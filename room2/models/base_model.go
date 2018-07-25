package models

import (
	"steve/room2/desk"
	"steve/room2/player"
	"steve/client_pb/room"
	"github.com/golang/protobuf/proto"
)

type BaseModel struct {
	desk  *desk.Desk
}

func (model *BaseModel) GetDesk() *desk.Desk {
	return model.desk
}

func (model *BaseModel) SetDesk(desk *desk.Desk) {
	model.desk = desk
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