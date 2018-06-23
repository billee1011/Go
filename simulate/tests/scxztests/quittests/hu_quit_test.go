package quittests

import (
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestHuQuit 玩家点击胡牌后,退出游戏,再次加如桌子提示加入成功,进入匹配队列
func TestHuQuit(t *testing.T) {
	params := global.NewCommonStartGameParams()
	params.GameID = room.GameId_GAMEID_XUEZHAN // 血战
	params.IsHsz = false
	params.PeiPaiGame = "scxz"
	params.WallCards = []uint32{31, 31, 31, 31, 32, 32, 32, 32}
	deskData, err := utils.StartGame(params)
	assert.Nil(t, err)
	assert.NotNil(t, deskData)
	//庄家等待自询状态,可以天胡
	assert.Nil(t, utils.WaitZixunNtf(deskData, params.BankerSeat))
	//庄家选择天胡,并且退出游戏
	assert.Nil(t, utils.SendHuReq(deskData, 0))
	utils.CheckHuNotify(t, deskData, []int{0}, 0, 11, room.HuType_HT_TIANHU)
	utils.SendQuitReq(deskData, 0)
	//此时离开的玩家可以加入新的队列,等待新的游戏
	p := utils.GetDeskPlayerBySeat(0, deskData)
	rsp, err := utils.ApplyJoinDesk(p.Player, room.GameId_GAMEID_XUEZHAN)
	assert.Nil(t, err)
	assert.Equal(t, room.RoomError_SUCCESS, rsp.GetErrCode())
}

// TestHuQuitRecover 玩家没有胡牌,没有认输,退出游戏后提示游戏进行中,需要进行恢复对局
func TestHuQuitRecover(t *testing.T) {
	params := global.NewCommonStartGameParams()
	params.GameID = room.GameId_GAMEID_XUEZHAN // 血战
	params.IsHsz = false
	params.PeiPaiGame = "scxz"
	params.WallCards = []uint32{31, 31, 31, 31, 32, 32, 32, 32}
	deskData, err := utils.StartGame(params)
	assert.Nil(t, err)
	assert.NotNil(t, deskData)
	//庄家等待自询状态,可以天胡
	assert.Nil(t, utils.WaitZixunNtf(deskData, params.BankerSeat))
	// 庄家选择天胡,并且退出游戏
	// assert.Nil(t, utils.SendHuReq(deskData, 0))
	// utils.CheckHuNotify(t, deskData, []int{0}, 0, 11, room.HuType_HT_TIANHU)
	utils.SendQuitReq(deskData, 0)
	//此时离开的玩家可以加入新的队列,等待新的游戏
	p := utils.GetDeskPlayerBySeat(0, deskData)
	rsp, err := utils.ApplyJoinDesk(p.Player, room.GameId_GAMEID_XUEZHAN)
	assert.Nil(t, err)
	assert.Equal(t, room.RoomError_DESK_GAME_PLAYING, rsp.GetErrCode())
}
