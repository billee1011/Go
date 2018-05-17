package tests

import (
	msgid "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/simulate/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Test_Peng 测试碰
func Test_Peng(t *testing.T) {
	var Int5w uint32 = 15
	params := commonStartGameParams
	// 0 号玩家的最后一张牌改成 5W， 打出后 1 号玩家可碰
	params.Cards[0][13] = &Card5W
	params.Cards[1][0] = &Card4W
	// 修改换三张的牌
	params.HszCards = [][]*room.Card{
		utils.MakeRoomCards(Card3W, Card3W, Card3W),
		utils.MakeRoomCards(Card7W, Card7W, Card7W),
		utils.MakeRoomCards(Card3T, Card3T, Card3T),
		utils.MakeRoomCards(Card7T, Card7T, Card7T),
	}

	deskData, err := utils.StartGame(params)
	assert.Nil(t, err)

	// 庄家出 5W
	assert.Nil(t, utils.SendChupaiReq(deskData, params.BankerSeat, Int5w))

	// 1 号玩家收到出牌问询通知，可以碰
	pengSeat := (commonStartGameParams.BankerSeat + 1) % len(deskData.Players)
	pengPlayer := utils.GetDeskPlayerBySeat(pengSeat, deskData)

	expector, _ := pengPlayer.Expectors[msgid.MsgID_ROOM_CHUPAIWENXUN_NTF]
	ntf := room.RoomChupaiWenxunNtf{}
	assert.Nil(t, expector.Recv(time.Second*1, &ntf))
	assert.Equal(t, Int5w, ntf.GetCard())
	assert.True(t, ntf.GetEnablePeng())
	assert.True(t, ntf.GetEnableQi())

	// 请求碰
	pengClient := pengPlayer.Player.GetClient()
	pengClient.SendPackage(utils.CreateMsgHead(msgid.MsgID_ROOM_XINGPAI_ACTION_REQ), &room.RoomXingpaiActionReq{
		ActionId: room.XingpaiAction_XA_PENG.Enum(),
	})

	from := utils.GetDeskPlayerBySeat(params.BankerSeat, deskData)
	// 所有玩家收到碰通知
	checkPengNotify(t, deskData, pengPlayer.Player.GetID(), from.Player.GetID(), Int5w)
}

func checkPengNotify(t *testing.T, deskData *utils.DeskData, to uint64, from uint64, card uint32) {
	for _, player := range deskData.Players {
		expector, _ := player.Expectors[msgid.MsgID_ROOM_PENG_NTF]
		pengNtf := room.RoomPengNtf{}
		expector.Recv(time.Second*1, &pengNtf)
		assert.Equal(t, to, pengNtf.GetToPlayerId())
		assert.Equal(t, from, pengNtf.GetFromPlayerId())
		assert.Equal(t, card, pengNtf.GetCard())
	}
}
