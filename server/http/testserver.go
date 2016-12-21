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

package http

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	TestServerIdHeader = "TestServer-ID"
)

type TestServer struct {
	ports         []string
	errorRate     float64
	latencyMean   float64
	latencyStdDev float64
	verbose       bool
}

type TestHandler struct {
	TestServer *TestServer
	Port       string
}

type TestServerControl struct {
	closer    chan<- bool
	waitgroup *sync.WaitGroup
}

// TestHandler functions

func (th TestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ts := th.TestServer

	w.Header().Set(TestServerIdHeader, th.Port)

	if ts.latencyMean > 0.0 {
		normLatency := (rand.NormFloat64() * ts.latencyStdDev) + ts.latencyMean
		if normLatency > 0 {
			ts.verbosef("sleeping for %f", normLatency)
			time.Sleep(time.Duration(normLatency) * time.Millisecond)
		}
	}

	if (rand.Float64() * 100.0) < ts.errorRate {
		ts.verbosef("failing")
		http.Error(w, "oopsies", 503)
	} else {
		ts.verbosef("succeeding")
		fmt.Fprintf(w, "Hi there, I love %s\n", r.URL.Path[1:])
	}
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

func (ts *TestServer) serveListener(addr string, listener net.Listener, wg *sync.WaitGroup) {
	defer func() {
		ts.logf("signaling completion for %s", addr)
		wg.Done()
	}()

	serveMux := http.NewServeMux()
	server := http.Server{Addr: addr, Handler: serveMux}
	th := TestHandler{ts, addr}
	serveMux.Handle("/", th)

	err := server.Serve(listener)
	if err != nil {
		ts.logf("failed to serve HTTP for %s: %v", addr, err)
	}
	ts.logf("server on port %s exited\n", addr)
}

func (ts *TestServer) closeListenerOnMessage(closer <-chan bool, listener net.Listener) {
	ok := true
	for ok {
		_, ok = <-closer
	}
	ts.logf("closing listener for %s", listener.Addr())
	listener.Close()
}

func (ts *TestServer) ServeAsync() *TestServerControl {
	var wg sync.WaitGroup
	closer := make(chan bool)
	wg.Add(len(ts.ports))
	for _, port := range ts.ports {
		ts.logf("launching server on port %s\n", port)

		addr := ":" + port

		listener, err := net.Listen("tcp", addr)
		if err != nil {
			ts.logf("failed to open listener for %s: %v", addr, err)
			wg.Done()
			continue
		}

		go ts.serveListener(addr, listener, &wg)
		go ts.closeListenerOnMessage(closer, listener)
	}
	ts.logf("servers started")

	return &TestServerControl{closer, &wg}
}

func (ts *TestServer) Serve() {
	ts.ServeAsync().Await()
}

// TestServerControl functions

// Stops the listeners and waits for the associated goroutines to exit
func (tsc *TestServerControl) Stop() {
	log.Printf("stopping servers")
	close(tsc.closer)
	tsc.Await()
}

func (tsc *TestServerControl) Await() {
	log.Printf("waiting for servers to stop")
	tsc.waitgroup.Wait()
}

func NewTestServer(
	ports []string,
	errorRate float64,
	latencyMean float64,
	latencyStdDev float64,
	verbose bool,
) (*TestServer, error) {
	rand.Seed(time.Now().UnixNano() ^ (int64(os.Getpid()) << 30))

	if errorRate < 0 || errorRate > 100 {
		return nil, fmt.Errorf("errorRate should be between 0 and 100")
	}

	if latencyMean < 0 {
		return nil, fmt.Errorf("latencyMean should be non-negative")
	}

	if latencyStdDev < 0 {
		return nil, fmt.Errorf("latencyStdDev should be non-negative")
	}

	m := map[string]bool{}
	// dedup in place; after the loop we're left with a slice all unique values at the front
	// in their original order so we re-slice to len(m) to get the unique values
	for _, v := range ports {
		if _, seen := m[v]; !seen {
			ports[len(m)] = v
			m[v] = true
		}
	}
	ports = ports[:len(m)]

	ts := TestServer{ports, errorRate, latencyMean, latencyStdDev, verbose}

	return &ts, nil
}

func NewTestServerFromFlagSet(f *flag.FlagSet, args []string) (*TestServer, error) {
	portsFlag := f.String("ports", "8889", "comma separated list of ports to listen on")
	errorRateFlag := f.Float64("error-rate", 0.01, "error rate as a percentage (100.0 = pure error)")
	latencyMeanFlag := f.Float64("latency-mean", 4, "mean latency in milliseconds")
	latencyStdDevFlag := f.Float64("latency-stddev", 1, "latency standard deviation in milliseconds")
	verboseFlag := f.Bool("verbose", false, "enable verbose logging")
	f.Parse(args)

	return NewTestServer(
		strings.Split(*portsFlag, ","),
		*errorRateFlag,
		*latencyMeanFlag,
		*latencyStdDevFlag,
		*verboseFlag)
}
