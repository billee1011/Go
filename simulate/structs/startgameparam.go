package structs

import "steve/client_pb/room"

// StartGameParams 启动游戏的参数
type StartGameParams struct {
	Cards      [][]*room.Card // 从庄家位置开始算起，每个位置的固定卡牌
	WallCards  []*room.Card   // 发完牌之后剩下的墙牌
	HszDir     room.Direction // 换三张的方向
	BankerSeat int            // 庄家座号
	ServerAddr string         // 服务器地址
	ClientVer  string         // 客户端版本号

	HszCards     [][]*room.Card   // 从庄家的位置算起，用来换三张的牌
	DingqueColor []room.CardColor // 定缺花色。 从庄家位置算起

}
