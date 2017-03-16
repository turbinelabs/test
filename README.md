
[//]: # ( Copyright 2017 Turbine Labs, Inc.                                   )
[//]: # ( you may not use this file except in compliance with the License.    )
[//]: # ( You may obtain a copy of the License at                             )
[//]: # (                                                                     )
[//]: # (     http://www.apache.org/licenses/LICENSE-2.0                      )
[//]: # (                                                                     )
[//]: # ( Unless required by applicable law or agreed to in writing, software )
[//]: # ( distributed under the License is distributed on an "AS IS" BASIS,   )
[//]: # ( WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or     )
[//]: # ( implied. See the License for the specific language governing        )
[//]: # ( permissions and limitations under the License.                      )

# turbinelabs/test

[![Apache 2.0](https://img.shields.io/hexpm/l/plug.svg)](LICENSE)
[![GoDoc](https://godoc.org/github.com/turbinelabs/test?status.svg)](https://godoc.org/github.com/turbinelabs/test)
[![CircleCI](https://circleci.com/gh/turbinelabs/test.svg?style=shield)](https://circleci.com/gh/turbinelabs/test)
[![Go Report Card](https://goreportcard.com/badge/github.com/turbinelabs/test)](https://goreportcard.com/report/github.com/turbinelabs/test)

The test project contains a variety of small helper packages to make writing
tests in go a little easier. The test project has no external dependencies
beyond the go standard library.

## Requirements

- Go 1.8 or later (previous versions may work, but we don't build or test against them)

## Install

```
go get -u github.com/turbinelabs/test/...
```

## Clone/Test

```
mkdir -p $GOPATH/src/turbinelabs
git clone https://github.com/turbinelabs/test.git > $GOPATH/src/turbinelabs/test
go test github.com/turbinelabs/test/...
```

## Packages

Each package is best described in its respective Godoc:

- [`assert`](https://godoc.org/github.com/turbinelabs/test/assert):
  a collection of useful test assertions
- [`category`](https://godoc.org/github.com/turbinelabs/test/category):
  categorization and selective execution of tests
- [`check`](https://godoc.org/github.com/turbinelabs/test/check):
  simple boolean checks, some of which are used by `assert`
- [`http`](https://godoc.org/github.com/turbinelabs/test/http):
  useful extensions to the net/http/httptest package
- [`io`](https://godoc.org/github.com/turbinelabs/test/io):
  useful [`io.Reader`](https://golang.org/pkg/io/#Reader) and
  [`io.Writer`](https://golang.org/pkg/io/#Writer) implementations
- [`log`](https://godoc.org/github.com/turbinelabs/test/log):
  useful [`log.Logger`](https://golang.org/pkg/log/#Logger) factories
- [`matcher`](https://godoc.org/github.com/turbinelabs/test/matcher):
  useful [`gomock.Matcher`](https://godoc.org/github.com/golang/mock/gomock#Matcher)
  implementations
- [`server`](https://godoc.org/github.com/turbinelabs/test/server):
  a command-line configurable test HTTP server
- [`stack`](https://godoc.org/github.com/turbinelabs/test/stack):
  produces user friendly stack traces
- [`strings`](https://godoc.org/github.com/turbinelabs/test/strings):
  provides string conversions that are useful in tests
- [`tempfile`](https://godoc.org/github.com/turbinelabs/test/tempfile):
  provides wrappers around ioutil.TempFile to easily create temporary files or
  file names
- [`testrunner`](https://godoc.org/github.com/turbinelabs/test/testrunner):
  a testrunner for `go test` that generates a junit-style test result XML file

## Versioning

Please see [Versioning of Turbine Labs Open Source Projects](http://github.com/turbinelabs/developer/blob/master/README.md#versioning).

## Pull Requests

Patches accepted! Please see [Contributing to Turbine Labs Open Source Projects](http://github.com/turbinelabs/developer/blob/master/README.md#contributing).

## Code of Conduct

All Turbine Labs open-sourced projects are released with a
[Contributor Code of Conduct](CODE_OF_CONDUCT.md). By participating in our
projects you agree to abide by its terms, which will be carefully enforced.
