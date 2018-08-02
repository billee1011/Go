package data

import (
	"steve/entity/cache"
	"strconv"
)

// PlayerInfo 玩家个人资料
type PlayerInfo struct {
	PlayerID   uint64
	NickName   string
	Gender     uint32
	Avator     string
	State      uint32
	GameID     uint32
	LevelID    uint32
	ChannelID  uint32
	ProvinceID uint32
	CityID     uint32
}

// PlayerState 玩家状态
type PlayerState struct {
	PlayerID  uint64
	State     uint32
	GameID    uint32
	IPAddr    string
	GateAddr  string
	MatchAddr string
	RoomAddr  string
}

func (p *PlayerInfo) generatePlayerInfo(info map[string]string) {
	// 性别
	gender, _ := strconv.ParseUint(info[cache.Gender], 10, 64)
	// 渠道ID
	channelID, _ := strconv.ParseUint(info[cache.ChannelID], 10, 64)
	// 省份ID
	provinceID, _ := strconv.ParseUint(info[cache.ProvinceID], 10, 64)
	// 城市ID
	cityID, _ := strconv.ParseUint(info[cache.CityID], 10, 64)

	p.NickName = info[cache.NickName]
	p.Avator = info[cache.Avatar]
	p.Gender = uint32(gender)
	p.ChannelID = uint32(channelID)
	p.ProvinceID = uint32(provinceID)
	p.CityID = uint32(cityID)
}

func (pState *PlayerState) generatePlayerState(info map[string]string) {
	// 游戏状态
	state, _ := strconv.ParseUint(info[cache.GameState], 10, 64)
	pState.State = uint32(state)
	// 游戏状态
	gameID, _ := strconv.ParseUint(info[cache.GameID], 10, 64)
	pState.GameID = uint32(gameID)
	// ip地址
	pState.IPAddr = info[cache.IPAddr]
	// 网关服地址
	pState.GateAddr = info[cache.GateAddr]
	// 匹配服地址
	pState.MatchAddr = info[cache.MatchAddr]
	// 房间服地址
	pState.RoomAddr = info[cache.RoomAddr]
}
