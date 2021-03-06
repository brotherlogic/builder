// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// BuildClient is the client API for Build service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BuildClient interface {
	Refresh(ctx context.Context, in *RefreshRequest, opts ...grpc.CallOption) (*RefreshResponse, error)
}

type buildClient struct {
	cc grpc.ClientConnInterface
}

func NewBuildClient(cc grpc.ClientConnInterface) BuildClient {
	return &buildClient{cc}
}

func (c *buildClient) Refresh(ctx context.Context, in *RefreshRequest, opts ...grpc.CallOption) (*RefreshResponse, error) {
	out := new(RefreshResponse)
	err := c.cc.Invoke(ctx, "/builder.Build/Refresh", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BuildServer is the server API for Build service.
// All implementations should embed UnimplementedBuildServer
// for forward compatibility
type BuildServer interface {
	Refresh(context.Context, *RefreshRequest) (*RefreshResponse, error)
}

// UnimplementedBuildServer should be embedded to have forward compatible implementations.
type UnimplementedBuildServer struct {
}

func (UnimplementedBuildServer) Refresh(context.Context, *RefreshRequest) (*RefreshResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Refresh not implemented")
}

// UnsafeBuildServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BuildServer will
// result in compilation errors.
type UnsafeBuildServer interface {
	mustEmbedUnimplementedBuildServer()
}

func RegisterBuildServer(s grpc.ServiceRegistrar, srv BuildServer) {
	s.RegisterService(&_Build_serviceDesc, srv)
}

func _Build_Refresh_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RefreshRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BuildServer).Refresh(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/builder.Build/Refresh",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BuildServer).Refresh(ctx, req.(*RefreshRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Build_serviceDesc = grpc.ServiceDesc{
	ServiceName: "builder.Build",
	HandlerType: (*BuildServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Refresh",
			Handler:    _Build_Refresh_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "builder.proto",
}
