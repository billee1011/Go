package models

import (
	"steve/client_pb/room"
	"steve/room2/desk"
	"steve/room2/player"

	"github.com/golang/protobuf/proto"
)

type BaseModel struct {
	desk *desk.Desk
}

func (model *BaseModel) GetDesk() *desk.Desk {
	return model.desk
}

func (model *BaseModel) SetDesk(desk *desk.Desk) {
	model.desk = desk
}

func (model *BaseModel) GetGameContext() interface{} {
	return model.desk.GetConfig().Context
}

// TranslateToRoomPlayer 将 deskPlayer 转换成 RoomPlayerInfo
func TranslateToRoomPlayer(player *player.Player) room.RoomPlayerInfo {
	return room.RoomPlayerInfo{
		PlayerId: proto.Uint64(player.GetPlayerID()),
		Name:     proto.String(""), // TODO
		Coin:     proto.Uint64(player.GetCoin()),
		Seat:     proto.Uint32(uint32(player.GetSeat())),
		// Location: TODO 没地方拿
	}
}
