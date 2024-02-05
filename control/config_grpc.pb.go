// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.2
// source: control/proto/config.proto

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

// ConfigClient is the client API for Config service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ConfigClient interface {
	Pull(ctx context.Context, in *Request, opts ...grpc.CallOption) (Config_PullClient, error)
	Push(ctx context.Context, opts ...grpc.CallOption) (Config_PushClient, error)
	PushPull(ctx context.Context, opts ...grpc.CallOption) (Config_PushPullClient, error)
	Control(ctx context.Context, in *Message, opts ...grpc.CallOption) (*Response, error)
}

type configClient struct {
	cc grpc.ClientConnInterface
}

func NewConfigClient(cc grpc.ClientConnInterface) ConfigClient {
	return &configClient{cc}
}

func (c *configClient) Pull(ctx context.Context, in *Request, opts ...grpc.CallOption) (Config_PullClient, error) {
	stream, err := c.cc.NewStream(ctx, &Config_ServiceDesc.Streams[0], "/control.Config/Pull", opts...)
	if err != nil {
		return nil, err
	}
	x := &configPullClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Config_PullClient interface {
	Recv() (*Entry, error)
	grpc.ClientStream
}

type configPullClient struct {
	grpc.ClientStream
}

func (x *configPullClient) Recv() (*Entry, error) {
	m := new(Entry)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *configClient) Push(ctx context.Context, opts ...grpc.CallOption) (Config_PushClient, error) {
	stream, err := c.cc.NewStream(ctx, &Config_ServiceDesc.Streams[1], "/control.Config/Push", opts...)
	if err != nil {
		return nil, err
	}
	x := &configPushClient{stream}
	return x, nil
}

type Config_PushClient interface {
	Send(*Entry) error
	CloseAndRecv() (*Response, error)
	grpc.ClientStream
}

type configPushClient struct {
	grpc.ClientStream
}

func (x *configPushClient) Send(m *Entry) error {
	return x.ClientStream.SendMsg(m)
}

func (x *configPushClient) CloseAndRecv() (*Response, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *configClient) PushPull(ctx context.Context, opts ...grpc.CallOption) (Config_PushPullClient, error) {
	stream, err := c.cc.NewStream(ctx, &Config_ServiceDesc.Streams[2], "/control.Config/PushPull", opts...)
	if err != nil {
		return nil, err
	}
	x := &configPushPullClient{stream}
	return x, nil
}

type Config_PushPullClient interface {
	Send(*Entry) error
	Recv() (*Entry, error)
	grpc.ClientStream
}

type configPushPullClient struct {
	grpc.ClientStream
}

func (x *configPushPullClient) Send(m *Entry) error {
	return x.ClientStream.SendMsg(m)
}

func (x *configPushPullClient) Recv() (*Entry, error) {
	m := new(Entry)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *configClient) Control(ctx context.Context, in *Message, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/control.Config/Control", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ConfigServer is the server API for Config service.
// All implementations must embed UnimplementedConfigServer
// for forward compatibility
type ConfigServer interface {
	Pull(*Request, Config_PullServer) error
	Push(Config_PushServer) error
	PushPull(Config_PushPullServer) error
	Control(context.Context, *Message) (*Response, error)
	mustEmbedUnimplementedConfigServer()
}

// UnimplementedConfigServer must be embedded to have forward compatible implementations.
type UnimplementedConfigServer struct {
}

func (UnimplementedConfigServer) Pull(*Request, Config_PullServer) error {
	return status.Errorf(codes.Unimplemented, "method Pull not implemented")
}
func (UnimplementedConfigServer) Push(Config_PushServer) error {
	return status.Errorf(codes.Unimplemented, "method Push not implemented")
}
func (UnimplementedConfigServer) PushPull(Config_PushPullServer) error {
	return status.Errorf(codes.Unimplemented, "method PushPull not implemented")
}
func (UnimplementedConfigServer) Control(context.Context, *Message) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Control not implemented")
}
func (UnimplementedConfigServer) mustEmbedUnimplementedConfigServer() {}

// UnsafeConfigServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ConfigServer will
// result in compilation errors.
type UnsafeConfigServer interface {
	mustEmbedUnimplementedConfigServer()
}

func RegisterConfigServer(s grpc.ServiceRegistrar, srv ConfigServer) {
	s.RegisterService(&Config_ServiceDesc, srv)
}

func _Config_Pull_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Request)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ConfigServer).Pull(m, &configPullServer{stream})
}

type Config_PullServer interface {
	Send(*Entry) error
	grpc.ServerStream
}

type configPullServer struct {
	grpc.ServerStream
}

func (x *configPullServer) Send(m *Entry) error {
	return x.ServerStream.SendMsg(m)
}

func _Config_Push_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ConfigServer).Push(&configPushServer{stream})
}

type Config_PushServer interface {
	SendAndClose(*Response) error
	Recv() (*Entry, error)
	grpc.ServerStream
}

type configPushServer struct {
	grpc.ServerStream
}

func (x *configPushServer) SendAndClose(m *Response) error {
	return x.ServerStream.SendMsg(m)
}

func (x *configPushServer) Recv() (*Entry, error) {
	m := new(Entry)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Config_PushPull_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ConfigServer).PushPull(&configPushPullServer{stream})
}

type Config_PushPullServer interface {
	Send(*Entry) error
	Recv() (*Entry, error)
	grpc.ServerStream
}

type configPushPullServer struct {
	grpc.ServerStream
}

func (x *configPushPullServer) Send(m *Entry) error {
	return x.ServerStream.SendMsg(m)
}

func (x *configPushPullServer) Recv() (*Entry, error) {
	m := new(Entry)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Config_Control_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Message)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConfigServer).Control(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/control.Config/Control",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConfigServer).Control(ctx, req.(*Message))
	}
	return interceptor(ctx, in, info, handler)
}

// Config_ServiceDesc is the grpc.ServiceDesc for Config service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Config_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "control.Config",
	HandlerType: (*ConfigServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Control",
			Handler:    _Config_Control_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Pull",
			Handler:       _Config_Pull_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "Push",
			Handler:       _Config_Push_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "PushPull",
			Handler:       _Config_PushPull_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "control/proto/config.proto",
}
