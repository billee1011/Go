package playermgr

import (
	playerdata "steve/common/data/player"
)

type player struct {
	playerID uint64
}

func (p *player) GetID() uint64 {
	return p.playerID
}

func (p *player) GetCoin() uint64 {
	return playerdata.GetPlayerCoin(p.playerID)
}

func (p *player) SetCoin(coin uint64) {
	playerdata.SetPlayerCoin(p.playerID, coin)
}

func (p *player) IsOnline() bool {
	return playerdata.GetPlayerGateAddr(p.playerID) != ""
}
