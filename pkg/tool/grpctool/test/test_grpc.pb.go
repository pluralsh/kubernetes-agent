// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.24.4
// source: pkg/tool/grpctool/test/test.proto

// If you make any changes make sure you run: make regenerate-proto

package test

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Testing_RequestResponse_FullMethodName          = "/plural.agent.grpctool.test.Testing/RequestResponse"
	Testing_StreamingRequestResponse_FullMethodName = "/plural.agent.grpctool.test.Testing/StreamingRequestResponse"
)

// TestingClient is the client API for Testing service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TestingClient interface {
	RequestResponse(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
	StreamingRequestResponse(ctx context.Context, opts ...grpc.CallOption) (Testing_StreamingRequestResponseClient, error)
}

type testingClient struct {
	cc grpc.ClientConnInterface
}

func NewTestingClient(cc grpc.ClientConnInterface) TestingClient {
	return &testingClient{cc}
}

func (c *testingClient) RequestResponse(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, Testing_RequestResponse_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *testingClient) StreamingRequestResponse(ctx context.Context, opts ...grpc.CallOption) (Testing_StreamingRequestResponseClient, error) {
	stream, err := c.cc.NewStream(ctx, &Testing_ServiceDesc.Streams[0], Testing_StreamingRequestResponse_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &testingStreamingRequestResponseClient{stream}
	return x, nil
}

type Testing_StreamingRequestResponseClient interface {
	Send(*Request) error
	Recv() (*Response, error)
	grpc.ClientStream
}

type testingStreamingRequestResponseClient struct {
	grpc.ClientStream
}

func (x *testingStreamingRequestResponseClient) Send(m *Request) error {
	return x.ClientStream.SendMsg(m)
}

func (x *testingStreamingRequestResponseClient) Recv() (*Response, error) {
	m := new(Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// TestingServer is the server API for Testing service.
// All implementations must embed UnimplementedTestingServer
// for forward compatibility
type TestingServer interface {
	RequestResponse(context.Context, *Request) (*Response, error)
	StreamingRequestResponse(Testing_StreamingRequestResponseServer) error
	mustEmbedUnimplementedTestingServer()
}

// UnimplementedTestingServer must be embedded to have forward compatible implementations.
type UnimplementedTestingServer struct {
}

func (UnimplementedTestingServer) RequestResponse(context.Context, *Request) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestResponse not implemented")
}
func (UnimplementedTestingServer) StreamingRequestResponse(Testing_StreamingRequestResponseServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamingRequestResponse not implemented")
}
func (UnimplementedTestingServer) mustEmbedUnimplementedTestingServer() {}

// UnsafeTestingServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TestingServer will
// result in compilation errors.
type UnsafeTestingServer interface {
	mustEmbedUnimplementedTestingServer()
}

func RegisterTestingServer(s grpc.ServiceRegistrar, srv TestingServer) {
	s.RegisterService(&Testing_ServiceDesc, srv)
}

func _Testing_RequestResponse_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TestingServer).RequestResponse(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Testing_RequestResponse_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TestingServer).RequestResponse(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Testing_StreamingRequestResponse_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(TestingServer).StreamingRequestResponse(&testingStreamingRequestResponseServer{stream})
}

type Testing_StreamingRequestResponseServer interface {
	Send(*Response) error
	Recv() (*Request, error)
	grpc.ServerStream
}

type testingStreamingRequestResponseServer struct {
	grpc.ServerStream
}

func (x *testingStreamingRequestResponseServer) Send(m *Response) error {
	return x.ServerStream.SendMsg(m)
}

func (x *testingStreamingRequestResponseServer) Recv() (*Request, error) {
	m := new(Request)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Testing_ServiceDesc is the grpc.ServiceDesc for Testing service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Testing_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "plural.agent.grpctool.test.Testing",
	HandlerType: (*TestingServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RequestResponse",
			Handler:    _Testing_RequestResponse_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamingRequestResponse",
			Handler:       _Testing_StreamingRequestResponse_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "pkg/tool/grpctool/test/test.proto",
}