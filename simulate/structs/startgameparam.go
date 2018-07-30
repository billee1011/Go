package structs

import (
	"steve/client_pb/common"
	"steve/client_pb/room"
)

// StartGameParams 启动游戏的参数
type StartGameParams struct {
	GameID     common.GameId  // 游戏ID
	Cards      [][]uint32     // 从庄家位置开始算起，每个位置的固定卡牌
	WallCards  []uint32       // 发完牌之后剩下的墙牌
	HszDir     room.Direction // 换三张的方向
	BankerSeat int            // 庄家座号

	IsHsz        bool             //是否换三张
	HszCards     [][]uint32       // 从庄家的位置算起，用来换三张的牌
	DingqueColor []room.CardColor // 定缺花色。 从庄家位置算起

	PlayerSeatGold map[int]uint64 //玩家座位对应的金币
	PeiPaiGame     string         // 配牌游戏名
	PlayerNum      int            //参与游戏的玩家人数
	IsDq           bool           //是否定缺
	DiFen          uint64         // 底分
}

// StartPukeGameParams 启动扑克游戏的参数
type StartPukeGameParams struct {
	GameID     common.GameId  // 游戏ID
	Cards      [][]uint32     // 从地主位置开始算起，每个位置的固定卡牌
	WallCards  []uint32       // 已废弃
	HszDir     room.Direction // 已废弃
	BankerSeat int            // 地主座位号
	ServerAddr string         // 服务器地址
	ClientVer  string         // 客户端版本号

	PlayerSeatGold map[int]uint64 //玩家座位对应的金币
	PeiPaiGame     string         // 配牌游戏名
}
