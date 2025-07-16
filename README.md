# insights-operator Utils

[![forthebadge made-with-go](http://ForTheBadge.com/images/badges/made-with-go.svg)](https://go.dev/)

[![GoDoc](https://godoc.org/github.com/RedHatInsights/insights-operator-utils?status.svg)](https://godoc.org/github.com/RedHatInsights/insights-operator-utils)
[![GitHub Pages](https://img.shields.io/badge/%20-GitHub%20Pages-informational)](https://redhatinsights.github.io/insights-operator-utils/)
[![Go Report Card](https://goreportcard.com/badge/github.com/RedHatInsights/insights-operator-utils)](https://goreportcard.com/report/github.com/RedHatInsights/insights-operator-utils)
[![Build Status](https://travis-ci.org/RedHatInsights/insights-operator-utils.svg?branch=master)](https://travis-ci.org/RedHatInsights/insights-operator-utils)
[![codecov](https://codecov.io/gh/RedHatInsights/insights-operator-utils/branch/master/graph/badge.svg)](https://codecov.io/gh/RedHatInsights/insights-operator-utils)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/RedHatInsights/insights-operator-utils)
[![License](https://img.shields.io/badge/license-Apache-blue)](https://github.com/RedHatInsights/insights-operator-utils/blob/master/LICENSE)

<!-- vim-markdown-toc GFM -->

- [insights-operator Utils](#insights-operator-utils)
  - [Description](#description)
  - [Sub-modules in this library](#sub-modules-in-this-library)
    - [`github.com/RedHatInsights/insights-operator-utils/collections`](#githubcomredhatinsightsinsights-operator-utilscollections)
    - [`github.com/RedHatInsights/insights-operator-utils/env`](#githubcomredhatinsightsinsights-operator-utilsenv)
    - [`github.com/RedHatInsights/insights-operator-utils/evaluator`](#githubcomredhatinsightsinsights-operator-utilsevaluator)
    - [`github.com/RedHatInsights/insights-operator-utils/generators`](#githubcomredhatinsightsinsights-operator-utilsgenerators)
    - [`github.com/RedHatInsights/insights-operator-utils/formatters`](#githubcomredhatinsightsinsights-operator-utilsformatters)
    - [`github.com/RedHatInsights/insights-operator-utils/http`](#githubcomredhatinsightsinsights-operator-utilshttp)
    - [`github.com/RedHatInsights/insights-operator-utils/logger`](#githubcomredhatinsightsinsights-operator-utilslogger)
    - [`github.com/RedHatInsights/insights-operator-utils/metrics`](#githubcomredhatinsightsinsights-operator-utilsmetrics)
    - [`github.com/RedHatInsights/insights-operator-utils/metrics/push`](#githubcomredhatinsightsinsights-operator-utilsmetricspush)
    - [`github.com/RedHatInsights/insights-operator-utils/migrations`](#githubcomredhatinsightsinsights-operator-utilsmigrations)
    - [`github.com/RedHatInsights/insights-operator-utils/parsers`](#githubcomredhatinsightsinsights-operator-utilsparsers)
    - [`github.com/RedHatInsights/insights-operator-utils/responses`](#githubcomredhatinsightsinsights-operator-utilsresponses)
    - [`github.com/RedHatInsights/insights-operator-utils/s3`](#githubcomredhatinsightsinsights-operator-utilss3)
    - [`github.com/RedHatInsights/insights-operator-utils/tls`](#githubcomredhatinsightsinsights-operator-utilstls)
    - [`github.com/RedHatInsights/insights-operator-utils/tests`](#githubcomredhatinsightsinsights-operator-utilstests)
    - [`github.com/RedHatInsights/insights-operator-utils/types`](#githubcomredhatinsightsinsights-operator-utilstypes)
  - [How to use this library](#how-to-use-this-library)
  - [Configuration](#configuration)
  - [Contribution](#contribution)
  - [Makefile targets](#makefile-targets)
  - [Testing](#testing)
  - [CI](#ci)
    - [Travis CI](#travis-ci)
    - [GolangCI](#golangci)
  - [Open Source Insights status](#open-source-insights-status)

<!-- vim-markdown-toc -->

## Description

Utility packages (written in Go) that are shared between different Insights
Operator, Insights Results Aggregator, and CCX Notification Service
repositories.

## Sub-modules in this library

### `github.com/RedHatInsights/insights-operator-utils/collections`

Helper functions to work with collections.

### `github.com/RedHatInsights/insights-operator-utils/env`

Functions to work with environment variables.

### `github.com/RedHatInsights/insights-operator-utils/evaluator`

Expression evaluator with ability to provide named values into expressions.

### `github.com/RedHatInsights/insights-operator-utils/generators`

Value generators - rule FQDNs etc.

### `github.com/RedHatInsights/insights-operator-utils/formatters`

Various text formatters utility functions.

### `github.com/RedHatInsights/insights-operator-utils/http`

HTTP-related utility functions.

### `github.com/RedHatInsights/insights-operator-utils/logger`

Configuration structures needed to configure the access to CloudWatch server to sending the log messages there.

### `github.com/RedHatInsights/insights-operator-utils/metrics`

Package metrics contains all metrics that needs to be exposed to Prometheus and indirectly to Grafana.

### `github.com/RedHatInsights/insights-operator-utils/metrics/push`

Package metrics/push contains some helping functions to push metrics to a Prometheus Pushgateway.

### `github.com/RedHatInsights/insights-operator-utils/migrations`

An implementation of a simple database migration mechanism that allows
semi-automatic transitions between various database versions as well as
building the latest version of the database from scratch.

### `github.com/RedHatInsights/insights-operator-utils/parsers`

Various text parser utility functions.

### `github.com/RedHatInsights/insights-operator-utils/responses`

Handlers for HTTP response.

### `github.com/RedHatInsights/insights-operator-utils/s3`

Helper functions to work with S3.

### `github.com/RedHatInsights/insights-operator-utils/tls`

Helper function to create [TLS configurations](https://pkg.go.dev/crypto/tls#Config).

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


## Makefile targets

```
Available targets are:

fmt                  Run go fmt -w for all sources
lint                 Run golint
vet                  Run go vet. Report likely mistakes in source code
cyclo                Run gocyclo
ineffassign          Run ineffassign checker
shellcheck           Run shellcheck
errcheck             Run errcheck
goconst              Run goconst checker
gosec                Run gosec checker
abcgo                Run ABC metrics checker
style                Run all the formatting related commands (fmt, vet, lint, cyclo) + check shell scripts
test                 Run the unit tests
benchmark            Run benchmarks
cover                Display test coverage on generated HTML pages
coverage             Display test coverage onto terminal
before_commit        Checks done before commit
license              Add license in every file in repository
help                 Show this help screen
```

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

### Travis CI

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

### GolangCI

Also GolangCI is configured for this repository and is run for all pull requests.

## Open Source Insights status

Open Source Insights status is available at [https://deps.dev/go/github.com%2Fredhatinsights%2Finsights-operator-utils/](https://deps.dev/go/github.com%2Fredhatinsights%2Finsights-operator-utils/)
