package fantests

import (
	"steve/client_pb/common"
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func kanzhang(t *testing.T) *utils.DeskData {
	params := global.NewCommonStartGameParams()
	params.GameID = common.GameId_GAMEID_ERRENMJ // 二人
	params.PeiPaiGame = "ermj"
	params.PlayerNum = 2
	params.BankerSeat = 0
	params.Cards = [][]uint32{
		//{41, 41, 46, 45, 45, 44, 44, 43, 43, 47, 47, 42, 42},
		{42, 42, 44, 44, 45, 45, 46, 46, 14, 12, 46, 19, 19},
		{11, 11, 16, 15, 15, 15, 15, 13, 13, 17, 17, 42, 13},
	}
	params.WallCards = []uint32{16, 42, 44, 14, 13, 45}
	params.IsHsz = false
	params.IsDq = false
	deskData, err2 := utils.StartGame(params)
	assert.NotNil(t, deskData)
	assert.Nil(t, err2)

	return deskData
}

//TestFan_Kanzhang_Zimo_ERM 坎张立即结算自摸测试
// 庄摸牌46,自摸
//期望赢分：14 = [2（箭刻） + 1（坎张） + 6（混一色）  +4（无花）+1（自摸）]* 1
func TestFan_Kanzhang_Zimo_ERM(t *testing.T) {
	deskData := kanzhang(t)
	//开局 0 自询
	assert.Nil(t, utils.WaitZixunNtf(deskData, 0))
	// 0玩家出16
	assert.Nil(t, utils.SendChupaiReq(deskData, 0, 16))
	//1玩家能碰,能胡13
	assert.Nil(t, utils.WaitChupaiWenxunNtf(deskData, 1, false, false, false))
	//开局 1 弃
	assert.Nil(t, utils.SendQiReq(deskData, 1))
	//开局 1 自询
	assert.Nil(t, utils.WaitZixunNtf(deskData, 1))
	// 1玩家出42
	assert.Nil(t, utils.SendChupaiReq(deskData, 1, 42))
	//0玩家能碰,能碰42
	assert.Nil(t, utils.WaitChupaiWenxunNtf(deskData, 0, true, false, false))
	// 0玩家发送碰请求42
	assert.Nil(t, utils.SendPengReq(deskData, 0))
	//0 自询
	assert.Nil(t, utils.WaitZixunNtf(deskData, 0))
	// 0玩家出19
	assert.Nil(t, utils.SendChupaiReq(deskData, 0, 19))
	//1 自询
	assert.Nil(t, utils.WaitZixunNtf(deskData, 1))
	// 1玩家出44
	assert.Nil(t, utils.SendChupaiReq(deskData, 1, 44))
	//0玩家能碰,能碰44
	assert.Nil(t, utils.WaitChupaiWenxunNtf(deskData, 0, true, false, false))
	// 0玩家发送碰请求44
	assert.Nil(t, utils.SendPengReq(deskData, 0))
	//0 自询
	assert.Nil(t, utils.WaitZixunNtf(deskData, 0))
	// 0玩家出19
	assert.Nil(t, utils.SendChupaiReq(deskData, 0, 19))
	//1 自询
	assert.Nil(t, utils.WaitZixunNtf(deskData, 1))
	// 1玩家出14
	assert.Nil(t, utils.SendChupaiReq(deskData, 1, 14))
	//0 自询
	assert.Nil(t, utils.WaitZixunNtf(deskData, 0))
	// 0胡
	assert.Nil(t, utils.SendHuReq(deskData, 0))
	// 检测分数
	winScro := 14 * (len(deskData.Players) - 1)

	utils.CheckFanSettle(t, deskData, 4, 0, int64(winScro), room.FanType_FT_KANZHANG)
}
