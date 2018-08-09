package registers

import (
	"steve/client_pb/msgid"
	"steve/client_pb/room"
	propclient "steve/common/data/prop"
	"steve/entity/constant"
	"steve/external/goldclient"
	"steve/gutils"
	modelmanager "steve/room/models"
	player2 "steve/room/player"
	"steve/room/util"
	"steve/structs/exchanger"
	"steve/structs/proto/gate_rpc"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

func HandleRoomChatReq(playerID uint64, header *steve_proto_gaterpc.Header, req room.RoomDeskChatReq) (ret []exchanger.ResponseMsg) {
	player := player2.GetPlayerMgr().GetPlayer(playerID)
	if player == nil {
		return
	}
	modelmanager.GetModelManager().GetChatModel(player.GetDesk().GetUid()).RoomChatMsgReq(player, header, req)
	return
}

// HandleRoomDeskQuitReq 处理玩家退出桌面请求
// 失败先不回复
func HandleRoomDeskQuitReq(playerID uint64, header *steve_proto_gaterpc.Header, req room.RoomDeskQuitReq) (rspMsg []exchanger.ResponseMsg) {
	response := room.RoomDeskQuitRsp{
		UserData: proto.Uint32(req.GetUserData()),
		ErrCode:  room.RoomError_SUCCESS.Enum(),
	}
	// 退出暂时总是回复成功
	// TODO : 血战的换对手需要等退出完成后才返回
	util.SendMessageToPlayer(playerID, msgid.MsgID_ROOM_DESK_QUIT_RSP, &response)

	player := player2.GetPlayerMgr().GetPlayer(playerID)
	if player == nil {
		return
	}
	desk := player.GetDesk()
	if desk == nil {
		return
	}
	modelmanager.GetModelManager().GetPlayerModel(desk.GetUid()).PlayerQuit(player)
	return
}

func noGamePlaying() []exchanger.ResponseMsg {
	body := &room.RoomResumeGameRsp{
		ResumeRes: room.RoomError_DESK_NO_GAME_PLAYING.Enum(),
	}
	return []exchanger.ResponseMsg{
		exchanger.ResponseMsg{
			MsgID: uint32(msgid.MsgID_ROOM_RESUME_GAME_RSP),
			Body:  body,
		},
	}
}

// HandleResumeGameReq 恢复对局请求
func HandleResumeGameReq(playerID uint64, header *steve_proto_gaterpc.Header, req room.RoomResumeGameReq) (ret []exchanger.ResponseMsg) {
	entry := logrus.WithField("player_id", playerID)
	player := player2.GetPlayerMgr().GetPlayer(playerID)
	if player == nil {
		entry.Debugln("玩家不存在")
		return noGamePlaying()
	}
	desk := player.GetDesk()
	if desk == nil {
		entry.Debugln("没有对应的牌桌")
		return noGamePlaying()
	}
	modelmanager.GetModelManager().GetPlayerModel(desk.GetUid()).PlayerEnter(player)
	return
}

// HandleCancelTuoGuanReq 处理取消托管请求
func HandleCancelTuoGuanReq(playerID uint64, header *steve_proto_gaterpc.Header, req room.RoomCancelTuoGuanReq) (ret []exchanger.ResponseMsg) {
	ret = []exchanger.ResponseMsg{}

	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "HandleCancelTuoGuanReq",
		"player_id": playerID,
	})
	player := player2.GetPlayerMgr().GetPlayer(playerID)
	if player == nil {
		logEntry.Debugln("获取玩家失败")
		return
	}
	desk := player.GetDesk()
	if desk == nil {
		logEntry.Debugln("玩家不在房间中")
		return
	}
	player.SetTuoguan(false, true)
	return
}

// HandleTuoGuanReq 处理取消托管请求
func HandleTuoGuanReq(playerID uint64, header *steve_proto_gaterpc.Header, req room.RoomTuoGuanReq) (ret []exchanger.ResponseMsg) {
	ret = []exchanger.ResponseMsg{}

	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "HandleTuoGuanReq",
		"player_id": playerID,
	})
	player := player2.GetPlayerMgr().GetPlayer(playerID)
	if player == nil {
		logEntry.Debugln("获取玩家失败")
		return
	}
	desk := player.GetDesk()
	if desk == nil {
		logEntry.Debugln("玩家不在房间中")
		return
	}
	player.SetTuoguan(req.GetTuoguan(), true)
	return
}

// HandleUsePropReq 使用道具请求处理
func HandleUsePropReq(playerID uint64, header *steve_proto_gaterpc.Header, req room.RoomUsePropReq) (ret []exchanger.ResponseMsg) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "HandleUsePropReq",
		"player_id": playerID,
	})

	errDesc := "不能使用该道具"
	rsp := room.RoomUsePropRsp{
		ErrCode: room.RoomError_FAILED.Enum(),
		ErrDesc: &errDesc,
	}
	ret = []exchanger.ResponseMsg{
		exchanger.ResponseMsg{
			MsgID: uint32(msgid.MsgID_ROOM_USE_PROP_RSP),
			Body:  &rsp,
		},
	}

	player := player2.GetPlayerMgr().GetPlayer(playerID)
	if player == nil {
		logEntry.Debugln("获取玩家失败")
		return
	}
	desk := player.GetDesk()
	if desk == nil {
		logEntry.Debugln("玩家不在房间中")
		return
	}

	propID := gutils.PropTypeClient2Server(req.GetPropId())
	prop, err := propclient.GetPlayerOneProp(playerID, propID)
	if err != nil {
		return
	}

	// 使用道具
	if prop.Count > 0 {
		err = propclient.AddPlayerProp(playerID, propID, -1)
		if err != nil {
			return
		}
	} else { // 使用金币购买
		propConfig, err := propclient.GetOnePropsConfig(propID)
		if err != nil {
			return
		}

		coin, err := goldclient.GetGold(playerID, constant.GOLD_COIN)
		if coin > propConfig.Limit {
			goldclient.AddGold(playerID, constant.GOLD_COIN, propConfig.Value, 0, 0, int32(desk.GetGameId()), desk.GetLevel())
		} else {
			return
		}
	}

	// 广播道具
	desPlayerID := req.GetPlayerId()
	ntf := room.RoomUsePropNtf{
		FromPlayerId: &playerID,
		ToPlayerId:   &desPlayerID,
		PropId:       req.PropId,
	}
	msgBody, err := proto.Marshal(&ntf)
	if err != nil {
		logEntry.WithError(err).Debugln("序列化失败")
		return
	}
	modelmanager.GetModelManager().GetMessageModel(playerID).BroadcastMessage([]uint64{}, msgid.MsgID_ROOM_USE_PROP_NTF, msgBody, true)
	return nil
}
