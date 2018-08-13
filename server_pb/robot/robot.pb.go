// Code generated by protoc-gen-go. DO NOT EDIT.
// source: robot.proto

package robot

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

// GameConfig 游戏玩法信息
type GameConfig struct {
	GameId  uint32 `protobuf:"varint,1,opt,name=game_id,json=gameId" json:"game_id,omitempty"`
	LevelId uint32 `protobuf:"varint,2,opt,name=level_id,json=levelId" json:"level_id,omitempty"`
}

func (m *GameConfig) Reset()                    { *m = GameConfig{} }
func (m *GameConfig) String() string            { return proto.CompactTextString(m) }
func (*GameConfig) ProtoMessage()               {}
func (*GameConfig) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

func (m *GameConfig) GetGameId() uint32 {
	if m != nil {
		return m.GameId
	}
	return 0
}

func (m *GameConfig) GetLevelId() uint32 {
	if m != nil {
		return m.LevelId
	}
	return 0
}

// GameWinRate 游戏对应的胜率
type GameWinRate struct {
	Game    *GameConfig `protobuf:"bytes,1,opt,name=game" json:"game,omitempty"`
	WinRate int32       `protobuf:"varint,2,opt,name=win_rate,json=winRate" json:"win_rate,omitempty"`
}

func (m *GameWinRate) Reset()                    { *m = GameWinRate{} }
func (m *GameWinRate) String() string            { return proto.CompactTextString(m) }
func (*GameWinRate) ProtoMessage()               {}
func (*GameWinRate) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{1} }

func (m *GameWinRate) GetGame() *GameConfig {
	if m != nil {
		return m.Game
	}
	return nil
}

func (m *GameWinRate) GetWinRate() int32 {
	if m != nil {
		return m.WinRate
	}
	return 0
}

// RobotPlayerInfo 机器人玩家信息
type RobotPlayerInfo struct {
	NickName    string           `protobuf:"bytes,1,opt,name=nick_name,json=nickName" json:"nick_name,omitempty"`
	Avatar      string           `protobuf:"bytes,2,opt,name=avatar" json:"avatar,omitempty"`
	Coin        uint64           `protobuf:"varint,3,opt,name=coin" json:"coin,omitempty"`
	State       RobotPlayerState `protobuf:"varint,4,opt,name=state,enum=robot.RobotPlayerState" json:"state,omitempty"`
	GameWinRate *GameWinRate     `protobuf:"bytes,5,opt,name=game_win_rate,json=gameWinRate" json:"game_win_rate,omitempty"`
}

func (m *RobotPlayerInfo) Reset()                    { *m = RobotPlayerInfo{} }
func (m *RobotPlayerInfo) String() string            { return proto.CompactTextString(m) }
func (*RobotPlayerInfo) ProtoMessage()               {}
func (*RobotPlayerInfo) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{2} }

func (m *RobotPlayerInfo) GetNickName() string {
	if m != nil {
		return m.NickName
	}
	return ""
}

func (m *RobotPlayerInfo) GetAvatar() string {
	if m != nil {
		return m.Avatar
	}
	return ""
}

func (m *RobotPlayerInfo) GetCoin() uint64 {
	if m != nil {
		return m.Coin
	}
	return 0
}

func (m *RobotPlayerInfo) GetState() RobotPlayerState {
	if m != nil {
		return m.State
	}
	return RobotPlayerState_RPS_IDIE
}

func (m *RobotPlayerInfo) GetGameWinRate() *GameWinRate {
	if m != nil {
		return m.GameWinRate
	}
	return nil
}

// WinRateRange 胜率范围
type WinRateRange struct {
	High int32 `protobuf:"varint,1,opt,name=high" json:"high,omitempty"`
	Low  int32 `protobuf:"varint,2,opt,name=low" json:"low,omitempty"`
}

func (m *WinRateRange) Reset()                    { *m = WinRateRange{} }
func (m *WinRateRange) String() string            { return proto.CompactTextString(m) }
func (*WinRateRange) ProtoMessage()               {}
func (*WinRateRange) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{3} }

func (m *WinRateRange) GetHigh() int32 {
	if m != nil {
		return m.High
	}
	return 0
}

func (m *WinRateRange) GetLow() int32 {
	if m != nil {
		return m.Low
	}
	return 0
}

// CoinsRange 金币范围
type CoinsRange struct {
	High int64 `protobuf:"varint,1,opt,name=high" json:"high,omitempty"`
	Low  int64 `protobuf:"varint,2,opt,name=low" json:"low,omitempty"`
}

func (m *CoinsRange) Reset()                    { *m = CoinsRange{} }
func (m *CoinsRange) String() string            { return proto.CompactTextString(m) }
func (*CoinsRange) ProtoMessage()               {}
func (*CoinsRange) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{4} }

func (m *CoinsRange) GetHigh() int64 {
	if m != nil {
		return m.High
	}
	return 0
}

func (m *CoinsRange) GetLow() int64 {
	if m != nil {
		return m.Low
	}
	return 0
}

// GetLeisureRobotInfoReq 获取空闲机器人信息请求
type GetLeisureRobotInfoReq struct {
	Game         *GameConfig      `protobuf:"bytes,1,opt,name=game" json:"game,omitempty"`
	WinRateRange *WinRateRange    `protobuf:"bytes,2,opt,name=win_rate_range,json=winRateRange" json:"win_rate_range,omitempty"`
	CoinsRange   *CoinsRange      `protobuf:"bytes,3,opt,name=coins_range,json=coinsRange" json:"coins_range,omitempty"`
	NewState     RobotPlayerState `protobuf:"varint,4,opt,name=new_state,json=newState,enum=robot.RobotPlayerState" json:"new_state,omitempty"`
}

func (m *GetLeisureRobotInfoReq) Reset()                    { *m = GetLeisureRobotInfoReq{} }
func (m *GetLeisureRobotInfoReq) String() string            { return proto.CompactTextString(m) }
func (*GetLeisureRobotInfoReq) ProtoMessage()               {}
func (*GetLeisureRobotInfoReq) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{5} }

func (m *GetLeisureRobotInfoReq) GetGame() *GameConfig {
	if m != nil {
		return m.Game
	}
	return nil
}

func (m *GetLeisureRobotInfoReq) GetWinRateRange() *WinRateRange {
	if m != nil {
		return m.WinRateRange
	}
	return nil
}

func (m *GetLeisureRobotInfoReq) GetCoinsRange() *CoinsRange {
	if m != nil {
		return m.CoinsRange
	}
	return nil
}

func (m *GetLeisureRobotInfoReq) GetNewState() RobotPlayerState {
	if m != nil {
		return m.NewState
	}
	return RobotPlayerState_RPS_IDIE
}

// GetLeisureRobotInfoRsp 获取空闲机器人信息响应
type GetLeisureRobotInfoRsp struct {
	RobotPlayerId uint64  `protobuf:"varint,1,opt,name=robot_player_id,json=robotPlayerId" json:"robot_player_id,omitempty"`
	Coin          int64   `protobuf:"varint,2,opt,name=coin" json:"coin,omitempty"`
	WinRate       float64 `protobuf:"fixed64,3,opt,name=win_rate,json=winRate" json:"win_rate,omitempty"`
	ErrCode       int32   `protobuf:"varint,4,opt,name=err_code,json=errCode" json:"err_code,omitempty"`
}

func (m *GetLeisureRobotInfoRsp) Reset()                    { *m = GetLeisureRobotInfoRsp{} }
func (m *GetLeisureRobotInfoRsp) String() string            { return proto.CompactTextString(m) }
func (*GetLeisureRobotInfoRsp) ProtoMessage()               {}
func (*GetLeisureRobotInfoRsp) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{6} }

func (m *GetLeisureRobotInfoRsp) GetRobotPlayerId() uint64 {
	if m != nil {
		return m.RobotPlayerId
	}
	return 0
}

func (m *GetLeisureRobotInfoRsp) GetCoin() int64 {
	if m != nil {
		return m.Coin
	}
	return 0
}

func (m *GetLeisureRobotInfoRsp) GetWinRate() float64 {
	if m != nil {
		return m.WinRate
	}
	return 0
}

func (m *GetLeisureRobotInfoRsp) GetErrCode() int32 {
	if m != nil {
		return m.ErrCode
	}
	return 0
}

// SetRobotPlayerStateReq 設置机器人玩家状态請求
type SetRobotPlayerStateReq struct {
	RobotPlayerId uint64           `protobuf:"varint,1,opt,name=robot_player_id,json=robotPlayerId" json:"robot_player_id,omitempty"`
	NewState      RobotPlayerState `protobuf:"varint,2,opt,name=new_state,json=newState,enum=robot.RobotPlayerState" json:"new_state,omitempty"`
	OldState      RobotPlayerState `protobuf:"varint,3,opt,name=old_state,json=oldState,enum=robot.RobotPlayerState" json:"old_state,omitempty"`
	ServerType    ServerType       `protobuf:"varint,4,opt,name=server_type,json=serverType,enum=robot.ServerType" json:"server_type,omitempty"`
	ServerAddr    string           `protobuf:"bytes,5,opt,name=server_addr,json=serverAddr" json:"server_addr,omitempty"`
}

func (m *SetRobotPlayerStateReq) Reset()                    { *m = SetRobotPlayerStateReq{} }
func (m *SetRobotPlayerStateReq) String() string            { return proto.CompactTextString(m) }
func (*SetRobotPlayerStateReq) ProtoMessage()               {}
func (*SetRobotPlayerStateReq) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{7} }

func (m *SetRobotPlayerStateReq) GetRobotPlayerId() uint64 {
	if m != nil {
		return m.RobotPlayerId
	}
	return 0
}

func (m *SetRobotPlayerStateReq) GetNewState() RobotPlayerState {
	if m != nil {
		return m.NewState
	}
	return RobotPlayerState_RPS_IDIE
}

func (m *SetRobotPlayerStateReq) GetOldState() RobotPlayerState {
	if m != nil {
		return m.OldState
	}
	return RobotPlayerState_RPS_IDIE
}

func (m *SetRobotPlayerStateReq) GetServerType() ServerType {
	if m != nil {
		return m.ServerType
	}
	return ServerType_ST_GATE
}

func (m *SetRobotPlayerStateReq) GetServerAddr() string {
	if m != nil {
		return m.ServerAddr
	}
	return ""
}

// SetRobotPlayerStateRsp 設置机器人玩家状态响应
type SetRobotPlayerStateRsp struct {
	Result  bool  `protobuf:"varint,1,opt,name=result" json:"result,omitempty"`
	ErrCode int32 `protobuf:"varint,2,opt,name=err_code,json=errCode" json:"err_code,omitempty"`
}

func (m *SetRobotPlayerStateRsp) Reset()                    { *m = SetRobotPlayerStateRsp{} }
func (m *SetRobotPlayerStateRsp) String() string            { return proto.CompactTextString(m) }
func (*SetRobotPlayerStateRsp) ProtoMessage()               {}
func (*SetRobotPlayerStateRsp) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{8} }

func (m *SetRobotPlayerStateRsp) GetResult() bool {
	if m != nil {
		return m.Result
	}
	return false
}

func (m *SetRobotPlayerStateRsp) GetErrCode() int32 {
	if m != nil {
		return m.ErrCode
	}
	return 0
}

// UpdataRobotGameWinRateReq 更新机器人胜率請求
type UpdataRobotGameWinRateReq struct {
	RobotPlayerId uint64  `protobuf:"varint,1,opt,name=robot_player_id,json=robotPlayerId" json:"robot_player_id,omitempty"`
	GameId        int32   `protobuf:"varint,2,opt,name=game_id,json=gameId" json:"game_id,omitempty"`
	OldWinRate    float64 `protobuf:"fixed64,3,opt,name=oldWinRate" json:"oldWinRate,omitempty"`
	NewWinRate    float64 `protobuf:"fixed64,4,opt,name=newWinRate" json:"newWinRate,omitempty"`
}

func (m *UpdataRobotGameWinRateReq) Reset()                    { *m = UpdataRobotGameWinRateReq{} }
func (m *UpdataRobotGameWinRateReq) String() string            { return proto.CompactTextString(m) }
func (*UpdataRobotGameWinRateReq) ProtoMessage()               {}
func (*UpdataRobotGameWinRateReq) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{9} }

func (m *UpdataRobotGameWinRateReq) GetRobotPlayerId() uint64 {
	if m != nil {
		return m.RobotPlayerId
	}
	return 0
}

func (m *UpdataRobotGameWinRateReq) GetGameId() int32 {
	if m != nil {
		return m.GameId
	}
	return 0
}

func (m *UpdataRobotGameWinRateReq) GetOldWinRate() float64 {
	if m != nil {
		return m.OldWinRate
	}
	return 0
}

func (m *UpdataRobotGameWinRateReq) GetNewWinRate() float64 {
	if m != nil {
		return m.NewWinRate
	}
	return 0
}

// UpdataRobotGameWinRateRsp 更新机器人胜率响应
type UpdataRobotGameWinRateRsp struct {
	Result  bool  `protobuf:"varint,1,opt,name=result" json:"result,omitempty"`
	ErrCode int32 `protobuf:"varint,2,opt,name=err_code,json=errCode" json:"err_code,omitempty"`
}

func (m *UpdataRobotGameWinRateRsp) Reset()                    { *m = UpdataRobotGameWinRateRsp{} }
func (m *UpdataRobotGameWinRateRsp) String() string            { return proto.CompactTextString(m) }
func (*UpdataRobotGameWinRateRsp) ProtoMessage()               {}
func (*UpdataRobotGameWinRateRsp) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{10} }

func (m *UpdataRobotGameWinRateRsp) GetResult() bool {
	if m != nil {
		return m.Result
	}
	return false
}

func (m *UpdataRobotGameWinRateRsp) GetErrCode() int32 {
	if m != nil {
		return m.ErrCode
	}
	return 0
}

// IsRobotPlayerReq 判断是否时机器人請求
type IsRobotPlayerReq struct {
	RobotPlayerId uint64 `protobuf:"varint,1,opt,name=robot_player_id,json=robotPlayerId" json:"robot_player_id,omitempty"`
}

func (m *IsRobotPlayerReq) Reset()                    { *m = IsRobotPlayerReq{} }
func (m *IsRobotPlayerReq) String() string            { return proto.CompactTextString(m) }
func (*IsRobotPlayerReq) ProtoMessage()               {}
func (*IsRobotPlayerReq) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{11} }

func (m *IsRobotPlayerReq) GetRobotPlayerId() uint64 {
	if m != nil {
		return m.RobotPlayerId
	}
	return 0
}

// IsRobotPlayerRsp 判断是否时机器人响应
type IsRobotPlayerRsp struct {
	Result bool `protobuf:"varint,1,opt,name=result" json:"result,omitempty"`
}

func (m *IsRobotPlayerRsp) Reset()                    { *m = IsRobotPlayerRsp{} }
func (m *IsRobotPlayerRsp) String() string            { return proto.CompactTextString(m) }
func (*IsRobotPlayerRsp) ProtoMessage()               {}
func (*IsRobotPlayerRsp) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{12} }

func (m *IsRobotPlayerRsp) GetResult() bool {
	if m != nil {
		return m.Result
	}
	return false
}

func init() {
	proto.RegisterType((*GameConfig)(nil), "robot.GameConfig")
	proto.RegisterType((*GameWinRate)(nil), "robot.GameWinRate")
	proto.RegisterType((*RobotPlayerInfo)(nil), "robot.RobotPlayerInfo")
	proto.RegisterType((*WinRateRange)(nil), "robot.WinRateRange")
	proto.RegisterType((*CoinsRange)(nil), "robot.CoinsRange")
	proto.RegisterType((*GetLeisureRobotInfoReq)(nil), "robot.GetLeisureRobotInfoReq")
	proto.RegisterType((*GetLeisureRobotInfoRsp)(nil), "robot.GetLeisureRobotInfoRsp")
	proto.RegisterType((*SetRobotPlayerStateReq)(nil), "robot.SetRobotPlayerStateReq")
	proto.RegisterType((*SetRobotPlayerStateRsp)(nil), "robot.SetRobotPlayerStateRsp")
	proto.RegisterType((*UpdataRobotGameWinRateReq)(nil), "robot.UpdataRobotGameWinRateReq")
	proto.RegisterType((*UpdataRobotGameWinRateRsp)(nil), "robot.UpdataRobotGameWinRateRsp")
	proto.RegisterType((*IsRobotPlayerReq)(nil), "robot.IsRobotPlayerReq")
	proto.RegisterType((*IsRobotPlayerRsp)(nil), "robot.IsRobotPlayerRsp")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for RobotService service

type RobotServiceClient interface {
	GetLeisureRobotInfoByInfo(ctx context.Context, in *GetLeisureRobotInfoReq, opts ...grpc.CallOption) (*GetLeisureRobotInfoRsp, error)
	SetRobotPlayerState(ctx context.Context, in *SetRobotPlayerStateReq, opts ...grpc.CallOption) (*SetRobotPlayerStateRsp, error)
	UpdataRobotGameWinRate(ctx context.Context, in *UpdataRobotGameWinRateReq, opts ...grpc.CallOption) (*UpdataRobotGameWinRateRsp, error)
	IsRobotPlayer(ctx context.Context, in *IsRobotPlayerReq, opts ...grpc.CallOption) (*IsRobotPlayerRsp, error)
}

type robotServiceClient struct {
	cc *grpc.ClientConn
}

func NewRobotServiceClient(cc *grpc.ClientConn) RobotServiceClient {
	return &robotServiceClient{cc}
}

func (c *robotServiceClient) GetLeisureRobotInfoByInfo(ctx context.Context, in *GetLeisureRobotInfoReq, opts ...grpc.CallOption) (*GetLeisureRobotInfoRsp, error) {
	out := new(GetLeisureRobotInfoRsp)
	err := grpc.Invoke(ctx, "/robot.RobotService/GetLeisureRobotInfoByInfo", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *robotServiceClient) SetRobotPlayerState(ctx context.Context, in *SetRobotPlayerStateReq, opts ...grpc.CallOption) (*SetRobotPlayerStateRsp, error) {
	out := new(SetRobotPlayerStateRsp)
	err := grpc.Invoke(ctx, "/robot.RobotService/SetRobotPlayerState", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *robotServiceClient) UpdataRobotGameWinRate(ctx context.Context, in *UpdataRobotGameWinRateReq, opts ...grpc.CallOption) (*UpdataRobotGameWinRateRsp, error) {
	out := new(UpdataRobotGameWinRateRsp)
	err := grpc.Invoke(ctx, "/robot.RobotService/UpdataRobotGameWinRate", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *robotServiceClient) IsRobotPlayer(ctx context.Context, in *IsRobotPlayerReq, opts ...grpc.CallOption) (*IsRobotPlayerRsp, error) {
	out := new(IsRobotPlayerRsp)
	err := grpc.Invoke(ctx, "/robot.RobotService/IsRobotPlayer", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for RobotService service

type RobotServiceServer interface {
	GetLeisureRobotInfoByInfo(context.Context, *GetLeisureRobotInfoReq) (*GetLeisureRobotInfoRsp, error)
	SetRobotPlayerState(context.Context, *SetRobotPlayerStateReq) (*SetRobotPlayerStateRsp, error)
	UpdataRobotGameWinRate(context.Context, *UpdataRobotGameWinRateReq) (*UpdataRobotGameWinRateRsp, error)
	IsRobotPlayer(context.Context, *IsRobotPlayerReq) (*IsRobotPlayerRsp, error)
}

func RegisterRobotServiceServer(s *grpc.Server, srv RobotServiceServer) {
	s.RegisterService(&_RobotService_serviceDesc, srv)
}

func _RobotService_GetLeisureRobotInfoByInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetLeisureRobotInfoReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RobotServiceServer).GetLeisureRobotInfoByInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/robot.RobotService/GetLeisureRobotInfoByInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RobotServiceServer).GetLeisureRobotInfoByInfo(ctx, req.(*GetLeisureRobotInfoReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _RobotService_SetRobotPlayerState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetRobotPlayerStateReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RobotServiceServer).SetRobotPlayerState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/robot.RobotService/SetRobotPlayerState",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RobotServiceServer).SetRobotPlayerState(ctx, req.(*SetRobotPlayerStateReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _RobotService_UpdataRobotGameWinRate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdataRobotGameWinRateReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RobotServiceServer).UpdataRobotGameWinRate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/robot.RobotService/UpdataRobotGameWinRate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RobotServiceServer).UpdataRobotGameWinRate(ctx, req.(*UpdataRobotGameWinRateReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _RobotService_IsRobotPlayer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IsRobotPlayerReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RobotServiceServer).IsRobotPlayer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/robot.RobotService/IsRobotPlayer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RobotServiceServer).IsRobotPlayer(ctx, req.(*IsRobotPlayerReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _RobotService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "robot.RobotService",
	HandlerType: (*RobotServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetLeisureRobotInfoByInfo",
			Handler:    _RobotService_GetLeisureRobotInfoByInfo_Handler,
		},
		{
			MethodName: "SetRobotPlayerState",
			Handler:    _RobotService_SetRobotPlayerState_Handler,
		},
		{
			MethodName: "UpdataRobotGameWinRate",
			Handler:    _RobotService_UpdataRobotGameWinRate_Handler,
		},
		{
			MethodName: "IsRobotPlayer",
			Handler:    _RobotService_IsRobotPlayer_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "robot.proto",
}

func init() { proto.RegisterFile("robot.proto", fileDescriptor1) }

var fileDescriptor1 = []byte{
	// 686 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x55, 0x5d, 0x6f, 0xd3, 0x30,
	0x14, 0x25, 0xfd, 0x58, 0xdb, 0x9b, 0x75, 0x1b, 0x9e, 0xd4, 0x6d, 0x45, 0x40, 0x15, 0x09, 0x34,
	0x21, 0xb1, 0x87, 0x30, 0x21, 0xc1, 0x13, 0xd0, 0x87, 0xa9, 0x02, 0x0d, 0xe4, 0x31, 0xf1, 0x80,
	0x44, 0x94, 0xd5, 0x77, 0x5d, 0x44, 0x16, 0x67, 0x4e, 0xb6, 0xa8, 0x8f, 0xbc, 0xf2, 0x1b, 0xf8,
	0x33, 0xfc, 0x17, 0x24, 0xfe, 0x06, 0xf2, 0x4d, 0xd2, 0xa6, 0x6b, 0xba, 0x75, 0xbc, 0xd9, 0x3e,
	0x39, 0xc7, 0xf7, 0x1c, 0x5f, 0x3b, 0x60, 0x2a, 0x79, 0x22, 0xe3, 0xbd, 0x50, 0xc9, 0x58, 0xb2,
	0x3a, 0x4d, 0xba, 0x80, 0xc1, 0xe5, 0x79, 0xba, 0x64, 0xbd, 0x01, 0x38, 0x70, 0xcf, 0xb1, 0x2f,
	0x83, 0x53, 0x6f, 0xc4, 0xb6, 0xa0, 0x31, 0x72, 0xcf, 0xd1, 0xf1, 0xc4, 0xb6, 0xd1, 0x33, 0x76,
	0xdb, 0x7c, 0x45, 0x4f, 0x07, 0x82, 0xed, 0x40, 0xd3, 0xc7, 0x2b, 0xf4, 0x35, 0x52, 0x21, 0xa4,
	0x41, 0xf3, 0x81, 0xb0, 0x3e, 0x82, 0xa9, 0x15, 0xbe, 0x78, 0x01, 0x77, 0x63, 0x64, 0x4f, 0xa0,
	0xa6, 0x39, 0xc4, 0x37, 0xed, 0xfb, 0x7b, 0xe9, 0xfe, 0xd3, 0x3d, 0x38, 0xc1, 0x5a, 0x30, 0xf1,
	0x02, 0x47, 0xb9, 0x31, 0x92, 0x60, 0x9d, 0x37, 0x92, 0x54, 0xc1, 0xfa, 0x6d, 0xc0, 0x3a, 0xd7,
	0xac, 0x4f, 0xbe, 0x3b, 0x46, 0x35, 0x08, 0x4e, 0x25, 0x7b, 0x00, 0xad, 0xc0, 0x1b, 0x7e, 0x77,
	0x82, 0x5c, 0xba, 0xc5, 0x9b, 0x7a, 0xe1, 0x50, 0x6b, 0x75, 0x60, 0xc5, 0xbd, 0x72, 0x63, 0x57,
	0x91, 0x52, 0x8b, 0x67, 0x33, 0xc6, 0xa0, 0x36, 0x94, 0x5e, 0xb0, 0x5d, 0xed, 0x19, 0xbb, 0x35,
	0x4e, 0x63, 0xf6, 0x1c, 0xea, 0x51, 0xac, 0x37, 0xad, 0xf5, 0x8c, 0xdd, 0x35, 0x7b, 0x2b, 0xab,
	0xaf, 0xb0, 0xdf, 0x91, 0x86, 0x79, 0xfa, 0x15, 0x7b, 0x09, 0x6d, 0x0a, 0x64, 0x52, 0x6b, 0x9d,
	0x6c, 0xb1, 0x82, 0xad, 0xcc, 0x38, 0x37, 0x47, 0xd3, 0x89, 0xb5, 0x0f, 0xab, 0xf9, 0xba, 0x1b,
	0x8c, 0x50, 0x97, 0x72, 0xe6, 0x8d, 0xce, 0xa8, 0xf4, 0x3a, 0xa7, 0x31, 0xdb, 0x80, 0xaa, 0x2f,
	0x93, 0xcc, 0xbd, 0x1e, 0x5a, 0x36, 0x40, 0x5f, 0x7a, 0x41, 0x34, 0xcf, 0xa9, 0xce, 0x73, 0xaa,
	0x29, 0xe7, 0x8f, 0x01, 0x9d, 0x03, 0x8c, 0x3f, 0xa0, 0x17, 0x5d, 0x2a, 0x24, 0x1f, 0x3a, 0x31,
	0x8e, 0x17, 0xcb, 0x1e, 0xc5, 0x2b, 0x58, 0xcb, 0xed, 0x39, 0x4a, 0xef, 0x4c, 0xf2, 0xa6, 0xbd,
	0x99, 0x11, 0x8a, 0x46, 0xf8, 0x6a, 0x52, 0xb4, 0x65, 0x83, 0xa9, 0x53, 0x8d, 0x32, 0x5e, 0x75,
	0x66, 0xa3, 0xa9, 0x15, 0x0e, 0xc3, 0xa9, 0xad, 0x7d, 0x68, 0x05, 0x98, 0x38, 0x4b, 0x9d, 0x42,
	0x33, 0xc0, 0x84, 0x46, 0xd6, 0xcf, 0x05, 0x36, 0xa3, 0x90, 0x3d, 0x85, 0x75, 0xa2, 0x3b, 0x21,
	0x31, 0xf3, 0xe6, 0xad, 0xf1, 0xb6, 0x2a, 0x74, 0x91, 0x98, 0xb4, 0x43, 0x1a, 0x5e, 0xda, 0x0e,
	0xc5, 0x36, 0xd4, 0xd5, 0x1b, 0x93, 0x36, 0xd4, 0x10, 0x2a, 0xe5, 0x0c, 0xa5, 0x48, 0xcb, 0xac,
	0xf3, 0x06, 0x2a, 0xd5, 0x97, 0x02, 0xad, 0x1f, 0x15, 0xe8, 0x1c, 0x61, 0x3c, 0x57, 0x2e, 0x5e,
	0x2c, 0x5d, 0xcc, 0x4c, 0x0a, 0x95, 0x25, 0x53, 0xd0, 0x2c, 0xe9, 0x8b, 0x8c, 0x55, 0xbd, 0x85,
	0x25, 0x7d, 0x91, 0xb2, 0x6c, 0x30, 0x23, 0x54, 0x57, 0xa8, 0x9c, 0x78, 0x1c, 0xe6, 0x99, 0xe7,
	0xa7, 0x74, 0x44, 0xc8, 0xe7, 0x71, 0x88, 0x1c, 0xa2, 0xc9, 0x98, 0x3d, 0x9e, 0x70, 0x5c, 0x21,
	0x14, 0xb5, 0x7d, 0x2b, 0xff, 0xe0, 0xad, 0x10, 0xca, 0x7a, 0x5f, 0x1e, 0x41, 0x14, 0xea, 0xeb,
	0xa8, 0x30, 0xba, 0xf4, 0x63, 0x72, 0xde, 0xe4, 0xd9, 0x6c, 0x26, 0xd0, 0xca, 0x6c, 0xa0, 0xbf,
	0x0c, 0xd8, 0x39, 0x0e, 0x85, 0x1b, 0xbb, 0x24, 0x58, 0xbc, 0x56, 0x77, 0xc8, 0xb4, 0xf0, 0x7a,
	0xa5, 0xfa, 0xf9, 0xeb, 0xf5, 0x08, 0x40, 0xfa, 0x22, 0x53, 0xcc, 0xce, 0xb9, 0xb0, 0xa2, 0xf1,
	0x00, 0x93, 0x1c, 0xaf, 0xa5, 0xf8, 0x74, 0xc5, 0x3a, 0x5c, 0x58, 0xdd, 0xff, 0xd9, 0x7d, 0x0d,
	0x1b, 0x83, 0xa8, 0x10, 0xdd, 0x1d, 0x4c, 0x5a, 0xcf, 0xae, 0x73, 0x17, 0x97, 0x60, 0xff, 0xad,
	0xc0, 0x2a, 0x7d, 0xaa, 0x0f, 0xd9, 0x1b, 0x22, 0xfb, 0x0a, 0x3b, 0x25, 0x97, 0xe8, 0xdd, 0x98,
	0xde, 0xd8, 0x87, 0xf9, 0x03, 0x51, 0xfa, 0x9a, 0x74, 0x6f, 0x82, 0xa3, 0xd0, 0xba, 0xc7, 0x8e,
	0x61, 0xb3, 0xa4, 0x23, 0x26, 0xb2, 0xe5, 0x17, 0xa6, 0x7b, 0x13, 0x4c, 0xb2, 0xdf, 0xa0, 0x53,
	0x1e, 0x3e, 0xeb, 0x65, 0xd4, 0x85, 0x9d, 0xd3, 0xbd, 0xe5, 0x0b, 0xd2, 0xef, 0x43, 0x7b, 0x26,
	0x50, 0x96, 0xdf, 0xa8, 0xeb, 0x47, 0xd4, 0x2d, 0x07, 0xb4, 0xc8, 0xc9, 0x0a, 0xfd, 0x4d, 0x5f,
	0xfc, 0x0b, 0x00, 0x00, 0xff, 0xff, 0xc1, 0x00, 0xb0, 0xfb, 0x6f, 0x07, 0x00, 0x00,
}
