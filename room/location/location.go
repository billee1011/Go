package location

import (
	"steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/room/interfaces"
	"steve/room/interfaces/global"
	"steve/structs/exchanger"
	"steve/structs/proto/gate_rpc"

	"github.com/golang/protobuf/proto"
)

// RoomPlayerLocationReq 处理地理位置请求
func RoomPlayerLocationReq(clientID uint64, header *steve_proto_gaterpc.Header, req room.RoomPlayerLocationReq) (ret []exchanger.ResponseMsg) {
	playerMgr := global.GetPlayerMgr()
	player := playerMgr.GetPlayerByClientID(clientID)
	deskMgr := global.GetDeskMgr()
	desk, err := deskMgr.GetRunDeskByPlayerID(player.GetID())
	if err != nil {
		return
	}
	deskPlayers := desk.GetDeskPlayers()
	playerLocations := getPlayerLocations(deskPlayers)

	rsp := &room.RoomPlayerLocationRsp{
		Locations: playerLocations,
	}
	ret = []exchanger.ResponseMsg{
		exchanger.ResponseMsg{
			MsgID: uint32(msgid.MsgID_ROOM_PLAYER_LOCATION_RSP),
			Body:  rsp,
		},
	}
	return
}

func getPlayerLocations(deskPlayers []interfaces.DeskPlayer) []*room.PlayerLocation {
	pls := make([]*room.PlayerLocation, 0)
	for _, deskPlayer := range deskPlayers {
		pl := &room.PlayerLocation{
			PlayerId: proto.Uint64(deskPlayer.GetPlayerID()),
			Location: deskPlayer.GetLocationInfos(),
		}

		pls = append(pls, pl)
	}
	return pls
}
