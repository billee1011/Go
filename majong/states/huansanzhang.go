package states

import (
	"errors"
	"fmt"
	"math/rand"
	"steve/client_pb/room"
	"steve/majong/global"
	"steve/majong/interfaces"
	"steve/majong/utils"
	"steve/peipai"
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
}

// ProcessEvent 处理换三张事件
func (s *HuansanzhangState) ProcessEvent(eventID majongpb.EventID, eventContext []byte, flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
	if eventID == majongpb.EventID_event_huansanzhang_request {
		huansanzhangReq := new(majongpb.HuansanzhangRequestEvent)
		err := proto.Unmarshal(eventContext, huansanzhangReq)
		if err != nil {
			return majongpb.StateID_state_huansanzhang, err
		}
		players := flow.GetMajongContext().Players
		player := utils.GetPlayerByID(players, huansanzhangReq.GetHead().PlayerId)
		huansnahzangCards := huansanzhangReq.Cards
		if len(huansnahzangCards) != 3 {
			return majongpb.StateID_state_huansanzhang, errors.New("换三张牌数不为3")
		}
		for i := 0; i < len(huansnahzangCards)-1; i++ {
			if huansnahzangCards[i].Color != huansnahzangCards[i+1].Color {
				return majongpb.StateID_state_huansanzhang, errors.New("换三张花色不一样")
			}
			if !utils.ContainCard(player.HandCards, huansnahzangCards[i]) {
				return majongpb.StateID_state_huansanzhang, fmt.Errorf("换三张的牌%v不存在玩家%v手牌", huansnahzangCards[i].Point, player.PalyerId)
			}

		}

		player.HuansanzhangCards = huansnahzangCards
		if huansanzhangReq.Sure {
			player.HuansanzhangSure = true
		}
		huansnazhangDone := 0
		for i := 0; i < len(players); i++ {
			if (len(players[i].HuansanzhangCards) == 3) && (players[i].HuansanzhangSure == true) {
				huansnazhangDone++
			}
		}
		if huansnazhangDone == len(players) { // 所有玩家换三张牌都收到，开始处理换三张
			rd := rand.New(rand.NewSource(time.Now().Unix())) // 生成换三张方向
			towards := rd.Intn(3)
			gameName := getGameName(flow)
			fx := peipai.GetHSZFangXiang(gameName)
			if fx != -1 {
				towards = fx
			}
			l := len(players)
			result := false
			for i, player := range players { // 根据不同方向处理换三张
				var pairPlayer *majongpb.Player
				if towards == int(room.Direction_ClockWise) {
					result = processHuansanzhang(player, players[(i+l-1)%l])
					pairPlayer = players[(i+l-1)%l]

				} else if towards == int(room.Direction_AntiClockWise) {
					result = processHuansanzhang(player, players[(i+l+1)%l])
					pairPlayer = players[(i+l+1)%l]

				} else if towards == int(room.Direction_Opposite) {
					result = processHuansanzhang(player, players[(i+l+2)%l])
					pairPlayer = players[(i+l+2)%l]
				}
				huansanzhangFinishNtf := room.RoomHuansanzhangFinishNtf{
					InCards:   utils.CardsToRoomCards(pairPlayer.HuansanzhangCards),
					OutCards:  utils.CardsToRoomCards(player.HuansanzhangCards),
					Direction: room.Direction(towards).Enum(),
				}
				if result {
					//TODO
					fmt.Println(huansanzhangFinishNtf)
					//flow.PushMessages([]uint64{player.PalyerId}，huansanzhangFinishNtf)
				}
			}
			return majongpb.StateID_state_dingque, nil
		}

		return majongpb.StateID_state_huansanzhang, errors.New("换三张尚有玩家未完成")
	}
	return majongpb.StateID_state_huansanzhang, global.ErrInvalidEvent
}

// OnExit 退出换三张状态
func (s *HuansanzhangState) OnExit(flow interfaces.MajongFlow) {

}

// processHuansanzhang 处理玩家换三张
func processHuansanzhang(playerIn, playerOut *majongpb.Player) bool {
	outCards := playerOut.HuansanzhangCards
	for _, outCard := range outCards {
		deleted := false
		playerOut.HandCards, deleted = utils.DeleteCardFromLast(playerOut.HandCards, outCard)
		if !deleted {
			logrus.Fatalf("huansanzhang err deleted fail cards[%v] deleteCard[%v]\n", playerOut.HandCards, outCard)
			return false
		}
		playerIn.HandCards = append(playerIn.HandCards, outCard)
	}
	return true
}
