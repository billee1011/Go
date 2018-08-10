package common

import (
	"fmt"
	"math/rand"
	"steve/client_pb/msgid"
	"steve/client_pb/room"
	"steve/common/mjoption"
	majongpb "steve/entity/majong"
	"steve/gutils"
	"steve/room/majong/interfaces"
	"steve/room/majong/utils"
	"time"

	"github.com/Sirupsen/logrus"
)

// HuansanzhangState 换三张状态
type HuansanzhangState struct {
}

// OnEntry 进入换三张状态
func (s *HuansanzhangState) OnEntry(flow interfaces.MajongFlow) {
	s.notifyPlayerHuangSanZhang(flow) //换三张推荐通知
}

// ProcessEvent 处理换三张事件
func (s *HuansanzhangState) ProcessEvent(eventID majongpb.EventID, eventContext interface{}, flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
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
	flow.GetMajongContext().TempData = new(majongpb.TempDatas) //清除临时数据

}

// nextState 下个状态
func (s *HuansanzhangState) nextState(flow interfaces.MajongFlow) majongpb.StateID {
	finished := flow.GetMajongContext().GetExcutedHuansanzhang()
	if !finished {
		return s.curState()
	}
	xpOption := mjoption.GetXingpaiOption(int(flow.GetMajongContext().GetXingpaiOptionId()))
	if xpOption.EnableDingque {
		return majongpb.StateID_state_dingque
	}
	return majongpb.StateID_state_zixun
}

// curState 当前状态
func (s *HuansanzhangState) curState() majongpb.StateID {
	return majongpb.StateID_state_huansanzhang
}

// onCartoonFinish 动画播放完毕
func (s *HuansanzhangState) onCartoonFinish(flow interfaces.MajongFlow, eventContext interface{}) (newState majongpb.StateID, err error) {
	cartoonFinishData := CartoonFinishData{
		CurState:        s.curState(),
		NextState:       s.nextState(flow),
		NeedCartoonType: room.CartoonType_CTNT_HUANSANZHANG,
		EventContext:    eventContext,
	}
	return OnCartoonFinish(cartoonFinishData, flow.GetMajongContext())
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
func (s *HuansanzhangState) onReq(eventID majongpb.EventID, eventContext interface{}, flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
	newState, err = majongpb.StateID_state_huansanzhang, nil

	logEntry := logrus.WithField("func_name", "HuansanzhangState.onReq")
	mjContext := flow.GetMajongContext()
	logEntry = utils.WithMajongContext(logEntry, mjContext)

	req := eventContext.(*majongpb.HuansanzhangRequestEvent)
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
	logrus.WithFields(logrus.Fields{
		"req_sure":      req.Sure,
		"req_player":    playerID,
		"req_hsz_cards": gutils.FmtMajongpbCards(reqCards),
		"sure":          player.GetHuansanzhangSure(),
	}).Infoln("--换三张中")
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
	HszFx := mjContext.GetOption().GetHszFx()
	if HszFx.GetNeedDeployFx() && HszFx.GetHuansanzhangFx() >= 0 && int(HszFx.GetHuansanzhangFx()) < len(directios) {
		towards = int(HszFx.GetHuansanzhangFx())
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
			return nil, fmt.Errorf("移除卡牌失败。玩家: %v, 卡牌：%v 手牌：%v", fromPlayer.GetPlayerId(), card, fromPlayer.GetHandCards())
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
		flow.PushMessages([]uint64{player.GetPlayerId()}, interfaces.ToClientMessage{
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
}

// notifyPlayerHuangSanZhang 通知玩家换三张
func (s *HuansanzhangState) notifyPlayerHuangSanZhang(flow interfaces.MajongFlow) {
	players := flow.GetMajongContext().GetPlayers()
	idHszMap := make(map[uint64]string)
	log := logrus.WithFields(logrus.Fields{})
	// 广播通知客户端进入定缺
	for _, player := range players {
		// 获取推荐换三张
		hszCards := gutils.GetRecommedHuanSanZhang(player.GetHandCards())
		// 检验换牌是否符合
		if !s.checkReq(log, player, hszCards) {
			log.WithFields(logrus.Fields{"hszCards": hszCards}).Infoln("换牌不符合")
			continue
		}
		// 先设置，用于超时AI
		player.HuansanzhangCards = hszCards
		hszNtf := &room.RoomHuansanzhangNtf{
			HszCard: utils.CardsToRoomCards(player.GetHuansanzhangCards()),
		}
		idHszMap[player.GetPlayerId()] = gutils.FmtMajongpbCards(player.GetHandCards())
		flow.PushMessages([]uint64{player.GetPlayerId()}, interfaces.ToClientMessage{
			MsgID: int(msgid.MsgID_ROOM_HUANSANZHANG_NTF),
			Msg:   hszNtf,
		})
	}
	// 日志
	log.WithFields(logrus.Fields{
		"idHszMap": idHszMap,
	}).Info("-----换三张开始-获取推荐换三张")
}
