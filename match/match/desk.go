package match

// Desk 匹配中的牌桌
type Desk struct {
	gameID     int      // 游戏ID
	needPlayer int      // 所需玩家人数
	players    []uint64 // 玩家列表
	finish     func()   // 匹配完成回调
}

// AddPlayer 向牌桌中加入玩家
func (d *Desk) AddPlayer(playerID uint64) {
	d.players = append(d.players, playerID)
	if len(d.players) == d.needPlayer {
		d.finish()
	}
}
