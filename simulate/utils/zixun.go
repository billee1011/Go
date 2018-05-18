package utils

import (
	msgid "steve/client_pb/msgId"
	"steve/client_pb/room"
	"time"
)

// WaitZixunNtf 等待自询通知
func WaitZixunNtf(desk *DeskData, seat int) error {
	player := GetDeskPlayerBySeat(seat, desk)
	expector, _ := player.Expectors[msgid.MsgID_ROOM_ZIXUN_NTF]

	ntf := room.RoomZixunNtf{}
	return expector.Recv(time.Second*2, &ntf)
}
