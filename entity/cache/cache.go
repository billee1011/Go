package cache

import "fmt"

// HallPlayer 大厅玩家
type HallPlayer struct {
	PlayerID uint64 `protobuf:"varint,1,opt,name=playerID" json:"playerID,omitempty"`
	NickName string `protobuf:"bytes,2,opt,name=nickName" json:"nickName,omitempty"`
	Avatar   string `protobuf:"bytes,3,opt,name=avatar" json:"avatar,omitempty"`
	Gender   uint64 `protobuf:"bytes,4,opt,name=gender" json:"gender,omitempty"`
	Name     string `protobuf:"bytes,5,opt,name=name" json:"name,omitempty"`
	Phone    string `protobuf:"bytes,6,opt,name=phone" json:"phone,omitempty"`
	IDdCard  string `protobuf:"bytes,7,opt,name=idCard" json:"idCard,omitempty"`
	Coin     uint64 `protobuf:"varint,8,opt,name=coin" json:"coin,omitempty"`
	GameID   uint64 `protobuf:"varint,9,opt,name=gameID" json:"gameID,omitempty"`
	State    uint64 `protobuf:"varint,10,opt,name=state" json:"state,omitempty"`
}

// RobotPlayer 机器人玩家
type RobotPlayer struct {
	PlayerID      uint64            `protobuf:"varint,1,opt,name=player_id,json=playerId" json:"player_id,omitempty"`
	NickName      string            `protobuf:"bytes,2,opt,name=nick_name,json=nickName" json:"nick_name,omitempty"`
	HeadImage     string            `protobuf:"bytes,3,opt,name=head_image,json=headImage" json:"head_image,omitempty"`
	Coin          uint64            `protobuf:"varint,4,opt,name=coin" json:"coin,omitempty"`
	State         uint64            `protobuf:"varint,5,opt,name=state" json:"state,omitempty"`
	GameIDWinRate map[uint64]uint64 `protobuf:"bytes,6,rep,name=game_id_win_rate,json=gameIdWinRate" json:"game_id_win_rate,omitempty" protobuf_key:"varint,1,opt,name=key" protobuf_val:"varint,2,opt,name=value"`
}

const (
	// AccountPlayerKey 账号关联的玩家
	AccountPlayerKey = "account:player:%v"
)

// Player 属性字段
const (
	// NickNameField ...昵称
	NickNameField = "nickname"
	// AvatarField ...头像
	AvatarField = "avatar"
	// GenderField  ...性别
	GenderField = "gender"
	// NameField  ...姓名
	NameField = "name"
	// PhoneField  ...联系电话
	PhoneField = "phone"
	// PlayerStateField ...玩家状态
	PlayerStateField = "player_state"
	// CoinField ...金币
	CoinField = "coin"
	// GateAddrField ...网关服地址
	GateAddrField = "gate_addr"
	// MatchAddrField ...匹配服地址
	MatchAddrField = "match_addr"
	// RoomAddrField ...房间服地址
	RoomAddrField = "room_addr"
)

// FmtAccountPlayerKey 账号所关联玩家 key
func FmtAccountPlayerKey(accountID uint64) string {
	return fmt.Sprintf(AccountPlayerKey, accountID)
}

// FmtPlayerIDKey 玩家ID的 key
func FmtPlayerIDKey(playerID uint64) string {
	return fmt.Sprintf("player:%v", playerID)
}

// FmtGameInfoKey 游戏配置信息的 key
func FmtGameInfoKey() string {
	return fmt.Sprintf("game:info")
}

// FmtPlayerStateKey 玩家State的 key
func FmtPlayerStateKey(playerID uint64) string {
	return fmt.Sprintf("playerState:%v", playerID)
}
