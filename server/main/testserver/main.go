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

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/turbinelabs/test/server"
)

// serves HTTP on a specified port with a specified error rate and latency
// distribution
func main() {
	ts, err := server.NewTestServerFromFlagSet(flag.CommandLine, os.Args[1:])
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println("Usage: ")
		flag.CommandLine.PrintDefaults()
		os.Exit(1)
	}
	ts.Serve()
}
