// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.7
// source: internal/module/reverse_tunnel/rpc/rpc.proto

package rpc

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

// ReverseTunnelClient is the client API for ReverseTunnel service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ReverseTunnelClient interface {
	Connect(ctx context.Context, opts ...grpc.CallOption) (ReverseTunnel_ConnectClient, error)
}

type reverseTunnelClient struct {
	cc grpc.ClientConnInterface
}

func NewReverseTunnelClient(cc grpc.ClientConnInterface) ReverseTunnelClient {
	return &reverseTunnelClient{cc}
}

func (c *reverseTunnelClient) Connect(ctx context.Context, opts ...grpc.CallOption) (ReverseTunnel_ConnectClient, error) {
	stream, err := c.cc.NewStream(ctx, &ReverseTunnel_ServiceDesc.Streams[0], "/gitlab.agent.reverse_tunnel.rpc.ReverseTunnel/Connect", opts...)
	if err != nil {
		return nil, err
	}
	x := &reverseTunnelConnectClient{stream}
	return x, nil
}

type ReverseTunnel_ConnectClient interface {
	Send(*ConnectRequest) error
	Recv() (*ConnectResponse, error)
	grpc.ClientStream
}

type reverseTunnelConnectClient struct {
	grpc.ClientStream
}

func (x *reverseTunnelConnectClient) Send(m *ConnectRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *reverseTunnelConnectClient) Recv() (*ConnectResponse, error) {
	m := new(ConnectResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ReverseTunnelServer is the server API for ReverseTunnel service.
// All implementations must embed UnimplementedReverseTunnelServer
// for forward compatibility
type ReverseTunnelServer interface {
	Connect(ReverseTunnel_ConnectServer) error
	mustEmbedUnimplementedReverseTunnelServer()
}

// UnimplementedReverseTunnelServer must be embedded to have forward compatible implementations.
type UnimplementedReverseTunnelServer struct {
}

func (UnimplementedReverseTunnelServer) Connect(ReverseTunnel_ConnectServer) error {
	return status.Errorf(codes.Unimplemented, "method Connect not implemented")
}
func (UnimplementedReverseTunnelServer) mustEmbedUnimplementedReverseTunnelServer() {}

// UnsafeReverseTunnelServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ReverseTunnelServer will
// result in compilation errors.
type UnsafeReverseTunnelServer interface {
	mustEmbedUnimplementedReverseTunnelServer()
}

func RegisterReverseTunnelServer(s grpc.ServiceRegistrar, srv ReverseTunnelServer) {
	s.RegisterService(&ReverseTunnel_ServiceDesc, srv)
}

func _ReverseTunnel_Connect_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ReverseTunnelServer).Connect(&reverseTunnelConnectServer{stream})
}

type ReverseTunnel_ConnectServer interface {
	Send(*ConnectResponse) error
	Recv() (*ConnectRequest, error)
	grpc.ServerStream
}

type reverseTunnelConnectServer struct {
	grpc.ServerStream
}

func (x *reverseTunnelConnectServer) Send(m *ConnectResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *reverseTunnelConnectServer) Recv() (*ConnectRequest, error) {
	m := new(ConnectRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ReverseTunnel_ServiceDesc is the grpc.ServiceDesc for ReverseTunnel service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ReverseTunnel_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "gitlab.agent.reverse_tunnel.rpc.ReverseTunnel",
	HandlerType: (*ReverseTunnelServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Connect",
			Handler:       _ReverseTunnel_Connect_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "internal/module/reverse_tunnel/rpc/rpc.proto",
}
