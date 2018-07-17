package game

import (
	playerdata "steve/common/data/player"
)

type Player struct {
	playerID uint64
	seatID   uint32 // 座号
	coin     uint64 // 进牌桌金币数

	desk *Desk //当前桌

	isQuit      bool // 是否已经退出牌桌
	overtime    int  // 超时计数
	maxOvertime int  // 最大超时次数
	isTuoGuan   bool // 是否在托管中
}

func NewPlayer(playerID uint64, seatID uint32, desk *Desk) *Player {
	return &Player{
		playerID: playerID,
		seatID:   seatID,
		desk:     desk,
		coin:     playerdata.GetPlayerCoin(playerID),

		isQuit:      false,
		overtime:    0,
		maxOvertime: 3,
		isTuoGuan:   false,
	}
}
