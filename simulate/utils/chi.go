package utils

import (
	 "steve/client_pb/msgid"
	"steve/client_pb/room"
	"steve/simulate/global"
	"testing"

	"github.com/stretchr/testify/assert"
)

// SendChiReq 请求吃
func SendChiReq(t *testing.T, deskData *DeskData, card []uint32, seat int) error {
	deskPlayer := GetDeskPlayerBySeat(seat, deskData)
	client := deskPlayer.Player.GetClient()
	_, err := client.SendPackage(CreateMsgHead(msgid.MsgID_ROOM_XINGPAI_ACTION_REQ), &room.RoomXingpaiActionReq{
		ActionId: room.XingpaiAction_XA_CHI.Enum(),
		ChiCards: card,
	})
	return err
}

// CheckChiNtfWithSeats 指定玩家检查吃牌广播
func CheckChiNtfWithSeats(t *testing.T, deskData *DeskData, card []uint32, to int, from int, expectedSeats []int) {
	chiPlayer := GetDeskPlayerBySeat(to, deskData)
	chupaiPlayer := GetDeskPlayerBySeat(from, deskData)
	for _, seat := range expectedSeats {
		deskPlayer := GetDeskPlayerBySeat(seat, deskData)
		expector := deskPlayer.Expectors[msgid.MsgID_ROOM_CHI_NTF]
		ntf := room.RoomChiNtf{}
		assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, &ntf))
		assert.Equal(t, card, ntf.Cards)
		assert.Equal(t, chiPlayer.Player.GetID(), ntf.GetToPlayerId())
		assert.Equal(t, chupaiPlayer.Player.GetID(), ntf.GetFromPlayerId())
	}
}
