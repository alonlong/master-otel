// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.7
// source: proto/stored/v1/stored.proto

package stored

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

// StoredServiceClient is the client API for StoredService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StoredServiceClient interface {
	CreateUser(ctx context.Context, in *v1.User, opts ...grpc.CallOption) (*v1.User, error)
	DeleteUser(ctx context.Context, in *v1.Identity, opts ...grpc.CallOption) (*v1.Empty, error)
}

type storedServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewStoredServiceClient(cc grpc.ClientConnInterface) StoredServiceClient {
	return &storedServiceClient{cc}
}

func (c *storedServiceClient) CreateUser(ctx context.Context, in *v1.User, opts ...grpc.CallOption) (*v1.User, error) {
	out := new(v1.User)
	err := c.cc.Invoke(ctx, "/proto.stored.v1.StoredService/CreateUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storedServiceClient) DeleteUser(ctx context.Context, in *v1.Identity, opts ...grpc.CallOption) (*v1.Empty, error) {
	out := new(v1.Empty)
	err := c.cc.Invoke(ctx, "/proto.stored.v1.StoredService/DeleteUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StoredServiceServer is the server API for StoredService service.
// All implementations should embed UnimplementedStoredServiceServer
// for forward compatibility
type StoredServiceServer interface {
	CreateUser(context.Context, *v1.User) (*v1.User, error)
	DeleteUser(context.Context, *v1.Identity) (*v1.Empty, error)
}

// UnimplementedStoredServiceServer should be embedded to have forward compatible implementations.
type UnimplementedStoredServiceServer struct {
}

func (UnimplementedStoredServiceServer) CreateUser(context.Context, *v1.User) (*v1.User, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUser not implemented")
}
func (UnimplementedStoredServiceServer) DeleteUser(context.Context, *v1.Identity) (*v1.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteUser not implemented")
}

// UnsafeStoredServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StoredServiceServer will
// result in compilation errors.
type UnsafeStoredServiceServer interface {
	mustEmbedUnimplementedStoredServiceServer()
}

func RegisterStoredServiceServer(s grpc.ServiceRegistrar, srv StoredServiceServer) {
	s.RegisterService(&StoredService_ServiceDesc, srv)
}

func _StoredService_CreateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v1.User)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StoredServiceServer).CreateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.stored.v1.StoredService/CreateUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StoredServiceServer).CreateUser(ctx, req.(*v1.User))
	}
	return interceptor(ctx, in, info, handler)
}

func _StoredService_DeleteUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v1.Identity)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StoredServiceServer).DeleteUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.stored.v1.StoredService/DeleteUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StoredServiceServer).DeleteUser(ctx, req.(*v1.Identity))
	}
	return interceptor(ctx, in, info, handler)
}

// StoredService_ServiceDesc is the grpc.ServiceDesc for StoredService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var StoredService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.stored.v1.StoredService",
	HandlerType: (*StoredServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateUser",
			Handler:    _StoredService_CreateUser_Handler,
		},
		{
			MethodName: "DeleteUser",
			Handler:    _StoredService_DeleteUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/stored/v1/stored.proto",
}
