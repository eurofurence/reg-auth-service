# reg-auth-service

<img src="https://github.com/eurofurence/reg-auth-service/actions/workflows/go.yml/badge.svg" alt="test status"/>
<img src="https://github.com/eurofurence/reg-auth-service/actions/workflows/codeql-analysis.yml/badge.svg" alt="code quality status"/>

## Overview

Authentication bridge service for the Eurofurence registration system.

## Installation

This service uses go modules to provide dependency management, see `go.mod`.

If you place this repository OUTSIDE of your gopath, `go build main.go` and 
`go test ./...` will download all required dependencies by default.

## Configuration

Copy `docs/config.example.yaml` to `config.yaml` in the main project
directory and edit to match your local development environment.

## Contract Tests

This microservice uses [pact-go](https://github.com/pact-foundation/pact-go) for contract tests
of its consumption of our OIDC identity provider.

As described in the [pact-go installation instructions](https://github.com/pact-foundation/pact-go#installation),
you will need to have the [pact ruby standalone binaries installed](https://raw.githubusercontent.com/pact-foundation/pact-ruby-standalone/master/install.sh).
They provide the local mock that receives the calls made during contract testing.

## Acceptance Tests

We aim for good coverage with BDD-style acceptance tests. These will be the
best starting point to understanding what this service does.

## Code Coverage

In order to include contract and acceptance tests in the coverage data, you must pass the
`-coverpkg=./...` argument to `go test`.

_If you use IntelliJ / GoLand, you can configure this as a
"Go Tool Argument" in the run configuration template for "Go Test"._
