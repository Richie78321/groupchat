// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.6.1
// source: chatservice/chat_service.proto

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

// ChatServiceClient is the client API for ChatService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ChatServiceClient interface {
	SubscribeChatroom(ctx context.Context, in *SubscribeChatroomRequest, opts ...grpc.CallOption) (ChatService_SubscribeChatroomClient, error)
	SendChat(ctx context.Context, in *SendChatRequest, opts ...grpc.CallOption) (*SendChatResponse, error)
	LikeChat(ctx context.Context, in *LikeChatRequest, opts ...grpc.CallOption) (*LikeChatResponse, error)
	MessageHistory(ctx context.Context, in *MessageHistoryRequest, opts ...grpc.CallOption) (*MessageHistoryResponse, error)
	ViewPeers(ctx context.Context, in *ViewPeersRequest, opts ...grpc.CallOption) (*ViewPeersResponse, error)
}

type chatServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewChatServiceClient(cc grpc.ClientConnInterface) ChatServiceClient {
	return &chatServiceClient{cc}
}

func (c *chatServiceClient) SubscribeChatroom(ctx context.Context, in *SubscribeChatroomRequest, opts ...grpc.CallOption) (ChatService_SubscribeChatroomClient, error) {
	stream, err := c.cc.NewStream(ctx, &ChatService_ServiceDesc.Streams[0], "/chatservice.ChatService/subscribe_chatroom", opts...)
	if err != nil {
		return nil, err
	}
	x := &chatServiceSubscribeChatroomClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ChatService_SubscribeChatroomClient interface {
	Recv() (*ChatroomSubscriptionUpdate, error)
	grpc.ClientStream
}

type chatServiceSubscribeChatroomClient struct {
	grpc.ClientStream
}

func (x *chatServiceSubscribeChatroomClient) Recv() (*ChatroomSubscriptionUpdate, error) {
	m := new(ChatroomSubscriptionUpdate)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *chatServiceClient) SendChat(ctx context.Context, in *SendChatRequest, opts ...grpc.CallOption) (*SendChatResponse, error) {
	out := new(SendChatResponse)
	err := c.cc.Invoke(ctx, "/chatservice.ChatService/send_chat", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatServiceClient) LikeChat(ctx context.Context, in *LikeChatRequest, opts ...grpc.CallOption) (*LikeChatResponse, error) {
	out := new(LikeChatResponse)
	err := c.cc.Invoke(ctx, "/chatservice.ChatService/like_chat", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatServiceClient) MessageHistory(ctx context.Context, in *MessageHistoryRequest, opts ...grpc.CallOption) (*MessageHistoryResponse, error) {
	out := new(MessageHistoryResponse)
	err := c.cc.Invoke(ctx, "/chatservice.ChatService/message_history", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatServiceClient) ViewPeers(ctx context.Context, in *ViewPeersRequest, opts ...grpc.CallOption) (*ViewPeersResponse, error) {
	out := new(ViewPeersResponse)
	err := c.cc.Invoke(ctx, "/chatservice.ChatService/view_peers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ChatServiceServer is the server API for ChatService service.
// All implementations must embed UnimplementedChatServiceServer
// for forward compatibility
type ChatServiceServer interface {
	SubscribeChatroom(*SubscribeChatroomRequest, ChatService_SubscribeChatroomServer) error
	SendChat(context.Context, *SendChatRequest) (*SendChatResponse, error)
	LikeChat(context.Context, *LikeChatRequest) (*LikeChatResponse, error)
	MessageHistory(context.Context, *MessageHistoryRequest) (*MessageHistoryResponse, error)
	ViewPeers(context.Context, *ViewPeersRequest) (*ViewPeersResponse, error)
	mustEmbedUnimplementedChatServiceServer()
}

// UnimplementedChatServiceServer must be embedded to have forward compatible implementations.
type UnimplementedChatServiceServer struct {
}

func (UnimplementedChatServiceServer) SubscribeChatroom(*SubscribeChatroomRequest, ChatService_SubscribeChatroomServer) error {
	return status.Errorf(codes.Unimplemented, "method SubscribeChatroom not implemented")
}
func (UnimplementedChatServiceServer) SendChat(context.Context, *SendChatRequest) (*SendChatResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendChat not implemented")
}
func (UnimplementedChatServiceServer) LikeChat(context.Context, *LikeChatRequest) (*LikeChatResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LikeChat not implemented")
}
func (UnimplementedChatServiceServer) MessageHistory(context.Context, *MessageHistoryRequest) (*MessageHistoryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MessageHistory not implemented")
}
func (UnimplementedChatServiceServer) ViewPeers(context.Context, *ViewPeersRequest) (*ViewPeersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ViewPeers not implemented")
}
func (UnimplementedChatServiceServer) mustEmbedUnimplementedChatServiceServer() {}

// UnsafeChatServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ChatServiceServer will
// result in compilation errors.
type UnsafeChatServiceServer interface {
	mustEmbedUnimplementedChatServiceServer()
}

func RegisterChatServiceServer(s grpc.ServiceRegistrar, srv ChatServiceServer) {
	s.RegisterService(&ChatService_ServiceDesc, srv)
}

func _ChatService_SubscribeChatroom_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SubscribeChatroomRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ChatServiceServer).SubscribeChatroom(m, &chatServiceSubscribeChatroomServer{stream})
}

type ChatService_SubscribeChatroomServer interface {
	Send(*ChatroomSubscriptionUpdate) error
	grpc.ServerStream
}

type chatServiceSubscribeChatroomServer struct {
	grpc.ServerStream
}

func (x *chatServiceSubscribeChatroomServer) Send(m *ChatroomSubscriptionUpdate) error {
	return x.ServerStream.SendMsg(m)
}

func _ChatService_SendChat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendChatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServiceServer).SendChat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chatservice.ChatService/send_chat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServiceServer).SendChat(ctx, req.(*SendChatRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatService_LikeChat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LikeChatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServiceServer).LikeChat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chatservice.ChatService/like_chat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServiceServer).LikeChat(ctx, req.(*LikeChatRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatService_MessageHistory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MessageHistoryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServiceServer).MessageHistory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chatservice.ChatService/message_history",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServiceServer).MessageHistory(ctx, req.(*MessageHistoryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatService_ViewPeers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ViewPeersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServiceServer).ViewPeers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chatservice.ChatService/view_peers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServiceServer).ViewPeers(ctx, req.(*ViewPeersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ChatService_ServiceDesc is the grpc.ServiceDesc for ChatService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ChatService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "chatservice.ChatService",
	HandlerType: (*ChatServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "send_chat",
			Handler:    _ChatService_SendChat_Handler,
		},
		{
			MethodName: "like_chat",
			Handler:    _ChatService_LikeChat_Handler,
		},
		{
			MethodName: "message_history",
			Handler:    _ChatService_MessageHistory_Handler,
		},
		{
			MethodName: "view_peers",
			Handler:    _ChatService_ViewPeers_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "subscribe_chatroom",
			Handler:       _ChatService_SubscribeChatroom_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "chatservice/chat_service.proto",
}
