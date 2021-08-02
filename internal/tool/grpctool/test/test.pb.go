// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.0
// source: internal/tool/grpctool/test/test.proto

package test

import (
	reflect "reflect"
	sync "sync"

	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	_ "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/grpctool/automata"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Enum1 int32

const (
	Enum1_v1 Enum1 = 0
	Enum1_v2 Enum1 = 1
)

// Enum value maps for Enum1.
var (
	Enum1_name = map[int32]string{
		0: "v1",
		1: "v2",
	}
	Enum1_value = map[string]int32{
		"v1": 0,
		"v2": 1,
	}
)

func (x Enum1) Enum() *Enum1 {
	p := new(Enum1)
	*p = x
	return p
}

func (x Enum1) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Enum1) Descriptor() protoreflect.EnumDescriptor {
	return file_internal_tool_grpctool_test_test_proto_enumTypes[0].Descriptor()
}

func (Enum1) Type() protoreflect.EnumType {
	return &file_internal_tool_grpctool_test_test_proto_enumTypes[0]
}

func (x Enum1) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Enum1.Descriptor instead.
func (Enum1) EnumDescriptor() ([]byte, []int) {
	return file_internal_tool_grpctool_test_test_proto_rawDescGZIP(), []int{0}
}

type Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	S1 string `protobuf:"bytes,1,opt,name=s1,proto3" json:"s1,omitempty"`
}

func (x *Request) Reset() {
	*x = Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_tool_grpctool_test_test_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Request) ProtoMessage() {}

func (x *Request) ProtoReflect() protoreflect.Message {
	mi := &file_internal_tool_grpctool_test_test_proto_msgTypes[0]
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
	return file_internal_tool_grpctool_test_test_proto_rawDescGZIP(), []int{0}
}

func (x *Request) GetS1() string {
	if x != nil {
		return x.S1
	}
	return ""
}

type Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Message:
	//	*Response_Scalar
	//	*Response_X1
	//	*Response_Data_
	//	*Response_Last_
	Message isResponse_Message `protobuf_oneof:"message"`
}

func (x *Response) Reset() {
	*x = Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_tool_grpctool_test_test_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response) ProtoMessage() {}

func (x *Response) ProtoReflect() protoreflect.Message {
	mi := &file_internal_tool_grpctool_test_test_proto_msgTypes[1]
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
	return file_internal_tool_grpctool_test_test_proto_rawDescGZIP(), []int{1}
}

func (m *Response) GetMessage() isResponse_Message {
	if m != nil {
		return m.Message
	}
	return nil
}

func (x *Response) GetScalar() int64 {
	if x, ok := x.GetMessage().(*Response_Scalar); ok {
		return x.Scalar
	}
	return 0
}

func (x *Response) GetX1() Enum1 {
	if x, ok := x.GetMessage().(*Response_X1); ok {
		return x.X1
	}
	return Enum1_v1
}

func (x *Response) GetData() *Response_Data {
	if x, ok := x.GetMessage().(*Response_Data_); ok {
		return x.Data
	}
	return nil
}

func (x *Response) GetLast() *Response_Last {
	if x, ok := x.GetMessage().(*Response_Last_); ok {
		return x.Last
	}
	return nil
}

type isResponse_Message interface {
	isResponse_Message()
}

type Response_Scalar struct {
	Scalar int64 `protobuf:"varint,1,opt,name=scalar,proto3,oneof"`
}

type Response_X1 struct {
	X1 Enum1 `protobuf:"varint,2,opt,name=x1,proto3,enum=gitlab.agent.grpctool.test.Enum1,oneof"`
}

type Response_Data_ struct {
	Data *Response_Data `protobuf:"bytes,3,opt,name=data,proto3,oneof"`
}

type Response_Last_ struct {
	Last *Response_Last `protobuf:"bytes,4,opt,name=last,proto3,oneof"`
}

func (*Response_Scalar) isResponse_Message() {}

func (*Response_X1) isResponse_Message() {}

func (*Response_Data_) isResponse_Message() {}

func (*Response_Last_) isResponse_Message() {}

type NoOneofs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *NoOneofs) Reset() {
	*x = NoOneofs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_tool_grpctool_test_test_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NoOneofs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NoOneofs) ProtoMessage() {}

func (x *NoOneofs) ProtoReflect() protoreflect.Message {
	mi := &file_internal_tool_grpctool_test_test_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NoOneofs.ProtoReflect.Descriptor instead.
func (*NoOneofs) Descriptor() ([]byte, []int) {
	return file_internal_tool_grpctool_test_test_proto_rawDescGZIP(), []int{2}
}

type TwoOneofs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Message1:
	//	*TwoOneofs_M11
	//	*TwoOneofs_M12
	Message1 isTwoOneofs_Message1 `protobuf_oneof:"message1"`
	// Types that are assignable to Message2:
	//	*TwoOneofs_M21
	//	*TwoOneofs_M22
	Message2 isTwoOneofs_Message2 `protobuf_oneof:"message2"`
}

func (x *TwoOneofs) Reset() {
	*x = TwoOneofs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_tool_grpctool_test_test_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TwoOneofs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TwoOneofs) ProtoMessage() {}

func (x *TwoOneofs) ProtoReflect() protoreflect.Message {
	mi := &file_internal_tool_grpctool_test_test_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TwoOneofs.ProtoReflect.Descriptor instead.
func (*TwoOneofs) Descriptor() ([]byte, []int) {
	return file_internal_tool_grpctool_test_test_proto_rawDescGZIP(), []int{3}
}

func (m *TwoOneofs) GetMessage1() isTwoOneofs_Message1 {
	if m != nil {
		return m.Message1
	}
	return nil
}

func (x *TwoOneofs) GetM11() int32 {
	if x, ok := x.GetMessage1().(*TwoOneofs_M11); ok {
		return x.M11
	}
	return 0
}

func (x *TwoOneofs) GetM12() int32 {
	if x, ok := x.GetMessage1().(*TwoOneofs_M12); ok {
		return x.M12
	}
	return 0
}

func (m *TwoOneofs) GetMessage2() isTwoOneofs_Message2 {
	if m != nil {
		return m.Message2
	}
	return nil
}

func (x *TwoOneofs) GetM21() int32 {
	if x, ok := x.GetMessage2().(*TwoOneofs_M21); ok {
		return x.M21
	}
	return 0
}

func (x *TwoOneofs) GetM22() int32 {
	if x, ok := x.GetMessage2().(*TwoOneofs_M22); ok {
		return x.M22
	}
	return 0
}

type isTwoOneofs_Message1 interface {
	isTwoOneofs_Message1()
}

type TwoOneofs_M11 struct {
	M11 int32 `protobuf:"varint,1,opt,name=m11,proto3,oneof"`
}

type TwoOneofs_M12 struct {
	M12 int32 `protobuf:"varint,2,opt,name=m12,proto3,oneof"`
}

func (*TwoOneofs_M11) isTwoOneofs_Message1() {}

func (*TwoOneofs_M12) isTwoOneofs_Message1() {}

type isTwoOneofs_Message2 interface {
	isTwoOneofs_Message2()
}

type TwoOneofs_M21 struct {
	M21 int32 `protobuf:"varint,3,opt,name=m21,proto3,oneof"`
}

type TwoOneofs_M22 struct {
	M22 int32 `protobuf:"varint,4,opt,name=m22,proto3,oneof"`
}

func (*TwoOneofs_M21) isTwoOneofs_Message2() {}

func (*TwoOneofs_M22) isTwoOneofs_Message2() {}

type TwoValidOneofs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Message1:
	//	*TwoValidOneofs_M11
	//	*TwoValidOneofs_M12
	Message1 isTwoValidOneofs_Message1 `protobuf_oneof:"message1"`
	// Types that are assignable to Message2:
	//	*TwoValidOneofs_M21
	//	*TwoValidOneofs_M22
	Message2 isTwoValidOneofs_Message2 `protobuf_oneof:"message2"`
}

func (x *TwoValidOneofs) Reset() {
	*x = TwoValidOneofs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_tool_grpctool_test_test_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TwoValidOneofs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TwoValidOneofs) ProtoMessage() {}

func (x *TwoValidOneofs) ProtoReflect() protoreflect.Message {
	mi := &file_internal_tool_grpctool_test_test_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TwoValidOneofs.ProtoReflect.Descriptor instead.
func (*TwoValidOneofs) Descriptor() ([]byte, []int) {
	return file_internal_tool_grpctool_test_test_proto_rawDescGZIP(), []int{4}
}

func (m *TwoValidOneofs) GetMessage1() isTwoValidOneofs_Message1 {
	if m != nil {
		return m.Message1
	}
	return nil
}

func (x *TwoValidOneofs) GetM11() int32 {
	if x, ok := x.GetMessage1().(*TwoValidOneofs_M11); ok {
		return x.M11
	}
	return 0
}

func (x *TwoValidOneofs) GetM12() int32 {
	if x, ok := x.GetMessage1().(*TwoValidOneofs_M12); ok {
		return x.M12
	}
	return 0
}

func (m *TwoValidOneofs) GetMessage2() isTwoValidOneofs_Message2 {
	if m != nil {
		return m.Message2
	}
	return nil
}

func (x *TwoValidOneofs) GetM21() int32 {
	if x, ok := x.GetMessage2().(*TwoValidOneofs_M21); ok {
		return x.M21
	}
	return 0
}

func (x *TwoValidOneofs) GetM22() int32 {
	if x, ok := x.GetMessage2().(*TwoValidOneofs_M22); ok {
		return x.M22
	}
	return 0
}

type isTwoValidOneofs_Message1 interface {
	isTwoValidOneofs_Message1()
}

type TwoValidOneofs_M11 struct {
	M11 int32 `protobuf:"varint,1,opt,name=m11,proto3,oneof"`
}

type TwoValidOneofs_M12 struct {
	M12 int32 `protobuf:"varint,2,opt,name=m12,proto3,oneof"`
}

func (*TwoValidOneofs_M11) isTwoValidOneofs_Message1() {}

func (*TwoValidOneofs_M12) isTwoValidOneofs_Message1() {}

type isTwoValidOneofs_Message2 interface {
	isTwoValidOneofs_Message2()
}

type TwoValidOneofs_M21 struct {
	M21 int32 `protobuf:"varint,3,opt,name=m21,proto3,oneof"`
}

type TwoValidOneofs_M22 struct {
	M22 int32 `protobuf:"varint,4,opt,name=m22,proto3,oneof"`
}

func (*TwoValidOneofs_M21) isTwoValidOneofs_Message2() {}

func (*TwoValidOneofs_M22) isTwoValidOneofs_Message2() {}

type OutOfOneof struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	X int32 `protobuf:"varint,1,opt,name=x,proto3" json:"x,omitempty"`
	// Types that are assignable to Message:
	//	*OutOfOneof_M1
	//	*OutOfOneof_M2
	Message isOutOfOneof_Message `protobuf_oneof:"message"`
}

func (x *OutOfOneof) Reset() {
	*x = OutOfOneof{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_tool_grpctool_test_test_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OutOfOneof) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OutOfOneof) ProtoMessage() {}

func (x *OutOfOneof) ProtoReflect() protoreflect.Message {
	mi := &file_internal_tool_grpctool_test_test_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OutOfOneof.ProtoReflect.Descriptor instead.
func (*OutOfOneof) Descriptor() ([]byte, []int) {
	return file_internal_tool_grpctool_test_test_proto_rawDescGZIP(), []int{5}
}

func (x *OutOfOneof) GetX() int32 {
	if x != nil {
		return x.X
	}
	return 0
}

func (m *OutOfOneof) GetMessage() isOutOfOneof_Message {
	if m != nil {
		return m.Message
	}
	return nil
}

func (x *OutOfOneof) GetM1() int32 {
	if x, ok := x.GetMessage().(*OutOfOneof_M1); ok {
		return x.M1
	}
	return 0
}

func (x *OutOfOneof) GetM2() int32 {
	if x, ok := x.GetMessage().(*OutOfOneof_M2); ok {
		return x.M2
	}
	return 0
}

type isOutOfOneof_Message interface {
	isOutOfOneof_Message()
}

type OutOfOneof_M1 struct {
	M1 int32 `protobuf:"varint,2,opt,name=m1,proto3,oneof"`
}

type OutOfOneof_M2 struct {
	M2 int32 `protobuf:"varint,3,opt,name=m2,proto3,oneof"`
}

func (*OutOfOneof_M1) isOutOfOneof_Message() {}

func (*OutOfOneof_M2) isOutOfOneof_Message() {}

type NotAllReachable struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Message:
	//	*NotAllReachable_M1
	//	*NotAllReachable_M2
	//	*NotAllReachable_M3
	Message isNotAllReachable_Message `protobuf_oneof:"message"`
}

func (x *NotAllReachable) Reset() {
	*x = NotAllReachable{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_tool_grpctool_test_test_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NotAllReachable) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NotAllReachable) ProtoMessage() {}

func (x *NotAllReachable) ProtoReflect() protoreflect.Message {
	mi := &file_internal_tool_grpctool_test_test_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NotAllReachable.ProtoReflect.Descriptor instead.
func (*NotAllReachable) Descriptor() ([]byte, []int) {
	return file_internal_tool_grpctool_test_test_proto_rawDescGZIP(), []int{6}
}

func (m *NotAllReachable) GetMessage() isNotAllReachable_Message {
	if m != nil {
		return m.Message
	}
	return nil
}

func (x *NotAllReachable) GetM1() int32 {
	if x, ok := x.GetMessage().(*NotAllReachable_M1); ok {
		return x.M1
	}
	return 0
}

func (x *NotAllReachable) GetM2() int32 {
	if x, ok := x.GetMessage().(*NotAllReachable_M2); ok {
		return x.M2
	}
	return 0
}

func (x *NotAllReachable) GetM3() int32 {
	if x, ok := x.GetMessage().(*NotAllReachable_M3); ok {
		return x.M3
	}
	return 0
}

type isNotAllReachable_Message interface {
	isNotAllReachable_Message()
}

type NotAllReachable_M1 struct {
	M1 int32 `protobuf:"varint,1,opt,name=m1,proto3,oneof"`
}

type NotAllReachable_M2 struct {
	M2 int32 `protobuf:"varint,2,opt,name=m2,proto3,oneof"`
}

type NotAllReachable_M3 struct {
	M3 int32 `protobuf:"varint,3,opt,name=m3,proto3,oneof"`
}

func (*NotAllReachable_M1) isNotAllReachable_Message() {}

func (*NotAllReachable_M2) isNotAllReachable_Message() {}

func (*NotAllReachable_M3) isNotAllReachable_Message() {}

type Response_Data struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data []byte `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *Response_Data) Reset() {
	*x = Response_Data{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_tool_grpctool_test_test_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Response_Data) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response_Data) ProtoMessage() {}

func (x *Response_Data) ProtoReflect() protoreflect.Message {
	mi := &file_internal_tool_grpctool_test_test_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response_Data.ProtoReflect.Descriptor instead.
func (*Response_Data) Descriptor() ([]byte, []int) {
	return file_internal_tool_grpctool_test_test_proto_rawDescGZIP(), []int{1, 0}
}

func (x *Response_Data) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

type Response_Last struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Response_Last) Reset() {
	*x = Response_Last{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_tool_grpctool_test_test_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Response_Last) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response_Last) ProtoMessage() {}

func (x *Response_Last) ProtoReflect() protoreflect.Message {
	mi := &file_internal_tool_grpctool_test_test_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response_Last.ProtoReflect.Descriptor instead.
func (*Response_Last) Descriptor() ([]byte, []int) {
	return file_internal_tool_grpctool_test_test_proto_rawDescGZIP(), []int{1, 1}
}

var File_internal_tool_grpctool_test_test_proto protoreflect.FileDescriptor

var file_internal_tool_grpctool_test_test_proto_rawDesc = []byte{
	0x0a, 0x26, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x74, 0x6f, 0x6f, 0x6c, 0x2f,
	0x67, 0x72, 0x70, 0x63, 0x74, 0x6f, 0x6f, 0x6c, 0x2f, 0x74, 0x65, 0x73, 0x74, 0x2f, 0x74, 0x65,
	0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1a, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62,
	0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x74, 0x6f, 0x6f, 0x6c, 0x2e,
	0x74, 0x65, 0x73, 0x74, 0x1a, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x74,
	0x6f, 0x6f, 0x6c, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x74, 0x6f, 0x6f, 0x6c, 0x2f, 0x61, 0x75, 0x74,
	0x6f, 0x6d, 0x61, 0x74, 0x61, 0x2f, 0x61, 0x75, 0x74, 0x6f, 0x6d, 0x61, 0x74, 0x61, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76,
	0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x22, 0x0a,
	0x07, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x02, 0x73, 0x31, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x20, 0x01, 0x52, 0x02, 0x73,
	0x31, 0x22, 0xb8, 0x02, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1e,
	0x0a, 0x06, 0x73, 0x63, 0x61, 0x6c, 0x61, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x42, 0x04,
	0x80, 0xf6, 0x2c, 0x02, 0x48, 0x00, 0x52, 0x06, 0x73, 0x63, 0x61, 0x6c, 0x61, 0x72, 0x12, 0x39,
	0x0a, 0x02, 0x78, 0x31, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x21, 0x2e, 0x67, 0x69, 0x74,
	0x6c, 0x61, 0x62, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x74, 0x6f,
	0x6f, 0x6c, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x65, 0x6e, 0x75, 0x6d, 0x31, 0x42, 0x04, 0x80,
	0xf6, 0x2c, 0x03, 0x48, 0x00, 0x52, 0x02, 0x78, 0x31, 0x12, 0x49, 0x0a, 0x04, 0x64, 0x61, 0x74,
	0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x29, 0x2e, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62,
	0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x74, 0x6f, 0x6f, 0x6c, 0x2e,
	0x74, 0x65, 0x73, 0x74, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x44, 0x61,
	0x74, 0x61, 0x42, 0x08, 0x80, 0xf6, 0x2c, 0x03, 0x80, 0xf6, 0x2c, 0x04, 0x48, 0x00, 0x52, 0x04,
	0x64, 0x61, 0x74, 0x61, 0x12, 0x4e, 0x0a, 0x04, 0x6c, 0x61, 0x73, 0x74, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x29, 0x2e, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x2e, 0x61, 0x67, 0x65, 0x6e,
	0x74, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x74, 0x6f, 0x6f, 0x6c, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x4c, 0x61, 0x73, 0x74, 0x42, 0x0d, 0x80,
	0xf6, 0x2c, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01, 0x48, 0x00, 0x52, 0x04,
	0x6c, 0x61, 0x73, 0x74, 0x1a, 0x1a, 0x0a, 0x04, 0x44, 0x61, 0x74, 0x61, 0x12, 0x12, 0x0a, 0x04,
	0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61,
	0x1a, 0x06, 0x0a, 0x04, 0x4c, 0x61, 0x73, 0x74, 0x42, 0x12, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x12, 0x07, 0x88, 0xf6, 0x2c, 0x01, 0xf8, 0x42, 0x01, 0x22, 0x0a, 0x0a, 0x08,
	0x4e, 0x6f, 0x4f, 0x6e, 0x65, 0x6f, 0x66, 0x73, 0x22, 0x7f, 0x0a, 0x09, 0x54, 0x77, 0x6f, 0x4f,
	0x6e, 0x65, 0x6f, 0x66, 0x73, 0x12, 0x12, 0x0a, 0x03, 0x6d, 0x31, 0x31, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x05, 0x48, 0x00, 0x52, 0x03, 0x6d, 0x31, 0x31, 0x12, 0x12, 0x0a, 0x03, 0x6d, 0x31, 0x32,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x48, 0x00, 0x52, 0x03, 0x6d, 0x31, 0x32, 0x12, 0x12, 0x0a,
	0x03, 0x6d, 0x32, 0x31, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x48, 0x01, 0x52, 0x03, 0x6d, 0x32,
	0x31, 0x12, 0x12, 0x0a, 0x03, 0x6d, 0x32, 0x32, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x48, 0x01,
	0x52, 0x03, 0x6d, 0x32, 0x32, 0x42, 0x10, 0x0a, 0x08, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x31, 0x12, 0x04, 0x88, 0xf6, 0x2c, 0x01, 0x42, 0x10, 0x0a, 0x08, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x32, 0x12, 0x04, 0x88, 0xf6, 0x2c, 0x03, 0x22, 0xae, 0x01, 0x0a, 0x0e, 0x54, 0x77,
	0x6f, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x4f, 0x6e, 0x65, 0x6f, 0x66, 0x73, 0x12, 0x18, 0x0a, 0x03,
	0x6d, 0x31, 0x31, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x42, 0x04, 0x80, 0xf6, 0x2c, 0x02, 0x48,
	0x00, 0x52, 0x03, 0x6d, 0x31, 0x31, 0x12, 0x21, 0x0a, 0x03, 0x6d, 0x31, 0x32, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x05, 0x42, 0x0d, 0x80, 0xf6, 0x2c, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0x01, 0x48, 0x00, 0x52, 0x03, 0x6d, 0x31, 0x32, 0x12, 0x18, 0x0a, 0x03, 0x6d, 0x32, 0x31,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x42, 0x04, 0x80, 0xf6, 0x2c, 0x04, 0x48, 0x01, 0x52, 0x03,
	0x6d, 0x32, 0x31, 0x12, 0x21, 0x0a, 0x03, 0x6d, 0x32, 0x32, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05,
	0x42, 0x0d, 0x80, 0xf6, 0x2c, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01, 0x48,
	0x01, 0x52, 0x03, 0x6d, 0x32, 0x32, 0x42, 0x10, 0x0a, 0x08, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x31, 0x12, 0x04, 0x88, 0xf6, 0x2c, 0x01, 0x42, 0x10, 0x0a, 0x08, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x32, 0x12, 0x04, 0x88, 0xf6, 0x2c, 0x03, 0x22, 0x64, 0x0a, 0x0a, 0x4f, 0x75,
	0x74, 0x4f, 0x66, 0x4f, 0x6e, 0x65, 0x6f, 0x66, 0x12, 0x0c, 0x0a, 0x01, 0x78, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x01, 0x78, 0x12, 0x16, 0x0a, 0x02, 0x6d, 0x31, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x05, 0x42, 0x04, 0x80, 0xf6, 0x2c, 0x01, 0x48, 0x00, 0x52, 0x02, 0x6d, 0x31, 0x12, 0x1f,
	0x0a, 0x02, 0x6d, 0x32, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x42, 0x0d, 0x80, 0xf6, 0x2c, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01, 0x48, 0x00, 0x52, 0x02, 0x6d, 0x32, 0x42,
	0x0f, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x04, 0x88, 0xf6, 0x2c, 0x02,
	0x22, 0x73, 0x0a, 0x0f, 0x4e, 0x6f, 0x74, 0x41, 0x6c, 0x6c, 0x52, 0x65, 0x61, 0x63, 0x68, 0x61,
	0x62, 0x6c, 0x65, 0x12, 0x16, 0x0a, 0x02, 0x6d, 0x31, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x42,
	0x04, 0x80, 0xf6, 0x2c, 0x02, 0x48, 0x00, 0x52, 0x02, 0x6d, 0x31, 0x12, 0x16, 0x0a, 0x02, 0x6d,
	0x32, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x42, 0x04, 0x80, 0xf6, 0x2c, 0x01, 0x48, 0x00, 0x52,
	0x02, 0x6d, 0x32, 0x12, 0x1f, 0x0a, 0x02, 0x6d, 0x33, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x42,
	0x0d, 0x80, 0xf6, 0x2c, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01, 0x48, 0x00,
	0x52, 0x02, 0x6d, 0x33, 0x42, 0x0f, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12,
	0x04, 0x88, 0xf6, 0x2c, 0x03, 0x2a, 0x17, 0x0a, 0x05, 0x65, 0x6e, 0x75, 0x6d, 0x31, 0x12, 0x06,
	0x0a, 0x02, 0x76, 0x31, 0x10, 0x00, 0x12, 0x06, 0x0a, 0x02, 0x76, 0x32, 0x10, 0x01, 0x32, 0xd6,
	0x01, 0x0a, 0x07, 0x54, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x12, 0x5e, 0x0a, 0x0f, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x23, 0x2e,
	0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x67, 0x72, 0x70,
	0x63, 0x74, 0x6f, 0x6f, 0x6c, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x24, 0x2e, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x2e, 0x61, 0x67, 0x65, 0x6e,
	0x74, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x74, 0x6f, 0x6f, 0x6c, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x6b, 0x0a, 0x18, 0x53, 0x74,
	0x72, 0x65, 0x61, 0x6d, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x23, 0x2e, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x2e,
	0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x74, 0x6f, 0x6f, 0x6c, 0x2e, 0x74,
	0x65, 0x73, 0x74, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x24, 0x2e, 0x67, 0x69,
	0x74, 0x6c, 0x61, 0x62, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x74,
	0x6f, 0x6f, 0x6c, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x28, 0x01, 0x30, 0x01, 0x42, 0x61, 0x5a, 0x5f, 0x67, 0x69, 0x74, 0x6c, 0x61,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x2d, 0x6f, 0x72, 0x67,
	0x2f, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x2d, 0x69, 0x6e, 0x74, 0x65, 0x67, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x2d, 0x61, 0x67, 0x65, 0x6e,
	0x74, 0x2f, 0x76, 0x31, 0x34, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x74,
	0x6f, 0x6f, 0x6c, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x74, 0x6f, 0x6f, 0x6c, 0x2f, 0x61, 0x75, 0x74,
	0x6f, 0x6d, 0x61, 0x74, 0x61, 0x2f, 0x74, 0x65, 0x73, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_internal_tool_grpctool_test_test_proto_rawDescOnce sync.Once
	file_internal_tool_grpctool_test_test_proto_rawDescData = file_internal_tool_grpctool_test_test_proto_rawDesc
)

func file_internal_tool_grpctool_test_test_proto_rawDescGZIP() []byte {
	file_internal_tool_grpctool_test_test_proto_rawDescOnce.Do(func() {
		file_internal_tool_grpctool_test_test_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_tool_grpctool_test_test_proto_rawDescData)
	})
	return file_internal_tool_grpctool_test_test_proto_rawDescData
}

var file_internal_tool_grpctool_test_test_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_internal_tool_grpctool_test_test_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_internal_tool_grpctool_test_test_proto_goTypes = []interface{}{
	(Enum1)(0),              // 0: gitlab.agent.grpctool.test.enum1
	(*Request)(nil),         // 1: gitlab.agent.grpctool.test.Request
	(*Response)(nil),        // 2: gitlab.agent.grpctool.test.Response
	(*NoOneofs)(nil),        // 3: gitlab.agent.grpctool.test.NoOneofs
	(*TwoOneofs)(nil),       // 4: gitlab.agent.grpctool.test.TwoOneofs
	(*TwoValidOneofs)(nil),  // 5: gitlab.agent.grpctool.test.TwoValidOneofs
	(*OutOfOneof)(nil),      // 6: gitlab.agent.grpctool.test.OutOfOneof
	(*NotAllReachable)(nil), // 7: gitlab.agent.grpctool.test.NotAllReachable
	(*Response_Data)(nil),   // 8: gitlab.agent.grpctool.test.Response.Data
	(*Response_Last)(nil),   // 9: gitlab.agent.grpctool.test.Response.Last
}
var file_internal_tool_grpctool_test_test_proto_depIdxs = []int32{
	0, // 0: gitlab.agent.grpctool.test.Response.x1:type_name -> gitlab.agent.grpctool.test.enum1
	8, // 1: gitlab.agent.grpctool.test.Response.data:type_name -> gitlab.agent.grpctool.test.Response.Data
	9, // 2: gitlab.agent.grpctool.test.Response.last:type_name -> gitlab.agent.grpctool.test.Response.Last
	1, // 3: gitlab.agent.grpctool.test.Testing.RequestResponse:input_type -> gitlab.agent.grpctool.test.Request
	1, // 4: gitlab.agent.grpctool.test.Testing.StreamingRequestResponse:input_type -> gitlab.agent.grpctool.test.Request
	2, // 5: gitlab.agent.grpctool.test.Testing.RequestResponse:output_type -> gitlab.agent.grpctool.test.Response
	2, // 6: gitlab.agent.grpctool.test.Testing.StreamingRequestResponse:output_type -> gitlab.agent.grpctool.test.Response
	5, // [5:7] is the sub-list for method output_type
	3, // [3:5] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_internal_tool_grpctool_test_test_proto_init() }
func file_internal_tool_grpctool_test_test_proto_init() {
	if File_internal_tool_grpctool_test_test_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_tool_grpctool_test_test_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
		file_internal_tool_grpctool_test_test_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
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
		file_internal_tool_grpctool_test_test_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NoOneofs); i {
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
		file_internal_tool_grpctool_test_test_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TwoOneofs); i {
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
		file_internal_tool_grpctool_test_test_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TwoValidOneofs); i {
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
		file_internal_tool_grpctool_test_test_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OutOfOneof); i {
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
		file_internal_tool_grpctool_test_test_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NotAllReachable); i {
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
		file_internal_tool_grpctool_test_test_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Response_Data); i {
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
		file_internal_tool_grpctool_test_test_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Response_Last); i {
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
	file_internal_tool_grpctool_test_test_proto_msgTypes[1].OneofWrappers = []interface{}{
		(*Response_Scalar)(nil),
		(*Response_X1)(nil),
		(*Response_Data_)(nil),
		(*Response_Last_)(nil),
	}
	file_internal_tool_grpctool_test_test_proto_msgTypes[3].OneofWrappers = []interface{}{
		(*TwoOneofs_M11)(nil),
		(*TwoOneofs_M12)(nil),
		(*TwoOneofs_M21)(nil),
		(*TwoOneofs_M22)(nil),
	}
	file_internal_tool_grpctool_test_test_proto_msgTypes[4].OneofWrappers = []interface{}{
		(*TwoValidOneofs_M11)(nil),
		(*TwoValidOneofs_M12)(nil),
		(*TwoValidOneofs_M21)(nil),
		(*TwoValidOneofs_M22)(nil),
	}
	file_internal_tool_grpctool_test_test_proto_msgTypes[5].OneofWrappers = []interface{}{
		(*OutOfOneof_M1)(nil),
		(*OutOfOneof_M2)(nil),
	}
	file_internal_tool_grpctool_test_test_proto_msgTypes[6].OneofWrappers = []interface{}{
		(*NotAllReachable_M1)(nil),
		(*NotAllReachable_M2)(nil),
		(*NotAllReachable_M3)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_internal_tool_grpctool_test_test_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_internal_tool_grpctool_test_test_proto_goTypes,
		DependencyIndexes: file_internal_tool_grpctool_test_test_proto_depIdxs,
		EnumInfos:         file_internal_tool_grpctool_test_test_proto_enumTypes,
		MessageInfos:      file_internal_tool_grpctool_test_test_proto_msgTypes,
	}.Build()
	File_internal_tool_grpctool_test_test_proto = out.File
	file_internal_tool_grpctool_test_test_proto_rawDesc = nil
	file_internal_tool_grpctool_test_test_proto_goTypes = nil
	file_internal_tool_grpctool_test_test_proto_depIdxs = nil
}
