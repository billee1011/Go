package utils

import (
	"steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/simulate/global"
	"testing"

	"github.com/stretchr/testify/assert"
)

// CheckBuhuaNtf 检查补花通知
func CheckBuhuaNtf(t *testing.T, buhuaSeats []int, huacards [][]uint32, bucards [][]uint32, recvSeats []int, deskData *DeskData) {
	for _, recvSeat := range recvSeats {
		player := GetDeskPlayerBySeat(recvSeat, deskData)
		expector := player.Expectors[msgId.MsgID_ROOM_BUHUA_NTF]
		ntf := room.RoomBuHuaNtf{}
		assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, &ntf))
		// if len(buhuaSeats) > 1 {
		buhuaInfos := ntf.GetBuhuaInfo()
		for index, info := range buhuaInfos {
			assert.Equal(t, huacards[index], info.OutHuaCards)
			buhuaPlayer := GetDeskPlayerBySeat(buhuaSeats[index], deskData)
			if buhuaPlayer.Player.GetID() == player.Player.GetID() {
				assert.Equal(t, bucards[index], info.BuCards)
			} else {
				if len(bucards[index]) == 0 {
					assert.Equal(t, *new([]uint32), info.GetBuCards())
				} else {
					assert.Equal(t, make([]uint32, len(bucards[index])), info.GetBuCards())
				}
			}
			assert.Equal(t, buhuaPlayer.Player.GetID(), info.GetPlayerId())
		}
		// }
	}
}
