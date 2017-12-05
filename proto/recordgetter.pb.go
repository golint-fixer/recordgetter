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
import recordcollection "github.com/brotherlogic/recordcollection/proto"

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
	CurrentPick *recordcollection.Record `protobuf:"bytes,1,opt,name=current_pick,json=currentPick" json:"current_pick,omitempty"`
}

func (m *State) Reset()                    { *m = State{} }
func (m *State) String() string            { return proto.CompactTextString(m) }
func (*State) ProtoMessage()               {}
func (*State) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *State) GetCurrentPick() *recordcollection.Record {
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
	GetRecord(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*recordcollection.Record, error)
	Listened(ctx context.Context, in *recordcollection.Record, opts ...grpc.CallOption) (*recordcollection.Record, error)
	Force(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
}

type recordGetterClient struct {
	cc *grpc.ClientConn
}

func NewRecordGetterClient(cc *grpc.ClientConn) RecordGetterClient {
	return &recordGetterClient{cc}
}

func (c *recordGetterClient) GetRecord(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*recordcollection.Record, error) {
	out := new(recordcollection.Record)
	err := grpc.Invoke(ctx, "/recordgetter.RecordGetter/GetRecord", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *recordGetterClient) Listened(ctx context.Context, in *recordcollection.Record, opts ...grpc.CallOption) (*recordcollection.Record, error) {
	out := new(recordcollection.Record)
	err := grpc.Invoke(ctx, "/recordgetter.RecordGetter/Listened", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *recordGetterClient) Force(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := grpc.Invoke(ctx, "/recordgetter.RecordGetter/Force", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for RecordGetter service

type RecordGetterServer interface {
	GetRecord(context.Context, *Empty) (*recordcollection.Record, error)
	Listened(context.Context, *recordcollection.Record) (*recordcollection.Record, error)
	Force(context.Context, *Empty) (*Empty, error)
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
	in := new(recordcollection.Record)
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
		return srv.(RecordGetterServer).Listened(ctx, req.(*recordcollection.Record))
	}
	return interceptor(ctx, in, info, handler)
}

func _RecordGetter_Force_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RecordGetterServer).Force(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/recordgetter.RecordGetter/Force",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RecordGetterServer).Force(ctx, req.(*Empty))
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
		{
			MethodName: "Force",
			Handler:    _RecordGetter_Force_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "recordgetter.proto",
}

func init() { proto.RegisterFile("recordgetter.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 213 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x90, 0xbd, 0x4e, 0x85, 0x40,
	0x10, 0x85, 0x43, 0x81, 0x3f, 0x73, 0xa9, 0xd6, 0xe6, 0xe6, 0x56, 0xe6, 0x56, 0x56, 0x4b, 0xc4,
	0x4e, 0x13, 0x2b, 0x91, 0xc6, 0xc2, 0xe0, 0x03, 0x18, 0x19, 0x26, 0xb0, 0xe1, 0x67, 0xc8, 0x38,
	0x14, 0x3e, 0x9c, 0xef, 0x66, 0xb2, 0x6b, 0xa1, 0x01, 0x6e, 0xfb, 0xe5, 0x9c, 0xef, 0xec, 0x0e,
	0x18, 0x21, 0x64, 0xa9, 0x1b, 0x52, 0x25, 0xb1, 0x93, 0xb0, 0xb2, 0x49, 0xfe, 0xb2, 0x43, 0xde,
	0x38, 0x6d, 0xe7, 0xca, 0x22, 0x0f, 0x69, 0x25, 0xac, 0x2d, 0x49, 0xcf, 0x8d, 0xc3, 0x34, 0xa4,
	0x90, 0xfb, 0x9e, 0x50, 0x1d, 0x8f, 0xa9, 0x6f, 0x2f, 0x70, 0x90, 0x1e, 0xcf, 0x21, 0xce, 0x87,
	0x49, 0xbf, 0x8e, 0x4f, 0x10, 0xbf, 0xe9, 0x87, 0x92, 0x79, 0x80, 0x04, 0x67, 0x11, 0x1a, 0xf5,
	0x7d, 0x72, 0xd8, 0xed, 0xa3, 0xeb, 0xe8, 0x66, 0x97, 0xed, 0xed, 0x42, 0x50, 0x7a, 0x50, 0xee,
	0x7e, 0xd3, 0xaf, 0x0e, 0xbb, 0xec, 0x3b, 0x82, 0x24, 0xf0, 0xc2, 0x3f, 0xd3, 0xdc, 0xc3, 0x65,
	0x41, 0x1a, 0x90, 0xb9, 0xb2, 0xff, 0xbe, 0xe5, 0x87, 0x0f, 0x9b, 0x66, 0xf3, 0x08, 0x17, 0x2f,
	0xee, 0x53, 0x69, 0xa4, 0xda, 0x6c, 0xa6, 0x4e, 0xf4, 0x6f, 0x21, 0x7e, 0x66, 0x41, 0x5a, 0xdf,
	0x5d, 0x83, 0xd5, 0x99, 0xbf, 0xca, 0xdd, 0x4f, 0x00, 0x00, 0x00, 0xff, 0xff, 0xc5, 0x4b, 0x37,
	0xe0, 0x80, 0x01, 0x00, 0x00,
}
