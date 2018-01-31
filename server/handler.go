package server

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	// TestServerIDHeader is the name of an HTTP header containing the
	// test server's ID.
	TestServerIDHeader = "TestServer-ID"

	// TestServerForceResponseCode is the name of an HTTP query
	// parameter indicating the HTTP response code that should set for
	// the request.
	TestServerForceResponseCode = "force-response-code"

	// TestServerEchoHeadersWithPrefix is the name of an HTTP query
	// parameter that causes HTTP headers in the request starting with
	// given prefix (the parameter's value) to be echoed in the
	// response. May be repeated to render multiple prefixes.
	TestServerEchoHeadersWithPrefix = "echo-headers-with-prefix"
)

// TestHandler is an http.Handler that implements the TestServer.
type TestHandler struct {
	TestServer *TestServer
	ID         string
}

func (th TestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ts := th.TestServer

	w.Header().Set(TestServerIDHeader, th.ID)

	if ts.latencyMean > 0 {
		normLatency := time.Duration(
			(ts.rand.NormFloat64() * float64(ts.latencyStdDev)) + float64(ts.latencyMean),
		)
		if normLatency > 0 {
			ts.verbosef("sleeping for %s", normLatency)
			time.Sleep(normLatency)
		}
	}

	respCode := -1
	values := r.URL.Query()

	if va, ok := values[TestServerForceResponseCode]; ok {
		if len(va) >= 1 {
			c, err := strconv.Atoi(va[0])
			if err != nil {
				log.Printf("Could not parse %v arg %q", TestServerForceResponseCode, va[0])
			} else {
				respCode = c
			}
		}
	}

	respCodeOrDefault := func(defaultRespCode int) int {
		if respCode == -1 {
			return defaultRespCode
		}
		return respCode
	}

	if ts.errorRate > 0.0 && ts.rand.Float64()*100.0 < ts.errorRate {
		ts.verbosef("failing")
		http.Error(w, "oopsies", respCodeOrDefault(ts.errorStatus))
		return
	}

	ts.verbosef("succeeding")
	w.WriteHeader(respCodeOrDefault(200))
	fmt.Fprintf(w, "Hi there, I love %s\n", r.URL.Path[1:])

	if prefixes, ok := values[TestServerEchoHeadersWithPrefix]; ok {
		if len(prefixes) >= 1 {
			for k, v := range r.Header {
				for _, prefix := range prefixes {
					if strings.HasPrefix(strings.ToLower(k), strings.ToLower(prefix)) {
						fmt.Fprintf(w, "Header %s = %s\n", k, strings.Join(v, ", "))
					}
				}
			}
		}
	}
}
