package utils

import (
	msgid "steve/client_pb/msgId"
	"steve/client_pb/room"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

// SendGangReq 发送杠请求
func SendGangReq(deskData *DeskData, seat int, card uint32, gangType room.GangType) error {
	zjPlayer := GetDeskPlayerBySeat(seat, deskData)
	zjClient := zjPlayer.Player.GetClient()
	_, err := zjClient.SendPackage(CreateMsgHead(msgid.MsgID_ROOM_XINGPAI_ACTION_REQ), &room.RoomXingpaiActionReq{
		ActionId: room.XingpaiAction_XA_GANG.Enum(),
		GangCard: proto.Uint32(card),
		GangType: gangType.Enum(),
	})
	return err
}

// CheckGangNotify 检查杠通知
func CheckGangNotify(t *testing.T, deskData *DeskData, to uint64, from uint64, card uint32, gangType room.GangType) {
	for _, player := range deskData.Players {
		expector, _ := player.Expectors[msgid.MsgID_ROOM_GANG_NTF]
		ntf := room.RoomGangNtf{}
		expector.Recv(time.Second*1, &ntf)
		assert.Equal(t, to, ntf.GetToPlayerId())
		assert.Equal(t, from, ntf.GetFromPlayerId())
		assert.Equal(t, card, ntf.GetCard())
		assert.Equal(t, gangType, ntf.GetGangType())
	}
}
