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

package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/turbinelabs/test/assert"
)

func TestNewTestServer(t *testing.T) {
	ts, err := NewTestServer(
		[]string{"1234", "1234", "1235", "1234", "1236", "1235"},
		0.0,
		0,
		0,
		false,
	)
	assert.Nil(t, err)
	assert.ArrayEqual(t, ts.ports, []string{"1234", "1235", "1236"})

	ts, err = NewTestServer(
		[]string{"9999", "1000", "9999"},
		50.0,
		time.Second,
		time.Millisecond,
		true,
	)
	assert.Nil(t, err)
	assert.ArrayEqual(t, ts.ports, []string{"9999", "1000"})
	assert.Equal(t, ts.errorRate, 50.0)
	assert.Equal(t, ts.latencyMean, time.Second)
	assert.Equal(t, ts.latencyStdDev, time.Millisecond)
	assert.True(t, ts.verbose)
	assert.NonNil(t, ts.rand)

	ts, err = NewTestServer([]string{"1234"}, -1.0, 0, 0, false)
	assert.ErrorContains(t, err, "error rate must be between 0 and 100")
	assert.Nil(t, ts)
}

func TestNewTestServerWithDynamicPorts(t *testing.T) {
	ts, err := NewTestServerWithDynamicPorts([]string{"a", "b", "c"}, 0.0, 0, 0, false)
	assert.Nil(t, err)
	assert.ArrayEqual(t, ts.ports, []string{"0", "0", "0"})
	assert.ArrayEqual(t, ts.listenerIDs, []string{"a", "b", "c"})

	ts, err = NewTestServerWithDynamicPorts(
		[]string{"a", "b"},
		50.0,
		time.Second,
		time.Millisecond,
		true,
	)
	assert.Nil(t, err)
	assert.ArrayEqual(t, ts.ports, []string{"0", "0"})
	assert.ArrayEqual(t, ts.listenerIDs, []string{"a", "b"})
	assert.Equal(t, ts.errorRate, 50.0)
	assert.Equal(t, ts.latencyMean, time.Second)
	assert.Equal(t, ts.latencyStdDev, time.Millisecond)
	assert.True(t, ts.verbose)
	assert.NonNil(t, ts.rand)

	ts, err = NewTestServerWithDynamicPorts([]string{}, 0.0, 0, 0, false)
	assert.ErrorContains(t, err, "must specify at least one listener ID")
	assert.Nil(t, ts)

	ts, err = NewTestServerWithDynamicPorts([]string{"x", "x"}, 0.0, 0, 0, false)
	assert.ErrorContains(t, err, "listener IDs must be unique")
	assert.Nil(t, ts)

	ts, err = NewTestServerWithDynamicPorts([]string{"a"}, -1.0, 0, 0, false)
	assert.ErrorContains(t, err, "error rate must be between 0 and 100")
	assert.Nil(t, ts)
}

func TestTestServer(t *testing.T) {
	ts, err := NewTestServerWithDynamicPorts([]string{"MY-ID"}, 0.0, 0, 0, false)
	assert.Nil(t, err)
	assert.NonNil(t, ts)

	tsc := ts.ServeAsync()
	defer tsc.Stop()

	ports := tsc.IDPortMap()
	assert.Equal(t, len(ports), 1)
	for id, port := range ports {
		assert.Equal(t, id, "MY-ID")
		assert.NotEqual(t, port, 0)
	}

	resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/this%%20test", ports["MY-ID"]))
	if assert.Nil(t, err) {
		defer resp.Body.Close()
	}
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)

	assert.Equal(t, string(body), "Hi there, I love this test\n")
	assert.Equal(t, resp.Header.Get(TestServerIDHeader), "MY-ID")
}
