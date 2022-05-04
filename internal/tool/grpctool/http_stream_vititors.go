package grpctool

import (
	"sync"

	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	HttpRequestHeaderFieldNumber  protoreflect.FieldNumber = 1
	HttpRequestDataFieldNumber    protoreflect.FieldNumber = 2
	HttpRequestTrailerFieldNumber protoreflect.FieldNumber = 3

	HttpResponseHeaderFieldNumber  protoreflect.FieldNumber = 1
	HttpResponseDataFieldNumber    protoreflect.FieldNumber = 2
	HttpResponseTrailerFieldNumber protoreflect.FieldNumber = 3
)

var (
	httpRequestSVOnce        sync.Once
	httpRequestStreamVisitor *StreamVisitor

	httpResponseSVOnce        sync.Once
	httpResponseStreamVisitor *StreamVisitor
)

func HttpRequestStreamVisitor() *StreamVisitor {
	httpRequestSVOnce.Do(func() {
		var err error
		httpRequestStreamVisitor, err = NewStreamVisitor(&HttpRequest{})
		if err != nil {
			panic(err) // this will never panic as long as the proto file is correct
		}
	})
	return httpRequestStreamVisitor
}

func HttpResponseStreamVisitor() *StreamVisitor {
	httpResponseSVOnce.Do(func() {
		var err error
		httpResponseStreamVisitor, err = NewStreamVisitor(&HttpResponse{})
		if err != nil {
			panic(err) // this will never panic as long as the proto file is correct
		}
	})
	return httpResponseStreamVisitor
}
