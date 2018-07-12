package utils

import (
	 "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/simulate/global"
	"testing"

	"github.com/stretchr/testify/assert"
)

// WaitZixunNtf 等待自询通知
func WaitZixunNtf(desk *DeskData, seat int) error {
	player := GetDeskPlayerBySeat(seat, desk)
	expector, _ := player.Expectors[msgid.MsgID_ROOM_ZIXUN_NTF]

	ntf := room.RoomZixunNtf{}
	return expector.Recv(global.DefaultWaitMessageTime, &ntf)
}

// CheckZixunNtf 检测自询通知,是否能补杠，暗杠，自摸
func CheckZixunNtf(t *testing.T, desk *DeskData, seat int, canBuGang, canAnGang, canZiMo bool) {
	player := GetDeskPlayerBySeat(seat, desk)
	expector, _ := player.Expectors[msgid.MsgID_ROOM_ZIXUN_NTF]

	ntf := room.RoomZixunNtf{}
	assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, &ntf))
	assert.Equal(t, len(ntf.GetEnableBugangCards()) > 0, canBuGang)
	assert.Equal(t, len(ntf.GetEnableAngangCards()) > 0, canAnGang)
	assert.Equal(t, ntf.GetEnableZimo(), canZiMo)
}

// CheckZixunNtfWithTing 自询检查听
func CheckZixunNtfWithTing(t *testing.T, desk *DeskData, seat int, canBuGang, canAnGang, canZiMo bool, canTing bool) {
	player := GetDeskPlayerBySeat(seat, desk)
	expector, _ := player.Expectors[msgid.MsgID_ROOM_ZIXUN_NTF]
	ntf := room.RoomZixunNtf{}
	assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, &ntf))
	assert.Equal(t, len(ntf.GetEnableBugangCards()) > 0, canBuGang)
	assert.Equal(t, len(ntf.GetEnableAngangCards()) > 0, canAnGang)
	assert.Equal(t, ntf.GetEnableZimo(), canZiMo)
	assert.Equal(t, canTing, ntf.GetEnableTing())
}
