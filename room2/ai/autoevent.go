package ai

import (
	"steve/entity/majong"
	"steve/room/config"
	"steve/room2/contexts"
	"steve/room2/desk"
	"steve/room2/fixed"
	playerPkg "steve/room2/player"
	"time"

	"github.com/spf13/viper"
)

// AutoEventGenerateParams 生成自动事件的参数
type AutoEventGenerateParams struct {
	MajongContext  *majong.MajongContext
	CurTime        time.Time
	StateTime      time.Time
	RobotLv        map[uint64]int
	TuoGuanPlayers []uint64
}

// AutoEventGenerateResult 自动事件生成结果
type AutoEventGenerateResult struct {
	Events []desk.DeskEvent
}

// DeskAutoEventGenerator 牌桌自动事件产生器
type DeskAutoEventGenerator interface {
	GenerateV2(params *AutoEventGenerateParams) AutoEventGenerateResult
	RegisterAI(gameID int, stateID majong.StateID, AI MajongAI) // 注册 AI
}

type AutoEventGenerator struct {
	majongAIs map[int](map[majong.StateID]MajongAI)
}

var atEvent *AutoEventGenerator

func init() {
	atEvent = &AutoEventGenerator{
		majongAIs: map[int](map[majong.StateID]MajongAI){},
	}
}

func GetAtEvent() *AutoEventGenerator {
	return atEvent
}

// getAI 根据状态和游戏 ID 获取 AI 对象
func (aeg *AutoEventGenerator) getAI(mjContext *majong.MajongContext) MajongAI {
	gameID := mjContext.GetGameId()
	gameAIs, ok := aeg.majongAIs[int(gameID)]
	if !ok {
		return nil
	}
	state := mjContext.GetCurState()
	AI, ok := gameAIs[state]
	if !ok {
		return nil
	}
	return AI
}

// getStateDuration 获取状态超时时间，通过config配置，随进程持续
func (aeg *AutoEventGenerator) getStateDuration() time.Duration {
	return time.Second * time.Duration(viper.GetInt(config.XingPaiTimeOut))
}

// addAIEvents 将 AI 产生的事件添加到结果中
func (aeg *AutoEventGenerator) addAIEvents(result *AutoEventGenerateResult, aiResult *AIEventGenerateResult, player *playerPkg.Player, eventType int) {
	for _, aiEvent := range aiResult.Events {
		event := desk.NewDeskEvent(int(aiEvent.ID), eventType, player.GetDesk(), desk.CreateEventParams(
			player.GetDesk().GetConfig().Context.(*contexts.MjContext).StateNumber, aiEvent.Context, player.GetPlayerID(),
		))
		result.Events = append(result.Events, event)
	}
}

// handlePlayerAI 处理玩家 AI
func (aeg *AutoEventGenerator) handlePlayerAI(result *AutoEventGenerateResult, AI MajongAI,
	player *majong.Player, mjContext *majong.MajongContext, aiType AIType, robotLv int) {
	playerID := player.GetPalyerId()
	aiResult, err := AI.GenerateAIEvent(AIEventGenerateParams{
		MajongContext: mjContext,
		PlayerID:      playerID,
		AIType:        aiType,
		RobotLv:       robotLv,
	})
	if err == nil {
		eventType := fixed.OverTimeEvent
		if aiType == RobotAI {
			eventType = fixed.RobotEvent
		} else if aiType == TuoGuangAI {
			eventType = fixed.TuoGuanEvent
		}
		player := playerPkg.GetPlayerMgr().GetPlayer(playerID)
		aeg.addAIEvents(result, &aiResult, player, eventType)
	}
}

// handlePlayerTuoGuan 处理玩家托管
func (aeg *AutoEventGenerator) handlePlayerTuoGuan(result *AutoEventGenerateResult, AI MajongAI,
	player *majong.Player, mjContext *majong.MajongContext) {
	aeg.handlePlayerAI(result, AI, player, mjContext, TuoGuangAI, 0)
}

// handlePlayerOverTime 处理玩家超时
func (aeg *AutoEventGenerator) handlePlayerOverTime(result *AutoEventGenerateResult, AI MajongAI,
	player *majong.Player, mjContext *majong.MajongContext) {
	aeg.handlePlayerAI(result, AI, player, mjContext, OverTimeAI, 0)
}

// handleOverTime 处理超时
func (aeg *AutoEventGenerator) handleOverTime(AI MajongAI, curTime time.Time, stateTime time.Time, mjContext *majong.MajongContext) (
	finish bool, result AutoEventGenerateResult) {

	finish, result = false, AutoEventGenerateResult{
		Events: []desk.DeskEvent{},
	}
	duration := aeg.getStateDuration()
	if duration == 0 || curTime.Sub(stateTime) < duration {
		return
	}
	players := mjContext.GetPlayers()
	for _, player := range players {
		aeg.handlePlayerOverTime(&result, AI, player, mjContext)
	}
	finish = true
	return
}

// isTuoGuan 玩家是否托管
func (aeg *AutoEventGenerator) isTuoGuan(playerID uint64, tuoGuanPlayers []uint64) bool {
	for _, player := range tuoGuanPlayers {
		if playerID == player {
			return true
		}
	}
	return false
}

// handleTuoGuan 执行所有玩家的托管
func (aeg *AutoEventGenerator) handleTuoGuan(tuoGuanPlayers []uint64, AI MajongAI, curTime time.Time, stateTime time.Time, mjContext *majong.MajongContext) AutoEventGenerateResult {
	result := AutoEventGenerateResult{
		Events: []desk.DeskEvent{},
	}
	tuoguanOprTime := 1 * time.Second
	if curTime.Sub(stateTime) < tuoguanOprTime {
		return result
	}
	players := mjContext.GetPlayers()
	for _, player := range players {
		playerID := player.GetPalyerId()
		if aeg.isTuoGuan(playerID, tuoGuanPlayers) {
			aeg.handlePlayerTuoGuan(&result, AI, player, mjContext)
		}
	}
	return result
}

// GenerateV2 利用 AI 生成自动事件
func (aeg *AutoEventGenerator) GenerateV2(params *AutoEventGenerateParams) (result AutoEventGenerateResult) {
	mjContext := params.MajongContext
	AI := aeg.getAI(mjContext)
	if AI == nil {
		return
	}
	if overTime, result := aeg.handleOverTime(AI, params.CurTime, params.StateTime, mjContext); overTime {
		return result
	}
	result = aeg.handleTuoGuan(params.TuoGuanPlayers, AI, params.CurTime, params.StateTime, mjContext)

	// 超过 1s 处理机器人事件
	if params.CurTime.Sub(params.StateTime) > 1*time.Second {
		players := mjContext.GetPlayers()
		for _, player := range players {
			playerID := player.GetPalyerId()
			if lv, exist := params.RobotLv[playerID]; exist && lv != 0 {
				aeg.handlePlayerAI(&result, AI, player, mjContext, RobotAI, lv)
			}
		}
	}
	return result
}

func (aeg *AutoEventGenerator) RegisterAI(gameID int, stateID majong.StateID, AI MajongAI) {
	if _, exist := aeg.majongAIs[gameID]; !exist {
		aeg.majongAIs[gameID] = make(map[majong.StateID]MajongAI)
	}
	aeg.majongAIs[gameID][stateID] = AI
}
