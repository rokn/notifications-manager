// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.27.3
// source: api/channels/api.proto

package api

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	ChannelsService_GetChannel_FullMethodName      = "/ChannelsService/GetChannel"
	ChannelsService_GetChannelNames_FullMethodName = "/ChannelsService/GetChannelNames"
)

// ChannelsServiceClient is the client API for ChannelsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ChannelsServiceClient interface {
	GetChannel(ctx context.Context, in *GetChannelRequest, opts ...grpc.CallOption) (*GetChannelResponse, error)
	GetChannelNames(ctx context.Context, in *GetChannelNamesRequest, opts ...grpc.CallOption) (*GetChannelNamesResponse, error)
}

type channelsServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewChannelsServiceClient(cc grpc.ClientConnInterface) ChannelsServiceClient {
	return &channelsServiceClient{cc}
}

func (c *channelsServiceClient) GetChannel(ctx context.Context, in *GetChannelRequest, opts ...grpc.CallOption) (*GetChannelResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetChannelResponse)
	err := c.cc.Invoke(ctx, ChannelsService_GetChannel_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *channelsServiceClient) GetChannelNames(ctx context.Context, in *GetChannelNamesRequest, opts ...grpc.CallOption) (*GetChannelNamesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetChannelNamesResponse)
	err := c.cc.Invoke(ctx, ChannelsService_GetChannelNames_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ChannelsServiceServer is the server API for ChannelsService service.
// All implementations must embed UnimplementedChannelsServiceServer
// for forward compatibility.
type ChannelsServiceServer interface {
	GetChannel(context.Context, *GetChannelRequest) (*GetChannelResponse, error)
	GetChannelNames(context.Context, *GetChannelNamesRequest) (*GetChannelNamesResponse, error)
	mustEmbedUnimplementedChannelsServiceServer()
}

// UnimplementedChannelsServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedChannelsServiceServer struct{}

func (UnimplementedChannelsServiceServer) GetChannel(context.Context, *GetChannelRequest) (*GetChannelResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetChannel not implemented")
}
func (UnimplementedChannelsServiceServer) GetChannelNames(context.Context, *GetChannelNamesRequest) (*GetChannelNamesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetChannelNames not implemented")
}
func (UnimplementedChannelsServiceServer) mustEmbedUnimplementedChannelsServiceServer() {}
func (UnimplementedChannelsServiceServer) testEmbeddedByValue()                         {}

// UnsafeChannelsServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ChannelsServiceServer will
// result in compilation errors.
type UnsafeChannelsServiceServer interface {
	mustEmbedUnimplementedChannelsServiceServer()
}

func RegisterChannelsServiceServer(s grpc.ServiceRegistrar, srv ChannelsServiceServer) {
	// If the following call pancis, it indicates UnimplementedChannelsServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&ChannelsService_ServiceDesc, srv)
}

func _ChannelsService_GetChannel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetChannelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChannelsServiceServer).GetChannel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ChannelsService_GetChannel_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChannelsServiceServer).GetChannel(ctx, req.(*GetChannelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChannelsService_GetChannelNames_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetChannelNamesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChannelsServiceServer).GetChannelNames(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ChannelsService_GetChannelNames_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChannelsServiceServer).GetChannelNames(ctx, req.(*GetChannelNamesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ChannelsService_ServiceDesc is the grpc.ServiceDesc for ChannelsService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ChannelsService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ChannelsService",
	HandlerType: (*ChannelsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetChannel",
			Handler:    _ChannelsService_GetChannel_Handler,
		},
		{
			MethodName: "GetChannelNames",
			Handler:    _ChannelsService_GetChannelNames_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/channels/api.proto",
}
