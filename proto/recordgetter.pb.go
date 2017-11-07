// Code generated by protoc-gen-go. DO NOT EDIT.
// source: recordgetter.proto

/*
Package recordgetter is a generated protocol buffer package.

It is generated from these files:
	recordgetter.proto

It has these top-level messages:
	Empty
	State
*/
package recordgetter

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import godiscogs "github.com/brotherlogic/godiscogs"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Empty struct {
}

func (m *Empty) Reset()                    { *m = Empty{} }
func (m *Empty) String() string            { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()               {}
func (*Empty) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type State struct {
	CurrentPick *godiscogs.Release `protobuf:"bytes,1,opt,name=current_pick,json=currentPick" json:"current_pick,omitempty"`
}

func (m *State) Reset()                    { *m = State{} }
func (m *State) String() string            { return proto.CompactTextString(m) }
func (*State) ProtoMessage()               {}
func (*State) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *State) GetCurrentPick() *godiscogs.Release {
	if m != nil {
		return m.CurrentPick
	}
	return nil
}

func init() {
	proto.RegisterType((*Empty)(nil), "recordgetter.Empty")
	proto.RegisterType((*State)(nil), "recordgetter.State")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for RecordGetter service

type RecordGetterClient interface {
	GetRecord(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*godiscogs.Release, error)
	Listened(ctx context.Context, in *godiscogs.Release, opts ...grpc.CallOption) (*godiscogs.Release, error)
}

type recordGetterClient struct {
	cc *grpc.ClientConn
}

func NewRecordGetterClient(cc *grpc.ClientConn) RecordGetterClient {
	return &recordGetterClient{cc}
}

func (c *recordGetterClient) GetRecord(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*godiscogs.Release, error) {
	out := new(godiscogs.Release)
	err := grpc.Invoke(ctx, "/recordgetter.RecordGetter/GetRecord", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *recordGetterClient) Listened(ctx context.Context, in *godiscogs.Release, opts ...grpc.CallOption) (*godiscogs.Release, error) {
	out := new(godiscogs.Release)
	err := grpc.Invoke(ctx, "/recordgetter.RecordGetter/Listened", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for RecordGetter service

type RecordGetterServer interface {
	GetRecord(context.Context, *Empty) (*godiscogs.Release, error)
	Listened(context.Context, *godiscogs.Release) (*godiscogs.Release, error)
}

func RegisterRecordGetterServer(s *grpc.Server, srv RecordGetterServer) {
	s.RegisterService(&_RecordGetter_serviceDesc, srv)
}

func _RecordGetter_GetRecord_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RecordGetterServer).GetRecord(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/recordgetter.RecordGetter/GetRecord",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RecordGetterServer).GetRecord(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _RecordGetter_Listened_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(godiscogs.Release)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RecordGetterServer).Listened(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/recordgetter.RecordGetter/Listened",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RecordGetterServer).Listened(ctx, req.(*godiscogs.Release))
	}
	return interceptor(ctx, in, info, handler)
}

var _RecordGetter_serviceDesc = grpc.ServiceDesc{
	ServiceName: "recordgetter.RecordGetter",
	HandlerType: (*RecordGetterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetRecord",
			Handler:    _RecordGetter_GetRecord_Handler,
		},
		{
			MethodName: "Listened",
			Handler:    _RecordGetter_Listened_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "recordgetter.proto",
}

func init() { proto.RegisterFile("recordgetter.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 198 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x8f, 0xbd, 0x6a, 0x85, 0x40,
	0x10, 0x85, 0xb1, 0x30, 0x3f, 0xab, 0xd5, 0xa4, 0x09, 0x56, 0xc1, 0x2a, 0xd5, 0x4a, 0x4c, 0xd2,
	0xa6, 0x0b, 0x36, 0x29, 0x82, 0x79, 0x80, 0xa0, 0xeb, 0xb0, 0x2e, 0xfe, 0x8c, 0xcc, 0x8e, 0x90,
	0xbc, 0xfd, 0x85, 0xf5, 0x82, 0xf7, 0x82, 0xdd, 0x70, 0x66, 0xce, 0x99, 0xef, 0x28, 0x60, 0x34,
	0xc4, 0x9d, 0x45, 0x11, 0x64, 0xbd, 0x30, 0x09, 0x41, 0x7a, 0xa9, 0x65, 0x2f, 0xd6, 0x49, 0xbf,
	0xb6, 0xda, 0xd0, 0x54, 0xb4, 0x4c, 0xd2, 0x23, 0x8f, 0x64, 0x9d, 0x29, 0x2c, 0x75, 0xce, 0x1b,
	0xb2, 0x7e, 0x9f, 0xb6, 0x80, 0xfc, 0x56, 0xc5, 0x9f, 0xd3, 0x22, 0xff, 0xf9, 0x87, 0x8a, 0x7f,
	0xa4, 0x11, 0x84, 0x77, 0x95, 0x9a, 0x95, 0x19, 0x67, 0xf9, 0x5d, 0x9c, 0x19, 0x1e, 0xa3, 0xa7,
	0xe8, 0x39, 0x29, 0x41, 0xef, 0xce, 0x1a, 0x47, 0x6c, 0x3c, 0xd6, 0xc9, 0xf9, 0xee, 0xdb, 0x99,
	0xa1, 0xfc, 0x53, 0x69, 0x1d, 0x58, 0xaa, 0xc0, 0x02, 0x6f, 0xea, 0xbe, 0x42, 0xd9, 0x24, 0x78,
	0xd0, 0x57, 0xec, 0xe1, 0x63, 0x76, 0x10, 0x09, 0xa5, 0xba, 0xfb, 0x72, 0x5e, 0x70, 0xc6, 0x0e,
	0x0e, 0xf6, 0x47, 0x9e, 0xf6, 0x26, 0x34, 0x79, 0x3d, 0x05, 0x00, 0x00, 0xff, 0xff, 0x72, 0xb3,
	0xb5, 0x2b, 0x20, 0x01, 0x00, 0x00,
}