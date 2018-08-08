package models

import (
	"steve/client_pb/msgid"
	"steve/client_pb/room"
	"steve/room/contexts"
	"steve/room/desk"
	"steve/room/fixed"
	playerpkg "steve/room/player"
	"sync"
)

type playerIDWithChannel struct {
	playerID      uint64
	finishChannel chan error
}

// PlayerModel ...
type PlayerModel struct {
	BaseModel
	players      []*playerpkg.Player
	enterChannel chan playerIDWithChannel
	leaveChannel chan playerIDWithChannel
	stopped      bool
	mu           sync.Mutex
}

// GetName get model name
func (model *PlayerModel) GetName() string {
	return fixed.PlayerModelName
}

// Active 激活 model
func (model *PlayerModel) Active() {
	model.players = make([]*playerpkg.Player, model.GetDesk().GetConfig().Num)
	ids := model.GetDesk().GetConfig().PlayerIds //GetModelManager().GetPlayerModel(model.GetDesk().GetUid()).GetDeskPlayerIDs()
	for i := 0; i < len(model.players); i++ {
		playerObj := playerpkg.GetPlayerMgr().GetPlayer(ids[i])
		playerObj.SetDesk(model.GetDesk())
		playerObj.SetQuit(false)
		playerObj.SetTuoguan(false, false)
		playerObj.SetEcoin(playerObj.GetCoin())

		model.players[i] = playerObj
	}
}

// Start start model
func (model *PlayerModel) Start() {
	model.stopped = false
}

// Stop stop model
func (model *PlayerModel) Stop() {
	model.mu.Lock()
	playerMgr := playerpkg.GetPlayerMgr()

	playerIDs := make([]uint64, 0, len(model.players))
	for _, pla := range model.players {
		if pla == nil {
			continue
		}
		if pla.GetDesk() == model.GetDesk() {
			pla.SetDesk(nil)
			playerIDs = append(playerIDs, pla.GetPlayerID())
		}
	}
	playerMgr.UnbindPlayerRoomAddr(playerIDs)

	// 消费完所有进入退出
forstart:
	for {
		select {
		case enterInfo := <-model.enterChannel:
			{
				close(enterInfo.finishChannel)
			}
		case leaveInfo := <-model.leaveChannel:
			{
				close(leaveInfo.finishChannel)
			}
		default:
			break forstart
		}
	}
	model.stopped = true
	model.mu.Unlock()
}

// getEnterChannel get enter channel
func (model *PlayerModel) getEnterChannel() chan playerIDWithChannel {
	return model.enterChannel
}

// getLeaveChannel get leave channel
func (model *PlayerModel) getLeaveChannel() chan playerIDWithChannel {
	return model.leaveChannel
}

// NewPlayertModel create player model
func NewPlayertModel(desk *desk.Desk) DeskModel {
	result := &PlayerModel{
		enterChannel: make(chan playerIDWithChannel, 4),
		leaveChannel: make(chan playerIDWithChannel, 4),
	}
	result.SetDesk(desk)
	return result
}

// PlayerEnter 玩家进入
func (model *PlayerModel) PlayerEnter(player *playerpkg.Player) {
	model.mu.Lock()
	if model.stopped {
		model.mu.Unlock()
		return
	}
	model.enterChannel <- playerIDWithChannel{
		playerID:      player.GetPlayerID(),
		finishChannel: make(chan error, 0),
	}
	model.mu.Unlock()
}

// handlePlayerEnter 处理玩家重入
func (model *PlayerModel) handlePlayerEnter(playerID uint64) {
	playerMgr := playerpkg.GetPlayerMgr()
	player := playerMgr.GetPlayer(playerID)
	player.SetQuit(false)
	player.SetTuoguan(false, false)
	model.playerQuitEnterDeskNtf(player, room.QuitEnterType_QET_ENTER)
}

// PlayerQuit 玩家退出
func (model *PlayerModel) PlayerQuit(player *playerpkg.Player) {
	model.mu.Lock()
	if model.stopped {
		model.mu.Unlock()
		return
	}
	model.leaveChannel <- playerIDWithChannel{
		playerID:      player.GetPlayerID(),
		finishChannel: make(chan error, 0),
	}
	model.mu.Unlock()
}

// handlePlayerLeave 处理玩家离开牌桌
func (model *PlayerModel) handlePlayerLeave(playerID uint64, needTuoguan bool) {
	playerMgr := playerpkg.GetPlayerMgr()
	player := playerMgr.GetPlayer(playerID)
	player.SetQuit(true)
	if !player.IsTuoguan() && needTuoguan {
		player.SetTuoguan(true, false)
	}
	model.playerQuitEnterDeskNtf(player, room.QuitEnterType_QET_QUIT)
}

func (model *PlayerModel) playerQuitEnterDeskNtf(player *playerpkg.Player, qeType room.QuitEnterType) {
	if player == nil {
		return
	}
	roomPlayer := TranslateToRoomPlayer(player)
	playerID := player.GetPlayerID()
	ntf := room.RoomDeskQuitEnterNtf{
		PlayerId:   &playerID,
		Type:       &qeType,
		PlayerInfo: &roomPlayer,
	}
	messageModel := GetModelManager().GetMessageModel(model.GetDesk().GetUid())
	messageModel.BroadCastDeskMessageExcept([]uint64{playerID}, true, msgid.MsgID_ROOM_DESK_QUIT_ENTER_NTF, &ntf)
}

func (model *PlayerModel) setContextPlayerQuit(player *playerpkg.Player, value bool) {
	for _, p := range model.GetDesk().GetConfig().Context.(*contexts.MajongDeskContext).MjContext.Players {
		if p.GetPalyerId() == player.GetPlayerID() {
			p.IsQuit = value
		}
	}
}

func (model *PlayerModel) GetDeskPlayers() []*playerpkg.Player {
	return model.players
}

func (model *PlayerModel) GetDeskPlayerByID(playerID uint64) *playerpkg.Player {
	for _, player := range model.players {
		if player.PlayerID == playerID {
			return player
		}
	}
	return nil
}

// GetDeskPlayerIDs 获取牌桌玩家 ID 列表， 座号作为索引
func (model *PlayerModel) GetDeskPlayerIDs() []uint64 {
	players := model.GetDeskPlayers()
	result := make([]uint64, len(players))
	for _, player := range players {
		result[player.GetSeat()] = player.GetPlayerID()
	}
	return result
}

// GetTuoguanPlayers 获取牌桌所有托管玩家
func (model *PlayerModel) GetTuoguanPlayers() []uint64 {
	players := GetModelManager().GetPlayerModel(model.GetDesk().GetUid()).GetDeskPlayers()
	result := make([]uint64, 0, len(players))
	for _, pla := range players {
		if pla.IsTuoguan() {
			result = append(result, pla.GetPlayerID())
		}
	}
	return result
}
