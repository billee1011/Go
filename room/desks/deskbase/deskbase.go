package deskbase

import "steve/room/interfaces"

// DeskBase 房间基类
type DeskBase struct {
	interfaces.DeskPlayerMgr
	uid    uint64
	gameID int
}

// NewDeskBase 创建房间基类对象
func NewDeskBase(uid uint64, gameID int, players []uint64) *DeskBase {
	deskPlayerMgr := createDeskPlayerMgr()
	deskPlayerMgr.setPlayers(players)
	return &DeskBase{
		DeskPlayerMgr: deskPlayerMgr,
		uid:           uid,
		gameID:        gameID,
	}
}

// GetUID 获取牌桌 UID
func (d *DeskBase) GetUID() uint64 {
	return d.uid
}

// GetGameID 获取游戏 ID
func (d *DeskBase) GetGameID() int {
	return d.gameID
}
