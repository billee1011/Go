package models

import (
	"context"
	client_match_pb "steve/client_pb/match"
	"steve/client_pb/msgid"
	"steve/room/desk"
	"steve/room/fixed"
	playerpkg "steve/room/player"
	"steve/server_pb/match"
	"steve/structs"

	"github.com/Sirupsen/logrus"
)

// ContinueModel 续局 model
type ContinueModel struct {
	BaseModel
}

// NewContinueModel 创建续局 model
func NewContinueModel(desk *desk.Desk) DeskModel {
	result := &ContinueModel{}
	result.SetDesk(desk)
	return result
}

// GetName 获取 model 名称
func (model *ContinueModel) GetName() string {
	return fixed.ChatModelName
}

// Active 激活 model
func (model *ContinueModel) Active() {}

// Start 启动 model
func (model *ContinueModel) Start() {

}

// Stop 停止 model
func (model *ContinueModel) Stop() {
}

// ContinueDesk 开始续局逻辑
func (model *ContinueModel) ContinueDesk(fixBanker bool, bankerSeat int, settleMap map[uint64]int64) {
	entry := logrus.WithFields(logrus.Fields{
		"func_name":   "DeskBase.ContinueDesk",
		"fix_banker":  fixBanker,
		"banker_seat": bankerSeat,
	})
	desk := model.GetDesk()
	playerIDs := desk.GetPlayerIds()
	continuePlayers := make([]*match.ContinuePlayer, 0, len(playerIDs))
	playerMgr := playerpkg.GetPlayerMgr()
	for _, playerID := range playerIDs {
		player := playerMgr.GetPlayer(playerID)
		if player == nil || player.GetDesk() != desk || player.IsQuit() { // 玩家已经退出牌桌，不续局
			entry.WithFields(logrus.Fields{
				"player_id": playerID,
				"quited":    player.IsQuit(),
			}).Debugln("玩家不满足续局条件")

			messageModel := GetModelManager().GetMessageModel(desk.GetUid())
			messageModel.BroadCastDeskMessage(nil, msgid.MsgID_MATCH_CONTINUE_DESK_DIMISS_NTF,
				&client_match_pb.MatchContinueDeskDimissNtf{}, true)
			return
		}
		continuePlayers = append(continuePlayers, &match.ContinuePlayer{
			PlayerId:   playerID,
			Seat:       int32(player.GetSeat()),
			Win:        settleMap[playerID] >= 0,
			RobotLevel: int32(player.GetRobotLv()),
		})
	}

	request := match.AddContinueDeskReq{
		Players:    continuePlayers,
		GameId:     int32(desk.GetGameId()),
		FixBanker:  fixBanker,
		BankerSeat: int32(bankerSeat),
	}
	exposer := structs.GetGlobalExposer()
	cc, err := exposer.RPCClient.GetConnectByServerName("match")
	if err != nil {
		entry.WithError(err).Errorln("获取 match 连接失败")
		return
	}
	mcc := match.NewMatchClient(cc)
	_, err = mcc.AddContinueDesk(context.Background(), &request)
	if err != nil {
		entry.WithError(err).Errorln("请求失败")
		return
	}
	entry.Debugln("添加续局牌桌完成")
}
