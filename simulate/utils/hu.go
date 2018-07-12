package utils

import (
	"fmt"
	"steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/simulate/global"
	"testing"

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
	rs := len(deskData.Players)
	playerAll := []int{0, 1, 2, 3}[:rs]
	CheckHuNotifyBySeats(t, deskData, huSeats, from, card, huType, playerAll)
}

// CheckHuNotifyBySeats 指定玩家检查胡通知
func CheckHuNotifyBySeats(t *testing.T, deskData *DeskData, huSeats []int, from int, card uint32, huType room.HuType, otherSeats []int) {
	huPlayers := []uint64{}
	for _, seat := range huSeats {
		huPlayers = append(huPlayers, GetDeskPlayerBySeat(seat, deskData).Player.GetID())
	}
	fromPlayer := GetDeskPlayerBySeat(from, deskData).Player.GetID()
	for _, oseat := range otherSeats {
		player := GetDeskPlayerBySeat(oseat, deskData)
		expector, _ := player.Expectors[msgid.MsgID_ROOM_HU_NTF]
		ntf := room.RoomHuNtf{}
		expector.Recv(global.DefaultWaitMessageTime, &ntf)
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
		expector.Recv(global.DefaultWaitMessageTime, &ntf)
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
		assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, &ntf))
		assert.Equal(t, len(huSeats)+1, len(ntf.BillPlayersInfo))
	}
}

// CheckInstantSettleScoreNotify 检查立即分数结算通知
func CheckInstantSettleScoreNotify(t *testing.T, deskData *DeskData, winSeat int, winScore int64) {
	winplayer := GetDeskPlayerBySeat(winSeat, deskData)
	winID := winplayer.Player.GetID()
	expector, _ := winplayer.Expectors[msgid.MsgID_ROOM_INSTANT_SETTLE]
	ntf := room.RoomSettleInstantRsp{}
	assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, &ntf))
	for _, billInfo := range ntf.BillPlayersInfo {
		// 赢的分数
		if billInfo.GetPid() == winID {
			assert.Equal(t, billInfo.GetScore(), winScore)
		}
		fmt.Println(billInfo)
	}
}

// CheckRoundSettleScoreNotify 检查单局分数结算通知
func CheckRoundSettleScoreNotify(t *testing.T, deskData *DeskData, winSeat int, winScore int64) {
	winplayer := GetDeskPlayerBySeat(winSeat, deskData)
	winID := winplayer.Player.GetID()
	expector, _ := winplayer.Expectors[msgid.MsgID_ROOM_ROUND_SETTLE]
	ntf := room.RoomBalanceInfoRsp{}
	assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, &ntf))
	for _, billInfo := range ntf.BillPlayersInfo {
		// 赢的分数
		if billInfo.GetPid() == winID {
			assert.Equal(t, billInfo.GetScore(), winScore)
		}
		fmt.Println(billInfo)
	}
}
