# reg-auth-service

<img src="https://github.com/eurofurence/reg-auth-service/actions/workflows/go.yml/badge.svg" alt="test status"/>

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

If you place this repository outside your GOPATH, build and test runs will download all required
dependencies by default.

## Running on localhost

Copy `docs/config.example.yaml` to `config.yaml` in the main project
directory and edit to match your local development environment.

Build using `go build cmd/main.go`.

Then run `./main -config config.yaml`.

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

## Limitations

The userinfo endpoint only works for configured applications whose cookie name matches the
global cookie name setting.

## Open Issues and Ideas

We track open issues as GitHub issues on this repository once it becomes clear what exactly needs to be done.
