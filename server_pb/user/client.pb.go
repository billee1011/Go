// Code generated by protoc-gen-go. DO NOT EDIT.
// source: client.proto

/*
Package user is a generated protocol buffer package.

It is generated from these files:
	client.proto
	errors.proto
	service.proto

It has these top-level messages:
	ClientDisconnect
	PlayerLogin
	GameConfig
	GameLevelConfig
	GetPlayerByAccountReq
	GetPlayerByAccountRsp
	GetPlayerInfoReq
	GetPlayerInfoRsp
	UpdatePlayerInfoReq
	UpdatePlayerInfoRsp
	GetPlayerStateReq
	GetPlayerStateRsp
	GetPlayerGameInfoReq
	GetPlayerGameInfoRsp
	UpdatePlayerStateReq
	UpdatePlayerStateRsp
	GetGameListInfoReq
	GetGameListInfoRsp
*/
package user

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// ClientDisconnect 客户端断开连接
type ClientDisconnect struct {
	ClientId uint64 `protobuf:"varint,1,opt,name=client_id,json=clientId" json:"client_id,omitempty"`
}

func (m *ClientDisconnect) Reset()                    { *m = ClientDisconnect{} }
func (m *ClientDisconnect) String() string            { return proto.CompactTextString(m) }
func (*ClientDisconnect) ProtoMessage()               {}
func (*ClientDisconnect) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *ClientDisconnect) GetClientId() uint64 {
	if m != nil {
		return m.ClientId
	}
	return 0
}

// PlayerLogin 玩家登录消息
type PlayerLogin struct {
	PlayerId uint64 `protobuf:"varint,1,opt,name=player_id,json=playerId" json:"player_id,omitempty"`
}

func (m *PlayerLogin) Reset()                    { *m = PlayerLogin{} }
func (m *PlayerLogin) String() string            { return proto.CompactTextString(m) }
func (*PlayerLogin) ProtoMessage()               {}
func (*PlayerLogin) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *PlayerLogin) GetPlayerId() uint64 {
	if m != nil {
		return m.PlayerId
	}
	return 0
}

func init() {
	proto.RegisterType((*ClientDisconnect)(nil), "user.ClientDisconnect")
	proto.RegisterType((*PlayerLogin)(nil), "user.PlayerLogin")
}

func init() { proto.RegisterFile("client.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 113 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x49, 0xce, 0xc9, 0x4c,
	0xcd, 0x2b, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x29, 0x2d, 0x4e, 0x2d, 0x52, 0xd2,
	0xe7, 0x12, 0x70, 0x06, 0x8b, 0xba, 0x64, 0x16, 0x27, 0xe7, 0xe7, 0xe5, 0xa5, 0x26, 0x97, 0x08,
	0x49, 0x73, 0x71, 0x42, 0x54, 0xc6, 0x67, 0xa6, 0x48, 0x30, 0x2a, 0x30, 0x6a, 0xb0, 0x04, 0x71,
	0x40, 0x04, 0x3c, 0x53, 0x94, 0xb4, 0xb8, 0xb8, 0x03, 0x72, 0x12, 0x2b, 0x53, 0x8b, 0x7c, 0xf2,
	0xd3, 0x33, 0xf3, 0x40, 0x6a, 0x0b, 0xc0, 0x5c, 0x24, 0xb5, 0x10, 0x01, 0xcf, 0x94, 0x24, 0x36,
	0xb0, 0x4d, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0x33, 0x87, 0x5b, 0xd1, 0x79, 0x00, 0x00,
	0x00,
}
