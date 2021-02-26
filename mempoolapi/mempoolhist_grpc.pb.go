// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package mempoolapi

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

// MempoolServiceClient is the client API for MempoolService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MempoolServiceClient interface {
	// You can get the confirmed raw transaction.
	GetMementry(ctx context.Context, in *TxidParam, opts ...grpc.CallOption) (*JsonResp, error)
	// You can get the confirmed mempool entry.
	GetMementryTime(ctx context.Context, in *TxidParam, opts ...grpc.CallOption) (*TimeResp, error)
	// You can get the confirmed tx and time.
	FindMempoolhist(ctx context.Context, in *TimerangeParam, opts ...grpc.CallOption) (*MemHistArray, error)
}

type mempoolServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMempoolServiceClient(cc grpc.ClientConnInterface) MempoolServiceClient {
	return &mempoolServiceClient{cc}
}

func (c *mempoolServiceClient) GetMementry(ctx context.Context, in *TxidParam, opts ...grpc.CallOption) (*JsonResp, error) {
	out := new(JsonResp)
	err := c.cc.Invoke(ctx, "/mempoolapi.MempoolService/GetMementry", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mempoolServiceClient) GetMementryTime(ctx context.Context, in *TxidParam, opts ...grpc.CallOption) (*TimeResp, error) {
	out := new(TimeResp)
	err := c.cc.Invoke(ctx, "/mempoolapi.MempoolService/GetMementryTime", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mempoolServiceClient) FindMempoolhist(ctx context.Context, in *TimerangeParam, opts ...grpc.CallOption) (*MemHistArray, error) {
	out := new(MemHistArray)
	err := c.cc.Invoke(ctx, "/mempoolapi.MempoolService/FindMempoolhist", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MempoolServiceServer is the server API for MempoolService service.
// All implementations must embed UnimplementedMempoolServiceServer
// for forward compatibility
type MempoolServiceServer interface {
	// You can get the confirmed raw transaction.
	GetMementry(context.Context, *TxidParam) (*JsonResp, error)
	// You can get the confirmed mempool entry.
	GetMementryTime(context.Context, *TxidParam) (*TimeResp, error)
	// You can get the confirmed tx and time.
	FindMempoolhist(context.Context, *TimerangeParam) (*MemHistArray, error)
	mustEmbedUnimplementedMempoolServiceServer()
}

// UnimplementedMempoolServiceServer must be embedded to have forward compatible implementations.
type UnimplementedMempoolServiceServer struct {
}

func (UnimplementedMempoolServiceServer) GetMementry(context.Context, *TxidParam) (*JsonResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMementry not implemented")
}
func (UnimplementedMempoolServiceServer) GetMementryTime(context.Context, *TxidParam) (*TimeResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMementryTime not implemented")
}
func (UnimplementedMempoolServiceServer) FindMempoolhist(context.Context, *TimerangeParam) (*MemHistArray, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindMempoolhist not implemented")
}
func (UnimplementedMempoolServiceServer) mustEmbedUnimplementedMempoolServiceServer() {}

// UnsafeMempoolServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MempoolServiceServer will
// result in compilation errors.
type UnsafeMempoolServiceServer interface {
	mustEmbedUnimplementedMempoolServiceServer()
}

func RegisterMempoolServiceServer(s grpc.ServiceRegistrar, srv MempoolServiceServer) {
	s.RegisterService(&MempoolService_ServiceDesc, srv)
}

func _MempoolService_GetMementry_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TxidParam)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MempoolServiceServer).GetMementry(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mempoolapi.MempoolService/GetMementry",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MempoolServiceServer).GetMementry(ctx, req.(*TxidParam))
	}
	return interceptor(ctx, in, info, handler)
}

func _MempoolService_GetMementryTime_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TxidParam)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MempoolServiceServer).GetMementryTime(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mempoolapi.MempoolService/GetMementryTime",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MempoolServiceServer).GetMementryTime(ctx, req.(*TxidParam))
	}
	return interceptor(ctx, in, info, handler)
}

func _MempoolService_FindMempoolhist_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TimerangeParam)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MempoolServiceServer).FindMempoolhist(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mempoolapi.MempoolService/FindMempoolhist",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MempoolServiceServer).FindMempoolhist(ctx, req.(*TimerangeParam))
	}
	return interceptor(ctx, in, info, handler)
}

// MempoolService_ServiceDesc is the grpc.ServiceDesc for MempoolService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MempoolService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "mempoolapi.MempoolService",
	HandlerType: (*MempoolServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetMementry",
			Handler:    _MempoolService_GetMementry_Handler,
		},
		{
			MethodName: "GetMementryTime",
			Handler:    _MempoolService_GetMementryTime_Handler,
		},
		{
			MethodName: "FindMempoolhist",
			Handler:    _MempoolService_FindMempoolhist_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "mempoolapi/mempoolhist.proto",
}
