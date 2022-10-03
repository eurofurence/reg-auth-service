# reg-auth-service

<img src="https://github.com/eurofurence/reg-auth-service/actions/workflows/go.yml/badge.svg" alt="test status"/>
<img src="https://github.com/eurofurence/reg-auth-service/actions/workflows/codeql-analysis.yml/badge.svg" alt="code quality status"/>

## Overview

Authentication bridge service between the Eurofurence registration system and an
[OpenID Connect](https://openid.net/developers/specs/) identity provider.

Implemented in go.

Command line arguments
```
-config <path-to-config-file> [-ecs-json-logging]
```

## Installation

This service uses go modules to provide dependency management, see `go.mod`.

If you place this repository outside of your GOPATH, build and test runs will download all required
dependencies by default.

## Running on localhost

Copy `docs/config.example.yaml` to `config.yaml` in the main project
directory and edit to match your local development environment.

Build using `go build cmd/main.go`.

Then run `./main -config config.yaml -migrate-database`.

## Installation on the server

See `install.sh`. This assumes a current build, and a valid configuration template in specific filenames.

## Test Coverage

In order to collect full test coverage, set go tool arguments to `-covermode=atomic -coverpkg=./internal/...`,
or manually run
```
go test -covermode=atomic -coverpkg=./internal/... ./...
```

## Acceptance Tests

We aim for good coverage with BDD-style acceptance tests. These will be the
best starting point to understanding what this service does.

## Contract Tests

This microservice uses [pact-go](https://github.com/pact-foundation/pact-go) for contract tests
of its consumption of our OIDC identity provider.

As described in the [pact-go installation instructions](https://github.com/pact-foundation/pact-go#installation),
you will need to have the [pact ruby standalone binaries installed](https://raw.githubusercontent.com/pact-foundation/pact-ruby-standalone/master/install.sh).
They provide the local mock that receives the calls made during contract testing.
