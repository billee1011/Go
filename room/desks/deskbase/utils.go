package deskbase

import (
	"steve/client_pb/room"
	"steve/room/interfaces"
	"steve/room/interfaces/global"

	"github.com/golang/protobuf/proto"
)

// TranslateToRoomPlayer 将 deskPlayer 转换成 RoomPlayerInfo
func TranslateToRoomPlayer(deskPlayer interfaces.DeskPlayer) room.RoomPlayerInfo {
	playerMgr := global.GetPlayerMgr()
	playerID := deskPlayer.GetPlayerID()
	player := playerMgr.GetPlayer(playerID)
	var coin uint64
	if player != nil {
		coin = player.GetCoin()
	}
	return room.RoomPlayerInfo{
		PlayerId: proto.Uint64(playerID),
		Name:     proto.String(""), // TODO
		Coin:     proto.Uint64(coin),
		Seat:     proto.Uint32(uint32(deskPlayer.GetSeat())),
		// Location: TODO 没地方拿
	}
}
