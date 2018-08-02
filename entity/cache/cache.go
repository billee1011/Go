package cache

import "fmt"

// RobotPlayer 机器人玩家
type RobotPlayer struct {
	PlayerID      uint64            `protobuf:"varint,1,opt,name=player_id,json=playerId" json:"player_id,omitempty"`
	NickName      string            `protobuf:"bytes,2,opt,name=nick_name,json=nickName" json:"nick_name,omitempty"`
	HeadImage     string            `protobuf:"bytes,3,opt,name=head_image,json=headImage" json:"head_image,omitempty"`
	Coin          uint64            `protobuf:"varint,4,opt,name=coin" json:"coin,omitempty"`
	State         uint64            `protobuf:"varint,5,opt,name=state" json:"state,omitempty"`
	GameIDWinRate map[uint64]uint64 `protobuf:"bytes,6,rep,name=game_id_win_rate,json=gameIdWinRate" json:"game_id_win_rate,omitempty" protobuf_key:"varint,1,opt,name=key" protobuf_val:"varint,2,opt,name=value"`
}

// key formats
const (
	// AccountPlayerKey 账号关联的玩家
	AccountPlayerKey = "account:player:%v"

	// playerTokenKeyFmt
	playerTokenKeyFmt = "playertoken:%d"
)

// Player 字段
const (
	// NickName ...昵称
	NickName = "nickname"
	// Avatar ...头像
	Avatar = "avatar"
	// Gender  ...性别
	Gender = "gender"
	// ChannelID ...渠道ID
	ChannelID = "channelID"
	// ProvinceID ...省份ID
	ProvinceID = "provinceID"
	// CityID ...城市ID
	CityID = "cityID"

	// GameState ...玩家游戏状态
	GameState = "game_state"
	// GameID ...正在进行的游戏id
	GameID = "game_id"
	// IPAddr ... 玩家ip地址
	IPAddr = "ip_addr"
	// GateAddr ...网关服地址
	GateAddr = "gate_addr"
	// MatchAddr ...匹配服地址
	MatchAddr = "match_addr"
	// RoomAddr ...房间服地址
	RoomAddr = "room_addr"
)

// FmtAccountPlayerKey 账号所关联玩家 key
func FmtAccountPlayerKey(accountID uint64) string {
	return fmt.Sprintf(AccountPlayerKey, accountID)
}

// FmtPlayerIDKey 玩家ID key
func FmtPlayerIDKey(playerID uint64) string {
	return fmt.Sprintf("player:%v", playerID)
}

// FmtPlayerTokenKey format player's token key
func FmtPlayerTokenKey(playerID uint64) string {
	return fmt.Sprintf(playerTokenKeyFmt, playerID)
}
