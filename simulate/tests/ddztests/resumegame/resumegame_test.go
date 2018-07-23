package play

import (
	"steve/client_pb/msgid"
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"steve/simulate/utils/doudizhu"
	"testing"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

//TestPlaycard1 出牌测试
//游戏过程中0号玩家发起叫地主
//期望：
//     1. 所有玩家都收到，0号玩家的叫地主广播
func TestPlaycard1(t *testing.T) {

	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "TestPlaycard1()",
	})

	// 配牌1
	params := doudizhu.NewStartDDZGameParamsTest1()

	deskData, err := utils.StartDDZGame(params)
	assert.NotNil(t, deskData)
	assert.Nil(t, err)

	// 当前状态的时间间隔
	logEntry.Infof("当前状态 = %v, 进入下一状态等待时间 = %v", deskData.DDZData.NextState.GetStage(), deskData.DDZData.NextState.GetTime())

	// 叫地主用例1
	assert.Nil(t, doudizhu.JiaodizhuTest1(deskData))

	// 加倍用例1
	assert.Nil(t, doudizhu.JiabeiTest1(deskData))

	// ------------------------------------------------ 	加倍后恢复对局	  ---------------------------------------
	// 最终地主的deskPlayer
	lordPlayer := deskData.Players[deskData.DDZData.ResultLordID]

	// 地主断开连接
	assert.Nil(t, lordPlayer.Player.GetClient().Stop())
	time.Sleep(time.Millisecond * 200) // 等200毫秒，确保连接断开

	// 重新登录
	accountID := lordPlayer.Player.GetAccountID()
	accountName := utils.GenerateAccountName(accountID)
	player, err := utils.LoginPlayer(accountID, accountName)
	assert.Nil(t, err)
	assert.NotNil(t, player)
	client := player.GetClient()

	// 步骤4
	utils.UpdateDDZPlayerClientInfo(client, player, deskData)

	// 监听恢复对局的回复消息
	resumeRspExpect, _ := client.ExpectMessage(msgid.MsgID_ROOM_DDZ_RESUME_RSP)
	assert.NotNil(t, resumeRspExpect)

	// 发出恢复对局请求
	_, resumeErr := client.SendPackage(utils.CreateMsgHead(msgid.MsgID_ROOM_DDZ_RESUME_REQ), &room.DDZResumeGameReq{})
	assert.Nil(t, resumeErr)

	resumeRsp := room.DDZResumeGameRsp{}
	if err := resumeRspExpect.Recv(global.DefaultWaitMessageTime, &resumeRsp); err != nil {
		logEntry.Errorf("玩家%d没有收到恢复对局的回复", player.GetID())
		assert.NotNil(t, nil)
		return
	}

	// 打印恢复对局回复的信息
	logEntry.Infof("玩家%d收到恢复对局的回复,resultCode = %d， resultStr = %s", resumeRsp.GetResult().GetErrCode(), resumeRsp.GetResult().GetErrDesc())

	// 成功时打印游戏信息
	if resumeRsp.GetResult().GetErrCode() == 0 {
		ddzDeskInfo := resumeRsp.GetGameInfo()
		// 桌子里面的每一个玩家
		for i := 0; i < len(ddzDeskInfo.GetPlayers()); i++ {
			ddzPlayrInfo := ddzDeskInfo.GetPlayers()[i]
			roomPlayerInfo := ddzPlayrInfo.GetPlayerInfo()
			logEntry.Infof("玩家:%d，名字:%s，金币数:%d，座位号:%d，打出的牌:%v，手中的牌:%v，是否为地主:%v，抢地主类型:%v，加倍类型:%v，手牌数量：%v",
				roomPlayerInfo.GetPlayerId(), roomPlayerInfo.GetName(), roomPlayerInfo.GetCoin(), roomPlayerInfo.GetSeat(),
				ddzPlayrInfo.GetOutCards(), ddzPlayrInfo.GetHandCards(), ddzPlayrInfo.GetLord(), ddzPlayrInfo.GetGrabLord(), ddzPlayrInfo.GetDouble(), ddzPlayrInfo.GetHandCardsCount())
		}

		// 当前状态
		logEntry.Infof("当前状态：%v，进入下一状态的等待时间:%d, 当前操作的玩家:%v, 底牌：%v, 抢庄倍数：%v, 加倍倍数：%v, 炸弹倍数：%v",
			ddzDeskInfo.GetStage().GetStage(), ddzDeskInfo.GetStage().GetTime(), ddzDeskInfo.GetCurPlayerId(),
			ddzDeskInfo.GetDipai(), ddzDeskInfo.GetTotalGrab(), ddzDeskInfo.GetTotalDouble(), ddzDeskInfo.GetTotalBomb())
	}
}
