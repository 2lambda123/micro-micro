// Code generated by protoc-gen-go. DO NOT EDIT.
// source: auth/api/proto/auth.proto

package go_micro_api_auth

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type ValidateRequest struct {
	Token                string   `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ValidateRequest) Reset()         { *m = ValidateRequest{} }
func (m *ValidateRequest) String() string { return proto.CompactTextString(m) }
func (*ValidateRequest) ProtoMessage()    {}
func (*ValidateRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_6a2d6e66004b67ef, []int{0}
}

func (m *ValidateRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ValidateRequest.Unmarshal(m, b)
}
func (m *ValidateRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ValidateRequest.Marshal(b, m, deterministic)
}
func (m *ValidateRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ValidateRequest.Merge(m, src)
}
func (m *ValidateRequest) XXX_Size() int {
	return xxx_messageInfo_ValidateRequest.Size(m)
}
func (m *ValidateRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ValidateRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ValidateRequest proto.InternalMessageInfo

func (m *ValidateRequest) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

type ValidateResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ValidateResponse) Reset()         { *m = ValidateResponse{} }
func (m *ValidateResponse) String() string { return proto.CompactTextString(m) }
func (*ValidateResponse) ProtoMessage()    {}
func (*ValidateResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_6a2d6e66004b67ef, []int{1}
}

func (m *ValidateResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ValidateResponse.Unmarshal(m, b)
}
func (m *ValidateResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ValidateResponse.Marshal(b, m, deterministic)
}
func (m *ValidateResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ValidateResponse.Merge(m, src)
}
func (m *ValidateResponse) XXX_Size() int {
	return xxx_messageInfo_ValidateResponse.Size(m)
}
func (m *ValidateResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ValidateResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ValidateResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*ValidateRequest)(nil), "go.micro.api.auth.ValidateRequest")
	proto.RegisterType((*ValidateResponse)(nil), "go.micro.api.auth.ValidateResponse")
}

func init() { proto.RegisterFile("auth/api/proto/auth.proto", fileDescriptor_6a2d6e66004b67ef) }

var fileDescriptor_6a2d6e66004b67ef = []byte{
	// 147 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x4c, 0x2c, 0x2d, 0xc9,
	0xd0, 0x4f, 0x2c, 0xc8, 0xd4, 0x2f, 0x28, 0xca, 0x2f, 0xc9, 0xd7, 0x07, 0x71, 0xf5, 0xc0, 0x4c,
	0x21, 0xc1, 0xf4, 0x7c, 0xbd, 0xdc, 0xcc, 0xe4, 0xa2, 0x7c, 0xbd, 0xc4, 0x82, 0x4c, 0x3d, 0x90,
	0x84, 0x92, 0x3a, 0x17, 0x7f, 0x58, 0x62, 0x4e, 0x66, 0x4a, 0x62, 0x49, 0x6a, 0x50, 0x6a, 0x61,
	0x69, 0x6a, 0x71, 0x89, 0x90, 0x08, 0x17, 0x6b, 0x49, 0x7e, 0x76, 0x6a, 0x9e, 0x04, 0xa3, 0x02,
	0xa3, 0x06, 0x67, 0x10, 0x84, 0xa3, 0x24, 0xc4, 0x25, 0x80, 0x50, 0x58, 0x5c, 0x90, 0x9f, 0x57,
	0x9c, 0x6a, 0x14, 0xcb, 0xc5, 0xe2, 0x58, 0x5a, 0x92, 0x21, 0x14, 0xca, 0xc5, 0x01, 0x93, 0x13,
	0x52, 0xd2, 0xc3, 0xb0, 0x44, 0x0f, 0xcd, 0x06, 0x29, 0x65, 0xbc, 0x6a, 0x20, 0x86, 0x2b, 0x31,
	0x24, 0xb1, 0x81, 0x5d, 0x6d, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0xfd, 0x61, 0x03, 0x27, 0xd2,
	0x00, 0x00, 0x00,
}
