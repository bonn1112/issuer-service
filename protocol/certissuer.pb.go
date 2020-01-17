// Code generated by protoc-gen-go. DO NOT EDIT.
// source: certissuer.proto

package protocol

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

type IssueBlockchainCertificateRequest struct {
	Issuer               string   `protobuf:"bytes,1,opt,name=Issuer,proto3" json:"Issuer,omitempty"`
	Filename             string   `protobuf:"bytes,2,opt,name=Filename,proto3" json:"Filename,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *IssueBlockchainCertificateRequest) Reset()         { *m = IssueBlockchainCertificateRequest{} }
func (m *IssueBlockchainCertificateRequest) String() string { return proto.CompactTextString(m) }
func (*IssueBlockchainCertificateRequest) ProtoMessage()    {}
func (*IssueBlockchainCertificateRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_abe078d4eb06409d, []int{0}
}

func (m *IssueBlockchainCertificateRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_IssueBlockchainCertificateRequest.Unmarshal(m, b)
}
func (m *IssueBlockchainCertificateRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_IssueBlockchainCertificateRequest.Marshal(b, m, deterministic)
}
func (m *IssueBlockchainCertificateRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IssueBlockchainCertificateRequest.Merge(m, src)
}
func (m *IssueBlockchainCertificateRequest) XXX_Size() int {
	return xxx_messageInfo_IssueBlockchainCertificateRequest.Size(m)
}
func (m *IssueBlockchainCertificateRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_IssueBlockchainCertificateRequest.DiscardUnknown(m)
}

var xxx_messageInfo_IssueBlockchainCertificateRequest proto.InternalMessageInfo

func (m *IssueBlockchainCertificateRequest) GetIssuer() string {
	if m != nil {
		return m.Issuer
	}
	return ""
}

func (m *IssueBlockchainCertificateRequest) GetFilename() string {
	if m != nil {
		return m.Filename
	}
	return ""
}

type IssueBlockchainCertificateReply struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *IssueBlockchainCertificateReply) Reset()         { *m = IssueBlockchainCertificateReply{} }
func (m *IssueBlockchainCertificateReply) String() string { return proto.CompactTextString(m) }
func (*IssueBlockchainCertificateReply) ProtoMessage()    {}
func (*IssueBlockchainCertificateReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_abe078d4eb06409d, []int{1}
}

func (m *IssueBlockchainCertificateReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_IssueBlockchainCertificateReply.Unmarshal(m, b)
}
func (m *IssueBlockchainCertificateReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_IssueBlockchainCertificateReply.Marshal(b, m, deterministic)
}
func (m *IssueBlockchainCertificateReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IssueBlockchainCertificateReply.Merge(m, src)
}
func (m *IssueBlockchainCertificateReply) XXX_Size() int {
	return xxx_messageInfo_IssueBlockchainCertificateReply.Size(m)
}
func (m *IssueBlockchainCertificateReply) XXX_DiscardUnknown() {
	xxx_messageInfo_IssueBlockchainCertificateReply.DiscardUnknown(m)
}

var xxx_messageInfo_IssueBlockchainCertificateReply proto.InternalMessageInfo

func init() {
	proto.RegisterType((*IssueBlockchainCertificateRequest)(nil), "protocol.IssueBlockchainCertificateRequest")
	proto.RegisterType((*IssueBlockchainCertificateReply)(nil), "protocol.IssueBlockchainCertificateReply")
}

func init() { proto.RegisterFile("certissuer.proto", fileDescriptor_abe078d4eb06409d) }

var fileDescriptor_abe078d4eb06409d = []byte{
	// 162 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x48, 0x4e, 0x2d, 0x2a,
	0xc9, 0x2c, 0x2e, 0x2e, 0x4d, 0x2d, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x00, 0x53,
	0xc9, 0xf9, 0x39, 0x4a, 0xe1, 0x5c, 0x8a, 0x9e, 0x20, 0x19, 0xa7, 0x9c, 0xfc, 0xe4, 0xec, 0xe4,
	0x8c, 0xc4, 0xcc, 0x3c, 0x67, 0x90, 0xe2, 0xb4, 0xcc, 0xe4, 0xc4, 0x92, 0xd4, 0xa0, 0xd4, 0xc2,
	0xd2, 0xd4, 0xe2, 0x12, 0x21, 0x31, 0x2e, 0x36, 0xb0, 0xa2, 0x22, 0x09, 0x46, 0x05, 0x46, 0x0d,
	0xce, 0x20, 0x28, 0x4f, 0x48, 0x8a, 0x8b, 0xc3, 0x2d, 0x33, 0x27, 0x35, 0x2f, 0x31, 0x37, 0x55,
	0x82, 0x09, 0x2c, 0x03, 0xe7, 0x2b, 0x29, 0x72, 0xc9, 0xe3, 0x33, 0xb8, 0x20, 0xa7, 0xd2, 0xa8,
	0x85, 0x91, 0x8b, 0x0b, 0x24, 0x08, 0x35, 0xad, 0x8c, 0x4b, 0x0a, 0xb7, 0x0e, 0x21, 0x6d, 0x3d,
	0x98, 0x9b, 0xf5, 0x08, 0x3a, 0x58, 0x4a, 0x93, 0x38, 0xc5, 0x05, 0x39, 0x95, 0x4a, 0x0c, 0x49,
	0x6c, 0x60, 0xb5, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0xb8, 0xd9, 0xcc, 0xef, 0x27, 0x01,
	0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// CertIssuerClient is the client API for CertIssuer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type CertIssuerClient interface {
	IssueBlockchainCertificate(ctx context.Context, in *IssueBlockchainCertificateRequest, opts ...grpc.CallOption) (*IssueBlockchainCertificateReply, error)
}

type certIssuerClient struct {
	cc *grpc.ClientConn
}

func NewCertIssuerClient(cc *grpc.ClientConn) CertIssuerClient {
	return &certIssuerClient{cc}
}

func (c *certIssuerClient) IssueBlockchainCertificate(ctx context.Context, in *IssueBlockchainCertificateRequest, opts ...grpc.CallOption) (*IssueBlockchainCertificateReply, error) {
	out := new(IssueBlockchainCertificateReply)
	err := c.cc.Invoke(ctx, "/protocol.CertIssuer/IssueBlockchainCertificate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CertIssuerServer is the server API for CertIssuer service.
type CertIssuerServer interface {
	IssueBlockchainCertificate(context.Context, *IssueBlockchainCertificateRequest) (*IssueBlockchainCertificateReply, error)
}

// UnimplementedCertIssuerServer can be embedded to have forward compatible implementations.
type UnimplementedCertIssuerServer struct {
}

func (*UnimplementedCertIssuerServer) IssueBlockchainCertificate(ctx context.Context, req *IssueBlockchainCertificateRequest) (*IssueBlockchainCertificateReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IssueBlockchainCertificate not implemented")
}

func RegisterCertIssuerServer(s *grpc.Server, srv CertIssuerServer) {
	s.RegisterService(&_CertIssuer_serviceDesc, srv)
}

func _CertIssuer_IssueBlockchainCertificate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IssueBlockchainCertificateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CertIssuerServer).IssueBlockchainCertificate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protocol.CertIssuer/IssueBlockchainCertificate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CertIssuerServer).IssueBlockchainCertificate(ctx, req.(*IssueBlockchainCertificateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _CertIssuer_serviceDesc = grpc.ServiceDesc{
	ServiceName: "protocol.CertIssuer",
	HandlerType: (*CertIssuerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "IssueBlockchainCertificate",
			Handler:    _CertIssuer_IssueBlockchainCertificate_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "certissuer.proto",
}
