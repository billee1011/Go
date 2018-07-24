package deskbase

import (
	"context"
	client_match_pb "steve/client_pb/match"
	"steve/client_pb/msgid"
	"steve/room/interfaces"
	"steve/room/interfaces/facade"
	"steve/server_pb/match"
	"steve/structs"

	"github.com/Sirupsen/logrus"
)

// DeskBase 房间基类
type DeskBase struct {
	interfaces.DeskPlayerMgr
	uid    uint64
	gameID int
}

// NewDeskBase 创建房间基类对象
func NewDeskBase(uid uint64, gameID int, deskPlayers []interfaces.DeskPlayer) *DeskBase {
	deskPlayerMgr := createDeskPlayerMgr()
	deskPlayerMgr.setPlayers(deskPlayers)
	return &DeskBase{
		DeskPlayerMgr: deskPlayerMgr,
		uid:           uid,
		gameID:        gameID,
	}
}

// GetUID 获取牌桌 UID
func (d *DeskBase) GetUID() uint64 {
	return d.uid
}

// GetGameID 获取游戏 ID
func (d *DeskBase) GetGameID() int {
	return d.gameID
}

func (d *DeskBase) isWinner(playerID uint64, winners []uint64) bool {
	if winners == nil {
		return false
	}
	for _, winner := range winners {
		if winner == playerID {
			return true
		}
	}
	return false
}

// ContinueDesk 续局牌桌
func (d *DeskBase) ContinueDesk(fixBanker bool, bankerSeat int, winners []uint64) {
	entry := logrus.WithFields(logrus.Fields{
		"func_name":   "DeskBase.ContinueDesk",
		"fix_banker":  fixBanker,
		"banker_seat": bankerSeat,
	})
	players := d.GetDeskPlayers()
	continuePlayers := make([]*match.ContinuePlayer, 0, len(players))
	for _, player := range players {
		playerID := player.GetPlayerID()
		if player.IsQuit() { // 玩家已经退出牌桌或者 玩家金币数为0，不续局
			entry.WithFields(logrus.Fields{
				"player_id": playerID,
				"quited":    player.IsQuit(),
			}).Debugln("玩家不满足续局条件")
			facade.BroadCastDeskMessage(d, nil, msgid.MsgID_MATCH_CONTINUE_DESK_DIMISS_NTF, &client_match_pb.MatchContinueDeskDimissNtf{}, true)
			return
		}
		continuePlayers = append(continuePlayers, &match.ContinuePlayer{
			PlayerId:   playerID,
			Seat:       int32(player.GetSeat()),
			Win:        d.isWinner(playerID, winners),
			RobotLevel: int32(player.GetRobotLv()),
		})
	}

	request := match.AddContinueDeskReq{
		Players:    continuePlayers,
		GameId:     int32(d.gameID),
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
