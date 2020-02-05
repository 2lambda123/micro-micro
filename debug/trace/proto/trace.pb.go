// Code generated by protoc-gen-go. DO NOT EDIT.
// source: micro/micro/debug/trace/proto/trace.proto

package go_micro_debug_trace

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

// Service describes a service running in the micro network.
type Service struct {
	// Service name, e.g. go.micro.service.greeter
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Version              string   `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
	Node                 *Node    `protobuf:"bytes,3,opt,name=node,proto3" json:"node,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Service) Reset()         { *m = Service{} }
func (m *Service) String() string { return proto.CompactTextString(m) }
func (*Service) ProtoMessage()    {}
func (*Service) Descriptor() ([]byte, []int) {
	return fileDescriptor_6510241a5452b8ec, []int{0}
}

func (m *Service) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Service.Unmarshal(m, b)
}
func (m *Service) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Service.Marshal(b, m, deterministic)
}
func (m *Service) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Service.Merge(m, src)
}
func (m *Service) XXX_Size() int {
	return xxx_messageInfo_Service.Size(m)
}
func (m *Service) XXX_DiscardUnknown() {
	xxx_messageInfo_Service.DiscardUnknown(m)
}

var xxx_messageInfo_Service proto.InternalMessageInfo

func (m *Service) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Service) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func (m *Service) GetNode() *Node {
	if m != nil {
		return m.Node
	}
	return nil
}

// Node describes a single instance of a service.
type Node struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Address              string   `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Node) Reset()         { *m = Node{} }
func (m *Node) String() string { return proto.CompactTextString(m) }
func (*Node) ProtoMessage()    {}
func (*Node) Descriptor() ([]byte, []int) {
	return fileDescriptor_6510241a5452b8ec, []int{1}
}

func (m *Node) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Node.Unmarshal(m, b)
}
func (m *Node) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Node.Marshal(b, m, deterministic)
}
func (m *Node) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Node.Merge(m, src)
}
func (m *Node) XXX_Size() int {
	return xxx_messageInfo_Node.Size(m)
}
func (m *Node) XXX_DiscardUnknown() {
	xxx_messageInfo_Node.DiscardUnknown(m)
}

var xxx_messageInfo_Node proto.InternalMessageInfo

func (m *Node) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Node) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

// Snapshot is a snapshot of Trace.Read from a particular service when called.
type Snapshot struct {
	// Source of the service where the snapshot was collected from
	Service *Service `protobuf:"bytes,1,opt,name=service,proto3" json:"service,omitempty"`
	// Unix timestamp, e.g. 1575561487
	Spans                []*Span  `protobuf:"bytes,2,rep,name=spans,proto3" json:"spans,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Snapshot) Reset()         { *m = Snapshot{} }
func (m *Snapshot) String() string { return proto.CompactTextString(m) }
func (*Snapshot) ProtoMessage()    {}
func (*Snapshot) Descriptor() ([]byte, []int) {
	return fileDescriptor_6510241a5452b8ec, []int{2}
}

func (m *Snapshot) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Snapshot.Unmarshal(m, b)
}
func (m *Snapshot) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Snapshot.Marshal(b, m, deterministic)
}
func (m *Snapshot) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Snapshot.Merge(m, src)
}
func (m *Snapshot) XXX_Size() int {
	return xxx_messageInfo_Snapshot.Size(m)
}
func (m *Snapshot) XXX_DiscardUnknown() {
	xxx_messageInfo_Snapshot.DiscardUnknown(m)
}

var xxx_messageInfo_Snapshot proto.InternalMessageInfo

func (m *Snapshot) GetService() *Service {
	if m != nil {
		return m.Service
	}
	return nil
}

func (m *Snapshot) GetSpans() []*Span {
	if m != nil {
		return m.Spans
	}
	return nil
}

type Span struct {
	// the trace id
	Trace string `protobuf:"bytes,1,opt,name=trace,proto3" json:"trace,omitempty"`
	// id of the span
	Id string `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	// parent span
	Parent string `protobuf:"bytes,3,opt,name=parent,proto3" json:"parent,omitempty"`
	// name of the resource
	Name string `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`
	// time of start in nanoseconds
	Started uint64 `protobuf:"varint,5,opt,name=started,proto3" json:"started,omitempty"`
	// duration of the execution in nanoseconds
	Duration uint64 `protobuf:"varint,6,opt,name=duration,proto3" json:"duration,omitempty"`
	// associated metadata
	Metadata             map[string]string `protobuf:"bytes,7,rep,name=metadata,proto3" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Span) Reset()         { *m = Span{} }
func (m *Span) String() string { return proto.CompactTextString(m) }
func (*Span) ProtoMessage()    {}
func (*Span) Descriptor() ([]byte, []int) {
	return fileDescriptor_6510241a5452b8ec, []int{3}
}

func (m *Span) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Span.Unmarshal(m, b)
}
func (m *Span) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Span.Marshal(b, m, deterministic)
}
func (m *Span) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Span.Merge(m, src)
}
func (m *Span) XXX_Size() int {
	return xxx_messageInfo_Span.Size(m)
}
func (m *Span) XXX_DiscardUnknown() {
	xxx_messageInfo_Span.DiscardUnknown(m)
}

var xxx_messageInfo_Span proto.InternalMessageInfo

func (m *Span) GetTrace() string {
	if m != nil {
		return m.Trace
	}
	return ""
}

func (m *Span) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Span) GetParent() string {
	if m != nil {
		return m.Parent
	}
	return ""
}

func (m *Span) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Span) GetStarted() uint64 {
	if m != nil {
		return m.Started
	}
	return 0
}

func (m *Span) GetDuration() uint64 {
	if m != nil {
		return m.Duration
	}
	return 0
}

func (m *Span) GetMetadata() map[string]string {
	if m != nil {
		return m.Metadata
	}
	return nil
}

type ReadRequest struct {
	// If set, only return services matching the filter
	Service *Service `protobuf:"bytes,1,opt,name=service,proto3" json:"service,omitempty"`
	// If false, only the current snapshots will be returned.
	// If true, all historical snapshots in memory will be returned.
	Past                 bool     `protobuf:"varint,2,opt,name=past,proto3" json:"past,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ReadRequest) Reset()         { *m = ReadRequest{} }
func (m *ReadRequest) String() string { return proto.CompactTextString(m) }
func (*ReadRequest) ProtoMessage()    {}
func (*ReadRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_6510241a5452b8ec, []int{4}
}

func (m *ReadRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReadRequest.Unmarshal(m, b)
}
func (m *ReadRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReadRequest.Marshal(b, m, deterministic)
}
func (m *ReadRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReadRequest.Merge(m, src)
}
func (m *ReadRequest) XXX_Size() int {
	return xxx_messageInfo_ReadRequest.Size(m)
}
func (m *ReadRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ReadRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ReadRequest proto.InternalMessageInfo

func (m *ReadRequest) GetService() *Service {
	if m != nil {
		return m.Service
	}
	return nil
}

func (m *ReadRequest) GetPast() bool {
	if m != nil {
		return m.Past
	}
	return false
}

type ReadResponse struct {
	Spans                []*Span  `protobuf:"bytes,1,rep,name=spans,proto3" json:"spans,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ReadResponse) Reset()         { *m = ReadResponse{} }
func (m *ReadResponse) String() string { return proto.CompactTextString(m) }
func (*ReadResponse) ProtoMessage()    {}
func (*ReadResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_6510241a5452b8ec, []int{5}
}

func (m *ReadResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReadResponse.Unmarshal(m, b)
}
func (m *ReadResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReadResponse.Marshal(b, m, deterministic)
}
func (m *ReadResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReadResponse.Merge(m, src)
}
func (m *ReadResponse) XXX_Size() int {
	return xxx_messageInfo_ReadResponse.Size(m)
}
func (m *ReadResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ReadResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ReadResponse proto.InternalMessageInfo

func (m *ReadResponse) GetSpans() []*Span {
	if m != nil {
		return m.Spans
	}
	return nil
}

type WriteRequest struct {
	// If set, only return services matching the filter
	Service *Service `protobuf:"bytes,1,opt,name=service,proto3" json:"service,omitempty"`
	// snapshot to write
	Stats                *Snapshot `protobuf:"bytes,2,opt,name=stats,proto3" json:"stats,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *WriteRequest) Reset()         { *m = WriteRequest{} }
func (m *WriteRequest) String() string { return proto.CompactTextString(m) }
func (*WriteRequest) ProtoMessage()    {}
func (*WriteRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_6510241a5452b8ec, []int{6}
}

func (m *WriteRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WriteRequest.Unmarshal(m, b)
}
func (m *WriteRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WriteRequest.Marshal(b, m, deterministic)
}
func (m *WriteRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WriteRequest.Merge(m, src)
}
func (m *WriteRequest) XXX_Size() int {
	return xxx_messageInfo_WriteRequest.Size(m)
}
func (m *WriteRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_WriteRequest.DiscardUnknown(m)
}

var xxx_messageInfo_WriteRequest proto.InternalMessageInfo

func (m *WriteRequest) GetService() *Service {
	if m != nil {
		return m.Service
	}
	return nil
}

func (m *WriteRequest) GetStats() *Snapshot {
	if m != nil {
		return m.Stats
	}
	return nil
}

type WriteResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *WriteResponse) Reset()         { *m = WriteResponse{} }
func (m *WriteResponse) String() string { return proto.CompactTextString(m) }
func (*WriteResponse) ProtoMessage()    {}
func (*WriteResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_6510241a5452b8ec, []int{7}
}

func (m *WriteResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WriteResponse.Unmarshal(m, b)
}
func (m *WriteResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WriteResponse.Marshal(b, m, deterministic)
}
func (m *WriteResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WriteResponse.Merge(m, src)
}
func (m *WriteResponse) XXX_Size() int {
	return xxx_messageInfo_WriteResponse.Size(m)
}
func (m *WriteResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_WriteResponse.DiscardUnknown(m)
}

var xxx_messageInfo_WriteResponse proto.InternalMessageInfo

type StreamRequest struct {
	// If set, only return services matching the filter
	Service *Service `protobuf:"bytes,1,opt,name=service,proto3" json:"service,omitempty"`
	// If set, only return services matching the namespace
	Namespace            string   `protobuf:"bytes,2,opt,name=namespace,proto3" json:"namespace,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StreamRequest) Reset()         { *m = StreamRequest{} }
func (m *StreamRequest) String() string { return proto.CompactTextString(m) }
func (*StreamRequest) ProtoMessage()    {}
func (*StreamRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_6510241a5452b8ec, []int{8}
}

func (m *StreamRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StreamRequest.Unmarshal(m, b)
}
func (m *StreamRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StreamRequest.Marshal(b, m, deterministic)
}
func (m *StreamRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StreamRequest.Merge(m, src)
}
func (m *StreamRequest) XXX_Size() int {
	return xxx_messageInfo_StreamRequest.Size(m)
}
func (m *StreamRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_StreamRequest.DiscardUnknown(m)
}

var xxx_messageInfo_StreamRequest proto.InternalMessageInfo

func (m *StreamRequest) GetService() *Service {
	if m != nil {
		return m.Service
	}
	return nil
}

func (m *StreamRequest) GetNamespace() string {
	if m != nil {
		return m.Namespace
	}
	return ""
}

type StreamResponse struct {
	Stats                []*Snapshot `protobuf:"bytes,1,rep,name=stats,proto3" json:"stats,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *StreamResponse) Reset()         { *m = StreamResponse{} }
func (m *StreamResponse) String() string { return proto.CompactTextString(m) }
func (*StreamResponse) ProtoMessage()    {}
func (*StreamResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_6510241a5452b8ec, []int{9}
}

func (m *StreamResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StreamResponse.Unmarshal(m, b)
}
func (m *StreamResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StreamResponse.Marshal(b, m, deterministic)
}
func (m *StreamResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StreamResponse.Merge(m, src)
}
func (m *StreamResponse) XXX_Size() int {
	return xxx_messageInfo_StreamResponse.Size(m)
}
func (m *StreamResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_StreamResponse.DiscardUnknown(m)
}

var xxx_messageInfo_StreamResponse proto.InternalMessageInfo

func (m *StreamResponse) GetStats() []*Snapshot {
	if m != nil {
		return m.Stats
	}
	return nil
}

func init() {
	proto.RegisterType((*Service)(nil), "go.micro.debug.trace.Service")
	proto.RegisterType((*Node)(nil), "go.micro.debug.trace.Node")
	proto.RegisterType((*Snapshot)(nil), "go.micro.debug.trace.Snapshot")
	proto.RegisterType((*Span)(nil), "go.micro.debug.trace.Span")
	proto.RegisterMapType((map[string]string)(nil), "go.micro.debug.trace.Span.MetadataEntry")
	proto.RegisterType((*ReadRequest)(nil), "go.micro.debug.trace.ReadRequest")
	proto.RegisterType((*ReadResponse)(nil), "go.micro.debug.trace.ReadResponse")
	proto.RegisterType((*WriteRequest)(nil), "go.micro.debug.trace.WriteRequest")
	proto.RegisterType((*WriteResponse)(nil), "go.micro.debug.trace.WriteResponse")
	proto.RegisterType((*StreamRequest)(nil), "go.micro.debug.trace.StreamRequest")
	proto.RegisterType((*StreamResponse)(nil), "go.micro.debug.trace.StreamResponse")
}

func init() {
	proto.RegisterFile("micro/micro/debug/trace/proto/trace.proto", fileDescriptor_6510241a5452b8ec)
}

var fileDescriptor_6510241a5452b8ec = []byte{
	// 507 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xa4, 0x54, 0xcf, 0x8b, 0xd3, 0x40,
	0x14, 0x26, 0x69, 0xd2, 0x1f, 0xaf, 0xdb, 0x55, 0x86, 0x45, 0x42, 0x50, 0xd1, 0xd1, 0xc3, 0x7a,
	0x49, 0x4b, 0x15, 0x14, 0xbd, 0x78, 0x50, 0x6f, 0x2b, 0x32, 0x55, 0x04, 0x6f, 0xb3, 0xcd, 0xb3,
	0x06, 0x6d, 0x12, 0x67, 0x26, 0x85, 0x3d, 0xf8, 0x47, 0xf8, 0x17, 0x6b, 0xe6, 0x47, 0xb2, 0x2d,
	0x34, 0xa5, 0xd0, 0x4b, 0x78, 0xdf, 0xcc, 0x37, 0xef, 0xbd, 0xef, 0x7b, 0x8f, 0xc0, 0xb3, 0x75,
	0xb6, 0x14, 0xc5, 0xd4, 0x7e, 0x53, 0xbc, 0xae, 0x56, 0x53, 0x25, 0xf8, 0x12, 0xa7, 0xa5, 0x28,
	0x54, 0x61, 0xe3, 0xc4, 0xc4, 0xe4, 0x62, 0x55, 0x24, 0x86, 0x97, 0x18, 0x5e, 0x62, 0xee, 0xe8,
	0x0a, 0x06, 0x0b, 0x14, 0x9b, 0x6c, 0x89, 0x84, 0x40, 0x90, 0xf3, 0x35, 0x46, 0xde, 0x23, 0xef,
	0x72, 0xc4, 0x4c, 0x4c, 0x22, 0x18, 0x6c, 0x50, 0xc8, 0xac, 0xc8, 0x23, 0xdf, 0x1c, 0x37, 0x90,
	0x24, 0x35, 0xbb, 0x48, 0x31, 0xea, 0xd5, 0xc7, 0xe3, 0x79, 0x9c, 0xec, 0xcb, 0x9e, 0x7c, 0xac,
	0x19, 0xcc, 0xf0, 0xe8, 0x0c, 0x02, 0x8d, 0xc8, 0x39, 0xf8, 0x59, 0xea, 0x6a, 0xd4, 0x91, 0xae,
	0xc0, 0xd3, 0x54, 0xa0, 0x94, 0x4d, 0x05, 0x07, 0x69, 0x05, 0xc3, 0x45, 0xce, 0x4b, 0xf9, 0xa3,
	0x50, 0xe4, 0x25, 0x0c, 0xa4, 0x6d, 0xd3, 0x3c, 0x1d, 0xcf, 0x1f, 0xec, 0x2f, 0xe8, 0xb4, 0xb0,
	0x86, 0x4d, 0x66, 0x10, 0xca, 0x92, 0xe7, 0x3a, 0x79, 0xaf, 0xbb, 0xcf, 0x45, 0x4d, 0x61, 0x96,
	0x48, 0xff, 0xfa, 0x10, 0x68, 0x4c, 0x2e, 0x20, 0x34, 0xb7, 0xae, 0x59, 0x0b, 0x5c, 0xff, 0x7e,
	0xdb, 0xff, 0x3d, 0xe8, 0x97, 0x5c, 0x60, 0xae, 0x8c, 0x13, 0x23, 0xe6, 0x50, 0xeb, 0x66, 0xb0,
	0xeb, 0xa6, 0x54, 0x5c, 0x28, 0x4c, 0xa3, 0xb0, 0x3e, 0x0e, 0x58, 0x03, 0x49, 0x0c, 0xc3, 0xb4,
	0x12, 0x5c, 0x69, 0xa3, 0xfb, 0xe6, 0xaa, 0xc5, 0xe4, 0x1d, 0x0c, 0xd7, 0xa8, 0x78, 0xca, 0x15,
	0x8f, 0x06, 0x46, 0xc5, 0x65, 0xb7, 0x8a, 0xe4, 0xca, 0x51, 0xdf, 0xe7, 0x4a, 0xdc, 0xb0, 0xf6,
	0x65, 0xfc, 0x06, 0x26, 0x3b, 0x57, 0xe4, 0x2e, 0xf4, 0x7e, 0xe2, 0x8d, 0x13, 0xa7, 0x43, 0x2d,
	0x78, 0xc3, 0x7f, 0x55, 0xe8, 0xd4, 0x59, 0xf0, 0xda, 0x7f, 0xe5, 0xd1, 0x6f, 0x30, 0x66, 0xc8,
	0x53, 0x86, 0xbf, 0x2b, 0x94, 0x27, 0x4c, 0xa3, 0x36, 0xa5, 0xe4, 0x52, 0x99, 0x02, 0x43, 0x66,
	0x62, 0xfa, 0x16, 0xce, 0x6c, 0x6e, 0x59, 0x16, 0xb9, 0xdc, 0x9a, 0x98, 0x77, 0xec, 0xc4, 0xfe,
	0xc0, 0xd9, 0x57, 0x91, 0x29, 0x3c, 0xb9, 0xbd, 0x17, 0x75, 0x69, 0xc5, 0x95, 0xdd, 0xc4, 0xf1,
	0xfc, 0x61, 0xc7, 0x33, 0xb7, 0x94, 0xcc, 0x92, 0xe9, 0x1d, 0x98, 0xb8, 0xf2, 0x56, 0x01, 0xfd,
	0x0e, 0x93, 0x85, 0x12, 0xc8, 0xd7, 0x27, 0x37, 0x74, 0x1f, 0x46, 0x7a, 0x71, 0x6a, 0x99, 0xcb,
	0x66, 0x2a, 0xb7, 0x07, 0xf4, 0x03, 0x9c, 0x37, 0x75, 0x9c, 0x77, 0xad, 0x00, 0xeb, 0xdd, 0x71,
	0x02, 0xe6, 0xff, 0x3c, 0x08, 0x3f, 0x9b, 0xe5, 0xbe, 0x82, 0x40, 0xcf, 0x82, 0x3c, 0xde, 0xff,
	0x70, 0x6b, 0x07, 0x62, 0x7a, 0x88, 0xe2, 0xda, 0xf9, 0x04, 0xa1, 0x71, 0x86, 0x74, 0x90, 0xb7,
	0xa7, 0x16, 0x3f, 0x39, 0xc8, 0x71, 0x19, 0xbf, 0x40, 0xdf, 0x4a, 0x26, 0x1d, 0xf4, 0x1d, 0xe3,
	0xe3, 0xa7, 0x87, 0x49, 0x36, 0xe9, 0xcc, 0xbb, 0xee, 0x9b, 0x5f, 0xe4, 0xf3, 0xff, 0x01, 0x00,
	0x00, 0xff, 0xff, 0x48, 0xf2, 0xb2, 0xe7, 0x4f, 0x05, 0x00, 0x00,
}
