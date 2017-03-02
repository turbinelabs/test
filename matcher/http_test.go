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

package matcher

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/turbinelabs/test/assert"
	"github.com/turbinelabs/test/io"
)

func mkResponse(statusCode int, body *string) *http.Response {
	rr := httptest.NewRecorder()
	rr.WriteHeader(statusCode)
	if body != nil {
		rr.WriteString(*body)
	} else {
		rr.Body = nil
	}
	return rr.Result()
}

func ok() *http.Response                   { s := ""; return mkResponse(200, &s) }
func okBody(s string) *http.Response       { return mkResponse(200, &s) }
func okNoBody() *http.Response             { return mkResponse(200, nil) }
func notFound() *http.Response             { s := ""; return mkResponse(404, &s) }
func notFoundBody(s string) *http.Response { return mkResponse(404, &s) }

func TestStatusCode(t *testing.T) {
	assert.True(t, StatusCode(200).Matches(ok()))
	assert.True(t, StatusCode(404).Matches(notFound()))
	assert.False(t, StatusCode(200).Matches(notFound()))
	assert.False(t, StatusCode(200).Matches("not a response"))
}

func TestStatusCodeString(t *testing.T) {
	assert.Equal(t, StatusCode(200).String(), "is an HTTP response with status code 200")
}

func TestBody(t *testing.T) {
	assert.True(t, Body([]byte("stuff")).Matches(okBody("stuff")))
	assert.True(t, Body([]byte{}).Matches(okBody("")))
	assert.True(t, Body([]byte("404")).Matches(notFoundBody("404")))
	assert.True(t, Body(nil).Matches(ok()))
	assert.True(t, Body(nil).Matches(okBody("")))
	assert.True(t, Body(nil).Matches(notFound()))
	assert.True(t, Body(nil).Matches(okNoBody()))

	assert.False(t, Body([]byte("stuff")).Matches(okBody("STUFF")))
	assert.False(t, Body(nil).Matches(okBody("abc")))
	assert.False(t, Body([]byte{}).Matches(okBody("abc")))
	assert.False(t, Body(nil).Matches(notFoundBody("abc")))

	assert.False(t, Body(nil).Matches("not a response"))
}

func TestBodyBadReader(t *testing.T) {
	rr := httptest.NewRecorder()
	rr.WriteHeader(200)
	rw := rr.Result()
	rw.Body = io.NewFailingReader()

	assert.False(t, Body(nil).Matches(rr.Result()))
}

func TestBodyString(t *testing.T) {
	assert.Equal(t, Body([]byte("xyz")).String(), `is an HTTP response with body "xyz"`)
	assert.Equal(t, Body([]byte("")).String(), `is an HTTP response with body ""`)
	assert.Equal(t, Body(nil).String(), `is an HTTP response with nil Body`)
}
