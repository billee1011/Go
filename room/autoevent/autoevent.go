package autoevent

import (
	"steve/room/interfaces"
	"steve/room/interfaces/global"
	"steve/server_pb/majong"
	"time"

	"github.com/spf13/viper"
	"steve/room/config"
	"steve/client_pb/room"
	"steve/server_pb/ddz"
)

type autoEventGenerator struct {
	commonAIs map[int](map[int32]interfaces.CommonAI)
}

func init() {
	global.SetDeskAutoEventGenerator(&autoEventGenerator{
		commonAIs: map[int](map[int32]interfaces.CommonAI){},
	})
}

// handlePlayerAI 处理玩家 AI
func (aeg *autoEventGenerator) handlePlayerAI(result *interfaces.AutoEventGenerateResult, AI interfaces.CommonAI,
	playerID uint64, mjContext *majong.MajongContext, aiType interfaces.AIType, robotLv int) {
	aiResult, err := AI.GenerateAIEvent(interfaces.AIEventGenerateParams{
		MajongContext: mjContext,
		PlayerID:      playerID,
		AIType:        aiType,
		RobotLv:       robotLv,
	})
	if err == nil {
		for _, aiEvent := range aiResult.Events {
			result.Events = append(result.Events, interfaces.Event{
				ID:        aiEvent.ID,
				Context:   aiEvent.Context,
				PlayerID:  playerID,
				EventType: interfaces.OverTimeEvent,
			})
		}
	}
}

// handlePlayerAI 处理玩家 AI
func (aeg *autoEventGenerator) handleDDZPlayerAI(result *interfaces.AutoEventGenerateResult, AI interfaces.CommonAI,
	playerID uint64, ddzContext *ddz.DDZContext, aiType interfaces.AIType, robotLv int) {
	aiResult, err := AI.GenerateAIEvent(interfaces.AIEventGenerateParams{
		DDZContext: ddzContext,
		PlayerID:      playerID,
		AIType:        aiType,
		RobotLv:       robotLv,
	})
	if err == nil {
		for _, aiEvent := range aiResult.Events {
			result.Events = append(result.Events, interfaces.Event{
				ID:        aiEvent.ID,
				Context:   aiEvent.Context,
				PlayerID:  playerID,
				EventType: interfaces.OverTimeEvent,
			})
		}
	}
}

// handleOverTime 处理超时
func (aeg *autoEventGenerator) handleOverTime(AI interfaces.CommonAI, stateTime time.Time, mjContext *majong.MajongContext) (
	finish bool, result interfaces.AutoEventGenerateResult) {

	finish, result = false, interfaces.AutoEventGenerateResult{
		Events: []interfaces.Event{},
	}
	duration := time.Second * time.Duration(viper.GetInt(config.XingPaiTimeOut))
	if duration == 0 || time.Now().Sub(stateTime) < duration {
		return
	}
	players := mjContext.GetPlayers()
	for _, player := range players {
		aeg.handlePlayerAI(&result, AI, player.GetPalyerId(), mjContext, interfaces.OverTimeAI, 0)
	}
	finish = true
	return
}

func (aeg *autoEventGenerator) handleDDZOverTime(AI interfaces.CommonAI, params *interfaces.AutoEventGenerateParams) (
	finish bool, result interfaces.AutoEventGenerateResult) {

	finish, result = false, interfaces.AutoEventGenerateResult{
		Events: []interfaces.Event{},
	}
	startTime := params.StartTime
	ddzContext := params.DDZContext
	duration := time.Second * time.Duration(ddzContext.Duration)
	if duration == 0 || time.Now().Sub(startTime) < duration {
		return
	}
	players := ddzContext.CountDownPlayers
	for _, player := range players {
		aeg.handleDDZPlayerAI(&result, AI, player, ddzContext, interfaces.OverTimeAI, 0)
	}
	finish = true
	return
}

// handleTuoGuan 执行所有玩家的托管
func (aeg *autoEventGenerator) handleTuoGuan(desk interfaces.Desk, AI interfaces.CommonAI, stateTime time.Time, mjContext *majong.MajongContext) interfaces.AutoEventGenerateResult {
	result := interfaces.AutoEventGenerateResult{
		Events: []interfaces.Event{},
	}
	tuoguanOprTime := 1 * time.Second
	if time.Now().Sub(stateTime) < tuoguanOprTime {
		return result
	}

	players := desk.GetDeskPlayers()
	for _, player := range players {
		playerID := player.GetPlayerID()
		if player.IsTuoguan() {
			aeg.handlePlayerAI(&result, AI, playerID, mjContext, interfaces.TuoGuangAI, 0)
		}
	}
	return result
}

func (aeg *autoEventGenerator) handleDDZTuoGuan(desk interfaces.Desk, AI interfaces.CommonAI, stateTime time.Time, ddzContext *ddz.DDZContext) interfaces.AutoEventGenerateResult {
	result := interfaces.AutoEventGenerateResult{
		Events: []interfaces.Event{},
	}
	tuoguanOprTime := 1 * time.Second
	if time.Now().Sub(stateTime) < tuoguanOprTime {
		return result
	}

	players := desk.GetDeskPlayers()
	for _, player := range players {
		playerID := player.GetPlayerID()
		if player.IsTuoguan() {
			aeg.handleDDZPlayerAI(&result, AI, playerID, ddzContext, interfaces.TuoGuangAI, 0)
		}
	}
	return result
}

// GenerateV2 利用 AI 生成自动事件
func (aeg *autoEventGenerator) GenerateV2(params *interfaces.AutoEventGenerateParams) (result interfaces.AutoEventGenerateResult) {
	desk := params.Desk
	gameID := desk.GetGameID()
	gameAIs, ok := aeg.commonAIs[int(gameID)]
	if !ok {
		//logrus.WithField("gameId", gameID).Debug("Can't find game AI")
		return
	}
	var state int32
	if gameID == int(room.GameId_GAMEID_DOUDIZHU) {
		state = int32(params.DDZContext.GetCurState())
	} else {
		state = int32(params.MajongContext.GetCurState())
	}
	AI, ok := gameAIs[int32(state)]
	if !ok {
		//logrus.WithField("gameId", gameID).WithField("state", state).Debug("Can't find state AI")
		return
	}

	if gameID == int(room.GameId_GAMEID_DOUDIZHU) {
		if overTime, result := aeg.handleDDZOverTime(AI, params); overTime {
			return result
		}
		result = aeg.handleDDZTuoGuan(params.Desk, AI, params.StartTime, params.DDZContext)
	} else {
		if overTime, result := aeg.handleOverTime(AI, params.StartTime, params.MajongContext); overTime {
			return result
		}
		result = aeg.handleTuoGuan(params.Desk, AI, params.StartTime, params.MajongContext)

		players := params.MajongContext.GetPlayers()
		for _, player := range players {
			playerID := player.GetPalyerId()
			if lv, exist := params.RobotLv[playerID]; exist {
				aeg.handlePlayerAI(&result, AI, player.GetPalyerId(), params.MajongContext, interfaces.RobotAI, lv)
			}
		}
	}
	return result
}

func (aeg *autoEventGenerator) RegisterAI(gameID int, stateID int32, AI interfaces.CommonAI) {
	//logrus.WithField("gameID", gameID).WithField("stateID", stateID).Debug("Register AI")
	if _, exist := aeg.commonAIs[gameID]; !exist {
		aeg.commonAIs[gameID] = make(map[int32]interfaces.CommonAI)
	}
	aeg.commonAIs[gameID][stateID] = AI
}
