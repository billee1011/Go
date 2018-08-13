package proptests

import (
	"steve/client_pb/common"
	"steve/client_pb/hall"
	"steve/client_pb/msgid"
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"steve/simulate/utils/doudizhu"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProp(t *testing.T) {
	// 配牌1
	params := doudizhu.NewStartDDZGameParamsTest1()
	params.PlayerSeatGold = map[int]uint64{0: 100000, 1: 100000, 2: 100000}
	seat := 0
	deskData, err := utils.StartDDZGame(params)
	assert.NotNil(t, deskData)
	assert.Nil(t, err)

	toPlayer := utils.GetDeskPlayerBySeat((seat+1)%len(deskData.Players), deskData)
	toPlayerID := toPlayer.Player.GetID()
	err = utils.SendGetPlayerGameInfoReq(seat, deskData, toPlayerID, params.GameID)
	assert.Nil(t, err)
	player := utils.GetDeskPlayerBySeat(seat, deskData)
	expector, _ := player.Expectors[msgid.MsgID_HALL_GET_PLAYER_GAME_INFO_RSP]
	rsp := hall.HallGetPlayerGameInfoRsp{}
	assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, &rsp))
	assert.NotEqual(t, 0, len(rsp.UserProperty))

	for _, prop := range rsp.UserProperty {
		err = utils.SendUsePropReq(seat, deskData, toPlayerID, common.PropType(*prop.PropId))
		assert.Nil(t, err)
		if common.PropType(*prop.PropId) == common.PropType_EGG_GUN {
			expector, _ = player.Expectors[msgid.MsgID_ROOM_USE_PROP_RSP]
			rsp1 := room.RoomUsePropRsp{}
			assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, &rsp1))
			assert.Equal(t, room.RoomError_FAILED, *rsp.ErrCode)
			break
		}
		for _, playert := range deskData.Players {
			expector, _ := playert.Expectors[msgid.MsgID_ROOM_USE_PROP_NTF]
			ntf := room.RoomUsePropNtf{}
			assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, &ntf))
			assert.Equal(t, player.Player.GetID(), ntf.FromPlayerId)
			assert.Equal(t, toPlayerID, ntf.ToPlayerId)
			assert.Equal(t, common.PropType(*prop.PropId), ntf.PropId)
		}
	}
}
