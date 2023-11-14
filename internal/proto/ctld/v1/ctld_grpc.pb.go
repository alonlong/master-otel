// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.7
// source: proto/ctld/v1/ctld.proto

package ctld

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	v1 "master-otel/internal/proto/common/v1"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// CtldServiceClient is the client API for CtldService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CtldServiceClient interface {
	CreateUser(ctx context.Context, in *v1.User, opts ...grpc.CallOption) (*v1.User, error)
	CreateEmail(ctx context.Context, in *v1.Email, opts ...grpc.CallOption) (*v1.Email, error)
}

type ctldServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCtldServiceClient(cc grpc.ClientConnInterface) CtldServiceClient {
	return &ctldServiceClient{cc}
}

func (c *ctldServiceClient) CreateUser(ctx context.Context, in *v1.User, opts ...grpc.CallOption) (*v1.User, error) {
	out := new(v1.User)
	err := c.cc.Invoke(ctx, "/proto.ctld.v1.CtldService/CreateUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ctldServiceClient) CreateEmail(ctx context.Context, in *v1.Email, opts ...grpc.CallOption) (*v1.Email, error) {
	out := new(v1.Email)
	err := c.cc.Invoke(ctx, "/proto.ctld.v1.CtldService/CreateEmail", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CtldServiceServer is the server API for CtldService service.
// All implementations should embed UnimplementedCtldServiceServer
// for forward compatibility
type CtldServiceServer interface {
	CreateUser(context.Context, *v1.User) (*v1.User, error)
	CreateEmail(context.Context, *v1.Email) (*v1.Email, error)
}

// UnimplementedCtldServiceServer should be embedded to have forward compatible implementations.
type UnimplementedCtldServiceServer struct {
}

func (UnimplementedCtldServiceServer) CreateUser(context.Context, *v1.User) (*v1.User, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUser not implemented")
}
func (UnimplementedCtldServiceServer) CreateEmail(context.Context, *v1.Email) (*v1.Email, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateEmail not implemented")
}

// UnsafeCtldServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CtldServiceServer will
// result in compilation errors.
type UnsafeCtldServiceServer interface {
	mustEmbedUnimplementedCtldServiceServer()
}

func RegisterCtldServiceServer(s grpc.ServiceRegistrar, srv CtldServiceServer) {
	s.RegisterService(&CtldService_ServiceDesc, srv)
}

func _CtldService_CreateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v1.User)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CtldServiceServer).CreateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.ctld.v1.CtldService/CreateUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CtldServiceServer).CreateUser(ctx, req.(*v1.User))
	}
	return interceptor(ctx, in, info, handler)
}

func _CtldService_CreateEmail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v1.Email)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CtldServiceServer).CreateEmail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.ctld.v1.CtldService/CreateEmail",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CtldServiceServer).CreateEmail(ctx, req.(*v1.Email))
	}
	return interceptor(ctx, in, info, handler)
}

// CtldService_ServiceDesc is the grpc.ServiceDesc for CtldService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CtldService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.ctld.v1.CtldService",
	HandlerType: (*CtldServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateUser",
			Handler:    _CtldService_CreateUser_Handler,
		},
		{
			MethodName: "CreateEmail",
			Handler:    _CtldService_CreateEmail_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/ctld/v1/ctld.proto",
}
