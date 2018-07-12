package fantests

import (
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

//TestDaSiXi_ZiMo 共同步骤

func TestDaSiXi_ZiMo(t *testing.T) {
	params := global.NewCommonStartGameParams()
	params.PlayerNum = 2
	params.BankerSeat = 0
	params.PeiPaiGame = "ermj"
	params.GameID = room.GameId_GAMEID_ERRENMJ
	params.IsDq = false
	params.IsHsz = false
	params.Cards = [][]uint32{
		{11, 11, 41, 41, 41, 42, 42, 42, 43, 43, 43, 44, 44},
		{14, 14, 15, 15, 15, 16, 16, 16, 17, 17, 17, 18, 18},
	}
	params.WallCards = []uint32{44, 12, 19, 56, 13, 14, 57, 58, 14, 19, 19, 19}
	deskData, err := utils.StartGame(params)
	assert.NotNil(t, deskData)
	assert.Nil(t, err)
	// 庄家
	banker := deskData.BankerSeat
	//开局 庄家 自询 能自摸
	utils.CheckZixunNtf(t, deskData, banker, false, false, true)
	// 发送胡请求
	// assert.Nil(t, utils.SendHuReq(deskData, banker))
	// 检测所有玩家收到自摸通知
	// utils.CheckHuNotify(t, deskData, []int{banker}, banker, 44, room.HuType_HT_TIANHU)

	// 检测对对胡自摸分数
	// winScro := 2 * 2 * (len(deskData.Players) - 1)
	// utils.CheckInstantSettleScoreNotify(t, deskData, banker, int64(winScro))
}
