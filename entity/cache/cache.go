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
)

// FmtAccountPlayerKey 账号所关联玩家 key
func FmtAccountPlayerKey(accountID uint64) string {
	return fmt.Sprintf(AccountPlayerKey, accountID)
}
