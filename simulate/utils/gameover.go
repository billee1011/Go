package utils

import (
	msgid "steve/client_pb/msgid"
	"steve/client_pb/room"
	"steve/simulate/global"
	"testing"

	"github.com/stretchr/testify/assert"
)

//WaitGameOverNtf 所有玩家等待游戏结束通知
func WaitGameOverNtf(t *testing.T, d *DeskData) {
	for i := 0; i < len(d.Players); i++ {
		player := GetDeskPlayerBySeat(i, d)
		expector, _ := player.Expectors[msgid.MsgID_ROOM_GAMEOVER_NTF]
		ntf := &room.RoomGameOverNtf{}
		assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, ntf))
	}
}
