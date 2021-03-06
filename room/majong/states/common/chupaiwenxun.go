package common

import (
	"errors"
	"steve/client_pb/msgid"
	"steve/client_pb/room"
	majongpb "steve/entity/majong"
	"steve/room/majong/global"
	"steve/room/majong/interfaces"
	"steve/room/majong/utils"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

// ChupaiwenxunState 出牌问询状态
type ChupaiwenxunState struct{}

// ProcessEvent 处理事件
func (s *ChupaiwenxunState) ProcessEvent(eventID majongpb.EventID, eventContext interface{}, flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
	switch eventID {
	case majongpb.EventID_event_hu_request,
		majongpb.EventID_event_gang_request,
		majongpb.EventID_event_peng_request,
		majongpb.EventID_event_chi_request,
		majongpb.EventID_event_qi_request:
		{
			return s.onActionRequestEvent(eventID, eventContext, flow)
		}
	}
	return majongpb.StateID_state_chupaiwenxun, global.ErrInvalidEvent
}

// OnEntry 进入状态
func (s *ChupaiwenxunState) OnEntry(flow interfaces.MajongFlow) {
	for _, player := range flow.GetMajongContext().GetPlayers() {
		player.HasSelected = false
	}
	s.notifyPossibleActions(flow)
}

// OnExit 退出状态 清除本状态数据
func (s *ChupaiwenxunState) OnExit(flow interfaces.MajongFlow) {
	s.clearActionRec(flow)
}

// notifyPossibleActions 通知出牌问询
func (s *ChupaiwenxunState) notifyPossibleActions(flow interfaces.MajongFlow) {
	mjContext := flow.GetMajongContext()
	players := mjContext.GetPlayers()
	card := mjContext.GetLastOutCard()

	for index, player := range players {
		actions := player.GetPossibleActions()
		if len(actions) == 0 {
			continue
		}
		ntf := room.RoomChupaiWenxunNtf{}
		ntf.Card = proto.Uint32(utils.ServerCard2Uint32(card))
		for _, action := range actions {
			switch action {
			case majongpb.Action_action_peng:
				{
					ntf.EnablePeng = proto.Bool(true)
				}
			case majongpb.Action_action_gang:
				{
					ntf.EnableMinggang = proto.Bool(true)
				}
			case majongpb.Action_action_hu:
				{
					ntf.EnableDianpao = proto.Bool(true)
				}
			case majongpb.Action_action_qi:
				{
					ntf.EnableQi = proto.Bool(true)
				}
			case majongpb.Action_action_chi:
				{
					ntf.ChiInfo = &room.RoomChiInfo{
						Cards:   player.GetEnbleChiCards(),
						ChiCard: proto.Uint32(utils.ServerCard2Uint32(mjContext.GetLastOutCard())),
					}
				}
			}
		}
		logrus.WithFields(logrus.Fields{
			"func_name":   "ChupaiwenxunState.notifyPossibleActions",
			"player_id":   player.GetPlayerId(),
			"player_seat": index,
			"actions":     actions,
		}).Debugln("发送问询通知")
		flow.PushMessages([]uint64{player.GetPlayerId()}, interfaces.ToClientMessage{
			MsgID: int(msgid.MsgID_ROOM_CHUPAIWENXUN_NTF),
			Msg:   &ntf,
		})
	}
}

// getMajongPlayer 获取玩家对象
func (s *ChupaiwenxunState) getMajongPlayer(playerID uint64, mjContext *majongpb.MajongContext) *majongpb.Player {
	return utils.GetMajongPlayer(playerID, mjContext)
}

// existAction 玩家是否存在对应的可选操作
func (s *ChupaiwenxunState) existAction(action majongpb.Action, player *majongpb.Player) bool {
	return utils.ExistPossibleAction(player, action)
}

// getPengRequestPlayer 获取碰请求的玩家
func (s *ChupaiwenxunState) getPengRequestPlayer(eventContext interface{}) (uint64, error) {
	pengRequest := eventContext.(*majongpb.PengRequestEvent)
	return pengRequest.GetHead().GetPlayerId(), nil
}

// getGangRequestPlayer 获取杠请求的玩家
func (s *ChupaiwenxunState) getGangRequestPlayer(eventContext interface{}) (uint64, error) {
	gangRequest := eventContext.(*majongpb.GangRequestEvent)
	return gangRequest.GetHead().GetPlayerId(), nil
}

// getHuRequestPlayer 获取胡请求的玩家
func (s *ChupaiwenxunState) getHuRequestPlayer(eventContext interface{}) (uint64, error) {
	huRequest := eventContext.(*majongpb.HuRequestEvent)
	return huRequest.GetHead().GetPlayerId(), nil
}

// getChiRequestPlayer 获取吃请求的玩家
func (s *ChupaiwenxunState) getChiRequestPlayer(eventContext interface{}) (uint64, error) {
	chiRequest := eventContext.(*majongpb.ChiRequestEvent)
	return chiRequest.GetHead().GetPlayerId(), nil
}

func (s *ChupaiwenxunState) getChiRequestInfo(eventContext interface{}, player *majongpb.Player) error {
	chiRequest := eventContext.(*majongpb.ChiRequestEvent)
	reqChiCards := chiRequest.GetCards()
	player.DesignChiCards = reqChiCards
	return nil
}

// getQiRequestPlayer 获取弃请求的玩家
func (s *ChupaiwenxunState) getQiRequestPlayer(eventContext interface{}) (uint64, error) {
	qiRequest := eventContext.(*majongpb.QiRequestEvent)
	return qiRequest.GetHead().GetPlayerId(), nil
}

// getRequestInfo 根据请求事件获取请求的基础信息
func (s *ChupaiwenxunState) getRequestInfo(eventID majongpb.EventID, eventContext interface{}, mjContext *majongpb.MajongContext) (
	player *majongpb.Player, action majongpb.Action, err error) {
	// 从 map 中查找对应的 action
	action, ok := map[majongpb.EventID]majongpb.Action{
		majongpb.EventID_event_peng_request: majongpb.Action_action_peng,
		majongpb.EventID_event_gang_request: majongpb.Action_action_gang,
		majongpb.EventID_event_hu_request:   majongpb.Action_action_hu,
		majongpb.EventID_event_qi_request:   majongpb.Action_action_qi,
		majongpb.EventID_event_chi_request:  majongpb.Action_action_chi,
	}[eventID]
	if !ok {
		err = global.ErrInvalidEvent
		return
	}

	// 从 map 中查找和调用对应的方法
	type getPlayerFunc func(eventContext interface{}) (uint64, error)
	f, ok := map[majongpb.EventID]getPlayerFunc{
		majongpb.EventID_event_peng_request: s.getPengRequestPlayer,
		majongpb.EventID_event_gang_request: s.getGangRequestPlayer,
		majongpb.EventID_event_hu_request:   s.getHuRequestPlayer,
		majongpb.EventID_event_qi_request:   s.getQiRequestPlayer,
		majongpb.EventID_event_chi_request:  s.getChiRequestPlayer,
	}[eventID]
	if !ok {
		err = global.ErrInvalidEvent
		return
	}
	playerID, err := f(eventContext)
	if err != nil {
		return
	}
	player = s.getMajongPlayer(playerID, mjContext)
	if player == nil {
		err = global.ErrInvalidRequestPlayer
		return
	}
	if eventID == majongpb.EventID_event_chi_request {
		err = s.getChiRequestInfo(eventContext, player)
		if err != nil {
			return
		}
	}
	return
}

// canPlayerAction 检测玩家是否可以执行指定行为
func (s *ChupaiwenxunState) canPlayerAction(player *majongpb.Player, action majongpb.Action) error {
	if !s.existAction(action, player) {
		err := errors.New("当前玩家不能执行该操作")
		return err
	}
	if player.GetHasSelected() {
		err := errors.New("玩家已经选择过了")
		return err
	}
	return nil
}

// onActionRequestEvent 处理玩家 action 请求事件
func (s *ChupaiwenxunState) onActionRequestEvent(eventID majongpb.EventID, eventContext interface{}, flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
	newState, err = majongpb.StateID_state_chupaiwenxun, nil
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "ChupaiwenxunState.onActionRequestEvent",
		"event_id":  eventID,
	})

	mjContext := flow.GetMajongContext()
	logEntry = utils.WithMajongContext(logEntry, mjContext)

	player, action, err := s.getRequestInfo(eventID, eventContext, mjContext)
	if err != nil {
		logEntry.WithError(err).Infoln("获取请求信息失败")
		return
	}
	logEntry = logEntry.WithField("player_id", player.GetPlayerId())
	if err = s.canPlayerAction(player, action); err != nil {
		logEntry.WithError(err).Infoln("玩家不能执行该行为")
		return
	}

	player.HasSelected, player.SelectedAction = true, action
	return s.makeDecision(flow)
}

func (s *ChupaiwenxunState) getActionPriority(action majongpb.Action) int {
	// priorityMap 行为的优先级， 数字越大代表优先级越高
	var priorityMap = map[majongpb.Action]int{
		majongpb.Action_action_hu:   100,
		majongpb.Action_action_gang: 90,
		majongpb.Action_action_peng: 80,
		majongpb.Action_action_chi:  70,
	}
	if p, ok := priorityMap[action]; ok {
		return p
	}
	return 0
}

// getMaxSelectedAction 获取选择的最高优先级 action， 以及选择的玩家列表
func (s *ChupaiwenxunState) getMaxSelectedAction(players []*majongpb.Player) (bool, majongpb.Action, []uint64) {
	hasMaxSelectedAction := false
	var maxSelectedAction majongpb.Action
	selectedPlayers := []uint64{}

	for _, player := range players {
		if !player.GetHasSelected() {
			continue
		}
		selectedAction := player.GetSelectedAction()
		selectedPriority := s.getActionPriority(selectedAction)

		if !hasMaxSelectedAction || selectedPriority > s.getActionPriority(maxSelectedAction) {
			hasMaxSelectedAction = true
			maxSelectedAction = selectedAction
			selectedPlayers = []uint64{player.GetPlayerId()}
			continue
		}
		if selectedPriority == s.getActionPriority(maxSelectedAction) {
			selectedPlayers = append(selectedPlayers, player.GetPlayerId())
		}
	}
	return hasMaxSelectedAction, maxSelectedAction, selectedPlayers
}

// getMaxNotSelectedAction 获取未选择的最高优先级 action
func (s *ChupaiwenxunState) getMaxNotSelectedAction(players []*majongpb.Player) (bool, majongpb.Action) {
	has := false
	var maxAction majongpb.Action

	for _, player := range players {
		possibles := player.GetPossibleActions()
		if len(possibles) == 0 || player.GetHasSelected() {
			continue
		}
		for _, a := range possibles {
			if !has || s.getActionPriority(a) > s.getActionPriority(maxAction) {
				has = true
				maxAction = a
			}
		}
	}
	return has, maxAction
}

// makeDecision 作决策
// step 1. 在没有选择的玩家中，找到他们之中能执行的最高优先级的动作， 称为动作 A
// step 2. 在已经选择的玩家中，查找到他们之中选择的最高优先级的动作， 称为动作 B，以及选择执行这个动作的玩家列表， 称为玩家列表 L
// step 3. 如果动作 A 的优先级不低于动作 B， 返回出牌问询状态。 否则对所有在 L 中的玩家执行 B， 并且返回对应的状态
func (s *ChupaiwenxunState) makeDecision(flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "ChupaiwenxunState.makeDecision",
	})
	newState, err = majongpb.StateID_state_chupaiwenxun, nil

	mjContext := flow.GetMajongContext()
	logEntry = utils.WithMajongContext(logEntry, mjContext)
	players := mjContext.GetPlayers()

	hasMaxSelected, maxSelected, maxSelPlayers := s.getMaxSelectedAction(players)
	hasMaxNotSelected, maxNotSelected := s.getMaxNotSelectedAction(players)

	if !hasMaxSelected && !hasMaxNotSelected {
		logEntry.WithField("players", players).Errorln("没有问询但是进入了问询状态")
		return
	}
	if !hasMaxSelected && hasMaxNotSelected {
		return
	}
	if hasMaxSelected && !hasMaxNotSelected {
		return s.doAction(flow, maxSelected, maxSelPlayers)
	}
	if hasMaxSelected && hasMaxNotSelected {
		if s.getActionPriority(maxSelected) > s.getActionPriority(maxNotSelected) {
			return s.doAction(flow, maxSelected, maxSelPlayers)
		}
	}
	return
}

// doAction 执行特定行为
func (s *ChupaiwenxunState) doAction(flow interfaces.MajongFlow, action majongpb.Action, playerIDs []uint64) (newState majongpb.StateID, err error) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name":      "ChupaiwenxunState.doAction",
		"action":         action,
		"action_players": playerIDs,
	})
	newState, err = majongpb.StateID_state_chupaiwenxun, nil

	type actionFunc func(interfaces.MajongFlow, []uint64) (newState majongpb.StateID, err error)

	f, ok := map[majongpb.Action]actionFunc{
		majongpb.Action_action_peng: s.doPeng,
		majongpb.Action_action_gang: s.doGang,
		majongpb.Action_action_hu:   s.doHu,
		majongpb.Action_action_qi:   s.doQi,
		majongpb.Action_action_chi:  s.doChi,
	}[action]

	if !ok {
		err = errors.New("不支持的 action")
		logEntry.Errorln(err)
		return
	}
	return f(flow, playerIDs)
}

// doPeng 执行碰操作
func (s *ChupaiwenxunState) doPeng(flow interfaces.MajongFlow, playerIDs []uint64) (newState majongpb.StateID, err error) {
	newState, err = majongpb.StateID_state_peng, nil

	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "ChupaiwenxunState.doPeng",
	})
	mjContext := flow.GetMajongContext()
	logEntry = utils.WithMajongContext(logEntry, mjContext)

	if len(playerIDs) != 1 {
		err := errors.New("执行碰的玩家数不为 1")
		logEntry.Errorln(err)
		return majongpb.StateID_state_chupaiwenxun, err
	}
	playerID := playerIDs[0]

	mjContext.LastPengPlayer = playerID
	return
}

// doGang 执行杠操作
func (s *ChupaiwenxunState) doGang(flow interfaces.MajongFlow, playerIDs []uint64) (newState majongpb.StateID, err error) {
	newState, err = majongpb.StateID_state_gang, nil

	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "ChupaiwenxunState.doGang",
	})
	mjContext := flow.GetMajongContext()
	logEntry = utils.WithMajongContext(logEntry, mjContext)

	if len(playerIDs) != 1 {
		err := errors.New("执行杠的玩家数不为 1")
		logEntry.Errorln(err)
		return majongpb.StateID_state_chupaiwenxun, err
	}
	playerID := playerIDs[0]

	card := mjContext.GetLastOutCard()
	mjContext.GangCard = card
	mjContext.LastGangPlayer = playerID
	return
}

// doHu 执行胡牌操作
func (s *ChupaiwenxunState) doHu(flow interfaces.MajongFlow, playerIDs []uint64) (newState majongpb.StateID, err error) {
	newState, err = majongpb.StateID_state_hu, nil

	logEntry := logrus.WithFields(logrus.Fields{
		"func_name":  "ChupaiwenxunState.doHu",
		"hu_players": playerIDs,
	})
	mjContext := flow.GetMajongContext()
	logEntry = utils.WithMajongContext(logEntry, mjContext)

	mjContext.LastHuPlayers = playerIDs
	return
}

// doQi 执行弃操作。 切换到下家摸牌状态
func (s *ChupaiwenxunState) doQi(flow interfaces.MajongFlow, playerIDs []uint64) (newState majongpb.StateID, err error) {
	newState, err = majongpb.StateID_state_mopai, nil

	logEntry := logrus.WithFields(logrus.Fields{
		"func_name":  "ChupaiwenxunState.doHu",
		"hu_players": playerIDs,
	})
	mjContext := flow.GetMajongContext()
	logEntry = utils.WithMajongContext(logEntry, mjContext)
	lastOutCardPlayer := mjContext.GetLastChupaiPlayer()

	players := mjContext.GetPlayers()
	for _, player := range players {
		if player.GetPlayerId() == lastOutCardPlayer {
			nextPlayer := utils.GetNextXpPlayerByID(lastOutCardPlayer, players, flow.GetMajongContext())
			mjContext.MopaiPlayer = nextPlayer.GetPlayerId()
			mjContext.MopaiType = majongpb.MopaiType_MT_NORMAL
			return
		}
	}
	err = errors.New("出牌玩家不存在")
	logEntry.Errorln(err)
	return
}

// doQi 执行吃操作。 切换到吃状态
func (s *ChupaiwenxunState) doChi(flow interfaces.MajongFlow, playerIDs []uint64) (newState majongpb.StateID, err error) {
	newState, err = majongpb.StateID_state_chi, nil

	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "ChupaiwenxunState.doChi",
	})
	mjContext := flow.GetMajongContext()
	logEntry = utils.WithMajongContext(logEntry, mjContext)

	if len(playerIDs) != 1 {
		err := errors.New("执行吃的玩家数不为 1")
		logEntry.Errorln(err)
		return majongpb.StateID_state_chupaiwenxun, err
	}
	playerID := playerIDs[0]

	mjContext.LastChiPlayer = playerID
	return
}

func (s *ChupaiwenxunState) clearActionRec(flow interfaces.MajongFlow) {
	mjContext := flow.GetMajongContext()
	for _, player := range mjContext.GetPlayers() {
		player.PossibleActions = []majongpb.Action{}
		player.HasSelected = false
		player.EnbleChiCards = []uint32{}
		player.SelectedAction = majongpb.Action(-1)
	}
}
