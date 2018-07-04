package states

import (
	"fmt"
	msgid "steve/client_pb/msgId"
	"steve/room/desks/ddzdesk/flow/ddz/ddzmachine"
	"steve/room/desks/ddzdesk/flow/machine"
	"steve/server_pb/ddz"
	"time"

	"github.com/golang/protobuf/proto"
	"strconv"
)

// PeiPai 配牌工具
func PeiPai(wallCards []uint32, value string) ([]uint32, error) {
	var cards []uint32
	for i := 0; i < len(value); i = i + 3 {
		card, err := strconv.Atoi(value[i : i+2])
		if err != nil {
			return nil, err
		}
		cards = append(cards, uint32(card))
	}
	for i := 0; i < len(cards); i++ {
		for j := len(wallCards) - 1; j >= 0; j-- {
			if cards[i] == wallCards[j] {
				wallCards[i], wallCards[j] = wallCards[j], wallCards[i]
				break
			}
		}
	}
	return wallCards, nil
}

// getDDZContext 从状态机中获取斗地主现场
func getDDZContext(m machine.Machine) *ddz.DDZContext {
	dm, ok := m.(*ddzmachine.DDZMachine)
	if !ok {
		return nil
	}
	return dm.GetDDZContext()
}

func getPlayers(m machine.Machine) []uint64 {
	dm, ok := m.(*ddzmachine.DDZMachine)
	if !ok {
		return nil
	}

	players := []uint64{}

	for _, player := range dm.GetDDZContext().GetPlayers() {
		players = append(players, player.GetPalyerId())
	}
	return players
}

// IsAllAbandon 判断是否三家都弃地主
func IsAllAbandon(players []*ddz.Player) bool {
	for _, player := range players {
		if player.Grab {
			return false
		}
	}
	return true
}

// GetTotalGrab 获取抢庄总倍数
func GetTotalGrab(players []*ddz.Player) (totalGrab uint32) {
	totalGrab = 1
	for _, player := range players {
		if player.Grab {
			totalGrab = totalGrab * 2
		}
	}
	return
}

// GetTotalDouble 获取加倍总倍数
func GetTotalDouble(players []*ddz.Player) (totalGrab uint32) {
	totalGrab = 1
	for _, player := range players {
		if player.IsDouble {
			totalGrab = totalGrab * 2
		}
	}
	return
}

//GetPlayerByID 根据玩家id获取玩家
func GetPlayerByID(players []*ddz.Player, id uint64) *ddz.Player {
	for _, player := range players {
		if player.PalyerId == id {
			return player
		}
	}
	return nil
}

//GetNextPlayerByID 根据玩家id获取下个玩家
func GetNextPlayerByID(players []*ddz.Player, id uint64) *ddz.Player {
	for k, player := range players {
		if player.PalyerId == id {
			index := (k + 1) % len(players)
			return players[index]
		}
	}
	return nil
}

// sendMessage 向玩家发送消息
func sendMessage(m machine.Machine, players []uint64, msgID msgid.MsgID, body proto.Message) error {
	dm, ok := m.(*ddzmachine.DDZMachine)
	if !ok {
		return fmt.Errorf("不是斗地主状态机")
	}
	return dm.SendMessage(players, msgID, body)
}

func sendToPlayer(m machine.Machine, playerID uint64, msgID msgid.MsgID, body proto.Message) error {
	dm, ok := m.(*ddzmachine.DDZMachine)
	if !ok {
		return fmt.Errorf("不是斗地主状态机")
	}
	return dm.SendMessage([]uint64{playerID}, msgID, body)
}

func broadcast(m machine.Machine, msgID msgid.MsgID, body proto.Message) error {
	dm, ok := m.(*ddzmachine.DDZMachine)
	if !ok {
		return fmt.Errorf("不是斗地主状态机")
	}
	return dm.SendMessage(getPlayers(m), msgID, body)
}

func broadcastExcept(m machine.Machine, playerID uint64, msgID msgid.MsgID, body proto.Message) error {
	dm, ok := m.(*ddzmachine.DDZMachine)
	if !ok {
		return fmt.Errorf("不是斗地主状态机")
	}
	allPlayers := getPlayers(m)
	players := []uint64{}
	for _, pid := range allPlayers {
		if pid != playerID {
			players = append(players, pid)
		}
	}
	return dm.SendMessage(players, msgID, body)
}

// setMachineAutoEvent 设置状态机自动事件
func setMachineAutoEvent(m machine.Machine, event machine.Event, duration time.Duration) {
	dm, ok := m.(*ddzmachine.DDZMachine)
	if !ok {
		return
	}
	dm.SetAutoEvent(event, duration)
}


// ContainsAll handCards是否包含所有outCards
func ContainsAll(handCards []DDZCard, outCards []DDZCard) bool {
	for _, outCard := range outCards {
		if(!Contains(handCards, outCard)){
			return false
		}
	}
	return true
}

// Contains cards是否包含card
func Contains(cards []DDZCard, card DDZCard) bool {
	for _, value := range cards {
		if value.equals(card) {
			return true
		}
	}
	return false
}

// ContainsPoint cards是否包含点数
func ContainsPoint(cards []DDZCard, point uint32) bool {
	for _, card := range cards {
		if card.point == point {
			return true
		}
	}
	return false
}

// RemovePoint 删除cards中所有点数为point的牌,并分别返回
func RemovePoint(cards []DDZCard, point uint32) (remain []DDZCard, deleted []DDZCard) {
	for _, card := range cards {
		if card.point == point {
			deleted = append(deleted, card)
		} else {
			remain = append(remain, card)
		}
	}
	return
}

// RemoveAll 从cards中删除removeCards
func RemoveAll(cards []DDZCard, removeCards []DDZCard) []DDZCard {
	var result []DDZCard
	for _, card := range cards {
		if !Contains(removeCards, card) {
			result = append(result, card)
		}
	}
	return result
}

// Remove 从cards中删除removeCard
func Remove(cards []DDZCard, removeCard DDZCard) []DDZCard {
	var result []DDZCard
	for _, card := range cards {
		if !card.equals(removeCard) {
			result = append(result, card)
		}
	}
	return result
}
