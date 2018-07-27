// Code generated by protoc-gen-go. DO NOT EDIT.
// source: service.proto

package user

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// GetPlayerByAccountReq 根据账号获取玩家请求
type GetPlayerByAccountReq struct {
	AccountId uint64 `protobuf:"varint,1,opt,name=account_id,json=accountId" json:"account_id,omitempty"`
}

func (m *GetPlayerByAccountReq) Reset()                    { *m = GetPlayerByAccountReq{} }
func (m *GetPlayerByAccountReq) String() string            { return proto.CompactTextString(m) }
func (*GetPlayerByAccountReq) ProtoMessage()               {}
func (*GetPlayerByAccountReq) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{0} }

func (m *GetPlayerByAccountReq) GetAccountId() uint64 {
	if m != nil {
		return m.AccountId
	}
	return 0
}

// GetPlayerByAccountRsp 根据账号获取玩家应答
type GetPlayerByAccountRsp struct {
	ErrCode  int32  `protobuf:"varint,1,opt,name=err_code,json=errCode" json:"err_code,omitempty"`
	PlayerId uint64 `protobuf:"varint,2,opt,name=player_id,json=playerId" json:"player_id,omitempty"`
}

func (m *GetPlayerByAccountRsp) Reset()                    { *m = GetPlayerByAccountRsp{} }
func (m *GetPlayerByAccountRsp) String() string            { return proto.CompactTextString(m) }
func (*GetPlayerByAccountRsp) ProtoMessage()               {}
func (*GetPlayerByAccountRsp) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{1} }

func (m *GetPlayerByAccountRsp) GetErrCode() int32 {
	if m != nil {
		return m.ErrCode
	}
	return 0
}

func (m *GetPlayerByAccountRsp) GetPlayerId() uint64 {
	if m != nil {
		return m.PlayerId
	}
	return 0
}

func init() {
	proto.RegisterType((*GetPlayerByAccountReq)(nil), "user.GetPlayerByAccountReq")
	proto.RegisterType((*GetPlayerByAccountRsp)(nil), "user.GetPlayerByAccountRsp")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for PlayerData service

type PlayerDataClient interface {
	// GetPlayerByAccount 根据账号获取玩家
	GetPlayerByAccount(ctx context.Context, in *GetPlayerByAccountReq, opts ...grpc.CallOption) (*GetPlayerByAccountRsp, error)
}

type playerDataClient struct {
	cc *grpc.ClientConn
}

func NewPlayerDataClient(cc *grpc.ClientConn) PlayerDataClient {
	return &playerDataClient{cc}
}

func (c *playerDataClient) GetPlayerByAccount(ctx context.Context, in *GetPlayerByAccountReq, opts ...grpc.CallOption) (*GetPlayerByAccountRsp, error) {
	out := new(GetPlayerByAccountRsp)
	err := grpc.Invoke(ctx, "/user.PlayerData/GetPlayerByAccount", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for PlayerData service

type PlayerDataServer interface {
	// GetPlayerByAccount 根据账号获取玩家
	GetPlayerByAccount(context.Context, *GetPlayerByAccountReq) (*GetPlayerByAccountRsp, error)
}

func RegisterPlayerDataServer(s *grpc.Server, srv PlayerDataServer) {
	s.RegisterService(&_PlayerData_serviceDesc, srv)
}

func _PlayerData_GetPlayerByAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPlayerByAccountReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlayerDataServer).GetPlayerByAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.PlayerData/GetPlayerByAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlayerDataServer).GetPlayerByAccount(ctx, req.(*GetPlayerByAccountReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _PlayerData_serviceDesc = grpc.ServiceDesc{
	ServiceName: "user.PlayerData",
	HandlerType: (*PlayerDataServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetPlayerByAccount",
			Handler:    _PlayerData_GetPlayerByAccount_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service.proto",
}

func init() { proto.RegisterFile("service.proto", fileDescriptor2) }

var fileDescriptor2 = []byte{
	// 175 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2d, 0x4e, 0x2d, 0x2a,
	0xcb, 0x4c, 0x4e, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x29, 0x2d, 0x4e, 0x2d, 0x52,
	0x32, 0xe3, 0x12, 0x75, 0x4f, 0x2d, 0x09, 0xc8, 0x49, 0xac, 0x4c, 0x2d, 0x72, 0xaa, 0x74, 0x4c,
	0x4e, 0xce, 0x2f, 0xcd, 0x2b, 0x09, 0x4a, 0x2d, 0x14, 0x92, 0xe5, 0xe2, 0x4a, 0x84, 0xf0, 0xe2,
	0x33, 0x53, 0x24, 0x18, 0x15, 0x18, 0x35, 0x58, 0x82, 0x38, 0xa1, 0x22, 0x9e, 0x29, 0x4a, 0xfe,
	0x58, 0xf5, 0x15, 0x17, 0x08, 0x49, 0x72, 0x71, 0xa4, 0x16, 0x15, 0xc5, 0x27, 0xe7, 0xa7, 0xa4,
	0x82, 0x75, 0xb1, 0x06, 0xb1, 0xa7, 0x16, 0x15, 0x39, 0xe7, 0xa7, 0xa4, 0x0a, 0x49, 0x73, 0x71,
	0x16, 0x80, 0x35, 0x80, 0x4c, 0x64, 0x02, 0x9b, 0xc8, 0x01, 0x11, 0xf0, 0x4c, 0x31, 0x8a, 0xe1,
	0xe2, 0x82, 0x98, 0xe6, 0x92, 0x58, 0x92, 0x28, 0xe4, 0xc7, 0x25, 0x84, 0x69, 0xbc, 0x90, 0xb4,
	0x1e, 0xc8, 0xcd, 0x7a, 0x58, 0x1d, 0x2c, 0x85, 0x5b, 0xb2, 0xb8, 0x20, 0x89, 0x0d, 0xec, 0x67,
	0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0xa3, 0xf3, 0x8b, 0xa9, 0x04, 0x01, 0x00, 0x00,
}
