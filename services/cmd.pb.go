// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cmd.proto

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

type EmptyReq struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EmptyReq) Reset()         { *m = EmptyReq{} }
func (m *EmptyReq) String() string { return proto.CompactTextString(m) }
func (*EmptyReq) ProtoMessage()    {}
func (*EmptyReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_7520252fb01eaf30, []int{0}
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
	return fileDescriptor_7520252fb01eaf30, []int{1}
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

type LaunchConfig struct {
	MemAlloc             int64    `protobuf:"varint,1,opt,name=mem_alloc,json=memAlloc,proto3" json:"mem_alloc,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LaunchConfig) Reset()         { *m = LaunchConfig{} }
func (m *LaunchConfig) String() string { return proto.CompactTextString(m) }
func (*LaunchConfig) ProtoMessage()    {}
func (*LaunchConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_7520252fb01eaf30, []int{2}
}

func (m *LaunchConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LaunchConfig.Unmarshal(m, b)
}
func (m *LaunchConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LaunchConfig.Marshal(b, m, deterministic)
}
func (m *LaunchConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LaunchConfig.Merge(m, src)
}
func (m *LaunchConfig) XXX_Size() int {
	return xxx_messageInfo_LaunchConfig.Size(m)
}
func (m *LaunchConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_LaunchConfig.DiscardUnknown(m)
}

var xxx_messageInfo_LaunchConfig proto.InternalMessageInfo

func (m *LaunchConfig) GetMemAlloc() int64 {
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
	return fileDescriptor_7520252fb01eaf30, []int{3}
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

func init() {
	proto.RegisterType((*EmptyReq)(nil), "services.EmptyReq")
	proto.RegisterType((*ServiceResp)(nil), "services.ServiceResp")
	proto.RegisterType((*LaunchConfig)(nil), "services.LaunchConfig")
	proto.RegisterType((*StatusResp)(nil), "services.StatusResp")
}

func init() { proto.RegisterFile("cmd.proto", fileDescriptor_7520252fb01eaf30) }

var fileDescriptor_7520252fb01eaf30 = []byte{
	// 267 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x51, 0xcf, 0x4b, 0xc3, 0x30,
	0x14, 0x5e, 0x37, 0xa9, 0xcd, 0xdb, 0x40, 0x78, 0xe8, 0x28, 0x8a, 0x30, 0x73, 0x1a, 0x08, 0x3d,
	0x38, 0x2f, 0x1e, 0x45, 0x3c, 0x0c, 0x3c, 0xa5, 0x77, 0x4b, 0x6d, 0x9f, 0xb3, 0xb0, 0x2c, 0x5d,
	0x92, 0x0a, 0xfe, 0x65, 0xfe, 0x7b, 0xd2, 0xa4, 0x75, 0x15, 0xdc, 0xf1, 0x7d, 0x79, 0xdf, 0xfb,
	0x7e, 0x04, 0x58, 0x21, 0xcb, 0xa4, 0xd6, 0xca, 0x2a, 0x8c, 0x0c, 0xe9, 0xcf, 0xaa, 0x20, 0xc3,
	0x01, 0xa2, 0x67, 0x59, 0xdb, 0x2f, 0x41, 0x7b, 0xfe, 0x0a, 0xd3, 0xd4, 0xe3, 0x82, 0x4c, 0x8d,
	0xd7, 0x00, 0x9a, 0xf6, 0x0d, 0x19, 0x9b, 0x55, 0x65, 0x1c, 0x2c, 0x82, 0xe5, 0x44, 0xb0, 0x0e,
	0x59, 0x97, 0x38, 0x87, 0xd0, 0xd8, 0xdc, 0x36, 0x26, 0x1e, 0xbb, 0xa7, 0x6e, 0xc2, 0x18, 0x4e,
	0x25, 0x19, 0x93, 0x6f, 0x28, 0x9e, 0x2c, 0x82, 0x25, 0x13, 0xfd, 0xc8, 0x6f, 0x61, 0xf6, 0x92,
	0x37, 0xbb, 0xe2, 0xe3, 0x49, 0xed, 0xde, 0xab, 0x0d, 0x5e, 0x01, 0x93, 0x24, 0xb3, 0x7c, 0xbb,
	0x55, 0x45, 0x77, 0x3f, 0x92, 0x24, 0x1f, 0xdb, 0x99, 0xaf, 0x01, 0x52, 0x77, 0xd0, 0x79, 0xb9,
	0x81, 0x59, 0x6b, 0x99, 0x74, 0xd6, 0xaa, 0x90, 0xdb, 0x66, 0x62, 0xea, 0xb1, 0x76, 0x8f, 0x86,
	0xba, 0xe3, 0x3f, 0xba, 0x77, 0xdf, 0x01, 0x9c, 0x79, 0x61, 0xd2, 0x5d, 0x40, 0x7c, 0x80, 0xd0,
	0x43, 0x38, 0x4f, 0xfa, 0x32, 0x92, 0xa1, 0xbb, 0xcb, 0x8b, 0x03, 0x3e, 0x68, 0x85, 0x8f, 0x70,
	0x05, 0x27, 0xa9, 0x55, 0x35, 0xe2, 0x61, 0xa1, 0xaf, 0xf0, 0x38, 0xe9, 0x1e, 0x42, 0x1f, 0xe7,
	0x5f, 0xda, 0xf9, 0x80, 0xf6, 0x1b, 0x9a, 0x8f, 0xde, 0x42, 0xf7, 0x5d, 0xab, 0x9f, 0x00, 0x00,
	0x00, 0xff, 0xff, 0x33, 0x69, 0xb8, 0x6b, 0xbb, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// LauncherServiceClient is the client API for LauncherService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type LauncherServiceClient interface {
	Launch(ctx context.Context, in *LaunchConfig, opts ...grpc.CallOption) (*ServiceResp, error)
	Stop(ctx context.Context, in *EmptyReq, opts ...grpc.CallOption) (*ServiceResp, error)
	Status(ctx context.Context, in *EmptyReq, opts ...grpc.CallOption) (*StatusResp, error)
}

type launcherServiceClient struct {
	cc *grpc.ClientConn
}

func NewLauncherServiceClient(cc *grpc.ClientConn) LauncherServiceClient {
	return &launcherServiceClient{cc}
}

func (c *launcherServiceClient) Launch(ctx context.Context, in *LaunchConfig, opts ...grpc.CallOption) (*ServiceResp, error) {
	out := new(ServiceResp)
	err := c.cc.Invoke(ctx, "/services.LauncherService/Launch", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *launcherServiceClient) Stop(ctx context.Context, in *EmptyReq, opts ...grpc.CallOption) (*ServiceResp, error) {
	out := new(ServiceResp)
	err := c.cc.Invoke(ctx, "/services.LauncherService/Stop", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *launcherServiceClient) Status(ctx context.Context, in *EmptyReq, opts ...grpc.CallOption) (*StatusResp, error) {
	out := new(StatusResp)
	err := c.cc.Invoke(ctx, "/services.LauncherService/Status", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LauncherServiceServer is the server API for LauncherService service.
type LauncherServiceServer interface {
	Launch(context.Context, *LaunchConfig) (*ServiceResp, error)
	Stop(context.Context, *EmptyReq) (*ServiceResp, error)
	Status(context.Context, *EmptyReq) (*StatusResp, error)
}

// UnimplementedLauncherServiceServer can be embedded to have forward compatible implementations.
type UnimplementedLauncherServiceServer struct {
}

func (*UnimplementedLauncherServiceServer) Launch(ctx context.Context, req *LaunchConfig) (*ServiceResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Launch not implemented")
}
func (*UnimplementedLauncherServiceServer) Stop(ctx context.Context, req *EmptyReq) (*ServiceResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Stop not implemented")
}
func (*UnimplementedLauncherServiceServer) Status(ctx context.Context, req *EmptyReq) (*StatusResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Status not implemented")
}

func RegisterLauncherServiceServer(s *grpc.Server, srv LauncherServiceServer) {
	s.RegisterService(&_LauncherService_serviceDesc, srv)
}

func _LauncherService_Launch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LaunchConfig)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LauncherServiceServer).Launch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/services.LauncherService/Launch",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LauncherServiceServer).Launch(ctx, req.(*LaunchConfig))
	}
	return interceptor(ctx, in, info, handler)
}

func _LauncherService_Stop_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LauncherServiceServer).Stop(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/services.LauncherService/Stop",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LauncherServiceServer).Stop(ctx, req.(*EmptyReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _LauncherService_Status_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LauncherServiceServer).Status(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/services.LauncherService/Status",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LauncherServiceServer).Status(ctx, req.(*EmptyReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _LauncherService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "services.LauncherService",
	HandlerType: (*LauncherServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Launch",
			Handler:    _LauncherService_Launch_Handler,
		},
		{
			MethodName: "Stop",
			Handler:    _LauncherService_Stop_Handler,
		},
		{
			MethodName: "Status",
			Handler:    _LauncherService_Status_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "cmd.proto",
}