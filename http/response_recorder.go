package http

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"testing"
)

// ResponseRecorder is an augmented version of net.http.httptest.ResponseRecorder
// with the ability to make various assertions concerning the HTTP response that
// was observed.
type ResponseRecorder struct {
	*httptest.ResponseRecorder            // embed the stdlib ResponseRecorder for API parity
	t                          *testing.T // the T that should be used to record errors
	bodyBuffer                 []byte     // holds the body response so that consecutive calls to GetBody work as expected
}

// Build a new ResponseRecorder that extends the net/http/httptest one by
// attaching some functions making it simple to assert some facts about the
// response or record an error if those assertions fail.
//
// Because the recording is parameterized by a single testing.T this should
// not be resude between test cases.
func NewResponseRecorder(t *testing.T) *ResponseRecorder {
	return &ResponseRecorder{
		&httptest.ResponseRecorder{Body: new(bytes.Buffer)},
		t,
		nil,
	}
}

// If the response did not have the indicated status this records a test error.
func (rr *ResponseRecorder) AssertStatus(want int) {
	if want != rr.Code {
		rr.t.Errorf("got: %d, want: %d", rr.Code, want)
	}
}

// If the response had any body content this will record a test error.
func (rr *ResponseRecorder) AssertNoBody() {
	if rr.Body.Len() != 0 {
		bstr := string(rr.GetBody())
		rr.t.Errorf("got: data written to response '%s', want: No data written to response", bstr)
	}
}

// Compares the wanted body against what was sent in the response and records
// an error if they don't match.
func (rr *ResponseRecorder) AssertBody(want string) {
	b := rr.GetBody()
	if want != string(b) {
		rr.t.Errorf("got: %s, want: %s", string(b), want)
	}
}

// Fetches the header name's value, compares to want and records
// an error if they don't match.
func (rr *ResponseRecorder) AssertHeader(name string, want string) {
	if got := rr.Header().Get(name); got != want {
		rr.t.Errorf("header %s: got: %s, want: %s", name, got, want)
	}
}

// Convenience method that renders an object as json and then compares that to
// the body sent in the response.  If they don't match a test error will be
// recorded.
func (rr *ResponseRecorder) AssertBodyJSON(want interface{}) {
	b, err := json.Marshal(want)
	if err != nil {
		log.Fatal(err)
	}
	rr.AssertBody(string(b))
}

// Convenience method that returns the contents of the body.
func (rr *ResponseRecorder) GetBody() []byte {
	if rr.bodyBuffer != nil {
		return rr.bodyBuffer
	}

	var err error
	rr.bodyBuffer, err = ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal("Failure reading response's recorded body")
	}

	return rr.bodyBuffer
}
