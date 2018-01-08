package server

import (
	"math/rand"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/turbinelabs/test/assert"
)

func TestHandlerReportsID(t *testing.T) {
	handler := TestHandler{
		TestServer: &TestServer{},
		ID:         "handler-id",
	}

	r := httptest.NewRequest("GET", "/foo", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	resp := w.Result()
	assert.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, resp.Header.Get(TestServerIDHeader), "handler-id")
	assert.Equal(t, w.Body.String(), "Hi there, I love foo\n")
}

func TestHandlerForceResponseCode(t *testing.T) {
	th := TestHandler{
		TestServer: &TestServer{},
		ID:         "1234",
	}

	r := httptest.NewRequest("GET", "/foo?force-response-code=599", nil)
	w := httptest.NewRecorder()

	th.ServeHTTP(w, r)

	resp := w.Result()
	assert.Equal(t, resp.StatusCode, 599)
	assert.Equal(t, w.Body.String(), "Hi there, I love foo\n")
}

func TestHandlerEchoHeadersWithPrefixOnSuccess(t *testing.T) {
	th := TestHandler{
		TestServer: &TestServer{},
		ID:         "1234",
	}

	r := httptest.NewRequest("GET", "/foo?echo-headers-with-prefix=x-", nil)
	r.Header.Add("x-show-me", "the-money")
	r.Header.Add("x-show-me", "state")
	r.Header.Add("y-me", "because")

	w := httptest.NewRecorder()

	th.ServeHTTP(w, r)

	resp := w.Result()
	assert.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, w.Body.String(), "Hi there, I love foo\nHeader X-Show-Me = the-money, state\n")
}

func TestHandlerIgnoreEchoHeadersWithPrefixOnFailure(t *testing.T) {
	th := TestHandler{
		TestServer: &TestServer{
			errorRate: 100.0,
			rand:      rand.New(rand.NewSource(1234)),
		},
		ID: "1234",
	}

	r := httptest.NewRequest("GET", "/foo?echo-headers-with-prefix=x-", nil)
	r.Header.Add("x-show-me", "the-money")

	w := httptest.NewRecorder()

	th.ServeHTTP(w, r)

	resp := w.Result()
	assert.Equal(t, resp.StatusCode, 503)
	assert.Equal(t, w.Body.String(), "oopsies\n")
}

func TestHandlerImplementsLatency(t *testing.T) {
	th := TestHandler{
		TestServer: &TestServer{
			latencyMean:   100 * time.Millisecond,
			latencyStdDev: 0, // no deviation
			rand:          rand.New(rand.NewSource(1234)),
		},
		ID: "latent",
	}

	r := httptest.NewRequest("GET", "/foo", nil)
	w := httptest.NewRecorder()

	start := time.Now()
	th.ServeHTTP(w, r)
	duration := time.Since(start)

	assert.GreaterThanEqual(t, duration, 100*time.Millisecond)
}

func TestHandlerImplementsLatencyWithStddev(t *testing.T) {
	rng := rand.New(rand.NewSource(1234))
	expectedDelay := time.Duration(rng.NormFloat64()*10.0+100.0) * time.Millisecond

	th := TestHandler{
		TestServer: &TestServer{
			latencyMean:   100 * time.Millisecond,
			latencyStdDev: 10 * time.Millisecond,
			rand:          rand.New(rand.NewSource(1234)),
		},
		ID: "latent",
	}

	r := httptest.NewRequest("GET", "/foo", nil)
	w := httptest.NewRecorder()

	start := time.Now()
	th.ServeHTTP(w, r)
	duration := time.Since(start)

	assert.GreaterThanEqual(t, duration, expectedDelay)
}

func TestHandlerImplementsErrorRate(t *testing.T) {
	rng := rand.New(rand.NewSource(1234))
	expectedFailures := []bool{}
	failures := 0
	for i := 0; i < 10; i++ {
		failure := rng.Float64()*100.0 < 50.0
		expectedFailures = append(expectedFailures, failure)
		if failure {
			failures++
		}
	}
	// at least one failure, at least one success
	assert.GreaterThan(t, failures, 0)
	assert.LessThan(t, failures, 10)

	th := TestHandler{
		TestServer: &TestServer{
			errorRate: 50.0,
			rand:      rand.New(rand.NewSource(1234)),
		},
		ID: "latent",
	}

	for _, expectedFailure := range expectedFailures {
		r := httptest.NewRequest("GET", "/foo", nil)
		w := httptest.NewRecorder()

		th.ServeHTTP(w, r)

		assert.Equal(t, w.Result().StatusCode != 200, expectedFailure)
	}
}
