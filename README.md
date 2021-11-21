# reg-auth-service

<img src="https://github.com/Fenrikur/reg-auth-service/actions/workflows/go.yml/badge.svg" alt="test status"/>
<img src="https://github.com/Fenrikur/reg-auth-service/actions/workflows/codeql-analysis.yml/badge.svg" alt="code quality status"/>

## Overview

Authentication service for the Eurofurence registration system.

## Installation

This service uses go modules to provide dependency management, see `go.mod`.

If you place this repository OUTSIDE of your gopath, `go build main.go` and `go test ./...` will download all required dependencies by default.
