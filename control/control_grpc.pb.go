// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.2
// source: control/proto/control.proto

package control

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

// ControlClient is the client API for Control service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ControlClient interface {
	Pull(ctx context.Context, in *Request, opts ...grpc.CallOption) (Control_PullClient, error)
	Push(ctx context.Context, opts ...grpc.CallOption) (Control_PushClient, error)
	PushPull(ctx context.Context, opts ...grpc.CallOption) (Control_PushPullClient, error)
	Control(ctx context.Context, in *Message, opts ...grpc.CallOption) (*Response, error)
}

type controlClient struct {
	cc grpc.ClientConnInterface
}

func NewControlClient(cc grpc.ClientConnInterface) ControlClient {
	return &controlClient{cc}
}

func (c *controlClient) Pull(ctx context.Context, in *Request, opts ...grpc.CallOption) (Control_PullClient, error) {
	stream, err := c.cc.NewStream(ctx, &Control_ServiceDesc.Streams[0], "/control.Control/Pull", opts...)
	if err != nil {
		return nil, err
	}
	x := &controlPullClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Control_PullClient interface {
	Recv() (*Entry, error)
	grpc.ClientStream
}

type controlPullClient struct {
	grpc.ClientStream
}

func (x *controlPullClient) Recv() (*Entry, error) {
	m := new(Entry)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *controlClient) Push(ctx context.Context, opts ...grpc.CallOption) (Control_PushClient, error) {
	stream, err := c.cc.NewStream(ctx, &Control_ServiceDesc.Streams[1], "/control.Control/Push", opts...)
	if err != nil {
		return nil, err
	}
	x := &controlPushClient{stream}
	return x, nil
}

type Control_PushClient interface {
	Send(*Entry) error
	CloseAndRecv() (*Response, error)
	grpc.ClientStream
}

type controlPushClient struct {
	grpc.ClientStream
}

func (x *controlPushClient) Send(m *Entry) error {
	return x.ClientStream.SendMsg(m)
}

func (x *controlPushClient) CloseAndRecv() (*Response, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *controlClient) PushPull(ctx context.Context, opts ...grpc.CallOption) (Control_PushPullClient, error) {
	stream, err := c.cc.NewStream(ctx, &Control_ServiceDesc.Streams[2], "/control.Control/PushPull", opts...)
	if err != nil {
		return nil, err
	}
	x := &controlPushPullClient{stream}
	return x, nil
}

type Control_PushPullClient interface {
	Send(*Entry) error
	Recv() (*Entry, error)
	grpc.ClientStream
}

type controlPushPullClient struct {
	grpc.ClientStream
}

func (x *controlPushPullClient) Send(m *Entry) error {
	return x.ClientStream.SendMsg(m)
}

func (x *controlPushPullClient) Recv() (*Entry, error) {
	m := new(Entry)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *controlClient) Control(ctx context.Context, in *Message, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/control.Control/Control", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ControlServer is the server API for Control service.
// All implementations must embed UnimplementedControlServer
// for forward compatibility
type ControlServer interface {
	Pull(*Request, Control_PullServer) error
	Push(Control_PushServer) error
	PushPull(Control_PushPullServer) error
	Control(context.Context, *Message) (*Response, error)
	mustEmbedUnimplementedControlServer()
}

// UnimplementedControlServer must be embedded to have forward compatible implementations.
type UnimplementedControlServer struct {
}

func (UnimplementedControlServer) Pull(*Request, Control_PullServer) error {
	return status.Errorf(codes.Unimplemented, "method Pull not implemented")
}
func (UnimplementedControlServer) Push(Control_PushServer) error {
	return status.Errorf(codes.Unimplemented, "method Push not implemented")
}
func (UnimplementedControlServer) PushPull(Control_PushPullServer) error {
	return status.Errorf(codes.Unimplemented, "method PushPull not implemented")
}
func (UnimplementedControlServer) Control(context.Context, *Message) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Control not implemented")
}
func (UnimplementedControlServer) mustEmbedUnimplementedControlServer() {}

// UnsafeControlServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ControlServer will
// result in compilation errors.
type UnsafeControlServer interface {
	mustEmbedUnimplementedControlServer()
}

func RegisterControlServer(s grpc.ServiceRegistrar, srv ControlServer) {
	s.RegisterService(&Control_ServiceDesc, srv)
}

func _Control_Pull_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Request)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ControlServer).Pull(m, &controlPullServer{stream})
}

type Control_PullServer interface {
	Send(*Entry) error
	grpc.ServerStream
}

type controlPullServer struct {
	grpc.ServerStream
}

func (x *controlPullServer) Send(m *Entry) error {
	return x.ServerStream.SendMsg(m)
}

func _Control_Push_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ControlServer).Push(&controlPushServer{stream})
}

type Control_PushServer interface {
	SendAndClose(*Response) error
	Recv() (*Entry, error)
	grpc.ServerStream
}

type controlPushServer struct {
	grpc.ServerStream
}

func (x *controlPushServer) SendAndClose(m *Response) error {
	return x.ServerStream.SendMsg(m)
}

func (x *controlPushServer) Recv() (*Entry, error) {
	m := new(Entry)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Control_PushPull_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ControlServer).PushPull(&controlPushPullServer{stream})
}

type Control_PushPullServer interface {
	Send(*Entry) error
	Recv() (*Entry, error)
	grpc.ServerStream
}

type controlPushPullServer struct {
	grpc.ServerStream
}

func (x *controlPushPullServer) Send(m *Entry) error {
	return x.ServerStream.SendMsg(m)
}

func (x *controlPushPullServer) Recv() (*Entry, error) {
	m := new(Entry)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Control_Control_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Message)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ControlServer).Control(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/control.Control/Control",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ControlServer).Control(ctx, req.(*Message))
	}
	return interceptor(ctx, in, info, handler)
}

// Control_ServiceDesc is the grpc.ServiceDesc for Control service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Control_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "control.Control",
	HandlerType: (*ControlServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Control",
			Handler:    _Control_Control_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Pull",
			Handler:       _Control_Pull_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "Push",
			Handler:       _Control_Push_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "PushPull",
			Handler:       _Control_PushPull_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "control/proto/control.proto",
}
