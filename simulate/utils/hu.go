package utils

import (
	"fmt"
	msgid "steve/client_pb/msgId"
	"steve/client_pb/room"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// SendHuReq 发送胡请求
func SendHuReq(deskData *DeskData, seat int) error {
	player := GetDeskPlayerBySeat(seat, deskData)
	client := player.Player.GetClient()
	_, err := client.SendPackage(CreateMsgHead(msgid.MsgID_ROOM_XINGPAI_ACTION_REQ), &room.RoomXingpaiActionReq{
		ActionId: room.XingpaiAction_XA_HU.Enum(),
	})
	return err
}

// CheckHuNotify 检查胡通知
func CheckHuNotify(t *testing.T, deskData *DeskData, huSeats []int, from int, card uint32, huType room.HuType) {
	huPlayers := []uint64{}
	for _, seat := range huSeats {
		huPlayers = append(huPlayers, GetDeskPlayerBySeat(seat, deskData).Player.GetID())
	}
	fromPlayer := GetDeskPlayerBySeat(from, deskData).Player.GetID()

	for _, player := range deskData.Players {
		expector, _ := player.Expectors[msgid.MsgID_ROOM_HU_NTF]
		ntf := room.RoomHuNtf{}
		expector.Recv(time.Second*1, &ntf)
		assert.Equal(t, huPlayers, ntf.GetPlayers())
		assert.Equal(t, fromPlayer, ntf.GetFromPlayerId())
		assert.Equal(t, card, ntf.GetCard())
		assert.Equal(t, huType, ntf.GetHuType())
	}
}

// CheckZiMoSettleNotify 检查自摸结算通知
func CheckZiMoSettleNotify(t *testing.T, deskData *DeskData, huSeats []int, from int, card uint32, huType room.HuType) {
	huPlayers := []uint64{}
	for _, seat := range huSeats {
		huPlayers = append(huPlayers, GetDeskPlayerBySeat(seat, deskData).Player.GetID())
	}

	for _, player := range deskData.Players {
		expector, _ := player.Expectors[msgid.MsgID_ROOM_INSTANT_SETTLE]
		ntf := room.RoomSettleInstantRsp{}
		expector.Recv(time.Second*1, &ntf)
		fmt.Println(ntf)
		assert.Equal(t, len(deskData.Players), len(ntf.BillPlayersInfo))
	}
}

// CheckDianPaoSettleNotify 检查点炮结算通知
func CheckDianPaoSettleNotify(t *testing.T, deskData *DeskData, huSeats []int, from int, card uint32, huType room.HuType) {
	huPlayers := []uint64{}
	for _, seat := range huSeats {
		huPlayers = append(huPlayers, GetDeskPlayerBySeat(seat, deskData).Player.GetID())
	}
	for _, player := range deskData.Players {
		expector, _ := player.Expectors[msgid.MsgID_ROOM_INSTANT_SETTLE]
		ntf := room.RoomSettleInstantRsp{}
		expector.Recv(time.Second*1, &ntf)
		fmt.Println(ntf)
		assert.Equal(t, len(deskData.Players), len(ntf.BillPlayersInfo))
	}
}
