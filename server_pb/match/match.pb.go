// Code generated by protoc-gen-go. DO NOT EDIT.
// source: match.proto

/*
Package match is a generated protocol buffer package.

It is generated from these files:
	match.proto

It has these top-level messages:
	ContinuePlayer
	AddContinueDeskReq
	AddContinueDeskRsp
*/
package match

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

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// ContinuePlayer 牌桌续局玩家
type ContinuePlayer struct {
	PlayerId   uint64 `protobuf:"varint,1,opt,name=player_id,json=playerId" json:"player_id,omitempty"`
	Seat       int32  `protobuf:"varint,2,opt,name=seat" json:"seat,omitempty"`
	Win        bool   `protobuf:"varint,3,opt,name=win" json:"win,omitempty"`
	RobotLevel int32  `protobuf:"varint,4,opt,name=robot_level,json=robotLevel" json:"robot_level,omitempty"`
}

func (m *ContinuePlayer) Reset()                    { *m = ContinuePlayer{} }
func (m *ContinuePlayer) String() string            { return proto.CompactTextString(m) }
func (*ContinuePlayer) ProtoMessage()               {}
func (*ContinuePlayer) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *ContinuePlayer) GetPlayerId() uint64 {
	if m != nil {
		return m.PlayerId
	}
	return 0
}

func (m *ContinuePlayer) GetSeat() int32 {
	if m != nil {
		return m.Seat
	}
	return 0
}

func (m *ContinuePlayer) GetWin() bool {
	if m != nil {
		return m.Win
	}
	return false
}

func (m *ContinuePlayer) GetRobotLevel() int32 {
	if m != nil {
		return m.RobotLevel
	}
	return 0
}

// AddContinueDeskReq 添加续局牌桌请求
type AddContinueDeskReq struct {
	Players    []*ContinuePlayer `protobuf:"bytes,1,rep,name=players" json:"players,omitempty"`
	GameId     int32             `protobuf:"varint,2,opt,name=game_id,json=gameId" json:"game_id,omitempty"`
	FixBanker  bool              `protobuf:"varint,3,opt,name=fix_banker,json=fixBanker" json:"fix_banker,omitempty"`
	BankerSeat int32             `protobuf:"varint,4,opt,name=banker_seat,json=bankerSeat" json:"banker_seat,omitempty"`
}

func (m *AddContinueDeskReq) Reset()                    { *m = AddContinueDeskReq{} }
func (m *AddContinueDeskReq) String() string            { return proto.CompactTextString(m) }
func (*AddContinueDeskReq) ProtoMessage()               {}
func (*AddContinueDeskReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *AddContinueDeskReq) GetPlayers() []*ContinuePlayer {
	if m != nil {
		return m.Players
	}
	return nil
}

func (m *AddContinueDeskReq) GetGameId() int32 {
	if m != nil {
		return m.GameId
	}
	return 0
}

func (m *AddContinueDeskReq) GetFixBanker() bool {
	if m != nil {
		return m.FixBanker
	}
	return false
}

func (m *AddContinueDeskReq) GetBankerSeat() int32 {
	if m != nil {
		return m.BankerSeat
	}
	return 0
}

// AddContinueDeskRsp 添加续局牌桌应答
type AddContinueDeskRsp struct {
}

func (m *AddContinueDeskRsp) Reset()                    { *m = AddContinueDeskRsp{} }
func (m *AddContinueDeskRsp) String() string            { return proto.CompactTextString(m) }
func (*AddContinueDeskRsp) ProtoMessage()               {}
func (*AddContinueDeskRsp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func init() {
	proto.RegisterType((*ContinuePlayer)(nil), "match.ContinuePlayer")
	proto.RegisterType((*AddContinueDeskReq)(nil), "match.AddContinueDeskReq")
	proto.RegisterType((*AddContinueDeskRsp)(nil), "match.AddContinueDeskRsp")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Match service

type MatchClient interface {
	AddContinueDesk(ctx context.Context, in *AddContinueDeskReq, opts ...grpc.CallOption) (*AddContinueDeskRsp, error)
}

type matchClient struct {
	cc *grpc.ClientConn
}

func NewMatchClient(cc *grpc.ClientConn) MatchClient {
	return &matchClient{cc}
}

func (c *matchClient) AddContinueDesk(ctx context.Context, in *AddContinueDeskReq, opts ...grpc.CallOption) (*AddContinueDeskRsp, error) {
	out := new(AddContinueDeskRsp)
	err := grpc.Invoke(ctx, "/match.Match/AddContinueDesk", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Match service

type MatchServer interface {
	AddContinueDesk(context.Context, *AddContinueDeskReq) (*AddContinueDeskRsp, error)
}

func RegisterMatchServer(s *grpc.Server, srv MatchServer) {
	s.RegisterService(&_Match_serviceDesc, srv)
}

func _Match_AddContinueDesk_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddContinueDeskReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MatchServer).AddContinueDesk(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/match.Match/AddContinueDesk",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MatchServer).AddContinueDesk(ctx, req.(*AddContinueDeskReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _Match_serviceDesc = grpc.ServiceDesc{
	ServiceName: "match.Match",
	HandlerType: (*MatchServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddContinueDesk",
			Handler:    _Match_AddContinueDesk_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "match.proto",
}

func init() { proto.RegisterFile("match.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 266 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x91, 0xdb, 0x4a, 0xc4, 0x30,
	0x10, 0x86, 0x8d, 0x6d, 0xf7, 0x30, 0x05, 0x95, 0x41, 0xb1, 0x2a, 0x62, 0xe9, 0x55, 0xaf, 0x56,
	0x58, 0x9f, 0xc0, 0xc3, 0x4d, 0x41, 0x41, 0xe2, 0x03, 0x94, 0xd4, 0xce, 0x6a, 0xd8, 0x6e, 0x53,
	0x9b, 0xa8, 0xeb, 0xcb, 0xf8, 0xac, 0x92, 0xa4, 0x7b, 0xa1, 0x8b, 0x77, 0x7f, 0xbe, 0x7f, 0x60,
	0xbe, 0x21, 0x10, 0xaf, 0x84, 0x79, 0x7e, 0x9d, 0x75, 0xbd, 0x32, 0x0a, 0x23, 0xf7, 0xc8, 0x0c,
	0xec, 0xdd, 0xaa, 0xd6, 0xc8, 0xf6, 0x9d, 0x1e, 0x1b, 0xf1, 0x45, 0x3d, 0x9e, 0xc1, 0xb4, 0x73,
	0xa9, 0x94, 0x75, 0xc2, 0x52, 0x96, 0x87, 0x7c, 0xe2, 0x41, 0x51, 0x23, 0x42, 0xa8, 0x49, 0x98,
	0x64, 0x37, 0x65, 0x79, 0xc4, 0x5d, 0xc6, 0x03, 0x08, 0x3e, 0x65, 0x9b, 0x04, 0x29, 0xcb, 0x27,
	0xdc, 0x46, 0xbc, 0x80, 0xb8, 0x57, 0x95, 0x32, 0x65, 0x43, 0x1f, 0xd4, 0x24, 0xa1, 0x1b, 0x06,
	0x87, 0xee, 0x2d, 0xc9, 0xbe, 0x19, 0xe0, 0x75, 0x5d, 0x6f, 0x36, 0xdf, 0x91, 0x5e, 0x72, 0x7a,
	0xc3, 0x4b, 0x18, 0xfb, 0x4d, 0x3a, 0x61, 0x69, 0x90, 0xc7, 0xf3, 0xa3, 0x99, 0x57, 0xfe, 0xad,
	0xc8, 0x37, 0x53, 0x78, 0x0c, 0xe3, 0x17, 0xb1, 0x22, 0x6b, 0xea, 0x8d, 0x46, 0xf6, 0x59, 0xd4,
	0x78, 0x0e, 0xb0, 0x90, 0xeb, 0xb2, 0x12, 0xed, 0x92, 0xfa, 0x41, 0x6d, 0xba, 0x90, 0xeb, 0x1b,
	0x07, 0xac, 0xa0, 0xaf, 0x4a, 0x77, 0xcd, 0x20, 0xe8, 0xd1, 0x13, 0x09, 0x93, 0x1d, 0x6e, 0xfb,
	0xe9, 0x6e, 0xce, 0x21, 0x7a, 0xb0, 0x3e, 0x58, 0xc0, 0xfe, 0x9f, 0x1a, 0x4f, 0x06, 0xd5, 0xed,
	0xb3, 0x4e, 0xff, 0xab, 0x74, 0x97, 0xed, 0x54, 0x23, 0xf7, 0x1d, 0x57, 0x3f, 0x01, 0x00, 0x00,
	0xff, 0xff, 0xb2, 0x30, 0xa9, 0xd6, 0x9d, 0x01, 0x00, 0x00,
}
