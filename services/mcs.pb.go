// Code generated by protoc-gen-go. DO NOT EDIT.
// source: mcs.proto

package services

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type BackupStatus int32

const (
	BackupStatus_FAILED    BackupStatus = 0
	BackupStatus_ZIPPING   BackupStatus = 1
	BackupStatus_UPLOADING BackupStatus = 2
	BackupStatus_DONE      BackupStatus = 3
)

var BackupStatus_name = map[int32]string{
	0: "FAILED",
	1: "ZIPPING",
	2: "UPLOADING",
	3: "DONE",
}

var BackupStatus_value = map[string]int32{
	"FAILED":    0,
	"ZIPPING":   1,
	"UPLOADING": 2,
	"DONE":      3,
}

func (x BackupStatus) String() string {
	return proto.EnumName(BackupStatus_name, int32(x))
}

func (BackupStatus) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_b92e0ece524b9239, []int{0}
}

type EmptyReq struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EmptyReq) Reset()         { *m = EmptyReq{} }
func (m *EmptyReq) String() string { return proto.CompactTextString(m) }
func (*EmptyReq) ProtoMessage()    {}
func (*EmptyReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_b92e0ece524b9239, []int{0}
}

func (m *EmptyReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EmptyReq.Unmarshal(m, b)
}
func (m *EmptyReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EmptyReq.Marshal(b, m, deterministic)
}
func (m *EmptyReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EmptyReq.Merge(m, src)
}
func (m *EmptyReq) XXX_Size() int {
	return xxx_messageInfo_EmptyReq.Size(m)
}
func (m *EmptyReq) XXX_DiscardUnknown() {
	xxx_messageInfo_EmptyReq.DiscardUnknown(m)
}

var xxx_messageInfo_EmptyReq proto.InternalMessageInfo

type PingReq struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PingReq) Reset()         { *m = PingReq{} }
func (m *PingReq) String() string { return proto.CompactTextString(m) }
func (*PingReq) ProtoMessage()    {}
func (*PingReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_b92e0ece524b9239, []int{1}
}

func (m *PingReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PingReq.Unmarshal(m, b)
}
func (m *PingReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PingReq.Marshal(b, m, deterministic)
}
func (m *PingReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PingReq.Merge(m, src)
}
func (m *PingReq) XXX_Size() int {
	return xxx_messageInfo_PingReq.Size(m)
}
func (m *PingReq) XXX_DiscardUnknown() {
	xxx_messageInfo_PingReq.DiscardUnknown(m)
}

var xxx_messageInfo_PingReq proto.InternalMessageInfo

type PongResp struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PongResp) Reset()         { *m = PongResp{} }
func (m *PongResp) String() string { return proto.CompactTextString(m) }
func (*PongResp) ProtoMessage()    {}
func (*PongResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_b92e0ece524b9239, []int{2}
}

func (m *PongResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PongResp.Unmarshal(m, b)
}
func (m *PongResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PongResp.Marshal(b, m, deterministic)
}
func (m *PongResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PongResp.Merge(m, src)
}
func (m *PongResp) XXX_Size() int {
	return xxx_messageInfo_PongResp.Size(m)
}
func (m *PongResp) XXX_DiscardUnknown() {
	xxx_messageInfo_PongResp.DiscardUnknown(m)
}

var xxx_messageInfo_PongResp proto.InternalMessageInfo

type ServiceResp struct {
	RequestId            int64    `protobuf:"varint,1,opt,name=request_id,json=requestId,proto3" json:"request_id,omitempty"`
	Status               int64    `protobuf:"varint,2,opt,name=status,proto3" json:"status,omitempty"`
	Message              string   `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ServiceResp) Reset()         { *m = ServiceResp{} }
func (m *ServiceResp) String() string { return proto.CompactTextString(m) }
func (*ServiceResp) ProtoMessage()    {}
func (*ServiceResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_b92e0ece524b9239, []int{3}
}

func (m *ServiceResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ServiceResp.Unmarshal(m, b)
}
func (m *ServiceResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ServiceResp.Marshal(b, m, deterministic)
}
func (m *ServiceResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ServiceResp.Merge(m, src)
}
func (m *ServiceResp) XXX_Size() int {
	return xxx_messageInfo_ServiceResp.Size(m)
}
func (m *ServiceResp) XXX_DiscardUnknown() {
	xxx_messageInfo_ServiceResp.DiscardUnknown(m)
}

var xxx_messageInfo_ServiceResp proto.InternalMessageInfo

func (m *ServiceResp) GetRequestId() int64 {
	if m != nil {
		return m.RequestId
	}
	return 0
}

func (m *ServiceResp) GetStatus() int64 {
	if m != nil {
		return m.Status
	}
	return 0
}

func (m *ServiceResp) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type StartConfig struct {
	MemAlloc             int64    `protobuf:"varint,1,opt,name=mem_alloc,json=memAlloc,proto3" json:"mem_alloc,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StartConfig) Reset()         { *m = StartConfig{} }
func (m *StartConfig) String() string { return proto.CompactTextString(m) }
func (*StartConfig) ProtoMessage()    {}
func (*StartConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_b92e0ece524b9239, []int{4}
}

func (m *StartConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StartConfig.Unmarshal(m, b)
}
func (m *StartConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StartConfig.Marshal(b, m, deterministic)
}
func (m *StartConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StartConfig.Merge(m, src)
}
func (m *StartConfig) XXX_Size() int {
	return xxx_messageInfo_StartConfig.Size(m)
}
func (m *StartConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_StartConfig.DiscardUnknown(m)
}

var xxx_messageInfo_StartConfig proto.InternalMessageInfo

func (m *StartConfig) GetMemAlloc() int64 {
	if m != nil {
		return m.MemAlloc
	}
	return 0
}

type StatusResp struct {
	ServerState          string   `protobuf:"bytes,1,opt,name=server_state,json=serverState,proto3" json:"server_state,omitempty"`
	Message              string   `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StatusResp) Reset()         { *m = StatusResp{} }
func (m *StatusResp) String() string { return proto.CompactTextString(m) }
func (*StatusResp) ProtoMessage()    {}
func (*StatusResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_b92e0ece524b9239, []int{5}
}

func (m *StatusResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StatusResp.Unmarshal(m, b)
}
func (m *StatusResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StatusResp.Marshal(b, m, deterministic)
}
func (m *StatusResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StatusResp.Merge(m, src)
}
func (m *StatusResp) XXX_Size() int {
	return xxx_messageInfo_StatusResp.Size(m)
}
func (m *StatusResp) XXX_DiscardUnknown() {
	xxx_messageInfo_StatusResp.DiscardUnknown(m)
}

var xxx_messageInfo_StatusResp proto.InternalMessageInfo

func (m *StatusResp) GetServerState() string {
	if m != nil {
		return m.ServerState
	}
	return ""
}

func (m *StatusResp) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type BackupResp struct {
	Status               BackupStatus `protobuf:"varint,1,opt,name=status,proto3,enum=services.BackupStatus" json:"status,omitempty"`
	Message              string       `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	LinkUrl              string       `protobuf:"bytes,3,opt,name=link_url,json=linkUrl,proto3" json:"link_url,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *BackupResp) Reset()         { *m = BackupResp{} }
func (m *BackupResp) String() string { return proto.CompactTextString(m) }
func (*BackupResp) ProtoMessage()    {}
func (*BackupResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_b92e0ece524b9239, []int{6}
}

func (m *BackupResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BackupResp.Unmarshal(m, b)
}
func (m *BackupResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BackupResp.Marshal(b, m, deterministic)
}
func (m *BackupResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BackupResp.Merge(m, src)
}
func (m *BackupResp) XXX_Size() int {
	return xxx_messageInfo_BackupResp.Size(m)
}
func (m *BackupResp) XXX_DiscardUnknown() {
	xxx_messageInfo_BackupResp.DiscardUnknown(m)
}

var xxx_messageInfo_BackupResp proto.InternalMessageInfo

func (m *BackupResp) GetStatus() BackupStatus {
	if m != nil {
		return m.Status
	}
	return BackupStatus_FAILED
}

func (m *BackupResp) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *BackupResp) GetLinkUrl() string {
	if m != nil {
		return m.LinkUrl
	}
	return ""
}

func init() {
	proto.RegisterEnum("services.BackupStatus", BackupStatus_name, BackupStatus_value)
	proto.RegisterType((*EmptyReq)(nil), "services.EmptyReq")
	proto.RegisterType((*PingReq)(nil), "services.PingReq")
	proto.RegisterType((*PongResp)(nil), "services.PongResp")
	proto.RegisterType((*ServiceResp)(nil), "services.ServiceResp")
	proto.RegisterType((*StartConfig)(nil), "services.StartConfig")
	proto.RegisterType((*StatusResp)(nil), "services.StatusResp")
	proto.RegisterType((*BackupResp)(nil), "services.BackupResp")
}

func init() { proto.RegisterFile("mcs.proto", fileDescriptor_b92e0ece524b9239) }

var fileDescriptor_b92e0ece524b9239 = []byte{
	// 396 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x52, 0x41, 0xcf, 0x93, 0x40,
	0x10, 0x05, 0x5a, 0x29, 0x0c, 0xd5, 0xd4, 0x89, 0x36, 0x58, 0x63, 0x52, 0x39, 0x35, 0x3d, 0x60,
	0xd2, 0x9a, 0x78, 0xb5, 0xda, 0x6a, 0x48, 0x6a, 0x4b, 0x20, 0xbd, 0x78, 0x90, 0x20, 0x5d, 0x1b,
	0x52, 0x28, 0x94, 0x5d, 0x4c, 0xfc, 0x25, 0xfe, 0x5d, 0xb3, 0xb0, 0x08, 0x4d, 0xfa, 0xe5, 0xbb,
	0xed, 0x9b, 0x9d, 0x79, 0x33, 0xf3, 0xde, 0x80, 0x9e, 0x46, 0xd4, 0xce, 0x8b, 0x8c, 0x65, 0xa8,
	0x51, 0x52, 0xfc, 0x8e, 0x23, 0x42, 0x2d, 0x00, 0x6d, 0x93, 0xe6, 0xec, 0x8f, 0x47, 0xae, 0x96,
	0x0e, 0x03, 0x37, 0xbe, 0x9c, 0xf8, 0x13, 0x40, 0x73, 0x33, 0xfe, 0xa4, 0xb9, 0xf5, 0x03, 0x0c,
	0xbf, 0x4e, 0xe7, 0x10, 0xdf, 0x00, 0x14, 0xe4, 0x5a, 0x12, 0xca, 0x82, 0xf8, 0x68, 0xca, 0x53,
	0x79, 0xd6, 0xf3, 0x74, 0x11, 0x71, 0x8e, 0x38, 0x06, 0x95, 0xb2, 0x90, 0x95, 0xd4, 0x54, 0xaa,
	0x2f, 0x81, 0xd0, 0x84, 0x41, 0x4a, 0x28, 0x0d, 0x4f, 0xc4, 0xec, 0x4d, 0xe5, 0x99, 0xee, 0x35,
	0xd0, 0x9a, 0x83, 0xe1, 0xb3, 0xb0, 0x60, 0x9f, 0xb3, 0xcb, 0xaf, 0xf8, 0x84, 0xaf, 0x41, 0x4f,
	0x49, 0x1a, 0x84, 0x49, 0x92, 0x45, 0x82, 0x5e, 0x4b, 0x49, 0xba, 0xe2, 0xd8, 0x72, 0x00, 0xfc,
	0x8a, 0xaf, 0x1a, 0xe5, 0x2d, 0x0c, 0xf9, 0x22, 0xa4, 0x08, 0x78, 0x13, 0x52, 0x65, 0xeb, 0x9e,
	0x51, 0xc7, 0x78, 0x1e, 0xe9, 0xb6, 0x55, 0x6e, 0xdb, 0x5e, 0x01, 0x3e, 0x85, 0xd1, 0xb9, 0xcc,
	0x2b, 0x2a, 0xfb, 0xff, 0xd8, 0x9c, 0xe4, 0xd9, 0x62, 0x6c, 0x37, 0x12, 0xd9, 0x75, 0x96, 0x68,
	0x7b, 0x67, 0x9d, 0x5b, 0x5e, 0x7c, 0x05, 0x5a, 0x12, 0x5f, 0xce, 0x41, 0x59, 0x24, 0xcd, 0xa6,
	0x1c, 0x1f, 0x8a, 0x64, 0xfe, 0x11, 0x86, 0x5d, 0x32, 0x04, 0x50, 0xbf, 0xac, 0x9c, 0xed, 0x66,
	0x3d, 0x92, 0xd0, 0x80, 0xc1, 0x77, 0xc7, 0x75, 0x9d, 0xdd, 0xd7, 0x91, 0x8c, 0x4f, 0x41, 0x3f,
	0xb8, 0xdb, 0xfd, 0x6a, 0xcd, 0xa1, 0x82, 0x1a, 0xf4, 0xd7, 0xfb, 0xdd, 0x66, 0xd4, 0x5b, 0xfc,
	0x55, 0x00, 0xbe, 0x45, 0x54, 0xf8, 0x81, 0xef, 0xa0, 0xcf, 0x1d, 0xc3, 0xe7, 0xed, 0xb4, 0xc2,
	0xc1, 0x09, 0x76, 0x42, 0x8d, 0x93, 0x12, 0x7e, 0x80, 0x27, 0x95, 0xd6, 0xf8, 0xb2, 0xfd, 0xee,
	0x88, 0x3f, 0xe9, 0x86, 0x5b, 0xcf, 0x2d, 0x09, 0x97, 0xd0, 0xf7, 0x59, 0x96, 0x63, 0x87, 0xb6,
	0xb9, 0x9b, 0x87, 0x8b, 0xde, 0x83, 0x2a, 0x36, 0xbd, 0x57, 0xf6, 0xe2, 0x66, 0x04, 0xe1, 0x69,
	0x5d, 0x55, 0xab, 0xf4, 0x58, 0x55, 0x6b, 0x9f, 0x25, 0xfd, 0x54, 0xab, 0xcb, 0x5e, 0xfe, 0x0b,
	0x00, 0x00, 0xff, 0xff, 0x20, 0xd8, 0x2c, 0x0e, 0xe6, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// McsServiceClient is the client API for McsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type McsServiceClient interface {
	Ping(ctx context.Context, in *PingReq, opts ...grpc.CallOption) (*PongResp, error)
	Start(ctx context.Context, in *StartConfig, opts ...grpc.CallOption) (*ServiceResp, error)
	Stop(ctx context.Context, in *EmptyReq, opts ...grpc.CallOption) (*ServiceResp, error)
	Status(ctx context.Context, in *EmptyReq, opts ...grpc.CallOption) (*StatusResp, error)
	Backup(ctx context.Context, in *EmptyReq, opts ...grpc.CallOption) (*BackupResp, error)
}

type mcsServiceClient struct {
	cc *grpc.ClientConn
}

func NewMcsServiceClient(cc *grpc.ClientConn) McsServiceClient {
	return &mcsServiceClient{cc}
}

func (c *mcsServiceClient) Ping(ctx context.Context, in *PingReq, opts ...grpc.CallOption) (*PongResp, error) {
	out := new(PongResp)
	err := c.cc.Invoke(ctx, "/services.McsService/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mcsServiceClient) Start(ctx context.Context, in *StartConfig, opts ...grpc.CallOption) (*ServiceResp, error) {
	out := new(ServiceResp)
	err := c.cc.Invoke(ctx, "/services.McsService/Start", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mcsServiceClient) Stop(ctx context.Context, in *EmptyReq, opts ...grpc.CallOption) (*ServiceResp, error) {
	out := new(ServiceResp)
	err := c.cc.Invoke(ctx, "/services.McsService/Stop", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mcsServiceClient) Status(ctx context.Context, in *EmptyReq, opts ...grpc.CallOption) (*StatusResp, error) {
	out := new(StatusResp)
	err := c.cc.Invoke(ctx, "/services.McsService/Status", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mcsServiceClient) Backup(ctx context.Context, in *EmptyReq, opts ...grpc.CallOption) (*BackupResp, error) {
	out := new(BackupResp)
	err := c.cc.Invoke(ctx, "/services.McsService/Backup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// McsServiceServer is the server API for McsService service.
type McsServiceServer interface {
	Ping(context.Context, *PingReq) (*PongResp, error)
	Start(context.Context, *StartConfig) (*ServiceResp, error)
	Stop(context.Context, *EmptyReq) (*ServiceResp, error)
	Status(context.Context, *EmptyReq) (*StatusResp, error)
	Backup(context.Context, *EmptyReq) (*BackupResp, error)
}

// UnimplementedMcsServiceServer can be embedded to have forward compatible implementations.
type UnimplementedMcsServiceServer struct {
}

func (*UnimplementedMcsServiceServer) Ping(ctx context.Context, req *PingReq) (*PongResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (*UnimplementedMcsServiceServer) Start(ctx context.Context, req *StartConfig) (*ServiceResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Start not implemented")
}
func (*UnimplementedMcsServiceServer) Stop(ctx context.Context, req *EmptyReq) (*ServiceResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Stop not implemented")
}
func (*UnimplementedMcsServiceServer) Status(ctx context.Context, req *EmptyReq) (*StatusResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Status not implemented")
}
func (*UnimplementedMcsServiceServer) Backup(ctx context.Context, req *EmptyReq) (*BackupResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Backup not implemented")
}

func RegisterMcsServiceServer(s *grpc.Server, srv McsServiceServer) {
	s.RegisterService(&_McsService_serviceDesc, srv)
}

func _McsService_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(McsServiceServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/services.McsService/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(McsServiceServer).Ping(ctx, req.(*PingReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _McsService_Start_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StartConfig)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(McsServiceServer).Start(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/services.McsService/Start",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(McsServiceServer).Start(ctx, req.(*StartConfig))
	}
	return interceptor(ctx, in, info, handler)
}

func _McsService_Stop_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(McsServiceServer).Stop(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/services.McsService/Stop",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(McsServiceServer).Stop(ctx, req.(*EmptyReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _McsService_Status_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(McsServiceServer).Status(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/services.McsService/Status",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(McsServiceServer).Status(ctx, req.(*EmptyReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _McsService_Backup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(McsServiceServer).Backup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/services.McsService/Backup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(McsServiceServer).Backup(ctx, req.(*EmptyReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _McsService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "services.McsService",
	HandlerType: (*McsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _McsService_Ping_Handler,
		},
		{
			MethodName: "Start",
			Handler:    _McsService_Start_Handler,
		},
		{
			MethodName: "Stop",
			Handler:    _McsService_Stop_Handler,
		},
		{
			MethodName: "Status",
			Handler:    _McsService_Status_Handler,
		},
		{
			MethodName: "Backup",
			Handler:    _McsService_Backup_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "mcs.proto",
}
