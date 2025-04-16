// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.19.6
// source: grpc/chatapp.proto

package grpc

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
	Chat_RoomChat_FullMethodName           = "/chat.Chat/RoomChat"
	Chat_SendPrivateMessage_FullMethodName = "/chat.Chat/SendPrivateMessage"
	Chat_LeaveRoom_FullMethodName          = "/chat.Chat/LeaveRoom"
)

// ChatClient is the client API for Chat service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ChatClient interface {
	RoomChat(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[ChatRoomMessage, ChatRoomMessage], error)
	SendPrivateMessage(ctx context.Context, in *PrivateMessage, opts ...grpc.CallOption) (*MessageResponse, error)
	LeaveRoom(ctx context.Context, in *LeaveRequest, opts ...grpc.CallOption) (*MessageResponse, error)
}

type chatClient struct {
	cc grpc.ClientConnInterface
}

func NewChatClient(cc grpc.ClientConnInterface) ChatClient {
	return &chatClient{cc}
}

func (c *chatClient) RoomChat(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[ChatRoomMessage, ChatRoomMessage], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &Chat_ServiceDesc.Streams[0], Chat_RoomChat_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[ChatRoomMessage, ChatRoomMessage]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Chat_RoomChatClient = grpc.BidiStreamingClient[ChatRoomMessage, ChatRoomMessage]

func (c *chatClient) SendPrivateMessage(ctx context.Context, in *PrivateMessage, opts ...grpc.CallOption) (*MessageResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MessageResponse)
	err := c.cc.Invoke(ctx, Chat_SendPrivateMessage_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatClient) LeaveRoom(ctx context.Context, in *LeaveRequest, opts ...grpc.CallOption) (*MessageResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MessageResponse)
	err := c.cc.Invoke(ctx, Chat_LeaveRoom_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ChatServer is the server API for Chat service.
// All implementations must embed UnimplementedChatServer
// for forward compatibility.
type ChatServer interface {
	RoomChat(grpc.BidiStreamingServer[ChatRoomMessage, ChatRoomMessage]) error
	SendPrivateMessage(context.Context, *PrivateMessage) (*MessageResponse, error)
	LeaveRoom(context.Context, *LeaveRequest) (*MessageResponse, error)
	mustEmbedUnimplementedChatServer()
}

// UnimplementedChatServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedChatServer struct{}

func (UnimplementedChatServer) RoomChat(grpc.BidiStreamingServer[ChatRoomMessage, ChatRoomMessage]) error {
	return status.Errorf(codes.Unimplemented, "method RoomChat not implemented")
}
func (UnimplementedChatServer) SendPrivateMessage(context.Context, *PrivateMessage) (*MessageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendPrivateMessage not implemented")
}
func (UnimplementedChatServer) LeaveRoom(context.Context, *LeaveRequest) (*MessageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LeaveRoom not implemented")
}
func (UnimplementedChatServer) mustEmbedUnimplementedChatServer() {}
func (UnimplementedChatServer) testEmbeddedByValue()              {}

// UnsafeChatServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ChatServer will
// result in compilation errors.
type UnsafeChatServer interface {
	mustEmbedUnimplementedChatServer()
}

func RegisterChatServer(s grpc.ServiceRegistrar, srv ChatServer) {
	// If the following call pancis, it indicates UnimplementedChatServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Chat_ServiceDesc, srv)
}

func _Chat_RoomChat_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ChatServer).RoomChat(&grpc.GenericServerStream[ChatRoomMessage, ChatRoomMessage]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Chat_RoomChatServer = grpc.BidiStreamingServer[ChatRoomMessage, ChatRoomMessage]

func _Chat_SendPrivateMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PrivateMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServer).SendPrivateMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Chat_SendPrivateMessage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServer).SendPrivateMessage(ctx, req.(*PrivateMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chat_LeaveRoom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LeaveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServer).LeaveRoom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Chat_LeaveRoom_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServer).LeaveRoom(ctx, req.(*LeaveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Chat_ServiceDesc is the grpc.ServiceDesc for Chat service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Chat_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "chat.Chat",
	HandlerType: (*ChatServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendPrivateMessage",
			Handler:    _Chat_SendPrivateMessage_Handler,
		},
		{
			MethodName: "LeaveRoom",
			Handler:    _Chat_LeaveRoom_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "RoomChat",
			Handler:       _Chat_RoomChat_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "grpc/chatapp.proto",
}
