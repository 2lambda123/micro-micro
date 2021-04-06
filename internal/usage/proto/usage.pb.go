// Code generated by protoc-gen-go. DO NOT EDIT.
// source: internal/usage/proto/usage.proto

package usage

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

type Usage struct {
	// service name
	Service string `protobuf:"bytes,1,opt,name=service,proto3" json:"service,omitempty"`
	// version of service
	Version string `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
	// unique service id
	Id string `protobuf:"bytes,3,opt,name=id,proto3" json:"id,omitempty"`
	// unix timestamp of report
	Timestamp uint64 `protobuf:"varint,4,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	// window of report in seconds
	Window uint64 `protobuf:"varint,5,opt,name=window,proto3" json:"window,omitempty"`
	// usage metrics
	Metrics              *Metrics `protobuf:"bytes,6,opt,name=metrics,proto3" json:"metrics,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Usage) Reset()         { *m = Usage{} }
func (m *Usage) String() string { return proto.CompactTextString(m) }
func (*Usage) ProtoMessage()    {}
func (*Usage) Descriptor() ([]byte, []int) {
	return fileDescriptor_57da1bcc49fd7460, []int{0}
}

func (m *Usage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Usage.Unmarshal(m, b)
}
func (m *Usage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Usage.Marshal(b, m, deterministic)
}
func (m *Usage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Usage.Merge(m, src)
}
func (m *Usage) XXX_Size() int {
	return xxx_messageInfo_Usage.Size(m)
}
func (m *Usage) XXX_DiscardUnknown() {
	xxx_messageInfo_Usage.DiscardUnknown(m)
}

var xxx_messageInfo_Usage proto.InternalMessageInfo

func (m *Usage) GetService() string {
	if m != nil {
		return m.Service
	}
	return ""
}

func (m *Usage) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func (m *Usage) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Usage) GetTimestamp() uint64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *Usage) GetWindow() uint64 {
	if m != nil {
		return m.Window
	}
	return 0
}

func (m *Usage) GetMetrics() *Metrics {
	if m != nil {
		return m.Metrics
	}
	return nil
}

type Metrics struct {
	// counts such as requests, services, etc
	Count                map[string]uint64 `protobuf:"bytes,1,rep,name=count,proto3" json:"count,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Metrics) Reset()         { *m = Metrics{} }
func (m *Metrics) String() string { return proto.CompactTextString(m) }
func (*Metrics) ProtoMessage()    {}
func (*Metrics) Descriptor() ([]byte, []int) {
	return fileDescriptor_57da1bcc49fd7460, []int{1}
}

func (m *Metrics) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Metrics.Unmarshal(m, b)
}
func (m *Metrics) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Metrics.Marshal(b, m, deterministic)
}
func (m *Metrics) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Metrics.Merge(m, src)
}
func (m *Metrics) XXX_Size() int {
	return xxx_messageInfo_Metrics.Size(m)
}
func (m *Metrics) XXX_DiscardUnknown() {
	xxx_messageInfo_Metrics.DiscardUnknown(m)
}

var xxx_messageInfo_Metrics proto.InternalMessageInfo

func (m *Metrics) GetCount() map[string]uint64 {
	if m != nil {
		return m.Count
	}
	return nil
}

func init() {
	proto.RegisterType((*Usage)(nil), "Usage")
	proto.RegisterType((*Metrics)(nil), "Metrics")
	proto.RegisterMapType((map[string]uint64)(nil), "Metrics.CountEntry")
}

func init() { proto.RegisterFile("internal/usage/proto/usage.proto", fileDescriptor_57da1bcc49fd7460) }

var fileDescriptor_57da1bcc49fd7460 = []byte{
	// 236 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x44, 0x90, 0xcd, 0x4a, 0xc4, 0x30,
	0x14, 0x85, 0x49, 0x7f, 0x9d, 0x3b, 0x20, 0x12, 0x45, 0x82, 0xb8, 0x08, 0x5d, 0xd5, 0x4d, 0x07,
	0xc6, 0xcd, 0xe0, 0x56, 0x5c, 0xba, 0x09, 0xf8, 0x00, 0xb5, 0xbd, 0xc8, 0xc5, 0x69, 0x32, 0x24,
	0x69, 0x87, 0x79, 0x20, 0xdf, 0x53, 0x9a, 0xb4, 0xcc, 0xee, 0x7c, 0xdf, 0xb9, 0x9b, 0x73, 0x41,
	0x92, 0xf6, 0x68, 0x75, 0x7b, 0xdc, 0x8d, 0xae, 0xfd, 0xc1, 0xdd, 0xc9, 0x1a, 0x6f, 0x62, 0x6e,
	0x42, 0xae, 0xfe, 0x18, 0xe4, 0x5f, 0x33, 0x73, 0x01, 0xa5, 0x43, 0x3b, 0x51, 0x87, 0x82, 0x49,
	0x56, 0x6f, 0xd4, 0x8a, 0x73, 0x33, 0xa1, 0x75, 0x64, 0xb4, 0x48, 0x62, 0xb3, 0x20, 0xbf, 0x85,
	0x84, 0x7a, 0x91, 0x06, 0x99, 0x50, 0xcf, 0x9f, 0x61, 0xe3, 0x69, 0x40, 0xe7, 0xdb, 0xe1, 0x24,
	0x32, 0xc9, 0xea, 0x4c, 0x5d, 0x05, 0x7f, 0x84, 0xe2, 0x4c, 0xba, 0x37, 0x67, 0x91, 0x87, 0x6a,
	0x21, 0x5e, 0x41, 0x39, 0xa0, 0xb7, 0xd4, 0x39, 0x51, 0x48, 0x56, 0x6f, 0xf7, 0x37, 0xcd, 0x67,
	0x64, 0xb5, 0x16, 0x95, 0x86, 0x72, 0x71, 0xfc, 0x05, 0xf2, 0xce, 0x8c, 0xda, 0x0b, 0x26, 0xd3,
	0x7a, 0xbb, 0xbf, 0x5f, 0x8f, 0x9b, 0xf7, 0xd9, 0x7e, 0x68, 0x6f, 0x2f, 0x2a, 0x5e, 0x3c, 0x1d,
	0x00, 0xae, 0x92, 0xdf, 0x41, 0xfa, 0x8b, 0x97, 0x65, 0xdd, 0x1c, 0xf9, 0x03, 0xe4, 0x53, 0x7b,
	0x1c, 0x31, 0xec, 0xca, 0x54, 0x84, 0xb7, 0xe4, 0xc0, 0xbe, 0x8b, 0xf0, 0x9e, 0xd7, 0xff, 0x00,
	0x00, 0x00, 0xff, 0xff, 0x8d, 0x2b, 0xa5, 0xb9, 0x42, 0x01, 0x00, 0x00,
}
