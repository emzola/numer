// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.23.4
// source: activity-service/proto/activity.proto

package proto

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
	ActivityService_GetUserActivities_FullMethodName    = "/activity.ActivityService/GetUserActivities"
	ActivityService_GetInvoiceActivities_FullMethodName = "/activity.ActivityService/GetInvoiceActivities"
)

// ActivityServiceClient is the client API for ActivityService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ActivityServiceClient interface {
	GetUserActivities(ctx context.Context, in *GetUserActivitiesRequest, opts ...grpc.CallOption) (*GetUserActivitiesResponse, error)
	GetInvoiceActivities(ctx context.Context, in *GetInvoiceActivitiesRequest, opts ...grpc.CallOption) (*GetInvoiceActivitiesResponse, error)
}

type activityServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewActivityServiceClient(cc grpc.ClientConnInterface) ActivityServiceClient {
	return &activityServiceClient{cc}
}

func (c *activityServiceClient) GetUserActivities(ctx context.Context, in *GetUserActivitiesRequest, opts ...grpc.CallOption) (*GetUserActivitiesResponse, error) {
	out := new(GetUserActivitiesResponse)
	err := c.cc.Invoke(ctx, ActivityService_GetUserActivities_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *activityServiceClient) GetInvoiceActivities(ctx context.Context, in *GetInvoiceActivitiesRequest, opts ...grpc.CallOption) (*GetInvoiceActivitiesResponse, error) {
	out := new(GetInvoiceActivitiesResponse)
	err := c.cc.Invoke(ctx, ActivityService_GetInvoiceActivities_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ActivityServiceServer is the server API for ActivityService service.
// All implementations must embed UnimplementedActivityServiceServer
// for forward compatibility
type ActivityServiceServer interface {
	GetUserActivities(context.Context, *GetUserActivitiesRequest) (*GetUserActivitiesResponse, error)
	GetInvoiceActivities(context.Context, *GetInvoiceActivitiesRequest) (*GetInvoiceActivitiesResponse, error)
	mustEmbedUnimplementedActivityServiceServer()
}

// UnimplementedActivityServiceServer must be embedded to have forward compatible implementations.
type UnimplementedActivityServiceServer struct {
}

func (UnimplementedActivityServiceServer) GetUserActivities(context.Context, *GetUserActivitiesRequest) (*GetUserActivitiesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserActivities not implemented")
}
func (UnimplementedActivityServiceServer) GetInvoiceActivities(context.Context, *GetInvoiceActivitiesRequest) (*GetInvoiceActivitiesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetInvoiceActivities not implemented")
}
func (UnimplementedActivityServiceServer) mustEmbedUnimplementedActivityServiceServer() {}

// UnsafeActivityServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ActivityServiceServer will
// result in compilation errors.
type UnsafeActivityServiceServer interface {
	mustEmbedUnimplementedActivityServiceServer()
}

func RegisterActivityServiceServer(s grpc.ServiceRegistrar, srv ActivityServiceServer) {
	s.RegisterService(&ActivityService_ServiceDesc, srv)
}

func _ActivityService_GetUserActivities_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserActivitiesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ActivityServiceServer).GetUserActivities(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ActivityService_GetUserActivities_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ActivityServiceServer).GetUserActivities(ctx, req.(*GetUserActivitiesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ActivityService_GetInvoiceActivities_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetInvoiceActivitiesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ActivityServiceServer).GetInvoiceActivities(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ActivityService_GetInvoiceActivities_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ActivityServiceServer).GetInvoiceActivities(ctx, req.(*GetInvoiceActivitiesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ActivityService_ServiceDesc is the grpc.ServiceDesc for ActivityService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ActivityService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "activity.ActivityService",
	HandlerType: (*ActivityServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetUserActivities",
			Handler:    _ActivityService_GetUserActivities_Handler,
		},
		{
			MethodName: "GetInvoiceActivities",
			Handler:    _ActivityService_GetInvoiceActivities_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "activity-service/proto/activity.proto",
}
