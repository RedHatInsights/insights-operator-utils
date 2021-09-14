# insights-operator Utils

[![GoDoc](https://godoc.org/github.com/RedHatInsights/insights-operator-utils?status.svg)](https://godoc.org/github.com/RedHatInsights/insights-operator-utils)
[![GitHub Pages](https://img.shields.io/badge/%20-GitHub%20Pages-informational)](https://redhatinsights.github.io/insights-operator-utils/)
[![Go Report Card](https://goreportcard.com/badge/github.com/RedHatInsights/insights-operator-utils)](https://goreportcard.com/report/github.com/RedHatInsights/insights-operator-utils)
[![Build Status](https://travis-ci.org/RedHatInsights/insights-operator-utils.svg?branch=master)](https://travis-ci.org/RedHatInsights/insights-operator-utils)
[![codecov](https://codecov.io/gh/RedHatInsights/insights-operator-utils/branch/master/graph/badge.svg)](https://codecov.io/gh/RedHatInsights/insights-operator-utils)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/RedHatInsights/insights-operator-utils)
[![License](https://img.shields.io/badge/license-Apache-blue)](https://github.com/RedHatInsights/insights-operator-utils/blob/master/LICENSE)

<!-- vim-markdown-toc GFM -->

* [Description](#description)
* [Sub-modules in this library](#sub-modules-in-this-library)
    * [`github.com/RedHatInsights/insights-operator-utils/collections`](#githubcomredhatinsightsinsights-operator-utilscollections)
    * [`github.com/RedHatInsights/insights-operator-utils/env`](#githubcomredhatinsightsinsights-operator-utilsenv)
    * [`github.com/RedHatInsights/insights-operator-utils/http`](#githubcomredhatinsightsinsights-operator-utilshttp)
    * [`github.com/RedHatInsights/insights-operator-utils/logger`](#githubcomredhatinsightsinsights-operator-utilslogger)
    * [`github.com/RedHatInsights/insights-operator-utils/metrics`](#githubcomredhatinsightsinsights-operator-utilsmetrics)
    * [`github.com/RedHatInsights/insights-operator-utils/parsers`](#githubcomredhatinsightsinsights-operator-utilsparsers)
    * [`github.com/RedHatInsights/insights-operator-utils/responses`](#githubcomredhatinsightsinsights-operator-utilsresponses)
    * [`github.com/RedHatInsights/insights-operator-utils/tests`](#githubcomredhatinsightsinsights-operator-utilstests)
    * [`github.com/RedHatInsights/insights-operator-utils/types`](#githubcomredhatinsightsinsights-operator-utilstypes)
* [How to use this library](#how-to-use-this-library)
* [Configuration](#configuration)
* [Contribution](#contribution)
* [Testing](#testing)
* [CI](#ci)

<!-- vim-markdown-toc -->

## Description

Utils that are shared between different insights-operator repositories.

## Sub-modules in this library

### `github.com/RedHatInsights/insights-operator-utils/collections`

Helper functions to work with collections.

### `github.com/RedHatInsights/insights-operator-utils/env`

Functions to work with environment variables.

### `github.com/RedHatInsights/insights-operator-utils/http`

HTTP-related utility functions.

### `github.com/RedHatInsights/insights-operator-utils/logger`

Configuration structures needed to configure the access to CloudWatch server to sending the log messages there.

### `github.com/RedHatInsights/insights-operator-utils/metrics`

Package metrics contains all metrics that needs to be exposed to Prometheus and indirectly to Grafana.

### `github.com/RedHatInsights/insights-operator-utils/parsers`

Various text parser utility functions.

### `github.com/RedHatInsights/insights-operator-utils/responses`

Handlers for HTTP response.

### `github.com/RedHatInsights/insights-operator-utils/tests`

Contains sub-modules to make unit tests easier to write.

### `github.com/RedHatInsights/insights-operator-utils/types`

Declaration of various data types (usually structures) used elsewhere in the aggregator code.



## How to use this library

Use selected sub-module from this library in your `import` statement. For example:

```
import (
	"encoding/json"
	"strings"

	"github.com/RedHatInsights/insights-operator-utils/types"
	"github.com/RedHatInsights/insights-results-aggregator-data/testdata"
)
```

or:

```
import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/RedHatInsights/insights-operator-utils/logger"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)
```

## Configuration

No further configuration is needed at this moment.

## Contribution

Please look into document [CONTRIBUTING.md](CONTRIBUTING.md) that contains all information about how to contribute to this project.

Please look also at [Definition of Done](DoD.md) document with further informations.


## Testing

Unit tests can be started by the following command:

```
./test.sh
```

It is also possible to specify CLI options for Go test. For example, if you need to disable test results caching, use the following command:

```
./test -count=1
```

## CI

[Travis CI](https://travis-ci.com/) is configured for this repository. Several tests and checks are started for all pull requests:

* Unit tests that use the standard tool `go test`
* `go fmt` tool to check code formatting. That tool is run with `-s` flag to perform [following transformations](https://golang.org/cmd/gofmt/#hdr-The_simplify_command)
* `go vet` to report likely mistakes in source code, for example suspicious constructs, such as Printf calls whose arguments do not align with the format string.
* `golint` as a linter for all Go sources stored in this repository
* `gocyclo` to report all functions and methods with too high cyclomatic complexity. The cyclomatic complexity of a function is calculated according to the following rules: 1 is the base complexity of a function +1 for each 'if', 'for', 'case', '&&' or '||' Go Report Card warns on functions with cyclomatic complexity > 9
* `goconst` to find repeated strings that could be replaced by a constant
* `gosec` to inspect source code for security problems by scanning the Go AST
* `ineffassign` to detect and print all ineffectual assignments in Go code
* `errcheck` for checking for all unchecked errors in go programs
* `shellcheck` to perform static analysis for all shell scripts used in this repository
* `abcgo` to measure ABC metrics for Go source code and check if the metrics does not exceed specified threshold

Please note that all checks mentioned above have to pass for the change to be merged into the main branch (look into Settings to check which branch has been set as main one).

History of checks performed by CI is available at [RedHatInsights / insights-operator-utils](https://travis-ci.org/RedHatInsights/insights-operator-utils).
