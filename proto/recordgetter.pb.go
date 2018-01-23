// Code generated by protoc-gen-go. DO NOT EDIT.
// source: recordgetter.proto

/*
Package recordgetter is a generated protocol buffer package.

It is generated from these files:
	recordgetter.proto

It has these top-level messages:
	Empty
	State
	GetRecordRequest
	GetRecordResponse
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

type GetRecordRequest struct {
	Refresh bool `protobuf:"varint,1,opt,name=refresh" json:"refresh,omitempty"`
}

func (m *GetRecordRequest) Reset()                    { *m = GetRecordRequest{} }
func (m *GetRecordRequest) String() string            { return proto.CompactTextString(m) }
func (*GetRecordRequest) ProtoMessage()               {}
func (*GetRecordRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *GetRecordRequest) GetRefresh() bool {
	if m != nil {
		return m.Refresh
	}
	return false
}

type GetRecordResponse struct {
	Record     *recordcollection.Record `protobuf:"bytes,1,opt,name=record" json:"record,omitempty"`
	NumListens int32                    `protobuf:"varint,2,opt,name=num_listens,json=numListens" json:"num_listens,omitempty"`
}

func (m *GetRecordResponse) Reset()                    { *m = GetRecordResponse{} }
func (m *GetRecordResponse) String() string            { return proto.CompactTextString(m) }
func (*GetRecordResponse) ProtoMessage()               {}
func (*GetRecordResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *GetRecordResponse) GetRecord() *recordcollection.Record {
	if m != nil {
		return m.Record
	}
	return nil
}

func (m *GetRecordResponse) GetNumListens() int32 {
	if m != nil {
		return m.NumListens
	}
	return 0
}

func init() {
	proto.RegisterType((*Empty)(nil), "recordgetter.Empty")
	proto.RegisterType((*State)(nil), "recordgetter.State")
	proto.RegisterType((*GetRecordRequest)(nil), "recordgetter.GetRecordRequest")
	proto.RegisterType((*GetRecordResponse)(nil), "recordgetter.GetRecordResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for RecordGetter service

type RecordGetterClient interface {
	GetRecord(ctx context.Context, in *GetRecordRequest, opts ...grpc.CallOption) (*GetRecordResponse, error)
	Listened(ctx context.Context, in *recordcollection.Record, opts ...grpc.CallOption) (*Empty, error)
	Force(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
}

type recordGetterClient struct {
	cc *grpc.ClientConn
}

func NewRecordGetterClient(cc *grpc.ClientConn) RecordGetterClient {
	return &recordGetterClient{cc}
}

func (c *recordGetterClient) GetRecord(ctx context.Context, in *GetRecordRequest, opts ...grpc.CallOption) (*GetRecordResponse, error) {
	out := new(GetRecordResponse)
	err := grpc.Invoke(ctx, "/recordgetter.RecordGetter/GetRecord", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *recordGetterClient) Listened(ctx context.Context, in *recordcollection.Record, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
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
	GetRecord(context.Context, *GetRecordRequest) (*GetRecordResponse, error)
	Listened(context.Context, *recordcollection.Record) (*Empty, error)
	Force(context.Context, *Empty) (*Empty, error)
}

func RegisterRecordGetterServer(s *grpc.Server, srv RecordGetterServer) {
	s.RegisterService(&_RecordGetter_serviceDesc, srv)
}

func _RecordGetter_GetRecord_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRecordRequest)
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
		return srv.(RecordGetterServer).GetRecord(ctx, req.(*GetRecordRequest))
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
	// 283 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x51, 0xbd, 0x4e, 0xc3, 0x30,
	0x10, 0x56, 0x90, 0xd2, 0x96, 0x4b, 0x06, 0x30, 0x4b, 0x94, 0x81, 0x56, 0x99, 0x3a, 0xa0, 0x04,
	0xca, 0x84, 0x58, 0x29, 0x5d, 0x3a, 0xa0, 0xf0, 0x00, 0x15, 0x71, 0xaf, 0x89, 0xd5, 0xc4, 0x0e,
	0xf6, 0x65, 0xe0, 0x15, 0x79, 0x2a, 0x24, 0x3b, 0xa0, 0x96, 0xb6, 0xea, 0xe8, 0xef, 0xbe, 0x3f,
	0x7d, 0x06, 0xa6, 0x91, 0x2b, 0xbd, 0x2e, 0x91, 0x08, 0x75, 0xda, 0x6a, 0x45, 0x8a, 0x85, 0xbb,
	0x58, 0x3c, 0x2f, 0x05, 0x55, 0x5d, 0x91, 0x72, 0xd5, 0x64, 0x85, 0x56, 0x54, 0xa1, 0xae, 0x55,
	0x29, 0x78, 0xe6, 0x58, 0x5c, 0xd5, 0x35, 0x72, 0x12, 0x4a, 0x66, 0x56, 0x7d, 0x00, 0x3b, 0xd3,
	0x64, 0x08, 0xfe, 0xbc, 0x69, 0xe9, 0x2b, 0x79, 0x01, 0xff, 0x9d, 0x3e, 0x08, 0xd9, 0x33, 0x84,
	0xbc, 0xd3, 0x1a, 0x25, 0xad, 0x5a, 0xc1, 0xb7, 0x91, 0x37, 0xf1, 0xa6, 0xc1, 0x2c, 0x4a, 0x0f,
	0x0c, 0x72, 0x0b, 0xe4, 0x41, 0xcf, 0x7e, 0x13, 0x7c, 0x9b, 0xdc, 0xc1, 0xd5, 0x02, 0xa9, 0xbf,
	0xe0, 0x67, 0x87, 0x86, 0x58, 0x04, 0x43, 0x8d, 0x1b, 0x8d, 0xa6, 0xb2, 0x5e, 0xa3, 0xfc, 0xf7,
	0x99, 0x6c, 0xe0, 0x7a, 0x87, 0x6d, 0x5a, 0x25, 0x0d, 0xb2, 0x7b, 0x18, 0xb8, 0xa8, 0xb3, 0xc9,
	0x3d, 0x8f, 0x8d, 0x21, 0x90, 0x5d, 0xb3, 0xaa, 0x85, 0x21, 0x94, 0x26, 0xba, 0x98, 0x78, 0x53,
	0x3f, 0x07, 0xd9, 0x35, 0x4b, 0x87, 0xcc, 0xbe, 0x3d, 0x08, 0x9d, 0x66, 0x61, 0xc7, 0x63, 0x4b,
	0xb8, 0xfc, 0x0b, 0x66, 0xb7, 0xe9, 0xde, 0xd8, 0xff, 0xfb, 0xc7, 0xe3, 0x93, 0xf7, 0xbe, 0xf1,
	0x13, 0x8c, 0x5c, 0x12, 0xae, 0xd9, 0xc9, 0xb6, 0xf1, 0xcd, 0xbe, 0x8d, 0x5d, 0x9d, 0x3d, 0x80,
	0xff, 0xaa, 0x34, 0x47, 0x76, 0xec, 0x7a, 0x54, 0x52, 0x0c, 0xec, 0xc7, 0x3d, 0xfe, 0x04, 0x00,
	0x00, 0xff, 0xff, 0x0f, 0xc6, 0xaf, 0x8f, 0x23, 0x02, 0x00, 0x00,
}
