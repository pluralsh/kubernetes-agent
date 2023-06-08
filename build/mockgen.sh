#!/usr/bin/env bash

# Make sure version matches go.mod
exec go run github.com/golang/mock/mockgen@v1.7.0-rc.1.0.20220812172401-5b455625bd2c "$@"
