// Code generated by protoc-gen-go. DO NOT EDIT.
// source: game_erren.proto

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

// 补花
type RoomBuHuaInfo struct {
	PlayerId             *uint64  `protobuf:"varint,1,opt,name=player_id,json=playerId" json:"player_id,omitempty"`
	OutHuaCards          []uint32 `protobuf:"varint,2,rep,name=out_hua_cards,json=outHuaCards" json:"out_hua_cards,omitempty"`
	BuCards              []uint32 `protobuf:"varint,3,rep,name=bu_cards,json=buCards" json:"bu_cards,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RoomBuHuaInfo) Reset()         { *m = RoomBuHuaInfo{} }
func (m *RoomBuHuaInfo) String() string { return proto.CompactTextString(m) }
func (*RoomBuHuaInfo) ProtoMessage()    {}
func (*RoomBuHuaInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_game_erren_aac2442421e12cfa, []int{0}
}
func (m *RoomBuHuaInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RoomBuHuaInfo.Unmarshal(m, b)
}
func (m *RoomBuHuaInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RoomBuHuaInfo.Marshal(b, m, deterministic)
}
func (dst *RoomBuHuaInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RoomBuHuaInfo.Merge(dst, src)
}
func (m *RoomBuHuaInfo) XXX_Size() int {
	return xxx_messageInfo_RoomBuHuaInfo.Size(m)
}
func (m *RoomBuHuaInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_RoomBuHuaInfo.DiscardUnknown(m)
}

var xxx_messageInfo_RoomBuHuaInfo proto.InternalMessageInfo

func (m *RoomBuHuaInfo) GetPlayerId() uint64 {
	if m != nil && m.PlayerId != nil {
		return *m.PlayerId
	}
	return 0
}

func (m *RoomBuHuaInfo) GetOutHuaCards() []uint32 {
	if m != nil {
		return m.OutHuaCards
	}
	return nil
}

func (m *RoomBuHuaInfo) GetBuCards() []uint32 {
	if m != nil {
		return m.BuCards
	}
	return nil
}

type RoomBuHuaNtf struct {
	BuhuaInfo            []*RoomBuHuaInfo `protobuf:"bytes,1,rep,name=buhua_info,json=buhuaInfo" json:"buhua_info,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *RoomBuHuaNtf) Reset()         { *m = RoomBuHuaNtf{} }
func (m *RoomBuHuaNtf) String() string { return proto.CompactTextString(m) }
func (*RoomBuHuaNtf) ProtoMessage()    {}
func (*RoomBuHuaNtf) Descriptor() ([]byte, []int) {
	return fileDescriptor_game_erren_aac2442421e12cfa, []int{1}
}
func (m *RoomBuHuaNtf) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RoomBuHuaNtf.Unmarshal(m, b)
}
func (m *RoomBuHuaNtf) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RoomBuHuaNtf.Marshal(b, m, deterministic)
}
func (dst *RoomBuHuaNtf) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RoomBuHuaNtf.Merge(dst, src)
}
func (m *RoomBuHuaNtf) XXX_Size() int {
	return xxx_messageInfo_RoomBuHuaNtf.Size(m)
}
func (m *RoomBuHuaNtf) XXX_DiscardUnknown() {
	xxx_messageInfo_RoomBuHuaNtf.DiscardUnknown(m)
}

var xxx_messageInfo_RoomBuHuaNtf proto.InternalMessageInfo

func (m *RoomBuHuaNtf) GetBuhuaInfo() []*RoomBuHuaInfo {
	if m != nil {
		return m.BuhuaInfo
	}
	return nil
}

func init() {
	proto.RegisterType((*RoomBuHuaInfo)(nil), "room.RoomBuHuaInfo")
	proto.RegisterType((*RoomBuHuaNtf)(nil), "room.RoomBuHuaNtf")
}

func init() { proto.RegisterFile("game_erren.proto", fileDescriptor_game_erren_aac2442421e12cfa) }

var fileDescriptor_game_erren_aac2442421e12cfa = []byte{
	// 200 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x8e, 0x41, 0x4b, 0x86, 0x40,
	0x10, 0x86, 0x31, 0x85, 0x74, 0x4c, 0x88, 0x2d, 0xc2, 0xe8, 0x22, 0x9e, 0x3c, 0x29, 0xf8, 0x13,
	0xec, 0xa2, 0x97, 0x0e, 0x7b, 0xec, 0xb2, 0xec, 0xea, 0x5a, 0x92, 0xee, 0xc8, 0xba, 0x13, 0xf4,
	0xef, 0x63, 0x2d, 0x3e, 0xf8, 0x6e, 0xc3, 0xfb, 0xcc, 0xcc, 0xfb, 0xc0, 0xfd, 0x87, 0xdc, 0xb4,
	0xd0, 0xd6, 0x6a, 0x53, 0xef, 0x16, 0x1d, 0xb2, 0xc8, 0x22, 0x6e, 0xe5, 0x17, 0x64, 0x1c, 0x71,
	0xeb, 0xa8, 0x27, 0x39, 0x98, 0x19, 0xd9, 0x0b, 0x24, 0xfb, 0x2a, 0x7f, 0xb4, 0x15, 0xcb, 0x94,
	0x07, 0x45, 0x50, 0x45, 0x3c, 0xfe, 0x0b, 0x86, 0x89, 0x95, 0x90, 0x21, 0x39, 0xf1, 0x49, 0x52,
	0x8c, 0xd2, 0x4e, 0x47, 0x7e, 0x53, 0x84, 0x55, 0xc6, 0x53, 0x24, 0xd7, 0x93, 0x7c, 0xf5, 0x11,
	0x7b, 0x86, 0x58, 0xd1, 0x3f, 0x0e, 0x4f, 0x7c, 0xab, 0xe8, 0x44, 0x65, 0x07, 0x77, 0x97, 0xb2,
	0x37, 0x37, 0xb3, 0x16, 0x40, 0x91, 0x7f, 0xb6, 0x98, 0x19, 0xf3, 0xa0, 0x08, 0xab, 0xb4, 0x7d,
	0xa8, 0xbd, 0x57, 0x7d, 0x25, 0xc5, 0x93, 0x73, 0xcd, 0x8f, 0xdd, 0xd3, 0xfb, 0xe3, 0xe1, 0xf4,
	0xb7, 0x6e, 0xc6, 0x75, 0xd1, 0xc6, 0x89, 0x5d, 0x35, 0xfe, 0xe0, 0x37, 0x00, 0x00, 0xff, 0xff,
	0xa8, 0x78, 0x76, 0x55, 0xe1, 0x00, 0x00, 0x00,
}
