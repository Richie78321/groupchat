// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.6.1
// source: chatservice/replication_service.proto

package chatservice

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

// ReplicationServiceClient is the client API for ReplicationService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ReplicationServiceClient interface {
	SubscribeUpdates(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (ReplicationService_SubscribeUpdatesClient, error)
}

type replicationServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewReplicationServiceClient(cc grpc.ClientConnInterface) ReplicationServiceClient {
	return &replicationServiceClient{cc}
}

func (c *replicationServiceClient) SubscribeUpdates(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (ReplicationService_SubscribeUpdatesClient, error) {
	stream, err := c.cc.NewStream(ctx, &ReplicationService_ServiceDesc.Streams[0], "/chatservice.ReplicationService/subscribe_updates", opts...)
	if err != nil {
		return nil, err
	}
	x := &replicationServiceSubscribeUpdatesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ReplicationService_SubscribeUpdatesClient interface {
	Recv() (*SubscriptionUpdate, error)
	grpc.ClientStream
}

type replicationServiceSubscribeUpdatesClient struct {
	grpc.ClientStream
}

func (x *replicationServiceSubscribeUpdatesClient) Recv() (*SubscriptionUpdate, error) {
	m := new(SubscriptionUpdate)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ReplicationServiceServer is the server API for ReplicationService service.
// All implementations must embed UnimplementedReplicationServiceServer
// for forward compatibility
type ReplicationServiceServer interface {
	SubscribeUpdates(*SubscribeRequest, ReplicationService_SubscribeUpdatesServer) error
	mustEmbedUnimplementedReplicationServiceServer()
}

// UnimplementedReplicationServiceServer must be embedded to have forward compatible implementations.
type UnimplementedReplicationServiceServer struct {
}

func (UnimplementedReplicationServiceServer) SubscribeUpdates(*SubscribeRequest, ReplicationService_SubscribeUpdatesServer) error {
	return status.Errorf(codes.Unimplemented, "method SubscribeUpdates not implemented")
}
func (UnimplementedReplicationServiceServer) mustEmbedUnimplementedReplicationServiceServer() {}

// UnsafeReplicationServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ReplicationServiceServer will
// result in compilation errors.
type UnsafeReplicationServiceServer interface {
	mustEmbedUnimplementedReplicationServiceServer()
}

func RegisterReplicationServiceServer(s grpc.ServiceRegistrar, srv ReplicationServiceServer) {
	s.RegisterService(&ReplicationService_ServiceDesc, srv)
}

func _ReplicationService_SubscribeUpdates_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SubscribeRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ReplicationServiceServer).SubscribeUpdates(m, &replicationServiceSubscribeUpdatesServer{stream})
}

type ReplicationService_SubscribeUpdatesServer interface {
	Send(*SubscriptionUpdate) error
	grpc.ServerStream
}

type replicationServiceSubscribeUpdatesServer struct {
	grpc.ServerStream
}

func (x *replicationServiceSubscribeUpdatesServer) Send(m *SubscriptionUpdate) error {
	return x.ServerStream.SendMsg(m)
}

// ReplicationService_ServiceDesc is the grpc.ServiceDesc for ReplicationService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ReplicationService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "chatservice.ReplicationService",
	HandlerType: (*ReplicationServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "subscribe_updates",
			Handler:       _ReplicationService_SubscribeUpdates_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "chatservice/replication_service.proto",
}
