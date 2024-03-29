// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v4.25.2
// source: control/proto/control.proto

package control

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Message_ActionType int32

const (
	Message_PING     Message_ActionType = 0
	Message_SHUTDOWN Message_ActionType = 1
)

// Enum value maps for Message_ActionType.
var (
	Message_ActionType_name = map[int32]string{
		0: "PING",
		1: "SHUTDOWN",
	}
	Message_ActionType_value = map[string]int32{
		"PING":     0,
		"SHUTDOWN": 1,
	}
)

func (x Message_ActionType) Enum() *Message_ActionType {
	p := new(Message_ActionType)
	*p = x
	return p
}

func (x Message_ActionType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Message_ActionType) Descriptor() protoreflect.EnumDescriptor {
	return file_control_proto_control_proto_enumTypes[0].Descriptor()
}

func (Message_ActionType) Type() protoreflect.EnumType {
	return &file_control_proto_control_proto_enumTypes[0]
}

func (x Message_ActionType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Message_ActionType.Descriptor instead.
func (Message_ActionType) EnumDescriptor() ([]byte, []int) {
	return file_control_proto_control_proto_rawDescGZIP(), []int{0, 0}
}

type Response_ResponseType int32

const (
	Response_OK    Response_ResponseType = 0
	Response_ERROR Response_ResponseType = 1
)

// Enum value maps for Response_ResponseType.
var (
	Response_ResponseType_name = map[int32]string{
		0: "OK",
		1: "ERROR",
	}
	Response_ResponseType_value = map[string]int32{
		"OK":    0,
		"ERROR": 1,
	}
)

func (x Response_ResponseType) Enum() *Response_ResponseType {
	p := new(Response_ResponseType)
	*p = x
	return p
}

func (x Response_ResponseType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Response_ResponseType) Descriptor() protoreflect.EnumDescriptor {
	return file_control_proto_control_proto_enumTypes[1].Descriptor()
}

func (Response_ResponseType) Type() protoreflect.EnumType {
	return &file_control_proto_control_proto_enumTypes[1]
}

func (x Response_ResponseType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Response_ResponseType.Descriptor instead.
func (Response_ResponseType) EnumDescriptor() ([]byte, []int) {
	return file_control_proto_control_proto_rawDescGZIP(), []int{1, 0}
}

type Request_RequestType int32

const (
	Request_CHANGES  Request_RequestType = 0
	Request_INIT     Request_RequestType = 1
	Request_SETTINGS Request_RequestType = 2
)

// Enum value maps for Request_RequestType.
var (
	Request_RequestType_name = map[int32]string{
		0: "CHANGES",
		1: "INIT",
		2: "SETTINGS",
	}
	Request_RequestType_value = map[string]int32{
		"CHANGES":  0,
		"INIT":     1,
		"SETTINGS": 2,
	}
)

func (x Request_RequestType) Enum() *Request_RequestType {
	p := new(Request_RequestType)
	*p = x
	return p
}

func (x Request_RequestType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Request_RequestType) Descriptor() protoreflect.EnumDescriptor {
	return file_control_proto_control_proto_enumTypes[2].Descriptor()
}

func (Request_RequestType) Type() protoreflect.EnumType {
	return &file_control_proto_control_proto_enumTypes[2]
}

func (x Request_RequestType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Request_RequestType.Descriptor instead.
func (Request_RequestType) EnumDescriptor() ([]byte, []int) {
	return file_control_proto_control_proto_rawDescGZIP(), []int{2, 0}
}

type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Action Message_ActionType `protobuf:"varint,1,opt,name=action,proto3,enum=control.Message_ActionType" json:"action,omitempty"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_control_proto_control_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_control_proto_control_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_control_proto_control_proto_rawDescGZIP(), []int{0}
}

func (x *Message) GetAction() Message_ActionType {
	if x != nil {
		return x.Action
	}
	return Message_PING
}

type Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type Response_ResponseType `protobuf:"varint,1,opt,name=type,proto3,enum=control.Response_ResponseType" json:"type,omitempty"`
}

func (x *Response) Reset() {
	*x = Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_control_proto_control_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response) ProtoMessage() {}

func (x *Response) ProtoReflect() protoreflect.Message {
	mi := &file_control_proto_control_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response.ProtoReflect.Descriptor instead.
func (*Response) Descriptor() ([]byte, []int) {
	return file_control_proto_control_proto_rawDescGZIP(), []int{1}
}

func (x *Response) GetType() Response_ResponseType {
	if x != nil {
		return x.Type
	}
	return Response_OK
}

type Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type Request_RequestType `protobuf:"varint,1,opt,name=type,proto3,enum=control.Request_RequestType" json:"type,omitempty"`
}

func (x *Request) Reset() {
	*x = Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_control_proto_control_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Request) ProtoMessage() {}

func (x *Request) ProtoReflect() protoreflect.Message {
	mi := &file_control_proto_control_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Request.ProtoReflect.Descriptor instead.
func (*Request) Descriptor() ([]byte, []int) {
	return file_control_proto_control_proto_rawDescGZIP(), []int{2}
}

func (x *Request) GetType() Request_RequestType {
	if x != nil {
		return x.Type
	}
	return Request_CHANGES
}

type Entry struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key    []*Key  `protobuf:"bytes,1,rep,name=Key,proto3" json:"Key,omitempty"`
	KeyI   int64   `protobuf:"varint,2,opt,name=KeyI,proto3" json:"KeyI,omitempty"`
	Value  *Object `protobuf:"bytes,3,opt,name=Value,proto3" json:"Value,omitempty"`
	Remove bool    `protobuf:"varint,4,opt,name=Remove,proto3" json:"Remove,omitempty"`
}

func (x *Entry) Reset() {
	*x = Entry{}
	if protoimpl.UnsafeEnabled {
		mi := &file_control_proto_control_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Entry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Entry) ProtoMessage() {}

func (x *Entry) ProtoReflect() protoreflect.Message {
	mi := &file_control_proto_control_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Entry.ProtoReflect.Descriptor instead.
func (*Entry) Descriptor() ([]byte, []int) {
	return file_control_proto_control_proto_rawDescGZIP(), []int{3}
}

func (x *Entry) GetKey() []*Key {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *Entry) GetKeyI() int64 {
	if x != nil {
		return x.KeyI
	}
	return 0
}

func (x *Entry) GetValue() *Object {
	if x != nil {
		return x.Value
	}
	return nil
}

func (x *Entry) GetRemove() bool {
	if x != nil {
		return x.Remove
	}
	return false
}

type Key struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key    string    `protobuf:"bytes,1,opt,name=Key,proto3" json:"Key,omitempty"`
	Index  []*Object `protobuf:"bytes,2,rep,name=Index,proto3" json:"Index,omitempty"`
	IndexI int64     `protobuf:"varint,3,opt,name=IndexI,proto3" json:"IndexI,omitempty"`
}

func (x *Key) Reset() {
	*x = Key{}
	if protoimpl.UnsafeEnabled {
		mi := &file_control_proto_control_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Key) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Key) ProtoMessage() {}

func (x *Key) ProtoReflect() protoreflect.Message {
	mi := &file_control_proto_control_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Key.ProtoReflect.Descriptor instead.
func (*Key) Descriptor() ([]byte, []int) {
	return file_control_proto_control_proto_rawDescGZIP(), []int{4}
}

func (x *Key) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *Key) GetIndex() []*Object {
	if x != nil {
		return x.Index
	}
	return nil
}

func (x *Key) GetIndexI() int64 {
	if x != nil {
		return x.IndexI
	}
	return 0
}

type Object struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	String_ *string  `protobuf:"bytes,1,opt,name=string,proto3,oneof" json:"string,omitempty"`
	Int64   *int64   `protobuf:"varint,2,opt,name=int64,proto3,oneof" json:"int64,omitempty"`
	Uint64  *uint64  `protobuf:"varint,3,opt,name=uint64,proto3,oneof" json:"uint64,omitempty"`
	Float32 *float32 `protobuf:"fixed32,4,opt,name=float32,proto3,oneof" json:"float32,omitempty"`
	Float64 *float64 `protobuf:"fixed64,5,opt,name=float64,proto3,oneof" json:"float64,omitempty"`
	Bool    *bool    `protobuf:"varint,6,opt,name=bool,proto3,oneof" json:"bool,omitempty"`
	Bytes   []byte   `protobuf:"bytes,7,opt,name=bytes,proto3,oneof" json:"bytes,omitempty"`
}

func (x *Object) Reset() {
	*x = Object{}
	if protoimpl.UnsafeEnabled {
		mi := &file_control_proto_control_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Object) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Object) ProtoMessage() {}

func (x *Object) ProtoReflect() protoreflect.Message {
	mi := &file_control_proto_control_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Object.ProtoReflect.Descriptor instead.
func (*Object) Descriptor() ([]byte, []int) {
	return file_control_proto_control_proto_rawDescGZIP(), []int{5}
}

func (x *Object) GetString_() string {
	if x != nil && x.String_ != nil {
		return *x.String_
	}
	return ""
}

func (x *Object) GetInt64() int64 {
	if x != nil && x.Int64 != nil {
		return *x.Int64
	}
	return 0
}

func (x *Object) GetUint64() uint64 {
	if x != nil && x.Uint64 != nil {
		return *x.Uint64
	}
	return 0
}

func (x *Object) GetFloat32() float32 {
	if x != nil && x.Float32 != nil {
		return *x.Float32
	}
	return 0
}

func (x *Object) GetFloat64() float64 {
	if x != nil && x.Float64 != nil {
		return *x.Float64
	}
	return 0
}

func (x *Object) GetBool() bool {
	if x != nil && x.Bool != nil {
		return *x.Bool
	}
	return false
}

func (x *Object) GetBytes() []byte {
	if x != nil {
		return x.Bytes
	}
	return nil
}

var File_control_proto_control_proto protoreflect.FileDescriptor

var file_control_proto_control_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f,
	0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x63,
	0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x22, 0x64, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x12, 0x33, 0x0a, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x1b, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x2e, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x2e, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x52, 0x06,
	0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x24, 0x0a, 0x0a, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x54, 0x79, 0x70, 0x65, 0x12, 0x08, 0x0a, 0x04, 0x50, 0x49, 0x4e, 0x47, 0x10, 0x00, 0x12, 0x0c,
	0x0a, 0x08, 0x53, 0x48, 0x55, 0x54, 0x44, 0x4f, 0x57, 0x4e, 0x10, 0x01, 0x22, 0x61, 0x0a, 0x08,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x32, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1e, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c,
	0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0x21, 0x0a, 0x0c,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x06, 0x0a, 0x02,
	0x4f, 0x4b, 0x10, 0x00, 0x12, 0x09, 0x0a, 0x05, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x10, 0x01, 0x22,
	0x6f, 0x0a, 0x07, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x30, 0x0a, 0x04, 0x74, 0x79,
	0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1c, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72,
	0x6f, 0x6c, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0x32, 0x0a, 0x0b,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x43,
	0x48, 0x41, 0x4e, 0x47, 0x45, 0x53, 0x10, 0x00, 0x12, 0x08, 0x0a, 0x04, 0x49, 0x4e, 0x49, 0x54,
	0x10, 0x01, 0x12, 0x0c, 0x0a, 0x08, 0x53, 0x45, 0x54, 0x54, 0x49, 0x4e, 0x47, 0x53, 0x10, 0x02,
	0x22, 0x7a, 0x0a, 0x05, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x1e, 0x0a, 0x03, 0x4b, 0x65, 0x79,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c,
	0x2e, 0x4b, 0x65, 0x79, 0x52, 0x03, 0x4b, 0x65, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x4b, 0x65, 0x79,
	0x49, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x4b, 0x65, 0x79, 0x49, 0x12, 0x25, 0x0a,
	0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x63,
	0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x2e, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x52, 0x05, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x22, 0x56, 0x0a, 0x03,
	0x4b, 0x65, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x4b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x4b, 0x65, 0x79, 0x12, 0x25, 0x0a, 0x05, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x02,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x2e, 0x4f,
	0x62, 0x6a, 0x65, 0x63, 0x74, 0x52, 0x05, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x16, 0x0a, 0x06,
	0x49, 0x6e, 0x64, 0x65, 0x78, 0x49, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x49, 0x6e,
	0x64, 0x65, 0x78, 0x49, 0x22, 0x9a, 0x02, 0x0a, 0x06, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x12,
	0x1b, 0x0a, 0x06, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x48,
	0x00, 0x52, 0x06, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x88, 0x01, 0x01, 0x12, 0x19, 0x0a, 0x05,
	0x69, 0x6e, 0x74, 0x36, 0x34, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x48, 0x01, 0x52, 0x05, 0x69,
	0x6e, 0x74, 0x36, 0x34, 0x88, 0x01, 0x01, 0x12, 0x1b, 0x0a, 0x06, 0x75, 0x69, 0x6e, 0x74, 0x36,
	0x34, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x48, 0x02, 0x52, 0x06, 0x75, 0x69, 0x6e, 0x74, 0x36,
	0x34, 0x88, 0x01, 0x01, 0x12, 0x1d, 0x0a, 0x07, 0x66, 0x6c, 0x6f, 0x61, 0x74, 0x33, 0x32, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x02, 0x48, 0x03, 0x52, 0x07, 0x66, 0x6c, 0x6f, 0x61, 0x74, 0x33, 0x32,
	0x88, 0x01, 0x01, 0x12, 0x1d, 0x0a, 0x07, 0x66, 0x6c, 0x6f, 0x61, 0x74, 0x36, 0x34, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x01, 0x48, 0x04, 0x52, 0x07, 0x66, 0x6c, 0x6f, 0x61, 0x74, 0x36, 0x34, 0x88,
	0x01, 0x01, 0x12, 0x17, 0x0a, 0x04, 0x62, 0x6f, 0x6f, 0x6c, 0x18, 0x06, 0x20, 0x01, 0x28, 0x08,
	0x48, 0x05, 0x52, 0x04, 0x62, 0x6f, 0x6f, 0x6c, 0x88, 0x01, 0x01, 0x12, 0x19, 0x0a, 0x05, 0x62,
	0x79, 0x74, 0x65, 0x73, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0c, 0x48, 0x06, 0x52, 0x05, 0x62, 0x79,
	0x74, 0x65, 0x73, 0x88, 0x01, 0x01, 0x42, 0x09, 0x0a, 0x07, 0x5f, 0x73, 0x74, 0x72, 0x69, 0x6e,
	0x67, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x69, 0x6e, 0x74, 0x36, 0x34, 0x42, 0x09, 0x0a, 0x07, 0x5f,
	0x75, 0x69, 0x6e, 0x74, 0x36, 0x34, 0x42, 0x0a, 0x0a, 0x08, 0x5f, 0x66, 0x6c, 0x6f, 0x61, 0x74,
	0x33, 0x32, 0x42, 0x0a, 0x0a, 0x08, 0x5f, 0x66, 0x6c, 0x6f, 0x61, 0x74, 0x36, 0x34, 0x42, 0x07,
	0x0a, 0x05, 0x5f, 0x62, 0x6f, 0x6f, 0x6c, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x62, 0x79, 0x74, 0x65,
	0x73, 0x32, 0xca, 0x01, 0x0a, 0x07, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x12, 0x2c, 0x0a,
	0x04, 0x50, 0x75, 0x6c, 0x6c, 0x12, 0x10, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x2e,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0e, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f,
	0x6c, 0x2e, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x22, 0x00, 0x30, 0x01, 0x12, 0x2d, 0x0a, 0x04, 0x50,
	0x75, 0x73, 0x68, 0x12, 0x0e, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x2e, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x1a, 0x11, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x2e, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x28, 0x01, 0x12, 0x30, 0x0a, 0x08, 0x50, 0x75,
	0x73, 0x68, 0x50, 0x75, 0x6c, 0x6c, 0x12, 0x0e, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c,
	0x2e, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x1a, 0x0e, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c,
	0x2e, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x22, 0x00, 0x28, 0x01, 0x30, 0x01, 0x12, 0x30, 0x0a, 0x07,
	0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x12, 0x10, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f,
	0x6c, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x11, 0x2e, 0x63, 0x6f, 0x6e, 0x74,
	0x72, 0x6f, 0x6c, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x0a,
	0x5a, 0x08, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_control_proto_control_proto_rawDescOnce sync.Once
	file_control_proto_control_proto_rawDescData = file_control_proto_control_proto_rawDesc
)

func file_control_proto_control_proto_rawDescGZIP() []byte {
	file_control_proto_control_proto_rawDescOnce.Do(func() {
		file_control_proto_control_proto_rawDescData = protoimpl.X.CompressGZIP(file_control_proto_control_proto_rawDescData)
	})
	return file_control_proto_control_proto_rawDescData
}

var file_control_proto_control_proto_enumTypes = make([]protoimpl.EnumInfo, 3)
var file_control_proto_control_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_control_proto_control_proto_goTypes = []interface{}{
	(Message_ActionType)(0),    // 0: control.Message.ActionType
	(Response_ResponseType)(0), // 1: control.Response.ResponseType
	(Request_RequestType)(0),   // 2: control.Request.RequestType
	(*Message)(nil),            // 3: control.Message
	(*Response)(nil),           // 4: control.Response
	(*Request)(nil),            // 5: control.Request
	(*Entry)(nil),              // 6: control.Entry
	(*Key)(nil),                // 7: control.Key
	(*Object)(nil),             // 8: control.Object
}
var file_control_proto_control_proto_depIdxs = []int32{
	0,  // 0: control.Message.action:type_name -> control.Message.ActionType
	1,  // 1: control.Response.type:type_name -> control.Response.ResponseType
	2,  // 2: control.Request.type:type_name -> control.Request.RequestType
	7,  // 3: control.Entry.Key:type_name -> control.Key
	8,  // 4: control.Entry.Value:type_name -> control.Object
	8,  // 5: control.Key.Index:type_name -> control.Object
	5,  // 6: control.Control.Pull:input_type -> control.Request
	6,  // 7: control.Control.Push:input_type -> control.Entry
	6,  // 8: control.Control.PushPull:input_type -> control.Entry
	3,  // 9: control.Control.Control:input_type -> control.Message
	6,  // 10: control.Control.Pull:output_type -> control.Entry
	4,  // 11: control.Control.Push:output_type -> control.Response
	6,  // 12: control.Control.PushPull:output_type -> control.Entry
	4,  // 13: control.Control.Control:output_type -> control.Response
	10, // [10:14] is the sub-list for method output_type
	6,  // [6:10] is the sub-list for method input_type
	6,  // [6:6] is the sub-list for extension type_name
	6,  // [6:6] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_control_proto_control_proto_init() }
func file_control_proto_control_proto_init() {
	if File_control_proto_control_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_control_proto_control_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Message); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_control_proto_control_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Response); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_control_proto_control_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Request); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_control_proto_control_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Entry); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_control_proto_control_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Key); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_control_proto_control_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Object); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_control_proto_control_proto_msgTypes[5].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_control_proto_control_proto_rawDesc,
			NumEnums:      3,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_control_proto_control_proto_goTypes,
		DependencyIndexes: file_control_proto_control_proto_depIdxs,
		EnumInfos:         file_control_proto_control_proto_enumTypes,
		MessageInfos:      file_control_proto_control_proto_msgTypes,
	}.Build()
	File_control_proto_control_proto = out.File
	file_control_proto_control_proto_rawDesc = nil
	file_control_proto_control_proto_goTypes = nil
	file_control_proto_control_proto_depIdxs = nil
}
