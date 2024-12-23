// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.26.1
// source: user.proto

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

const (
	UserRPCService_CreateUser_FullMethodName = "/rpc.UserRPCService/CreateUser"
	UserRPCService_GetUser_FullMethodName    = "/rpc.UserRPCService/GetUser"
)

// UserRPCServiceClient is the client API for UserRPCService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserRPCServiceClient interface {
	CreateUser(ctx context.Context, in *NewUserMessage, opts ...grpc.CallOption) (*UserMessage, error)
	GetUser(ctx context.Context, in *IdMessage, opts ...grpc.CallOption) (*UserMessage, error)
}

type userRPCServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewUserRPCServiceClient(cc grpc.ClientConnInterface) UserRPCServiceClient {
	return &userRPCServiceClient{cc}
}

func (c *userRPCServiceClient) CreateUser(ctx context.Context, in *NewUserMessage, opts ...grpc.CallOption) (*UserMessage, error) {
	out := new(UserMessage)
	err := c.cc.Invoke(ctx, UserRPCService_CreateUser_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userRPCServiceClient) GetUser(ctx context.Context, in *IdMessage, opts ...grpc.CallOption) (*UserMessage, error) {
	out := new(UserMessage)
	err := c.cc.Invoke(ctx, UserRPCService_GetUser_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserRPCServiceServer is the server API for UserRPCService service.
// All implementations must embed UnimplementedUserRPCServiceServer
// for forward compatibility
type UserRPCServiceServer interface {
	CreateUser(context.Context, *NewUserMessage) (*UserMessage, error)
	GetUser(context.Context, *IdMessage) (*UserMessage, error)
	mustEmbedUnimplementedUserRPCServiceServer()
}

// UnimplementedUserRPCServiceServer must be embedded to have forward compatible implementations.
type UnimplementedUserRPCServiceServer struct {
}

func (UnimplementedUserRPCServiceServer) CreateUser(context.Context, *NewUserMessage) (*UserMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUser not implemented")
}
func (UnimplementedUserRPCServiceServer) GetUser(context.Context, *IdMessage) (*UserMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUser not implemented")
}
func (UnimplementedUserRPCServiceServer) mustEmbedUnimplementedUserRPCServiceServer() {}

// UnsafeUserRPCServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserRPCServiceServer will
// result in compilation errors.
type UnsafeUserRPCServiceServer interface {
	mustEmbedUnimplementedUserRPCServiceServer()
}

func RegisterUserRPCServiceServer(s grpc.ServiceRegistrar, srv UserRPCServiceServer) {
	s.RegisterService(&UserRPCService_ServiceDesc, srv)
}

func _UserRPCService_CreateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewUserMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserRPCServiceServer).CreateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserRPCService_CreateUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserRPCServiceServer).CreateUser(ctx, req.(*NewUserMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserRPCService_GetUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IdMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserRPCServiceServer).GetUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserRPCService_GetUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserRPCServiceServer).GetUser(ctx, req.(*IdMessage))
	}
	return interceptor(ctx, in, info, handler)
}

// UserRPCService_ServiceDesc is the grpc.ServiceDesc for UserRPCService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UserRPCService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "rpc.UserRPCService",
	HandlerType: (*UserRPCServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateUser",
			Handler:    _UserRPCService_CreateUser_Handler,
		},
		{
			MethodName: "GetUser",
			Handler:    _UserRPCService_GetUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "user.proto",
}
