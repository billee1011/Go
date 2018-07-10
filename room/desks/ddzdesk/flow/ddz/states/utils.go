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
	"github.com/Sirupsen/logrus"
)

// PeiPai 配牌工具
func PeiPai(wallCards []uint32, value string) error {
	var cards []uint32
	for i := 0; i < len(value); i = i + 3 {
		card, err := strconv.Atoi(value[i : i+2])
		if err != nil {
			return err
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
	logrus.WithFields(logrus.Fields{"wallCards":wallCards, "peipai:":value}).Debug("斗地主配牌成功")
	return nil
}

// getDDZContext 从状态机中获取斗地主现场
func getDDZContext(m machine.Machine) *ddz.DDZContext {
	dm, ok := m.(*ddzmachine.DDZMachine)
	if !ok {
		return nil
	}
	return dm.GetDDZContext()
}

func remove(playerIds []uint64, removeId uint64) []uint64 {
	var result []uint64
	for _, playerId := range playerIds {
		if playerId != removeId {
			result = append(result, playerId)
		}
	}
	return result
}

func getPlayerIds(m machine.Machine) []uint64 {
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
	return sendMessage(m, []uint64{playerID}, msgID, body)
}

func broadcast(m machine.Machine, msgID msgid.MsgID, body proto.Message) error {
	return sendMessage(m, getPlayerIds(m), msgID, body)
}

func broadcastExcept(m machine.Machine, playerID uint64, msgID msgid.MsgID, body proto.Message) error {
	allPlayers := getPlayerIds(m)
	players := []uint64{}
	for _, pid := range allPlayers {
		if pid != playerID {
			players = append(players, pid)
		}
	}
	return sendMessage(m, players, msgID, body)
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
func ContainsAll(handCards []Poker, outCards []Poker) bool {
	for _, outCard := range outCards {
		if !Contains(handCards, outCard) {
			return false
		}
	}
	return true
}

// Contains cards是否包含card
func Contains(cards []Poker, card Poker) bool {
	for _, value := range cards {
		if value.equals(card) {
			return true
		}
	}
	return false
}

// ContainsPoint cards是否包含点数
func ContainsPoint(cards []Poker, point uint32) bool {
	for _, card := range cards {
		if card.Point == point {
			return true
		}
	}
	return false
}

// RemoveAll 从cards中删除removeCards
func RemoveAll(cards []Poker, removeCards []Poker) []Poker {
	var result []Poker
	for _, card := range cards {
		if !Contains(removeCards, card) {
			result = append(result, card)
		}
	}
	return result
}

// AppendAll 从cards中添加addCards
func AppendAll(cards []Poker, addCards []Poker) []Poker {
	for _, card := range addCards {
		cards = append(cards, card)
	}
	return cards
}

func If(judge bool, trueReturn interface{}, falseReturn interface{}) interface{} {
	if judge {
		return trueReturn
	} else {
		return falseReturn
	}
}