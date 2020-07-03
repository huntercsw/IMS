// Code generated by protoc-gen-go. DO NOT EDIT.
// source: imsLb.proto

package imsPb

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

type GetServiceRequest struct {
	ServiceGroup         string   `protobuf:"bytes,1,opt,name=serviceGroup,proto3" json:"serviceGroup,omitempty"`
	ServiceName          string   `protobuf:"bytes,2,opt,name=serviceName,proto3" json:"serviceName,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetServiceRequest) Reset()         { *m = GetServiceRequest{} }
func (m *GetServiceRequest) String() string { return proto.CompactTextString(m) }
func (*GetServiceRequest) ProtoMessage()    {}
func (*GetServiceRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_533f9b36b1cce5b0, []int{0}
}

func (m *GetServiceRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetServiceRequest.Unmarshal(m, b)
}
func (m *GetServiceRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetServiceRequest.Marshal(b, m, deterministic)
}
func (m *GetServiceRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetServiceRequest.Merge(m, src)
}
func (m *GetServiceRequest) XXX_Size() int {
	return xxx_messageInfo_GetServiceRequest.Size(m)
}
func (m *GetServiceRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetServiceRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetServiceRequest proto.InternalMessageInfo

func (m *GetServiceRequest) GetServiceGroup() string {
	if m != nil {
		return m.ServiceGroup
	}
	return ""
}

func (m *GetServiceRequest) GetServiceName() string {
	if m != nil {
		return m.ServiceName
	}
	return ""
}

type GetServiceResponse struct {
	ServiceHost          string   `protobuf:"bytes,1,opt,name=serviceHost,proto3" json:"serviceHost,omitempty"`
	ServicePort          string   `protobuf:"bytes,2,opt,name=servicePort,proto3" json:"servicePort,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetServiceResponse) Reset()         { *m = GetServiceResponse{} }
func (m *GetServiceResponse) String() string { return proto.CompactTextString(m) }
func (*GetServiceResponse) ProtoMessage()    {}
func (*GetServiceResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_533f9b36b1cce5b0, []int{1}
}

func (m *GetServiceResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetServiceResponse.Unmarshal(m, b)
}
func (m *GetServiceResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetServiceResponse.Marshal(b, m, deterministic)
}
func (m *GetServiceResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetServiceResponse.Merge(m, src)
}
func (m *GetServiceResponse) XXX_Size() int {
	return xxx_messageInfo_GetServiceResponse.Size(m)
}
func (m *GetServiceResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetServiceResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetServiceResponse proto.InternalMessageInfo

func (m *GetServiceResponse) GetServiceHost() string {
	if m != nil {
		return m.ServiceHost
	}
	return ""
}

func (m *GetServiceResponse) GetServicePort() string {
	if m != nil {
		return m.ServicePort
	}
	return ""
}

func init() {
	proto.RegisterType((*GetServiceRequest)(nil), "imsPb.GetServiceRequest")
	proto.RegisterType((*GetServiceResponse)(nil), "imsPb.GetServiceResponse")
}

func init() { proto.RegisterFile("imsLb.proto", fileDescriptor_533f9b36b1cce5b0) }

var fileDescriptor_533f9b36b1cce5b0 = []byte{
	// 187 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xce, 0xcc, 0x2d, 0xf6,
	0x49, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0xcd, 0xcc, 0x2d, 0x0e, 0x48, 0x52, 0x8a,
	0xe4, 0x12, 0x74, 0x4f, 0x2d, 0x09, 0x4e, 0x2d, 0x2a, 0xcb, 0x4c, 0x4e, 0x0d, 0x4a, 0x2d, 0x2c,
	0x4d, 0x2d, 0x2e, 0x11, 0x52, 0xe2, 0xe2, 0x29, 0x86, 0x88, 0xb8, 0x17, 0xe5, 0x97, 0x16, 0x48,
	0x30, 0x2a, 0x30, 0x6a, 0x70, 0x06, 0xa1, 0x88, 0x09, 0x29, 0x70, 0x71, 0x43, 0xf9, 0x7e, 0x89,
	0xb9, 0xa9, 0x12, 0x4c, 0x60, 0x25, 0xc8, 0x42, 0x4a, 0x11, 0x5c, 0x42, 0xc8, 0x46, 0x17, 0x17,
	0xe4, 0xe7, 0x15, 0xa7, 0x22, 0xe9, 0xf3, 0xc8, 0x2f, 0x2e, 0x81, 0x1a, 0x8d, 0x2c, 0x84, 0xa4,
	0x22, 0x20, 0xbf, 0xa8, 0x04, 0xcd, 0x64, 0x90, 0x90, 0x51, 0x34, 0x17, 0x9f, 0x67, 0x6e, 0xb1,
	0x4f, 0x7e, 0x62, 0x8a, 0x53, 0x62, 0x4e, 0x62, 0x5e, 0x72, 0xaa, 0x90, 0x27, 0x97, 0x00, 0xc2,
	0xae, 0xf0, 0xcc, 0x92, 0x8c, 0xa0, 0x20, 0x21, 0x09, 0x3d, 0xb0, 0x17, 0xf5, 0x30, 0xfc, 0x27,
	0x25, 0x89, 0x45, 0x06, 0xe2, 0x3c, 0x25, 0x86, 0x24, 0x36, 0x70, 0xf8, 0x18, 0x03, 0x02, 0x00,
	0x00, 0xff, 0xff, 0x77, 0x0a, 0xa6, 0x94, 0x2e, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ImsLoadBalanceClient is the client API for ImsLoadBalance service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ImsLoadBalanceClient interface {
	GetServiceWithRR(ctx context.Context, in *GetServiceRequest, opts ...grpc.CallOption) (*GetServiceResponse, error)
}

type imsLoadBalanceClient struct {
	cc *grpc.ClientConn
}

func NewImsLoadBalanceClient(cc *grpc.ClientConn) ImsLoadBalanceClient {
	return &imsLoadBalanceClient{cc}
}

func (c *imsLoadBalanceClient) GetServiceWithRR(ctx context.Context, in *GetServiceRequest, opts ...grpc.CallOption) (*GetServiceResponse, error) {
	out := new(GetServiceResponse)
	err := c.cc.Invoke(ctx, "/imsPb.ImsLoadBalance/GetServiceWithRR", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ImsLoadBalanceServer is the server API for ImsLoadBalance service.
type ImsLoadBalanceServer interface {
	GetServiceWithRR(context.Context, *GetServiceRequest) (*GetServiceResponse, error)
}

// UnimplementedImsLoadBalanceServer can be embedded to have forward compatible implementations.
type UnimplementedImsLoadBalanceServer struct {
}

func (*UnimplementedImsLoadBalanceServer) GetServiceWithRR(ctx context.Context, req *GetServiceRequest) (*GetServiceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetServiceWithRR not implemented")
}

func RegisterImsLoadBalanceServer(s *grpc.Server, srv ImsLoadBalanceServer) {
	s.RegisterService(&_ImsLoadBalance_serviceDesc, srv)
}

func _ImsLoadBalance_GetServiceWithRR_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetServiceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ImsLoadBalanceServer).GetServiceWithRR(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/imsPb.ImsLoadBalance/GetServiceWithRR",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ImsLoadBalanceServer).GetServiceWithRR(ctx, req.(*GetServiceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ImsLoadBalance_serviceDesc = grpc.ServiceDesc{
	ServiceName: "imsPb.ImsLoadBalance",
	HandlerType: (*ImsLoadBalanceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetServiceWithRR",
			Handler:    _ImsLoadBalance_GetServiceWithRR_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "imsLb.proto",
}