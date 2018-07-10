package autoevent

import (
	"steve/room/interfaces"
	"steve/room/interfaces/global"
	"steve/server_pb/majong"
	"time"

	"github.com/spf13/viper"
	"steve/room/config"
)

type autoEventGenerator struct {
	majongAIs map[int](map[int32]interfaces.CommonAI)
}

func init() {
	global.SetDeskAutoEventGenerator(&autoEventGenerator{
		majongAIs: map[int](map[int32]interfaces.CommonAI){},
	})
}

// getAI 根据状态和游戏 ID 获取 AI 对象
func (aeg *autoEventGenerator) getAI(mjContext *majong.MajongContext) interfaces.CommonAI {
	gameID := mjContext.GetGameId()
	gameAIs, ok := aeg.majongAIs[int(gameID)]
	if !ok {
		return nil
	}
	state := mjContext.GetCurState()
	AI, ok := gameAIs[int32(state)]
	if !ok {
		return nil
	}
	return AI
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

// handleTuoGuan 执行所有玩家的托管
func (aeg *autoEventGenerator) handleTuoGuan(desk interfaces.Desk, AI interfaces.CommonAI, stateTime time.Time, mjContext *majong.MajongContext) interfaces.AutoEventGenerateResult {
	result := interfaces.AutoEventGenerateResult{
		Events: []interfaces.Event{},
	}
	tuoguanOprTime := 1 * time.Second
	if time.Now().Sub(stateTime) < tuoguanOprTime {
		return result
	}
	//players := mjContext.GetPlayers()
	players := desk.GetDeskPlayers()
	for _, player := range players {
		playerID := player.GetPlayerID()
		if player.IsTuoguan() {
			aeg.handlePlayerAI(&result, AI, playerID, mjContext, interfaces.TuoGuangAI, 0)
		}
	}
	return result
}

// GenerateV2 利用 AI 生成自动事件
func (aeg *autoEventGenerator) GenerateV2(params *interfaces.AutoEventGenerateParams) (result interfaces.AutoEventGenerateResult) {
	mjContext := params.MajongContext
	AI := aeg.getAI(mjContext)
	if AI == nil {
		return
	}
	if overTime, result := aeg.handleOverTime(AI, params.StateTime, mjContext); overTime {
		return result
	}
	result = aeg.handleTuoGuan(params.Desk, AI, params.StateTime, mjContext)

	players := mjContext.GetPlayers()
	for _, player := range players {
		playerID := player.GetPalyerId()
		if lv, exist := params.RobotLv[playerID]; exist {
			aeg.handlePlayerAI(&result, AI, player.GetPalyerId(), mjContext, interfaces.RobotAI, lv)
		}
	}
	return result
}

func (aeg *autoEventGenerator) RegisterAI(gameID int, stateID int32, AI interfaces.CommonAI) {
	if _, exist := aeg.majongAIs[gameID]; !exist {
		aeg.majongAIs[gameID] = make(map[int32]interfaces.CommonAI)
	}
	aeg.majongAIs[gameID][stateID] = AI
}
