package utils

import (
	"fmt"
	 "steve/client_pb/msgid"
	"steve/client_pb/room"
	"steve/simulate/global"
	"testing"

	"github.com/stretchr/testify/assert"
)

// WaitTingInfoNtf 等待听牌信息通知
func WaitTingInfoNtf(t *testing.T, desk *DeskData, seat int, canTingCards ...uint32) error {
	player := GetDeskPlayerBySeat(seat, desk)
	expector, _ := player.Expectors[msgid.MsgID_ROOM_TINGINFO_NTF]

	ntf := room.RoomTingInfoNtf{}
	if err := expector.Recv(global.DefaultWaitMessageTime, &ntf); err != nil {
		return err
	}
	tingCardInfo := ntf.GetTingCardInfos()
	canTingSum := len(canTingCards)
	if canTingSum > 0 {
		assert.Equal(t, len(tingCardInfo), canTingSum)
		for _, tingCard := range tingCardInfo {
			flag := false
			for _, canTingCard := range canTingCards {
				if canTingCard == tingCard.GetTingCard() {
					flag = true
				}
			}
			assert.True(t, flag)
		}
	}
	fmt.Println(tingCardInfo)
	return nil
}
