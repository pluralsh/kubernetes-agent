#!/usr/bin/env bash

# Make sure version matches go.mod
exec go run go.uber.org/mock/mockgen@v0.3.0 "$@"
