package cache

import "fmt"

// HallPlayer 大厅玩家
type HallPlayer struct {
	PlayerID  uint64 `protobuf:"varint,1,opt,name=playerID" json:"playerID,omitempty"`
	NickName  string `protobuf:"bytes,2,opt,name=nickName" json:"nickName,omitempty"`
	HeadImage string `protobuf:"bytes,3,opt,name=headImage" json:"headImage,omitempty"`
	Coin      uint64 `protobuf:"varint,4,opt,name=coin" json:"coin,omitempty"`
	GameID    uint64 `protobuf:"varint,5,opt,name=gameID" json:"gameID,omitempty"`
	State     uint64 `protobuf:"varint,6,opt,name=state" json:"state,omitempty"`
}

const (
	// AccountPlayerKey 账号关联的玩家
	AccountPlayerKey = "account:player:%v"
	// NickNameField ...昵称
	NickNameField = "nick_name"
	// HeadImageField ...头像
	HeadImageField = "head_image"
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

// FmtPlayerKey 返回玩家的 key
func FmtPlayerKey(playerID uint64) string {
	return fmt.Sprintf("player:%v", playerID)
}
