package match

import "time"

// DeskPlayer 匹配中的牌桌玩家
type DeskPlayer struct {
	playerID uint64 // 玩家 ID
	robotLv  int    // 机器人等级, 为0时表示普通玩家
}

// GetPlayerID 获取 PlayerID
func (dp *DeskPlayer) GetPlayerID() uint64 {
	return dp.playerID
}

// GetRobotLv 获取机器人等级
func (dp *DeskPlayer) GetRobotLv() int {
	return dp.robotLv
}

// Desk 匹配中的牌桌
type Desk struct {
	ID                uint64       // ID
	gameID            int          // 游戏ID
	needPlayer        int          // 所需玩家人数
	players           []DeskPlayer // 玩家列表
	finish            func(*Desk)  // 匹配完成回调
	lastAddPlayerTime time.Time    // 最近一次添加玩家的时间
}

// CreateMatchDesk 创建匹配牌桌
func CreateMatchDesk(ID uint64, gameID int, needPlayer int, finish func(*Desk)) *Desk {
	return &Desk{
		ID:         ID,
		gameID:     gameID,
		needPlayer: needPlayer,
		finish:     finish,
		players:    make([]DeskPlayer, 0, needPlayer),
	}
}

// AddPlayer 向牌桌中加入玩家
func (d *Desk) AddPlayer(playerID uint64, robotLv int) {
	d.players = append(d.players, DeskPlayer{playerID: playerID, robotLv: robotLv})
	if len(d.players) == d.needPlayer {
		d.finish(d)
		return
	}
	d.lastAddPlayerTime = time.Now()
}

// GetLastAddPlayerTime 获取最近一次添加玩家的时间
func (d *Desk) GetLastAddPlayerTime() time.Time {
	return d.lastAddPlayerTime
}

// GetID 获取 ID
func (d *Desk) GetID() uint64 {
	return d.ID
}

// GetGameID 获取游戏 ID
func (d *Desk) GetGameID() int {
	return d.gameID
}

// GetPlayers 获取牌桌玩家
func (d *Desk) GetPlayers() []DeskPlayer {
	return d.players
}
