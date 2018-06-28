package common

import (
	"fmt"
	"math/rand"
	"steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/common/mjoption"
	"steve/majong/global"
	"steve/majong/interfaces"
	"steve/majong/utils"
	"steve/room/peipai/handle"
	majongpb "steve/server_pb/majong"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

// HuansanzhangState 换三张状态
type HuansanzhangState struct {
}

// OnEntry 进入换三张状态
func (s *HuansanzhangState) OnEntry(flow interfaces.MajongFlow) {
	// 客户端强烈要求不要这个通知
	// facade.BroadcaseMessage(flow, msgid.MsgID_ROOM_HUANSANZHANG_NTF, &room.RoomHuansanzhangNtf{})
}

// ProcessEvent 处理换三张事件
func (s *HuansanzhangState) ProcessEvent(eventID majongpb.EventID, eventContext []byte, flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
	switch eventID {
	case majongpb.EventID_event_huansanzhang_request:
		{
			return s.onReq(eventID, eventContext, flow)
		}
	case majongpb.EventID_event_huansanzhang_finish:
		{
			return s.nextState(flow), nil
		}
	case majongpb.EventID_event_cartoon_finish_request:
		{
			return s.onCartoonFinish(flow, eventContext)
		}
	}
	return majongpb.StateID_state_huansanzhang, nil
}

// OnExit 退出换三张状态
func (s *HuansanzhangState) OnExit(flow interfaces.MajongFlow) {

}

// nextState 下个状态
func (s *HuansanzhangState) nextState(flow interfaces.MajongFlow) majongpb.StateID {
	xpOption := mjoption.GetXingpaiOption(int(flow.GetMajongContext().GetXingpaiOptionId()))
	if xpOption.NeedDingque {
		return majongpb.StateID_state_dingque
	}
	return majongpb.StateID_state_zixun
}

// curState 当前状态
func (s *HuansanzhangState) curState() majongpb.StateID {
	return majongpb.StateID_state_huansanzhang
}

// onCartoonFinish 动画播放完毕
func (s *HuansanzhangState) onCartoonFinish(flow interfaces.MajongFlow, eventContext []byte) (newState majongpb.StateID, err error) {
	finished := flow.GetMajongContext().GetExcutedHuansanzhang()
	if !finished {
		return s.curState(), global.ErrInvalidEvent
	}
	return OnCartoonFinish(s.curState(), s.nextState(flow), room.CartoonType_CTNT_HUANSANZHANG, eventContext)
}

// checkReq 检测玩家请求是否合法
func (s *HuansanzhangState) checkReq(logEntry *logrus.Entry, player *majongpb.Player, cards []*majongpb.Card) bool {
	if len(cards) != 3 {
		logEntry.Infoln("请求牌数不为3")
		return false
	}
	for i := 0; i < len(cards)-1; i++ {
		if cards[i].Color != cards[i+1].Color {
			logEntry.Infoln("花色不一致")
			return false
		}
		if !utils.ContainCard(player.HandCards, cards[i]) {
			logEntry.WithField("card", cards[i]).Infoln("玩家手牌不存在该牌")
			return false
		}
	}
	// TODO : 判断重复情况
	return true
}

// onReq 处理换三张请求事件
func (s *HuansanzhangState) onReq(eventID majongpb.EventID, eventContext []byte, flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
	newState, err = majongpb.StateID_state_huansanzhang, nil

	logEntry := logrus.WithField("func_name", "HuansanzhangState.onReq")
	mjContext := flow.GetMajongContext()
	logEntry = utils.WithMajongContext(logEntry, mjContext)

	req := new(majongpb.HuansanzhangRequestEvent)

	if marshalErr := proto.Unmarshal(eventContext, req); marshalErr != nil {
		logEntry.WithError(marshalErr).Errorln(global.ErrUnmarshalEvent)
		return majongpb.StateID_state_huansanzhang, global.ErrUnmarshalEvent
	}
	playerID := req.GetHead().GetPlayerId()

	player := utils.GetPlayerByID(mjContext.GetPlayers(), playerID)
	reqCards := req.GetCards()

	logEntry = logEntry.WithFields(logrus.Fields{
		"req_player": playerID,
		"req_cards":  reqCards,
		"hand_cards": player.GetHandCards(),
	})
	if !s.checkReq(logEntry, player, reqCards) {
		return
	}
	player.HuansanzhangCards = reqCards
	if req.Sure {
		// 应答
		onHuanSanZhangRsq(playerID, flow)
		player.HuansanzhangSure = true
	}
	return s.execute(flow)
}

// checkDone 检测换三张是否可执行
func (s *HuansanzhangState) checkDone(players []*majongpb.Player) bool {
	for _, player := range players {
		if len(player.GetHuansanzhangCards()) != 3 || !player.GetHuansanzhangSure() {
			return false
		}
	}
	return true
}

// randDirection 随机换三张方向
func (s *HuansanzhangState) randDirection(flow interfaces.MajongFlow) room.Direction {
	rd := rand.New(rand.NewSource(time.Now().Unix())) // 生成换三张方向
	directios := []room.Direction{room.Direction_ClockWise, room.Direction_Opposite, room.Direction_AntiClockWise}
	towards := rd.Intn(len(directios))

	mjContext := flow.GetMajongContext()
	fx := handle.GetHSZFangXiang(int(mjContext.GetGameId()))
	if fx >= 0 && fx < len(directios) {
		towards = fx
	}
	return directios[towards]
}

// makePairs 生成换三张配对数据。
// 返回： 座号->从哪个座号拿牌
func (s *HuansanzhangState) makePairs(playerCount int, dir room.Direction) map[int]int {
	switch dir {
	case room.Direction_AntiClockWise:
		{
			return map[int]int{0: 3, 1: 0, 2: 1, 3: 2}
		}
	case room.Direction_ClockWise:
		{
			return map[int]int{0: 1, 1: 2, 2: 3, 3: 0}
		}
	case room.Direction_Opposite:
		{
			return map[int]int{0: 2, 1: 3, 2: 0, 3: 1}
		}
	default:
		return map[int]int{}
	}
}

// fetchCardFrom seat 从 from 获取换三张的牌
func (s *HuansanzhangState) fetchCardFrom(flow interfaces.MajongFlow, seat int, from int) ([]*majongpb.Card, error) {
	mjContext := flow.GetMajongContext()
	players := mjContext.GetPlayers()

	fromPlayer := players[from]
	curPlayer := players[seat]

	cards := fromPlayer.GetHuansanzhangCards()

	var ok bool
	for _, card := range cards {
		if fromPlayer.HandCards, ok = utils.RemoveCards(fromPlayer.GetHandCards(), card, 1); !ok {
			return nil, fmt.Errorf("移除卡牌失败。玩家: %v, 卡牌：%v 手牌：%v", fromPlayer.GetPalyerId(), card, fromPlayer.GetHandCards())
		}
		curPlayer.HandCards = append(curPlayer.HandCards, card)
	}
	return cards, nil
}

// notifyFinish 通知换三张完成
func (s *HuansanzhangState) notifyFinish(flow interfaces.MajongFlow, dir room.Direction, exchangesIn map[int][]*majongpb.Card, exchangesOut map[int][]*majongpb.Card) {
	mjContext := flow.GetMajongContext()
	players := mjContext.GetPlayers()

	for index, player := range players {
		inCards, ok := exchangesIn[index]
		if !ok {
			continue
		}
		outCards, ok := exchangesOut[index]
		if !ok {
			continue
		}
		notify := room.RoomHuansanzhangFinishNtf{
			InCards:   utils.ServerCards2Uint32(inCards),
			OutCards:  utils.ServerCards2Uint32(outCards),
			Direction: dir.Enum(),
		}
		flow.PushMessages([]uint64{player.GetPalyerId()}, interfaces.ToClientMessage{
			MsgID: int(msgid.MsgID_ROOM_HUANSANZHANG_FINISH_NTF),
			Msg:   &notify,
		})
	}
}

// execute 执行换三张
func (s *HuansanzhangState) execute(flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
	newState, err = majongpb.StateID_state_huansanzhang, nil

	logEntry := logrus.WithField("func_name", "HuansanzhangState.execute")
	mjContext := flow.GetMajongContext()
	logEntry = utils.WithMajongContext(logEntry, mjContext)

	players := mjContext.GetPlayers()

	if !s.checkDone(players) {
		return
	}

	dir := s.randDirection(flow)
	logEntry = logEntry.WithField("direction", dir)
	pairs := s.makePairs(len(players), dir)

	exchangesIn := map[int][]*majongpb.Card{}
	exchangesOut := map[int][]*majongpb.Card{}

	for seat, from := range pairs {
		if cards, err := s.fetchCardFrom(flow, seat, from); err != nil {
			logEntry.Errorln(err)
			return newState, err
		} else {
			exchangesIn[seat] = cards
			exchangesOut[from] = cards
		}
	}
	mjContext.ExcutedHuansanzhang = true
	s.notifyFinish(flow, dir, exchangesIn, exchangesOut)
	flow.SetAutoEvent(majongpb.AutoEvent{
		EventId:      majongpb.EventID_event_huansanzhang_finish,
		EventContext: nil,
		WaitTime:     mjContext.GetOption().GetMaxHuansanzhangCartoonTime(),
	})
	return
}

//onRsq 换三张应答
func onHuanSanZhangRsq(playerID uint64, flow interfaces.MajongFlow) {
	// 错误码-成功
	errCode := room.RoomError_SUCCESS
	// 定缺应答 请求-响应
	toClientRsq := interfaces.ToClientMessage{
		MsgID: int(msgid.MsgID_ROOM_HUANSANZHANG_RSP),
		Msg: &room.RoomHuansanzhangRsp{
			ErrCode: &errCode,
		},
	}
	// 推送消息应答
	flow.PushMessages([]uint64{playerID}, toClientRsq)
	logrus.WithFields(logrus.Fields{
		"msgID":      msgid.MsgID_ROOM_HUANSANZHANG_RSP,
		"dingQueNtf": toClientRsq,
	}).Info("-----换三张成功应答")
}
