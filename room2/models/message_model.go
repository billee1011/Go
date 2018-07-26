package models

import (
	"steve/client_pb/msgid"
	"github.com/golang/protobuf/proto"
	"fmt"
	"github.com/Sirupsen/logrus"
	"steve/room2/util"
	"steve/room2/player"
	"steve/room2/fixed"
	"steve/room2/desk"
)

var gMessageSender util.MessageSender

type MessageModel struct {
	BaseModel
}

func (model MessageModel) GetName() string {
	return fixed.Message
}
func (model MessageModel) Start() {

}
func (model MessageModel) Stop() {

}

func NewMessageModel(desk *desk.Desk) DeskModel {
	result := &MessageModel{}
	result.SetDesk(desk)
	return result
}

// BroadCastDeskMessage 广播消息给牌卓玩家
func (model *MessageModel) BroadCastDeskMessage(playerIDs []uint64, msgID msgid.MsgID, body proto.Message, exceptQuit bool) error {
	msgBody, err := proto.Marshal(body)
	if err != nil {
		return err
	}
	model.BroadcastMessage(playerIDs, msgID, msgBody, exceptQuit)
	return nil
}

func find(datas []uint64, data uint64) bool {
	for _, d := range datas {
		if d == data {
			return true
		}
	}
	return false
}

func (model *MessageModel) BroadcastMessage(playerIDs []uint64, msgID msgid.MsgID, body []byte, exceptQuit bool) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name":       "deskPlayerMgr.BroadcastMessage",
		"dest_player_ids": playerIDs,
		"msg_id":          msgID,
	})
	// 是否针对所有玩家
	if playerIDs == nil || len(playerIDs) == 0 {
		deskId := model.GetDesk().GetUid()
		m:=GetModelManager().GetPlayerModel(deskId)
		playerIDs = m.GetDeskPlayerIDs()
		logEntry = logEntry.WithField("all_player_ids", playerIDs)
	}
	playerIDs = model.removeQuit(playerIDs)
	logEntry = logEntry.WithField("real_dest_player_ids", playerIDs)

	if len(playerIDs) == 0 {
		return
	}
	util.BroadCastMessageBare(playerIDs, msgID, body)
	logEntry.Debugln("广播消息")
}

func (model *MessageModel) removeQuit(playerIDs []uint64) []uint64 {
	result := []uint64{}
	for _, playerID := range playerIDs {
		if quited := player.GetPlayerMgr().GetPlayer(playerID).IsQuit(); !quited {
			result = append(result, playerID)
		}
	}
	return result
}

// BroadCastDeskMessageExcept 广播消息给牌桌玩家
func  (model *MessageModel) BroadCastDeskMessageExcept(expcetPlayers []uint64, exceptQuit bool, msgID msgid.MsgID, body proto.Message) error {
	playerIDs := []uint64{}
	deskPlayers := GetModelManager().GetPlayerModel(model.GetDesk().GetUid()).GetDeskPlayers()
	for _, deskPlayer := range deskPlayers {
		playerID := deskPlayer.GetPlayerID()
		if find(expcetPlayers, playerID) {
			continue
		}
		playerIDs = append(playerIDs, playerID)
	}
	if len(playerIDs) == 0 {
		return fmt.Errorf("没有广播玩家")
	}
	err := model.BroadCastDeskMessage(playerIDs, msgID, body, exceptQuit)
	return err
}
