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

// PlayerState 玩家状态
type PlayerState int32

const (
	PlayerState_PS_IDIE     PlayerState = 0
	PlayerState_PS_MATCHING PlayerState = 1
	PlayerState_PS_GAMEING  PlayerState = 2
)

var PlayerState_name = map[int32]string{
	0: "PS_IDIE",
	1: "PS_MATCHING",
	2: "PS_GAMEING",
}
var PlayerState_value = map[string]int32{
	"PS_IDIE":     0,
	"PS_MATCHING": 1,
	"PS_GAMEING":  2,
}

func (x PlayerState) String() string {
	return proto.EnumName(PlayerState_name, int32(x))
}
func (PlayerState) EnumDescriptor() ([]byte, []int) { return fileDescriptor2, []int{0} }

// ServerType 服务类型
type ServerType int32

const (
	ServerType_ST_GATE  ServerType = 0
	ServerType_ST_MATCH ServerType = 1
	ServerType_ST_ROOM  ServerType = 2
)

var ServerType_name = map[int32]string{
	0: "ST_GATE",
	1: "ST_MATCH",
	2: "ST_ROOM",
}
var ServerType_value = map[string]int32{
	"ST_GATE":  0,
	"ST_MATCH": 1,
	"ST_ROOM":  2,
}

func (x ServerType) String() string {
	return proto.EnumName(ServerType_name, int32(x))
}
func (ServerType) EnumDescriptor() ([]byte, []int) { return fileDescriptor2, []int{1} }

// GameConfig 游戏玩法信息
type GameConfig struct {
	GameId   uint32 `protobuf:"varint,1,opt,name=game_id,json=gameId" json:"game_id,omitempty"`
	GameName string `protobuf:"bytes,2,opt,name=game_name,json=gameName" json:"game_name,omitempty"`
	GameType uint32 `protobuf:"varint,3,opt,name=game_type,json=gameType" json:"game_type,omitempty"`
}

func (m *GameConfig) Reset()                    { *m = GameConfig{} }
func (m *GameConfig) String() string            { return proto.CompactTextString(m) }
func (*GameConfig) ProtoMessage()               {}
func (*GameConfig) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{0} }

func (m *GameConfig) GetGameId() uint32 {
	if m != nil {
		return m.GameId
	}
	return 0
}

func (m *GameConfig) GetGameName() string {
	if m != nil {
		return m.GameName
	}
	return ""
}

func (m *GameConfig) GetGameType() uint32 {
	if m != nil {
		return m.GameType
	}
	return 0
}

// GameConfigLevel 游戏场次信息
type GameLevelConfig struct {
	GameId     uint32 `protobuf:"varint,1,opt,name=game_id,json=gameId" json:"game_id,omitempty"`
	LevelId    uint32 `protobuf:"varint,2,opt,name=level_id,json=levelId" json:"level_id,omitempty"`
	LevelName  string `protobuf:"bytes,3,opt,name=level_name,json=levelName" json:"level_name,omitempty"`
	BaseScores uint32 `protobuf:"varint,4,opt,name=base_scores,json=baseScores" json:"base_scores,omitempty"`
	LowScores  uint32 `protobuf:"varint,5,opt,name=low_scores,json=lowScores" json:"low_scores,omitempty"`
	HighScores uint32 `protobuf:"varint,6,opt,name=high_scores,json=highScores" json:"high_scores,omitempty"`
	MinPeople  uint32 `protobuf:"varint,7,opt,name=min_people,json=minPeople" json:"min_people,omitempty"`
	MaxPeople  uint32 `protobuf:"varint,8,opt,name=max_people,json=maxPeople" json:"max_people,omitempty"`
}

func (m *GameLevelConfig) Reset()                    { *m = GameLevelConfig{} }
func (m *GameLevelConfig) String() string            { return proto.CompactTextString(m) }
func (*GameLevelConfig) ProtoMessage()               {}
func (*GameLevelConfig) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{1} }

func (m *GameLevelConfig) GetGameId() uint32 {
	if m != nil {
		return m.GameId
	}
	return 0
}

func (m *GameLevelConfig) GetLevelId() uint32 {
	if m != nil {
		return m.LevelId
	}
	return 0
}

func (m *GameLevelConfig) GetLevelName() string {
	if m != nil {
		return m.LevelName
	}
	return ""
}

func (m *GameLevelConfig) GetBaseScores() uint32 {
	if m != nil {
		return m.BaseScores
	}
	return 0
}

func (m *GameLevelConfig) GetLowScores() uint32 {
	if m != nil {
		return m.LowScores
	}
	return 0
}

func (m *GameLevelConfig) GetHighScores() uint32 {
	if m != nil {
		return m.HighScores
	}
	return 0
}

func (m *GameLevelConfig) GetMinPeople() uint32 {
	if m != nil {
		return m.MinPeople
	}
	return 0
}

func (m *GameLevelConfig) GetMaxPeople() uint32 {
	if m != nil {
		return m.MaxPeople
	}
	return 0
}

// GetPlayerByAccountReq 根据账号获取玩家请求
type GetPlayerByAccountReq struct {
	AccountId uint64 `protobuf:"varint,1,opt,name=account_id,json=accountId" json:"account_id,omitempty"`
}

func (m *GetPlayerByAccountReq) Reset()                    { *m = GetPlayerByAccountReq{} }
func (m *GetPlayerByAccountReq) String() string            { return proto.CompactTextString(m) }
func (*GetPlayerByAccountReq) ProtoMessage()               {}
func (*GetPlayerByAccountReq) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{2} }

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
func (*GetPlayerByAccountRsp) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{3} }

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

// GetPlayerInfoReq 获取玩家信息
type GetPlayerInfoReq struct {
	PlayerId uint64 `protobuf:"varint,1,opt,name=player_id,json=playerId" json:"player_id,omitempty"`
}

func (m *GetPlayerInfoReq) Reset()                    { *m = GetPlayerInfoReq{} }
func (m *GetPlayerInfoReq) String() string            { return proto.CompactTextString(m) }
func (*GetPlayerInfoReq) ProtoMessage()               {}
func (*GetPlayerInfoReq) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{4} }

func (m *GetPlayerInfoReq) GetPlayerId() uint64 {
	if m != nil {
		return m.PlayerId
	}
	return 0
}

// GetPlayerInfoRsp 获取玩家信息应答
type GetPlayerInfoRsp struct {
	ErrCode  int32  `protobuf:"varint,1,opt,name=err_code,json=errCode" json:"err_code,omitempty"`
	PlayerId uint64 `protobuf:"varint,2,opt,name=player_id,json=playerId" json:"player_id,omitempty"`
	NickName string `protobuf:"bytes,3,opt,name=nick_name,json=nickName" json:"nick_name,omitempty"`
	Avatar   string `protobuf:"bytes,4,opt,name=avatar" json:"avatar,omitempty"`
	Gender   uint32 `protobuf:"varint,5,opt,name=gender" json:"gender,omitempty"`
	Name     string `protobuf:"bytes,6,opt,name=name" json:"name,omitempty"`
	Phone    string `protobuf:"bytes,7,opt,name=phone" json:"phone,omitempty"`
	IpAddr   string `protobuf:"bytes,8,opt,name=ip_addr,json=ipAddr" json:"ip_addr,omitempty"`
}

func (m *GetPlayerInfoRsp) Reset()                    { *m = GetPlayerInfoRsp{} }
func (m *GetPlayerInfoRsp) String() string            { return proto.CompactTextString(m) }
func (*GetPlayerInfoRsp) ProtoMessage()               {}
func (*GetPlayerInfoRsp) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{5} }

func (m *GetPlayerInfoRsp) GetErrCode() int32 {
	if m != nil {
		return m.ErrCode
	}
	return 0
}

func (m *GetPlayerInfoRsp) GetPlayerId() uint64 {
	if m != nil {
		return m.PlayerId
	}
	return 0
}

func (m *GetPlayerInfoRsp) GetNickName() string {
	if m != nil {
		return m.NickName
	}
	return ""
}

func (m *GetPlayerInfoRsp) GetAvatar() string {
	if m != nil {
		return m.Avatar
	}
	return ""
}

func (m *GetPlayerInfoRsp) GetGender() uint32 {
	if m != nil {
		return m.Gender
	}
	return 0
}

func (m *GetPlayerInfoRsp) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *GetPlayerInfoRsp) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *GetPlayerInfoRsp) GetIpAddr() string {
	if m != nil {
		return m.IpAddr
	}
	return ""
}

// UpdatePlayerInfoReq 修改玩家信息
type UpdatePlayerInfoReq struct {
	PlayerId uint64 `protobuf:"varint,1,opt,name=player_id,json=playerId" json:"player_id,omitempty"`
	NickName string `protobuf:"bytes,2,opt,name=nick_name,json=nickName" json:"nick_name,omitempty"`
	Avatar   string `protobuf:"bytes,3,opt,name=avatar" json:"avatar,omitempty"`
	Gender   uint32 `protobuf:"varint,4,opt,name=gender" json:"gender,omitempty"`
	Name     string `protobuf:"bytes,5,opt,name=name" json:"name,omitempty"`
	Phone    string `protobuf:"bytes,6,opt,name=phone" json:"phone,omitempty"`
}

func (m *UpdatePlayerInfoReq) Reset()                    { *m = UpdatePlayerInfoReq{} }
func (m *UpdatePlayerInfoReq) String() string            { return proto.CompactTextString(m) }
func (*UpdatePlayerInfoReq) ProtoMessage()               {}
func (*UpdatePlayerInfoReq) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{6} }

func (m *UpdatePlayerInfoReq) GetPlayerId() uint64 {
	if m != nil {
		return m.PlayerId
	}
	return 0
}

func (m *UpdatePlayerInfoReq) GetNickName() string {
	if m != nil {
		return m.NickName
	}
	return ""
}

func (m *UpdatePlayerInfoReq) GetAvatar() string {
	if m != nil {
		return m.Avatar
	}
	return ""
}

func (m *UpdatePlayerInfoReq) GetGender() uint32 {
	if m != nil {
		return m.Gender
	}
	return 0
}

func (m *UpdatePlayerInfoReq) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *UpdatePlayerInfoReq) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

// UpdatePlayerInfoRsp 修改玩家信息应答
type UpdatePlayerInfoRsp struct {
	ErrCode int32 `protobuf:"varint,1,opt,name=err_code,json=errCode" json:"err_code,omitempty"`
	Result  bool  `protobuf:"varint,2,opt,name=result" json:"result,omitempty"`
}

func (m *UpdatePlayerInfoRsp) Reset()                    { *m = UpdatePlayerInfoRsp{} }
func (m *UpdatePlayerInfoRsp) String() string            { return proto.CompactTextString(m) }
func (*UpdatePlayerInfoRsp) ProtoMessage()               {}
func (*UpdatePlayerInfoRsp) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{7} }

func (m *UpdatePlayerInfoRsp) GetErrCode() int32 {
	if m != nil {
		return m.ErrCode
	}
	return 0
}

func (m *UpdatePlayerInfoRsp) GetResult() bool {
	if m != nil {
		return m.Result
	}
	return false
}

// GetPlayerStateReq 获取玩家状态
type GetPlayerStateReq struct {
	PlayerId uint64 `protobuf:"varint,1,opt,name=player_id,json=playerId" json:"player_id,omitempty"`
}

func (m *GetPlayerStateReq) Reset()                    { *m = GetPlayerStateReq{} }
func (m *GetPlayerStateReq) String() string            { return proto.CompactTextString(m) }
func (*GetPlayerStateReq) ProtoMessage()               {}
func (*GetPlayerStateReq) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{8} }

func (m *GetPlayerStateReq) GetPlayerId() uint64 {
	if m != nil {
		return m.PlayerId
	}
	return 0
}

// GetPlayerStateRsp 获取玩家状态应答
type GetPlayerStateRsp struct {
	ErrCode int32       `protobuf:"varint,1,opt,name=err_code,json=errCode" json:"err_code,omitempty"`
	State   PlayerState `protobuf:"varint,2,opt,name=state,enum=user.PlayerState" json:"state,omitempty"`
}

func (m *GetPlayerStateRsp) Reset()                    { *m = GetPlayerStateRsp{} }
func (m *GetPlayerStateRsp) String() string            { return proto.CompactTextString(m) }
func (*GetPlayerStateRsp) ProtoMessage()               {}
func (*GetPlayerStateRsp) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{9} }

func (m *GetPlayerStateRsp) GetErrCode() int32 {
	if m != nil {
		return m.ErrCode
	}
	return 0
}

func (m *GetPlayerStateRsp) GetState() PlayerState {
	if m != nil {
		return m.State
	}
	return PlayerState_PS_IDIE
}

// GetPlayerGameInfoReq 获取玩家游戏信息
type GetPlayerGameInfoReq struct {
	PlayerId uint64 `protobuf:"varint,1,opt,name=player_id,json=playerId" json:"player_id,omitempty"`
	GameId   uint32 `protobuf:"varint,2,opt,name=game_id,json=gameId" json:"game_id,omitempty"`
}

func (m *GetPlayerGameInfoReq) Reset()                    { *m = GetPlayerGameInfoReq{} }
func (m *GetPlayerGameInfoReq) String() string            { return proto.CompactTextString(m) }
func (*GetPlayerGameInfoReq) ProtoMessage()               {}
func (*GetPlayerGameInfoReq) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{10} }

func (m *GetPlayerGameInfoReq) GetPlayerId() uint64 {
	if m != nil {
		return m.PlayerId
	}
	return 0
}

func (m *GetPlayerGameInfoReq) GetGameId() uint32 {
	if m != nil {
		return m.GameId
	}
	return 0
}

// GetPlayerGameInfoRsp 获取玩家游戏信息应答
type GetPlayerGameInfoRsp struct {
	ErrCode          int32  `protobuf:"varint,1,opt,name=err_code,json=errCode" json:"err_code,omitempty"`
	GameId           uint32 `protobuf:"varint,2,opt,name=game_id,json=gameId" json:"game_id,omitempty"`
	WinningRate      uint32 `protobuf:"varint,3,opt,name=winning_rate,json=winningRate" json:"winning_rate,omitempty"`
	TotalBurea       uint32 `protobuf:"varint,4,opt,name=total_burea,json=totalBurea" json:"total_burea,omitempty"`
	MaxWinningStream uint32 `protobuf:"varint,5,opt,name=max_winning_stream,json=maxWinningStream" json:"max_winning_stream,omitempty"`
	MaxMultiple      uint32 `protobuf:"varint,6,opt,name=max_multiple,json=maxMultiple" json:"max_multiple,omitempty"`
}

func (m *GetPlayerGameInfoRsp) Reset()                    { *m = GetPlayerGameInfoRsp{} }
func (m *GetPlayerGameInfoRsp) String() string            { return proto.CompactTextString(m) }
func (*GetPlayerGameInfoRsp) ProtoMessage()               {}
func (*GetPlayerGameInfoRsp) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{11} }

func (m *GetPlayerGameInfoRsp) GetErrCode() int32 {
	if m != nil {
		return m.ErrCode
	}
	return 0
}

func (m *GetPlayerGameInfoRsp) GetGameId() uint32 {
	if m != nil {
		return m.GameId
	}
	return 0
}

func (m *GetPlayerGameInfoRsp) GetWinningRate() uint32 {
	if m != nil {
		return m.WinningRate
	}
	return 0
}

func (m *GetPlayerGameInfoRsp) GetTotalBurea() uint32 {
	if m != nil {
		return m.TotalBurea
	}
	return 0
}

func (m *GetPlayerGameInfoRsp) GetMaxWinningStream() uint32 {
	if m != nil {
		return m.MaxWinningStream
	}
	return 0
}

func (m *GetPlayerGameInfoRsp) GetMaxMultiple() uint32 {
	if m != nil {
		return m.MaxMultiple
	}
	return 0
}

// UpdatePlayerStateReq 更新玩家状态
type UpdatePlayerStateReq struct {
	PlayerId   uint64      `protobuf:"varint,1,opt,name=player_id,json=playerId" json:"player_id,omitempty"`
	OldState   PlayerState `protobuf:"varint,2,opt,name=old_state,json=oldState,enum=user.PlayerState" json:"old_state,omitempty"`
	NewState   PlayerState `protobuf:"varint,3,opt,name=new_state,json=newState,enum=user.PlayerState" json:"new_state,omitempty"`
	ServerType ServerType  `protobuf:"varint,4,opt,name=server_type,json=serverType,enum=user.ServerType" json:"server_type,omitempty"`
	ServerAddr string      `protobuf:"bytes,5,opt,name=server_addr,json=serverAddr" json:"server_addr,omitempty"`
}

func (m *UpdatePlayerStateReq) Reset()                    { *m = UpdatePlayerStateReq{} }
func (m *UpdatePlayerStateReq) String() string            { return proto.CompactTextString(m) }
func (*UpdatePlayerStateReq) ProtoMessage()               {}
func (*UpdatePlayerStateReq) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{12} }

func (m *UpdatePlayerStateReq) GetPlayerId() uint64 {
	if m != nil {
		return m.PlayerId
	}
	return 0
}

func (m *UpdatePlayerStateReq) GetOldState() PlayerState {
	if m != nil {
		return m.OldState
	}
	return PlayerState_PS_IDIE
}

func (m *UpdatePlayerStateReq) GetNewState() PlayerState {
	if m != nil {
		return m.NewState
	}
	return PlayerState_PS_IDIE
}

func (m *UpdatePlayerStateReq) GetServerType() ServerType {
	if m != nil {
		return m.ServerType
	}
	return ServerType_ST_GATE
}

func (m *UpdatePlayerStateReq) GetServerAddr() string {
	if m != nil {
		return m.ServerAddr
	}
	return ""
}

// UpdatePlayerStateRsp  更新玩家状态应答
type UpdatePlayerStateRsp struct {
	ErrCode int32 `protobuf:"varint,1,opt,name=err_code,json=errCode" json:"err_code,omitempty"`
	Result  bool  `protobuf:"varint,2,opt,name=result" json:"result,omitempty"`
}

func (m *UpdatePlayerStateRsp) Reset()                    { *m = UpdatePlayerStateRsp{} }
func (m *UpdatePlayerStateRsp) String() string            { return proto.CompactTextString(m) }
func (*UpdatePlayerStateRsp) ProtoMessage()               {}
func (*UpdatePlayerStateRsp) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{13} }

func (m *UpdatePlayerStateRsp) GetErrCode() int32 {
	if m != nil {
		return m.ErrCode
	}
	return 0
}

func (m *UpdatePlayerStateRsp) GetResult() bool {
	if m != nil {
		return m.Result
	}
	return false
}

// GetGameListInfoReq 获取游戏列表信息
type GetGameListInfoReq struct {
}

func (m *GetGameListInfoReq) Reset()                    { *m = GetGameListInfoReq{} }
func (m *GetGameListInfoReq) String() string            { return proto.CompactTextString(m) }
func (*GetGameListInfoReq) ProtoMessage()               {}
func (*GetGameListInfoReq) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{14} }

// GetGameListInfoRsp 获取游戏列表信息应答
type GetGameListInfoRsp struct {
	ErrCode         int32              `protobuf:"varint,1,opt,name=err_code,json=errCode" json:"err_code,omitempty"`
	GameConfig      []*GameConfig      `protobuf:"bytes,2,rep,name=game_config,json=gameConfig" json:"game_config,omitempty"`
	GameLevelConfig []*GameLevelConfig `protobuf:"bytes,3,rep,name=game_level_config,json=gameLevelConfig" json:"game_level_config,omitempty"`
}

func (m *GetGameListInfoRsp) Reset()                    { *m = GetGameListInfoRsp{} }
func (m *GetGameListInfoRsp) String() string            { return proto.CompactTextString(m) }
func (*GetGameListInfoRsp) ProtoMessage()               {}
func (*GetGameListInfoRsp) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{15} }

func (m *GetGameListInfoRsp) GetErrCode() int32 {
	if m != nil {
		return m.ErrCode
	}
	return 0
}

func (m *GetGameListInfoRsp) GetGameConfig() []*GameConfig {
	if m != nil {
		return m.GameConfig
	}
	return nil
}

func (m *GetGameListInfoRsp) GetGameLevelConfig() []*GameLevelConfig {
	if m != nil {
		return m.GameLevelConfig
	}
	return nil
}

// GetPlayerServerInfoReq
type GetPlayerServerInfoReq struct {
	PlayerId uint64 `protobuf:"varint,1,opt,name=player_id,json=playerId" json:"player_id,omitempty"`
}

func (m *GetPlayerServerInfoReq) Reset()                    { *m = GetPlayerServerInfoReq{} }
func (m *GetPlayerServerInfoReq) String() string            { return proto.CompactTextString(m) }
func (*GetPlayerServerInfoReq) ProtoMessage()               {}
func (*GetPlayerServerInfoReq) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{16} }

func (m *GetPlayerServerInfoReq) GetPlayerId() uint64 {
	if m != nil {
		return m.PlayerId
	}
	return 0
}

// GetPlayerServerInfoRsp 获取玩家服务端信息
type GetPlayerServerInfoRsp struct {
	ErrCode   int32  `protobuf:"varint,1,opt,name=err_code,json=errCode" json:"err_code,omitempty"`
	PlayerId  uint64 `protobuf:"varint,2,opt,name=player_id,json=playerId" json:"player_id,omitempty"`
	ClientId  string `protobuf:"bytes,3,opt,name=client_id,json=clientId" json:"client_id,omitempty"`
	MatchAddr string `protobuf:"bytes,4,opt,name=match_addr,json=matchAddr" json:"match_addr,omitempty"`
	GateAddr  string `protobuf:"bytes,5,opt,name=gate_addr,json=gateAddr" json:"gate_addr,omitempty"`
	RoomAddr  string `protobuf:"bytes,6,opt,name=room_addr,json=roomAddr" json:"room_addr,omitempty"`
}

func (m *GetPlayerServerInfoRsp) Reset()                    { *m = GetPlayerServerInfoRsp{} }
func (m *GetPlayerServerInfoRsp) String() string            { return proto.CompactTextString(m) }
func (*GetPlayerServerInfoRsp) ProtoMessage()               {}
func (*GetPlayerServerInfoRsp) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{17} }

func (m *GetPlayerServerInfoRsp) GetErrCode() int32 {
	if m != nil {
		return m.ErrCode
	}
	return 0
}

func (m *GetPlayerServerInfoRsp) GetPlayerId() uint64 {
	if m != nil {
		return m.PlayerId
	}
	return 0
}

func (m *GetPlayerServerInfoRsp) GetClientId() string {
	if m != nil {
		return m.ClientId
	}
	return ""
}

func (m *GetPlayerServerInfoRsp) GetMatchAddr() string {
	if m != nil {
		return m.MatchAddr
	}
	return ""
}

func (m *GetPlayerServerInfoRsp) GetGateAddr() string {
	if m != nil {
		return m.GateAddr
	}
	return ""
}

func (m *GetPlayerServerInfoRsp) GetRoomAddr() string {
	if m != nil {
		return m.RoomAddr
	}
	return ""
}

func init() {
	proto.RegisterType((*GameConfig)(nil), "user.GameConfig")
	proto.RegisterType((*GameLevelConfig)(nil), "user.GameLevelConfig")
	proto.RegisterType((*GetPlayerByAccountReq)(nil), "user.GetPlayerByAccountReq")
	proto.RegisterType((*GetPlayerByAccountRsp)(nil), "user.GetPlayerByAccountRsp")
	proto.RegisterType((*GetPlayerInfoReq)(nil), "user.GetPlayerInfoReq")
	proto.RegisterType((*GetPlayerInfoRsp)(nil), "user.GetPlayerInfoRsp")
	proto.RegisterType((*UpdatePlayerInfoReq)(nil), "user.UpdatePlayerInfoReq")
	proto.RegisterType((*UpdatePlayerInfoRsp)(nil), "user.UpdatePlayerInfoRsp")
	proto.RegisterType((*GetPlayerStateReq)(nil), "user.GetPlayerStateReq")
	proto.RegisterType((*GetPlayerStateRsp)(nil), "user.GetPlayerStateRsp")
	proto.RegisterType((*GetPlayerGameInfoReq)(nil), "user.GetPlayerGameInfoReq")
	proto.RegisterType((*GetPlayerGameInfoRsp)(nil), "user.GetPlayerGameInfoRsp")
	proto.RegisterType((*UpdatePlayerStateReq)(nil), "user.UpdatePlayerStateReq")
	proto.RegisterType((*UpdatePlayerStateRsp)(nil), "user.UpdatePlayerStateRsp")
	proto.RegisterType((*GetGameListInfoReq)(nil), "user.GetGameListInfoReq")
	proto.RegisterType((*GetGameListInfoRsp)(nil), "user.GetGameListInfoRsp")
	proto.RegisterType((*GetPlayerServerInfoReq)(nil), "user.GetPlayerServerInfoReq")
	proto.RegisterType((*GetPlayerServerInfoRsp)(nil), "user.GetPlayerServerInfoRsp")
	proto.RegisterEnum("user.PlayerState", PlayerState_name, PlayerState_value)
	proto.RegisterEnum("user.ServerType", ServerType_name, ServerType_value)
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
	// GetPlayerInfo 获取玩家信息
	GetPlayerInfo(ctx context.Context, in *GetPlayerInfoReq, opts ...grpc.CallOption) (*GetPlayerInfoRsp, error)
	// UpdatePlayerInfo 更新玩家信息
	UpdatePlayerInfo(ctx context.Context, in *UpdatePlayerInfoReq, opts ...grpc.CallOption) (*UpdatePlayerInfoRsp, error)
	// GetPlayerState 获取玩家状态
	GetPlayerState(ctx context.Context, in *GetPlayerStateReq, opts ...grpc.CallOption) (*GetPlayerStateRsp, error)
	// GetPlayerGameInfo 获取玩家游戏信息
	GetPlayerGameInfo(ctx context.Context, in *GetPlayerGameInfoReq, opts ...grpc.CallOption) (*GetPlayerGameInfoRsp, error)
	// UpdatePlayerState 更新玩家状态
	UpdatePlayerState(ctx context.Context, in *UpdatePlayerStateReq, opts ...grpc.CallOption) (*UpdatePlayerStateRsp, error)
	// GetGameListInfo 获取游戏列表
	GetGameListInfo(ctx context.Context, in *GetGameListInfoReq, opts ...grpc.CallOption) (*GetGameListInfoRsp, error)
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

func (c *playerDataClient) GetPlayerInfo(ctx context.Context, in *GetPlayerInfoReq, opts ...grpc.CallOption) (*GetPlayerInfoRsp, error) {
	out := new(GetPlayerInfoRsp)
	err := grpc.Invoke(ctx, "/user.PlayerData/GetPlayerInfo", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *playerDataClient) UpdatePlayerInfo(ctx context.Context, in *UpdatePlayerInfoReq, opts ...grpc.CallOption) (*UpdatePlayerInfoRsp, error) {
	out := new(UpdatePlayerInfoRsp)
	err := grpc.Invoke(ctx, "/user.PlayerData/UpdatePlayerInfo", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *playerDataClient) GetPlayerState(ctx context.Context, in *GetPlayerStateReq, opts ...grpc.CallOption) (*GetPlayerStateRsp, error) {
	out := new(GetPlayerStateRsp)
	err := grpc.Invoke(ctx, "/user.PlayerData/GetPlayerState", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *playerDataClient) GetPlayerGameInfo(ctx context.Context, in *GetPlayerGameInfoReq, opts ...grpc.CallOption) (*GetPlayerGameInfoRsp, error) {
	out := new(GetPlayerGameInfoRsp)
	err := grpc.Invoke(ctx, "/user.PlayerData/GetPlayerGameInfo", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *playerDataClient) UpdatePlayerState(ctx context.Context, in *UpdatePlayerStateReq, opts ...grpc.CallOption) (*UpdatePlayerStateRsp, error) {
	out := new(UpdatePlayerStateRsp)
	err := grpc.Invoke(ctx, "/user.PlayerData/UpdatePlayerState", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *playerDataClient) GetGameListInfo(ctx context.Context, in *GetGameListInfoReq, opts ...grpc.CallOption) (*GetGameListInfoRsp, error) {
	out := new(GetGameListInfoRsp)
	err := grpc.Invoke(ctx, "/user.PlayerData/GetGameListInfo", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for PlayerData service

type PlayerDataServer interface {
	// GetPlayerByAccount 根据账号获取玩家
	GetPlayerByAccount(context.Context, *GetPlayerByAccountReq) (*GetPlayerByAccountRsp, error)
	// GetPlayerInfo 获取玩家信息
	GetPlayerInfo(context.Context, *GetPlayerInfoReq) (*GetPlayerInfoRsp, error)
	// UpdatePlayerInfo 更新玩家信息
	UpdatePlayerInfo(context.Context, *UpdatePlayerInfoReq) (*UpdatePlayerInfoRsp, error)
	// GetPlayerState 获取玩家状态
	GetPlayerState(context.Context, *GetPlayerStateReq) (*GetPlayerStateRsp, error)
	// GetPlayerGameInfo 获取玩家游戏信息
	GetPlayerGameInfo(context.Context, *GetPlayerGameInfoReq) (*GetPlayerGameInfoRsp, error)
	// UpdatePlayerState 更新玩家状态
	UpdatePlayerState(context.Context, *UpdatePlayerStateReq) (*UpdatePlayerStateRsp, error)
	// GetGameListInfo 获取游戏列表
	GetGameListInfo(context.Context, *GetGameListInfoReq) (*GetGameListInfoRsp, error)
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

func _PlayerData_GetPlayerInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPlayerInfoReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlayerDataServer).GetPlayerInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.PlayerData/GetPlayerInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlayerDataServer).GetPlayerInfo(ctx, req.(*GetPlayerInfoReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _PlayerData_UpdatePlayerInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdatePlayerInfoReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlayerDataServer).UpdatePlayerInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.PlayerData/UpdatePlayerInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlayerDataServer).UpdatePlayerInfo(ctx, req.(*UpdatePlayerInfoReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _PlayerData_GetPlayerState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPlayerStateReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlayerDataServer).GetPlayerState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.PlayerData/GetPlayerState",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlayerDataServer).GetPlayerState(ctx, req.(*GetPlayerStateReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _PlayerData_GetPlayerGameInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPlayerGameInfoReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlayerDataServer).GetPlayerGameInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.PlayerData/GetPlayerGameInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlayerDataServer).GetPlayerGameInfo(ctx, req.(*GetPlayerGameInfoReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _PlayerData_UpdatePlayerState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdatePlayerStateReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlayerDataServer).UpdatePlayerState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.PlayerData/UpdatePlayerState",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlayerDataServer).UpdatePlayerState(ctx, req.(*UpdatePlayerStateReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _PlayerData_GetGameListInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetGameListInfoReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlayerDataServer).GetGameListInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.PlayerData/GetGameListInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlayerDataServer).GetGameListInfo(ctx, req.(*GetGameListInfoReq))
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
		{
			MethodName: "GetPlayerInfo",
			Handler:    _PlayerData_GetPlayerInfo_Handler,
		},
		{
			MethodName: "UpdatePlayerInfo",
			Handler:    _PlayerData_UpdatePlayerInfo_Handler,
		},
		{
			MethodName: "GetPlayerState",
			Handler:    _PlayerData_GetPlayerState_Handler,
		},
		{
			MethodName: "GetPlayerGameInfo",
			Handler:    _PlayerData_GetPlayerGameInfo_Handler,
		},
		{
			MethodName: "UpdatePlayerState",
			Handler:    _PlayerData_UpdatePlayerState_Handler,
		},
		{
			MethodName: "GetGameListInfo",
			Handler:    _PlayerData_GetGameListInfo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service.proto",
}

func init() { proto.RegisterFile("service.proto", fileDescriptor2) }

var fileDescriptor2 = []byte{
	// 1007 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x56, 0xcb, 0x6e, 0xdb, 0x46,
	0x14, 0x2d, 0x25, 0x59, 0x22, 0x2f, 0x63, 0x5b, 0x9e, 0x3a, 0x8e, 0xa2, 0xa0, 0x88, 0xcb, 0x4d,
	0x8d, 0xa0, 0x70, 0x9b, 0xf4, 0xb1, 0xe9, 0xa2, 0x90, 0x1d, 0x43, 0x16, 0x1a, 0xdb, 0x02, 0xa5,
	0x22, 0x4b, 0x62, 0x2c, 0x4e, 0x64, 0xa2, 0x7c, 0x75, 0x38, 0xb2, 0xec, 0x75, 0xb7, 0xfd, 0x8d,
	0x2e, 0xfa, 0x23, 0xfd, 0x88, 0x6e, 0xba, 0xec, 0x6f, 0x14, 0x73, 0x67, 0x28, 0x51, 0x32, 0x69,
	0x1b, 0xe9, 0x8e, 0x73, 0xce, 0x3d, 0xf3, 0xba, 0x77, 0xce, 0x25, 0x6c, 0x66, 0x8c, 0x5f, 0x07,
	0x13, 0x76, 0x98, 0xf2, 0x44, 0x24, 0xa4, 0x31, 0xcb, 0x18, 0x77, 0x3c, 0x80, 0x3e, 0x8d, 0xd8,
	0x71, 0x12, 0x7f, 0x08, 0xa6, 0xe4, 0x19, 0xb4, 0xa6, 0x34, 0x62, 0x5e, 0xe0, 0x77, 0x8c, 0x7d,
	0xe3, 0x60, 0xd3, 0x6d, 0xca, 0xe1, 0xc0, 0x27, 0x2f, 0xc0, 0x42, 0x22, 0xa6, 0x11, 0xeb, 0xd4,
	0xf6, 0x8d, 0x03, 0xcb, 0x35, 0x25, 0x70, 0x4e, 0x23, 0xb6, 0x20, 0xc5, 0x6d, 0xca, 0x3a, 0x75,
	0xd4, 0x21, 0x39, 0xbe, 0x4d, 0x99, 0xf3, 0x5b, 0x0d, 0xb6, 0xe5, 0x0a, 0xef, 0xd8, 0x35, 0x0b,
	0x1f, 0x5a, 0xe6, 0x39, 0x98, 0xa1, 0x8c, 0x93, 0x4c, 0x0d, 0x99, 0x16, 0x8e, 0x07, 0x3e, 0xf9,
	0x0c, 0x40, 0x51, 0xb8, 0x85, 0x3a, 0x6e, 0xc1, 0x42, 0x04, 0xf7, 0xf0, 0x12, 0xec, 0x4b, 0x9a,
	0x31, 0x2f, 0x9b, 0x24, 0x9c, 0x65, 0x9d, 0x06, 0x8a, 0x41, 0x42, 0x23, 0x44, 0x50, 0x9f, 0xcc,
	0x73, 0x7e, 0x03, 0x79, 0x2b, 0x4c, 0xe6, 0x9a, 0x7e, 0x09, 0xf6, 0x55, 0x30, 0xbd, 0xca, 0xf9,
	0xa6, 0xd2, 0x4b, 0x68, 0xa9, 0x8f, 0x82, 0xd8, 0x4b, 0x59, 0x92, 0x86, 0xac, 0xd3, 0x52, 0xfa,
	0x28, 0x88, 0x87, 0x08, 0x20, 0x4d, 0x6f, 0x72, 0xda, 0xd4, 0x34, 0xbd, 0x51, 0xb4, 0xf3, 0x3d,
	0x3c, 0xed, 0x33, 0x31, 0x0c, 0xe9, 0x2d, 0xe3, 0x47, 0xb7, 0xbd, 0xc9, 0x24, 0x99, 0xc5, 0xc2,
	0x65, 0xbf, 0x4a, 0x1d, 0x55, 0xa3, 0xfc, 0x36, 0x1a, 0xae, 0xa5, 0x91, 0x81, 0xef, 0x5c, 0x94,
	0xea, 0xb2, 0x54, 0xde, 0x14, 0xe3, 0xdc, 0x9b, 0x24, 0x3e, 0x43, 0xd5, 0x86, 0xdb, 0x62, 0x9c,
	0x1f, 0x27, 0x3e, 0xa6, 0x23, 0x45, 0x41, 0x7e, 0x8b, 0x0d, 0xd7, 0x54, 0xc0, 0xc0, 0x77, 0xbe,
	0x82, 0xf6, 0x62, 0xc2, 0x41, 0xfc, 0x21, 0x91, 0x7b, 0x58, 0x11, 0x18, 0x6b, 0x82, 0xbf, 0x8d,
	0x75, 0xc5, 0xc7, 0xaf, 0x2e, 0xc9, 0x38, 0x98, 0xfc, 0x52, 0xcc, 0xa1, 0x29, 0x01, 0x4c, 0xe1,
	0x1e, 0x34, 0xe9, 0x35, 0x15, 0x94, 0x63, 0xf6, 0x2c, 0x57, 0x8f, 0x24, 0x3e, 0x65, 0xb1, 0xcf,
	0xb8, 0xce, 0x9a, 0x1e, 0x11, 0x02, 0x0d, 0x9c, 0xa7, 0x89, 0xd1, 0xf8, 0x4d, 0x76, 0x61, 0x23,
	0xbd, 0x4a, 0x62, 0x95, 0x20, 0xcb, 0x55, 0x03, 0x59, 0x6f, 0x41, 0xea, 0x51, 0xdf, 0xe7, 0x98,
	0x19, 0xcb, 0x6d, 0x06, 0x69, 0xcf, 0xf7, 0xb9, 0xf3, 0xa7, 0x01, 0x9f, 0xfe, 0x9c, 0xfa, 0x54,
	0xb0, 0xc7, 0xdf, 0xc8, 0xea, 0x21, 0x6a, 0x95, 0x87, 0xa8, 0x57, 0x1c, 0xa2, 0x51, 0x7a, 0x88,
	0x8d, 0xb2, 0x43, 0x34, 0x0b, 0x87, 0x70, 0x4e, 0x4b, 0xb6, 0x7a, 0x7f, 0x2a, 0xf6, 0xa0, 0xc9,
	0x59, 0x36, 0x0b, 0x05, 0xee, 0xd2, 0x74, 0xf5, 0xc8, 0xf9, 0x1a, 0x76, 0x16, 0x19, 0x1d, 0x09,
	0x2a, 0xd8, 0x83, 0x45, 0xf0, 0xfe, 0x8e, 0xe2, 0xfe, 0x95, 0xbf, 0x80, 0x8d, 0x4c, 0x86, 0xe1,
	0xc2, 0x5b, 0x6f, 0x76, 0x0e, 0xa5, 0xd7, 0x1c, 0x16, 0xf5, 0x8a, 0x77, 0xde, 0xc1, 0xee, 0x62,
	0x62, 0xe9, 0x12, 0x8f, 0x4a, 0x40, 0xc1, 0x3e, 0x6a, 0x45, 0xfb, 0x70, 0xfe, 0x31, 0xca, 0xa6,
	0xbb, 0x7f, 0xab, 0x55, 0x93, 0x91, 0xcf, 0xe1, 0xc9, 0x3c, 0x88, 0xe3, 0x20, 0x9e, 0x7a, 0x5c,
	0x1e, 0x45, 0x19, 0x9b, 0xad, 0x31, 0x97, 0x0a, 0x34, 0x1d, 0x91, 0x08, 0x1a, 0x7a, 0x97, 0x33,
	0xce, 0x68, 0x6e, 0x3a, 0x08, 0x1d, 0x49, 0x84, 0x7c, 0x09, 0x44, 0xba, 0x42, 0x3e, 0x4f, 0x26,
	0x38, 0xa3, 0x91, 0x2e, 0xe3, 0x76, 0x44, 0x6f, 0xde, 0x2b, 0x62, 0x84, 0xb8, 0x5c, 0x51, 0x46,
	0x47, 0xb3, 0x50, 0x04, 0xd2, 0x45, 0x94, 0x09, 0xd9, 0x11, 0xbd, 0x39, 0xd3, 0x90, 0xf3, 0xaf,
	0x01, 0xbb, 0xc5, 0x2a, 0x78, 0x54, 0xfa, 0xc8, 0x21, 0x58, 0x49, 0xe8, 0x7b, 0x0f, 0xa4, 0xc4,
	0x4c, 0x42, 0x1f, 0xbf, 0x64, 0x7c, 0xcc, 0xe6, 0x3a, 0xbe, 0x5e, 0x19, 0x1f, 0xb3, 0xb9, 0x8a,
	0x7f, 0x0d, 0xb6, 0xec, 0x2d, 0x8c, 0xab, 0x16, 0xd0, 0x40, 0x45, 0x5b, 0x29, 0x46, 0x48, 0xc8,
	0x56, 0xe0, 0x42, 0xb6, 0xf8, 0x96, 0x57, 0xa7, 0x25, 0xf8, 0x2c, 0x55, 0xf9, 0xeb, 0x00, 0x7c,
	0x9a, 0x83, 0xb2, 0x83, 0x7e, 0x5c, 0xbd, 0xef, 0x02, 0xe9, 0x33, 0x81, 0x4d, 0x28, 0xc8, 0x84,
	0x2e, 0x31, 0xe7, 0x0f, 0xe3, 0x2e, 0x7c, 0xff, 0xfc, 0xaf, 0xc1, 0xc6, 0x52, 0x99, 0x60, 0x17,
	0xeb, 0xd4, 0xf6, 0xeb, 0x07, 0x76, 0x7e, 0xcc, 0x65, 0x13, 0x75, 0x61, 0xba, 0x6c, 0xa8, 0x3d,
	0xd8, 0x41, 0x89, 0x6a, 0x5d, 0x5a, 0x58, 0x47, 0xe1, 0xd3, 0xa5, 0xb0, 0xd0, 0x1b, 0xdd, 0xed,
	0xe9, 0x2a, 0xe0, 0x7c, 0x07, 0x7b, 0xcb, 0xb7, 0x87, 0xf7, 0xf3, 0x28, 0xdf, 0xfe, 0xcb, 0x28,
	0xd7, 0xfd, 0x3f, 0xf7, 0x9e, 0x84, 0x01, 0x53, 0xad, 0x4a, 0xbb, 0xb7, 0x02, 0x54, 0x7f, 0x8e,
	0xa8, 0x98, 0x5c, 0xa9, 0x7c, 0x2a, 0x07, 0xb7, 0x10, 0x91, 0xe9, 0x54, 0xff, 0x08, 0x82, 0x15,
	0xb3, 0x6d, 0x4a, 0x20, 0x27, 0x79, 0x92, 0x44, 0x8a, 0x54, 0xa6, 0x67, 0x4a, 0x40, 0x92, 0xaf,
	0x7e, 0x00, 0xbb, 0x50, 0x02, 0xc4, 0x86, 0xd6, 0x70, 0xe4, 0x0d, 0xde, 0x0e, 0x4e, 0xda, 0x9f,
	0x90, 0x6d, 0xb0, 0x87, 0x23, 0xef, 0xac, 0x37, 0x3e, 0x3e, 0x1d, 0x9c, 0xf7, 0xdb, 0x06, 0xd9,
	0x02, 0x18, 0x8e, 0xbc, 0x7e, 0xef, 0xec, 0x44, 0x8e, 0x6b, 0xaf, 0xbe, 0x05, 0x58, 0x16, 0xa0,
	0xd4, 0x8e, 0xc6, 0x5e, 0xbf, 0x37, 0x96, 0xda, 0x27, 0x60, 0x8e, 0xc6, 0x4a, 0xdb, 0x36, 0x34,
	0xe5, 0x5e, 0x5c, 0x9c, 0xb5, 0x6b, 0x6f, 0x7e, 0x6f, 0x00, 0xa8, 0x35, 0xdf, 0x52, 0x41, 0xc9,
	0x39, 0x16, 0xca, 0x5a, 0x13, 0x26, 0x2f, 0x74, 0xfe, 0xca, 0xda, 0x7a, 0xb7, 0x9a, 0xcc, 0x52,
	0xf2, 0x23, 0x6c, 0xae, 0x74, 0x54, 0xb2, 0xb7, 0x16, 0xad, 0x13, 0xdc, 0x2d, 0xc5, 0xb3, 0x94,
	0x9c, 0x42, 0x7b, 0xbd, 0x15, 0x90, 0xe7, 0x2a, 0xb6, 0xa4, 0x9b, 0x75, 0xab, 0xa8, 0x2c, 0x25,
	0x47, 0xb0, 0xb5, 0x6a, 0xec, 0xe4, 0xd9, 0xda, 0x9a, 0xb9, 0xc3, 0x74, 0xcb, 0x89, 0x2c, 0x25,
	0x3f, 0x15, 0x9a, 0x43, 0x6e, 0xba, 0xa4, 0xbb, 0x16, 0x5d, 0x30, 0xf7, 0x6e, 0x25, 0xa7, 0x26,
	0xbb, 0xf3, 0xec, 0xf3, 0xc9, 0xca, 0x8c, 0xaf, 0x5b, 0xc9, 0x65, 0x29, 0x39, 0x81, 0xed, 0xb5,
	0x17, 0x4e, 0x3a, 0x8b, 0xb5, 0xd7, 0xfc, 0xa0, 0x5b, 0xc1, 0x64, 0xe9, 0x65, 0x13, 0x7f, 0x98,
	0xbf, 0xf9, 0x2f, 0x00, 0x00, 0xff, 0xff, 0xae, 0xd5, 0x83, 0x84, 0x41, 0x0b, 0x00, 0x00,
}
