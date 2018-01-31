/*
Copyright 2017 Turbine Labs, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package server produces a binary called testserver, which serves HTTP
// on a specified port with a specified error rate and latency distribution.
package server

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	DefaultErrorStatus = 503
)

type closerChan chan struct{}

// TestServer represents one or more HTTP listeners.
type TestServer struct {
	ports         []string
	listenerIDs   []string
	errorStatus   int
	errorRate     float64
	latencyMean   time.Duration
	latencyStdDev time.Duration
	verbose       bool
	rand          *rand.Rand
}

// TestServerControl provides the ability to control a TestServer. It
// provides a mechanism for stopping the server and awaiting the
// termination of all listeners.
type TestServerControl struct {
	idPortMap map[string]int
	closer    closerChan
	waitgroup *sync.WaitGroup
}

// TestServer functions

func (ts *TestServer) logf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (ts *TestServer) verbosef(format string, v ...interface{}) {
	if ts.verbose {
		ts.logf(format, v...)
	}
}

func (ts *TestServer) serveListener(
	addr,
	listenerID string,
	listener net.Listener,
	wg *sync.WaitGroup,
) {
	defer func() {
		ts.logf("signaling completion for %s", addr)
		wg.Done()
	}()

	serveMux := http.NewServeMux()
	server := http.Server{Addr: addr, Handler: serveMux}
	th := TestHandler{ts, listenerID}
	serveMux.Handle("/", th)

	err := server.Serve(listener)
	if err != nil {
		ts.logf("failed to serve HTTP for %s: %v", addr, err)
	}
	ts.logf("server on port %s exited\n", addr)
}

func (ts *TestServer) closeListenerOnMessage(closer closerChan, listener net.Listener) {
	ok := true
	for ok {
		_, ok = <-closer
	}
	ts.logf("closing listener for %s", listener.Addr())
	listener.Close()
}

// SetErrorStatus configures the error code returned when the server's
// error rate is non-zero. It defaults to 503 (service
// unavailable). The value must be at least 500 and less than 600.
func (ts *TestServer) SetErrorStatus(code int) error {
	if code < 400 || code >= 600 {
		return fmt.Errorf("status code %d: out of range", code)
	}

	ts.errorStatus = code
	return nil
}

// ServeAsync starts the configured listeners for this TestServer and
// returns a TestServerControl which may be used to stop the listeners
// at a later point in time.
func (ts *TestServer) ServeAsync() *TestServerControl {
	if len(ts.ports) != len(ts.listenerIDs) {
		// Only tests should be able to cause this since these fields are not public.
		panic("failed invariant: list of ports and listener IDs must be the same length")
	}

	closer := closerChan(make(chan struct{}))

	wg := &sync.WaitGroup{}
	wg.Add(len(ts.ports))

	idPortMap := map[string]int{}
	for idx, port := range ts.ports {
		addr := ":" + port
		listenerID := ts.listenerIDs[idx]

		listener, err := net.Listen("tcp", addr)
		if err != nil {
			ts.logf("failed to open listener for %s: %v", addr, err)
			wg.Done()
			continue
		}

		// Port may have been dynamically selected, so retrieve it.
		resolvedPort := listener.Addr().(*net.TCPAddr).Port
		idPortMap[listenerID] = resolvedPort

		addr = fmt.Sprintf(":%d", resolvedPort)
		ts.logf("launching server on port %s\n", addr)

		go ts.serveListener(addr, listenerID, listener, wg)
		go ts.closeListenerOnMessage(closer, listener)
	}
	ts.logf("servers started")

	return &TestServerControl{idPortMap, closer, wg}
}

// TestServerControl functions

// Stop halts the listeners and waits for their associated goroutines to exit.
func (tsc *TestServerControl) Stop() {
	log.Printf("stopping servers")
	close(tsc.closer)
	tsc.Await()
}

// Await waits for all listeners to exit.
func (tsc *TestServerControl) Await() {
	log.Printf("waiting for servers to stop")
	tsc.waitgroup.Wait()
}

// IDPortMap returns a map of the ports used by the TestServer to
// their respective TestServerIDHeader values. For hard-coded ports
// the ID is always the port prefixed with a colon.
func (tsc *TestServerControl) IDPortMap() map[string]int {
	return tsc.idPortMap
}

// NewTestServer creates a new TestServer with the given
// configuration. The error rate is expressed as a percentage and must
// be between 0 and 100, inclusive. Duplicate ports are ignored.
func NewTestServer(
	ports []string,
	errorRate float64,
	latencyMean time.Duration,
	latencyStdDev time.Duration,
	verbose bool,
) (*TestServer, error) {
	if errorRate < 0 || errorRate > 100 {
		return nil, fmt.Errorf("error rate must be between 0 and 100")
	}

	// dedupe in place: after the loop we're left with a slice all unique values at the front
	// in their original order so we re-slice to len(m) to get the unique values
	m := map[string]struct{}{}
	for _, v := range ports {
		if _, seen := m[v]; !seen {
			ports[len(m)] = v
			m[v] = struct{}{}
		}
	}
	ports = ports[:len(m)]
	listenerIDs := make([]string, len(ports))
	for i, port := range ports {
		listenerIDs[i] = ":" + port
	}

	ts := TestServer{
		ports,
		listenerIDs,
		DefaultErrorStatus,
		errorRate,
		latencyMean,
		latencyStdDev,
		verbose,
		mkRand(),
	}

	return &ts, nil
}

// NewTestServerWithDynamicPorts creates a new TestServer with the
// given configuration. The error rate is expressed as a percentage
// and must be between 0 and 100, inclusive. A port is selected for
// each listener dynamically. ListenerIDs must be unique. Responses
// from a given port will contain the TestServerIDHeader with the
// corresponding value from listenerIDs. A mapping of IDs to their
// ports can be obtained via the TestServerControl object returned
// from ServeAsync.
func NewTestServerWithDynamicPorts(
	listenerIDs []string,
	errorRate float64,
	latencyMean time.Duration,
	latencyStdDev time.Duration,
	verbose bool,
) (*TestServer, error) {
	if len(listenerIDs) == 0 {
		return nil, errors.New("must specify at least one listener ID")
	}

	seenMap := map[string]struct{}{}
	for _, id := range listenerIDs {
		if _, seen := seenMap[id]; seen {
			return nil, errors.New("listener IDs must be unique")
		}
		seenMap[id] = struct{}{}
	}

	if errorRate < 0 || errorRate > 100 {
		return nil, errors.New("error rate must be between 0 and 100")
	}

	ports := make([]string, len(listenerIDs))
	for i := range ports {
		ports[i] = "0"
	}

	ts := TestServer{
		ports,
		listenerIDs,
		DefaultErrorStatus,
		errorRate,
		latencyMean,
		latencyStdDev,
		verbose,
		mkRand(),
	}

	return &ts, nil
}

func mkRand() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano() ^ (int64(os.Getpid()) << 30)))
}
