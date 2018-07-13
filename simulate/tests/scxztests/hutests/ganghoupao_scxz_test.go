package hutests

import (
	"fmt"
	"steve/client_pb/msgid"
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Test_SCXZ_Duo_Qiangganghu 杠后炮测试
// 0 号玩家暗杠，摸1筒，并且打出9万， 然后1号玩家可胡
// 期望：
// 1. 1号玩家收到出牌问询通知，且可胡
// 2. 1号玩家请求胡，所有玩家收到胡通知，2号玩家摸牌后收到自询通知
func Test_SCXZ_Ganghoupao(t *testing.T) {
	params := global.NewCommonStartGameParams()
	params.GameID = room.GameId_GAMEID_XUEZHAN // 血战
	params.PeiPaiGame = "scxz"
	params.BankerSeat = 0
	huSeat := 1
	gangSeat := params.BankerSeat
	// 0 号玩家手牌改成 Card1W, Card1W, Card1W, Card1W, Card2W, Card2W, Card2W, Card2W, Card3W, Card3W, Card3W, Card3W, Card4W, Card9W
	// 换三张后手牌为 Card5T, Card5T, Card5T， Card1W, Card2W, Card2W, Card2W, Card2W, Card3W, Card3W, Card3W, Card3W, Card4W, Card9W
	params.Cards[0][13] = 19
	// 1 号玩家手牌改成 Card5W, Card5W, Card5W, Card5W, Card6W, Card6W, Card6W, Card6W, Card7W, Card7W, Card7W, Card7W, Card9W
	params.Cards[1][12] = 19
	// 1号玩家换三张后手牌为  Card1W, Card1W, Card1W, Card5W, Card5W, Card5W,  Card6W, Card6W, Card6W,  Card7W, Card7W, Card7W, Card9W
	params.HszCards[1] = []uint32{15, 16, 17}

	// 墙牌改为 1筒
	params.WallCards = []uint32{31}

	deskData, err := utils.StartGame(params)
	assert.Nil(t, err)
	assert.NotNil(t, deskData)
	// 庄家杠 2W
	assert.Nil(t, utils.WaitZixunNtf(deskData, 0))
	assert.Nil(t, utils.SendGangReq(deskData, 0, 12, room.GangType_AnGang))
	// 庄家等到自询通知后出 9W
	assert.Nil(t, utils.WaitZixunNtf(deskData, 0))
	assert.Nil(t, utils.SendChupaiReq(deskData, 0, 19))

	// 1 号玩家收到出牌问询通知， 可以胡
	huPlayer := utils.GetDeskPlayerBySeat(1, deskData)
	expector, _ := huPlayer.Expectors[msgid.MsgID_ROOM_CHUPAIWENXUN_NTF]
	ntf := room.RoomChupaiWenxunNtf{}
	assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, &ntf))
	assert.True(t, ntf.GetEnableDianpao())
	assert.True(t, ntf.GetEnableQi())

	// 发送胡请求
	assert.Nil(t, utils.SendHuReq(deskData, 1))
	// 检测所有玩家收到杆后炮通知
	utils.CheckHuNotify(t, deskData, []int{1}, 0, 19, room.HuType_HT_GANGHOUPAO)

	// 检测所有玩家是否收到呼叫转移的通知
	checkGangHouPaoSettleScoreNotify(t, deskData, gangSeat, huSeat)

}

// checkGangHouPaoSettleScoreNotify 检查杠后炮分数结算通知
func checkGangHouPaoSettleScoreNotify(t *testing.T, deskData *utils.DeskData, gangSeat int, huSeat int) {
	gangplayer := utils.GetDeskPlayerBySeat(gangSeat, deskData)
	gangID := gangplayer.Player.GetID()
	huPlayer := utils.GetDeskPlayerBySeat(huSeat, deskData)
	huPlayerID := huPlayer.Player.GetID()
	expector, _ := gangplayer.Expectors[msgid.MsgID_ROOM_INSTANT_SETTLE]
	ntf := room.RoomSettleInstantRsp{}
	expector.Recv(global.DefaultWaitMessageTime, &ntf)
	assert.Equal(t, len(deskData.Players), len(ntf.BillPlayersInfo))
	gangWinScore := 6
	for _, billInfo := range ntf.BillPlayersInfo {
		if billInfo.GetPid() == gangID {
			assert.Equal(t, billInfo.GetScore(), int64(gangWinScore))
		} else {
			assert.Equal(t, billInfo.GetScore(), -int64((gangWinScore / 3)))
		}
	}
	expector, _ = gangplayer.Expectors[msgid.MsgID_ROOM_INSTANT_SETTLE]
	ntf = room.RoomSettleInstantRsp{}
	expector.Recv(global.DefaultWaitMessageTime, &ntf)
	dianpaoWinScore := 16
	for _, billInfo := range ntf.BillPlayersInfo {
		if billInfo.GetPid() == gangID {
			assert.Equal(t, billInfo.GetScore(), int64(-dianpaoWinScore))
		} else if billInfo.GetPid() == huPlayerID {
			assert.Equal(t, billInfo.GetScore(), int64((dianpaoWinScore)))
		} else {
			assert.Equal(t, billInfo.GetScore(), int64(0))
		}
	}

	expector, _ = gangplayer.Expectors[msgid.MsgID_ROOM_INSTANT_SETTLE]
	ntf = room.RoomSettleInstantRsp{}
	expector.Recv(global.DefaultWaitMessageTime, &ntf)
	fmt.Println(ntf)
	callTransferScore := 6
	for _, billInfo := range ntf.BillPlayersInfo {
		if billInfo.GetPid() == gangID {
			assert.Equal(t, billInfo.GetScore(), int64(-callTransferScore))
		} else if billInfo.GetPid() == huPlayerID {
			assert.Equal(t, billInfo.GetScore(), int64((callTransferScore)))
		} else {
			assert.Equal(t, billInfo.GetScore(), int64(0))
		}
	}

	expector, _ = gangplayer.Expectors[msgid.MsgID_ROOM_ROUND_SETTLE]
	ntf2 := room.RoomBalanceInfoRsp{}
	expector.Recv(time.Second*5, &ntf2)
	fmt.Println(ntf2)
}
