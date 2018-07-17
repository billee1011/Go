package game

type PlayerManager struct {
	players map[uint64]*Player // playerID - *Player
}

var DefaultPlayManager = NewPlayerManager()

func NewPlayerManager() *PlayerManager {
	return new(PlayerManager)
}

func (pm *PlayerManager) AddPlayer(player *Player) {
	pm.players[player.playerID] = player
}

func (pm *PlayerManager) RemovePlayer(playerID uint64) {
	delete(pm.players, playerID)
}

func (pm *PlayerManager) GetPlayer(playerID uint64) *Player {
	return pm.players[playerID]
}
