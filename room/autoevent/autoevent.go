package autoevent

import (
	"steve/entity/majong"
	"steve/room/interfaces"
	"steve/room/interfaces/global"
	"time"

	"steve/client_pb/room"
	"steve/room/config"
	"steve/server_pb/ddz"

	"github.com/spf13/viper"
)

type autoEventGenerator struct {
	commonAIs map[int](map[int32]interfaces.CommonAI) // Key：游戏ID，对应枚举GameId，Value:该游戏各个状态对应的AI产生器
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
		eventType := interfaces.OverTimeEvent
		if aiType == interfaces.RobotAI {
			eventType = interfaces.RobotEvent
		} else if aiType == interfaces.TuoGuangAI {
			eventType = interfaces.TuoGuanEvent
		}
		aeg.addAIEvents(result, &aiResult, playerID, eventType)
	}
}

// addAIEvents 将 AI 产生的事件添加到结果中
func (aeg *autoEventGenerator) addAIEvents(result *interfaces.AutoEventGenerateResult, aiResult *interfaces.AIEventGenerateResult, playerID uint64, eventType interfaces.EventType) {
	for _, aiEvent := range aiResult.Events {
		result.Events = append(result.Events, interfaces.Event{
			ID:        aiEvent.ID,
			Context:   aiEvent.Context,
			PlayerID:  playerID,
			EventType: eventType,
		})
	}
}

// handlePlayerAI 处理玩家 AI
// result 		: 存放AI事件的结果
// AI			: 具体的AI产生器
// playerID 	: AI针对的玩家的playerID
// ddzContext 	: 牌局信息
// aiType		: 托管，超时等，对应枚举 AIType
// robotLv		: 机器人级别
func (aeg *autoEventGenerator) handleDDZPlayerAI(result *interfaces.AutoEventGenerateResult, AI interfaces.CommonAI,
	playerID uint64, ddzContext *ddz.DDZContext, aiType interfaces.AIType, robotLv int) {

	// 由该AI产生具体的AI事件
	aiResult, err := AI.GenerateAIEvent(interfaces.AIEventGenerateParams{
		DDZContext: ddzContext,
		PlayerID:   playerID,
		AIType:     aiType,
		RobotLv:    robotLv,
	})

	// 未出错时，把产生的每一个AI事件压入结果集
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

// handleDDZOverTime 斗地主超时处理
// finish : 是否处理完成
// result : 产生的AI事件结果集合
func (aeg *autoEventGenerator) handleDDZOverTime(AI interfaces.CommonAI, params *interfaces.AutoEventGenerateParams) (
	finish bool, result interfaces.AutoEventGenerateResult) {

	finish, result = false, interfaces.AutoEventGenerateResult{
		Events: []interfaces.Event{},
	}

	// 开始时间
	startTime := params.StartTime

	ddzContext := params.DDZContext

	// 倒计时的时长
	duration := time.Second * time.Duration(ddzContext.Duration)

	// 未到倒计时，不处理
	if duration == 0 || time.Now().Sub(startTime) < duration {
		return
	}

	// 处理每一个处于倒计时的玩家，产生具体的AI事件，并把事件存入result
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

// handleDDZTuoGuan 斗地主托管处理
// finish : 是否处理完成
// result : 产生的AI事件结果集合
func (aeg *autoEventGenerator) handleDDZTuoGuan(desk interfaces.Desk, AI interfaces.CommonAI, stateTime time.Time, ddzContext *ddz.DDZContext) interfaces.AutoEventGenerateResult {
	result := interfaces.AutoEventGenerateResult{
		Events: []interfaces.Event{},
	}

	// 托管时的操作等待时间
	tuoguanOprTime := 2 * time.Second

	if time.Now().Sub(stateTime) < tuoguanOprTime {
		return result
	}

	// 遍历桌子的所有玩家
	players := desk.GetDeskPlayers()
	for _, player := range players {
		playerID := player.GetPlayerID()

		// 若处于托管中，则产生具体的AI事件，并把事件存入result
		if player.IsTuoguan() {
			aeg.handleDDZPlayerAI(&result, AI, playerID, ddzContext, interfaces.TuoGuangAI, 0)
		}
	}
	return result
}

// GenerateV2 利用 AI 生成自动事件
func (aeg *autoEventGenerator) GenerateV2(params *interfaces.AutoEventGenerateParams) (result interfaces.AutoEventGenerateResult) {
	desk := params.Desk

	// 该桌子所属的游戏ID
	gameID := desk.GetGameID()

	// 该游戏各个状态对应的AI产生器
	gameAIs, ok := aeg.commonAIs[int(gameID)]
	if !ok {
		//logrus.WithField("gameId", gameID).Debug("Can't find game AI")
		return
	}

	// 当前状态ID
	var state int32
	if gameID == int(room.GameId_GAMEID_DOUDIZHU) {
		state = int32(params.DDZContext.GetCurState())
	} else {
		state = int32(params.MajongContext.GetCurState())
	}

	// 当前状态的AI产生器
	AI, ok := gameAIs[int32(state)]
	if !ok {
		//logrus.WithField("gameId", gameID).WithField("state", state).Debug("Can't find state AI")
		return
	}

	// 斗地主的特殊处理
	if gameID == int(room.GameId_GAMEID_DOUDIZHU) {

		// 先处理超时
		if overTime, result := aeg.handleDDZOverTime(AI, params); overTime {
			return result
		}

		// 未超时，则处理托管
		result = aeg.handleDDZTuoGuan(params.Desk, AI, params.StartTime, params.DDZContext)
	} else {
		if overTime, result := aeg.handleOverTime(AI, params.StartTime, params.MajongContext); overTime {
			return result
		}
		result = aeg.handleTuoGuan(params.Desk, AI, params.StartTime, params.MajongContext)
		// 超过 1s 处理机器人事件
		if time.Now().Sub(params.StartTime) > 1*time.Second {
			players := params.MajongContext.GetPlayers()
			for _, player := range players {
				playerID := player.GetPalyerId()
				if lv, exist := params.RobotLv[playerID]; exist && lv != 0 {
					aeg.handlePlayerAI(&result, AI, player.GetPalyerId(), params.MajongContext, interfaces.RobotAI, lv)
				}
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
