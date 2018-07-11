package user

import (
	"steve/client_pb/common"
	"steve/client_pb/hall"
	"steve/client_pb/msgid"
	"steve/common/data/player"
	"steve/structs/exchanger"
	"steve/structs/proto/gate_rpc"

	"github.com/golang/protobuf/proto"
)

// getPlayerState 获取玩家状态
func getPlayerState(playerID uint64) common.PlayerState {
	if player.GetPlayerRoomAddr(playerID) == "" {
		return common.PlayerState_PS_IDLE
	}
	return common.PlayerState_PS_GAMEING
}

// HandleGetPlayerInfoReq 处理获取玩家信息请求
func HandleGetPlayerInfoReq(playerID uint64, header *steve_proto_gaterpc.Header, req hall.HallGetPlayerInfoReq) (rspMsg []exchanger.ResponseMsg) {
	response := &hall.HallGetPlayerInfoRsp{}
	rspMsg = []exchanger.ResponseMsg{
		exchanger.ResponseMsg{
			MsgID: uint32(msgid.MsgID_HALL_GET_PLAYER_INFO_RSP),
			Body:  response,
		},
	}
	response.ErrCode = proto.Uint32(0)
	response.Coin = proto.Uint64(player.GetPlayerCoin(playerID))
	response.NickName = proto.String(player.GetPlayerNickName(playerID))
	response.PlayerState = getPlayerState(playerID).Enum()
	return
}

// HandleGetPlayerStateReq 获取玩家是否正在游戏中
func HandleGetPlayerStateReq(playerID uint64, header *steve_proto_gaterpc.Header, req hall.HallGetPlayerStateReq) (rspMsg []exchanger.ResponseMsg) {
	response := &hall.HallGetPlayerStateRsp{}
	rspMsg = []exchanger.ResponseMsg{
		exchanger.ResponseMsg{
			MsgID: uint32(msgid.MsgID_HALL_GET_PLAYER_STATE_RSP),
			Body:  response,
		},
	}
	response.PlayerState = getPlayerState(playerID).Enum()
	return
}
