package utils

import (
	msgid "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/interfaces"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createMsgHead(msgID msgid.MsgID) interfaces.SendHead {
	return interfaces.SendHead{
		Head: interfaces.Head{
			MsgID: uint32(msgID),
		},
	}
}

// CreateMsgHead 创建消息头
func CreateMsgHead(msgID msgid.MsgID) interfaces.SendHead {
	return createMsgHead(msgID)
}

// MakeRoomCards 构造牌切片
func MakeRoomCards(card ...room.Card) []*room.Card {
	result := []*room.Card{}
	for i := range card {
		result = append(result, &card[i])
	}
	return result
}

//CheckChuPaiNotify 检查出牌广播
func CheckChuPaiNotify(t *testing.T, deskData *DeskData, card uint32, activePlayer uint64) {
	for _, player := range deskData.Players {
		messageExpector := player.Expectors[msgid.MsgID_ROOM_CHUPAI_NTF]
		ntf := &room.RoomChupaiNtf{}
		assert.Nil(t, messageExpector.Recv(global.DefaultWaitMessageTime, ntf))
		assert.Equal(t, card, ntf.GetCard())
		assert.Equal(t, activePlayer, ntf.GetPlayer())
	}
}

//CheckMoPaiNotify 检查摸牌广播
func CheckMoPaiNotify(t *testing.T, deskData *DeskData, mopaiSeat int) *DeskPlayer {
	player := GetDeskPlayerBySeat(mopaiSeat, deskData)
	for _, deskPlayer := range deskData.Players {
		expector, _ := deskPlayer.Expectors[msgid.MsgID_ROOM_MOPAI_NTF]
		ntf := &room.RoomMopaiNtf{}
		assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, ntf))
		assert.Equal(t, false, ntf.GetBack())
		if player.Seat == deskPlayer.Seat {
			assert.Equal(t, uint32(31), ntf.GetCard())
		} else {
			assert.Equal(t, uint32(0), ntf.GetCard())
		}
	}
	return player
}
