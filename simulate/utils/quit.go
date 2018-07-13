package utils

import (
	"steve/client_pb/msgid"
	"steve/client_pb/room"
	"steve/simulate/global"
	"testing"

	"github.com/stretchr/testify/assert"
)

// SendQuitReq 发送退出牌桌请求
func SendQuitReq(deskData *DeskData, seat int) error {
	player := GetDeskPlayerBySeat(seat, deskData)
	client := player.Player.GetClient()
	_, err := client.SendPackage(CreateMsgHead(msgid.MsgID_ROOM_DESK_QUIT_REQ), &room.RoomDeskQuitReq{})
	return err
}

// RecvQuitNtf 退出广播
func RecvQuitNtf(t *testing.T, deskData *DeskData, seats []int) {
	for _, seat := range seats {
		player := GetDeskPlayerBySeat(seat, deskData)
		expector, _ := player.Expectors[msgid.MsgID_ROOM_DESK_QUIT_ENTER_NTF]
		ntf := room.RoomDeskQuitEnterNtf{}
		assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, &ntf))
		assert.Equal(t, room.QuitEnterType_QET_QUIT, ntf.GetType())
	}
}
