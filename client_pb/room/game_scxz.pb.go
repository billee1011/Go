// Code generated by protoc-gen-go. DO NOT EDIT.
// source: game_scxz.proto

package room // import "steve/client_pb/room"

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

// RoomGiveUpReq 认输请求
type RoomGiveUpReq struct {
	Reserve              *uint32  `protobuf:"varint,1,opt,name=reserve" json:"reserve,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RoomGiveUpReq) Reset()         { *m = RoomGiveUpReq{} }
func (m *RoomGiveUpReq) String() string { return proto.CompactTextString(m) }
func (*RoomGiveUpReq) ProtoMessage()    {}
func (*RoomGiveUpReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_game_scxz_d7a71545d2f997cc, []int{0}
}
func (m *RoomGiveUpReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RoomGiveUpReq.Unmarshal(m, b)
}
func (m *RoomGiveUpReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RoomGiveUpReq.Marshal(b, m, deterministic)
}
func (dst *RoomGiveUpReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RoomGiveUpReq.Merge(dst, src)
}
func (m *RoomGiveUpReq) XXX_Size() int {
	return xxx_messageInfo_RoomGiveUpReq.Size(m)
}
func (m *RoomGiveUpReq) XXX_DiscardUnknown() {
	xxx_messageInfo_RoomGiveUpReq.DiscardUnknown(m)
}

var xxx_messageInfo_RoomGiveUpReq proto.InternalMessageInfo

func (m *RoomGiveUpReq) GetReserve() uint32 {
	if m != nil && m.Reserve != nil {
		return *m.Reserve
	}
	return 0
}

// RoomGiveUpRsp 认输响应
type RoomGiveUpRsp struct {
	ErrCode              *RoomError `protobuf:"varint,1,opt,name=err_code,json=errCode,enum=room.RoomError" json:"err_code,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *RoomGiveUpRsp) Reset()         { *m = RoomGiveUpRsp{} }
func (m *RoomGiveUpRsp) String() string { return proto.CompactTextString(m) }
func (*RoomGiveUpRsp) ProtoMessage()    {}
func (*RoomGiveUpRsp) Descriptor() ([]byte, []int) {
	return fileDescriptor_game_scxz_d7a71545d2f997cc, []int{1}
}
func (m *RoomGiveUpRsp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RoomGiveUpRsp.Unmarshal(m, b)
}
func (m *RoomGiveUpRsp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RoomGiveUpRsp.Marshal(b, m, deterministic)
}
func (dst *RoomGiveUpRsp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RoomGiveUpRsp.Merge(dst, src)
}
func (m *RoomGiveUpRsp) XXX_Size() int {
	return xxx_messageInfo_RoomGiveUpRsp.Size(m)
}
func (m *RoomGiveUpRsp) XXX_DiscardUnknown() {
	xxx_messageInfo_RoomGiveUpRsp.DiscardUnknown(m)
}

var xxx_messageInfo_RoomGiveUpRsp proto.InternalMessageInfo

func (m *RoomGiveUpRsp) GetErrCode() RoomError {
	if m != nil && m.ErrCode != nil {
		return *m.ErrCode
	}
	return RoomError_SUCCESS
}

// RoomGiveUpNtf 认输通知
type RoomGiveUpNtf struct {
	PlayerId             []uint64 `protobuf:"varint,1,rep,name=player_id,json=playerId" json:"player_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RoomGiveUpNtf) Reset()         { *m = RoomGiveUpNtf{} }
func (m *RoomGiveUpNtf) String() string { return proto.CompactTextString(m) }
func (*RoomGiveUpNtf) ProtoMessage()    {}
func (*RoomGiveUpNtf) Descriptor() ([]byte, []int) {
	return fileDescriptor_game_scxz_d7a71545d2f997cc, []int{2}
}
func (m *RoomGiveUpNtf) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RoomGiveUpNtf.Unmarshal(m, b)
}
func (m *RoomGiveUpNtf) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RoomGiveUpNtf.Marshal(b, m, deterministic)
}
func (dst *RoomGiveUpNtf) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RoomGiveUpNtf.Merge(dst, src)
}
func (m *RoomGiveUpNtf) XXX_Size() int {
	return xxx_messageInfo_RoomGiveUpNtf.Size(m)
}
func (m *RoomGiveUpNtf) XXX_DiscardUnknown() {
	xxx_messageInfo_RoomGiveUpNtf.DiscardUnknown(m)
}

var xxx_messageInfo_RoomGiveUpNtf proto.InternalMessageInfo

func (m *RoomGiveUpNtf) GetPlayerId() []uint64 {
	if m != nil {
		return m.PlayerId
	}
	return nil
}

func init() {
	proto.RegisterType((*RoomGiveUpReq)(nil), "room.RoomGiveUpReq")
	proto.RegisterType((*RoomGiveUpRsp)(nil), "room.RoomGiveUpRsp")
	proto.RegisterType((*RoomGiveUpNtf)(nil), "room.RoomGiveUpNtf")
}

func init() { proto.RegisterFile("game_scxz.proto", fileDescriptor_game_scxz_d7a71545d2f997cc) }

var fileDescriptor_game_scxz_d7a71545d2f997cc = []byte{
	// 187 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4f, 0x4f, 0xcc, 0x4d,
	0x8d, 0x2f, 0x4e, 0xae, 0xa8, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x29, 0xca, 0xcf,
	0xcf, 0x95, 0xe2, 0x4e, 0x2d, 0x2a, 0xca, 0x2f, 0x82, 0x08, 0x29, 0x69, 0x72, 0xf1, 0x06, 0xe5,
	0xe7, 0xe7, 0xba, 0x67, 0x96, 0xa5, 0x86, 0x16, 0x04, 0xa5, 0x16, 0x0a, 0x49, 0x70, 0xb1, 0x17,
	0xa5, 0x16, 0xa7, 0x16, 0x95, 0xa5, 0x4a, 0x30, 0x2a, 0x30, 0x6a, 0xf0, 0x06, 0xc1, 0xb8, 0x4a,
	0xd6, 0x28, 0x4a, 0x8b, 0x0b, 0x84, 0xb4, 0xb8, 0x38, 0x52, 0x8b, 0x8a, 0xe2, 0x93, 0xf3, 0x53,
	0x20, 0x6a, 0xf9, 0x8c, 0xf8, 0xf5, 0x40, 0x36, 0xe8, 0x81, 0x94, 0xb9, 0x82, 0x2c, 0x09, 0x62,
	0x4f, 0x2d, 0x2a, 0x72, 0xce, 0x4f, 0x49, 0x55, 0xd2, 0x41, 0xd6, 0xec, 0x57, 0x92, 0x26, 0x24,
	0xcd, 0xc5, 0x59, 0x90, 0x93, 0x58, 0x99, 0x5a, 0x14, 0x9f, 0x99, 0x22, 0xc1, 0xa8, 0xc0, 0xac,
	0xc1, 0x12, 0xc4, 0x01, 0x11, 0xf0, 0x4c, 0x71, 0x12, 0x8b, 0x12, 0x29, 0x2e, 0x49, 0x2d, 0x4b,
	0xd5, 0x4f, 0xce, 0xc9, 0x4c, 0xcd, 0x2b, 0x89, 0x2f, 0x48, 0xd2, 0x07, 0x19, 0x0c, 0x08, 0x00,
	0x00, 0xff, 0xff, 0x95, 0xc8, 0x06, 0xd3, 0xd2, 0x00, 0x00, 0x00,
}
