/*
Copyright 2018 Turbine Labs, Inc.

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
	"fmt"
	"io/ioutil"
	"net/http"
)

// StatusCode matches an *http.Response with the given StatusCode.
func StatusCode(code int) Matcher {
	return &statusCode{code}
}

// Body matches an *http.Response with the given body bytes. If bytes
// is nil it will match an http.Response with Body == nil or a
// zero-length Body. The response's Body, if any, is closed after
// matching.
func Body(bytes []byte) Matcher {
	return &body{bytes}
}

type statusCode struct {
	statusCode int
}

func (s *statusCode) Matches(i interface{}) bool {
	if r, ok := i.(*http.Response); ok {
		return r.StatusCode == s.statusCode
	}

	return false
}

func (s *statusCode) String() string {
	return fmt.Sprintf("is an HTTP response with status code %d", s.statusCode)
}

type body struct {
	body []byte
}

func (b *body) Matches(i interface{}) bool {
	if r, ok := i.(*http.Response); ok {
		if r.Body == nil {
			return len(b.body) == 0
		}

		defer r.Body.Close()
		bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return false
		}

		if len(bytes) != len(b.body) {
			return false
		}

		for i := range bytes {
			if bytes[i] != b.body[i] {
				return false
			}
		}

		return true
	}

	return false
}

func (b *body) String() string {
	if b.body == nil {
		return "is an HTTP response with nil Body"
	}
	return fmt.Sprintf("is an HTTP response with body %q", string(b.body))
}
