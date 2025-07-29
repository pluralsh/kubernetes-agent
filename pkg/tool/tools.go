//go:build tools

package tool

import (
	_ "github.com/envoyproxy/protoc-gen-validate"
	_ "github.com/golangci/golangci-lint/v2/cmd/golangci-lint"
	_ "github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
