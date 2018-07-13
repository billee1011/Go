package desks

import (
	"steve/client_pb/msgid"
	"steve/client_pb/room"
	"steve/room/interfaces"
	"steve/room/interfaces/global"
	"steve/structs/exchanger"
	"steve/structs/proto/gate_rpc"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func setupPlayerMgr(mc *gomock.Controller) {
	pm := interfaces.NewMockPlayerMgr(mc)
	global.SetPlayerMgr(pm)

	pm.EXPECT().GetPlayerByClientID(gomock.Any()).DoAndReturn(func(clientID uint64) interfaces.Player {
		mockPlayer := interfaces.NewMockPlayer(mc)
		mockPlayer.EXPECT().GetID().Return(clientID).AnyTimes()
		return mockPlayer
	}).AnyTimes()

	pm.EXPECT().GetPlayer(gomock.Any()).DoAndReturn(func(playerID uint64) interfaces.Player {
		mockPlayer := interfaces.NewMockPlayer(mc)
		mockPlayer.EXPECT().GetClientID().Return(playerID).AnyTimes()
		mockPlayer.EXPECT().GetID().Return(playerID).AnyTimes()
		return mockPlayer
	}).AnyTimes()
}

func setup(mc *gomock.Controller) {
	// gJoinApplyMgr = newApplyMgr(true)
	setupPlayerMgr(mc)
}

func apply(clientID uint64) []exchanger.ResponseMsg {
	header := steve_proto_gaterpc.Header{MsgId: uint32(msgid.MsgID_ROOM_JOIN_DESK_REQ)}
	req := room.RoomJoinDeskReq{}

	return HandleRoomJoinDeskReq(clientID, &header, req)
}

// TestHandleRoomJoinDeskReq 测试正常情况申请加入
func TestHandleRoomJoinDeskReq(t *testing.T) {
	mc := gomock.NewController(t)
	defer mc.Finish()
	setup(mc)

	var clientID uint64 = 10
	rspMsgs := apply(clientID)

	assert.NotNil(t, rspMsgs)
	assert.Equal(t, 1, len(rspMsgs))
	assert.Equal(t, uint32(msgid.MsgID_ROOM_JOIN_DESK_RSP), rspMsgs[0].MsgID)

	rspBody, ok := rspMsgs[0].Body.(*room.RoomJoinDeskRsp)
	assert.True(t, ok)
	assert.Equal(t, room.RoomError_SUCCESS, rspBody.GetErrCode())
}

// TestHandleRoomJoinDeskReq_NotLogin 测试用户未登录
func TestHandleRoomJoinDeskReq_NotLogin(t *testing.T) {
	mc := gomock.NewController(t)
	defer mc.Finish()
	setup(mc)

	pm := interfaces.NewMockPlayerMgr(mc)
	global.SetPlayerMgr(pm)
	pm.EXPECT().GetPlayerByClientID(gomock.Any()).Return(nil).AnyTimes()

	rspMsgs := apply(10)

	assert.NotNil(t, rspMsgs)
	assert.Equal(t, 1, len(rspMsgs))
	assert.Equal(t, uint32(msgid.MsgID_ROOM_JOIN_DESK_RSP), rspMsgs[0].MsgID)

	rspBody, ok := rspMsgs[0].Body.(*room.RoomJoinDeskRsp)
	assert.True(t, ok)
	assert.Equal(t, room.RoomError_NOT_LOGIN, rspBody.GetErrCode())
}

func Test_joinApplyManager_checkMatch(t *testing.T) {
	groupCount := 6

	mc := gomock.NewController(t)
	defer mc.Finish()
	setupPlayerMgr(mc)

	deskFactory := interfaces.NewMockDeskFactory(mc)
	global.SetDeskFactory(deskFactory)

	deskFactory.EXPECT().CreateDesk(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(players []uint64, gameID int, opt interfaces.CreateDeskOptions) (interfaces.CreateDeskResult, error) {
			desk := interfaces.NewMockDesk(mc)

			roomPlayers := []*room.RoomPlayerInfo{}
			for _, playerID := range players {
				roomPlayers = append(roomPlayers, &room.RoomPlayerInfo{
					PlayerId: proto.Uint64(playerID),
				})
			}
			desk.EXPECT().GetPlayers().Return(roomPlayers).AnyTimes()

			return interfaces.CreateDeskResult{
				Desk: desk,
			}, nil
		}).AnyTimes()

	messageSender := interfaces.NewMockMessageSender(mc)
	global.SetMessageSender(messageSender)

	deskMgr := interfaces.NewMockDeskMgr(mc)
	deskMgr.EXPECT().RunDesk(gomock.Any()).Times(groupCount)
	global.SetDeskMgr(deskMgr)

	for i := 0; i < groupCount; i++ {
		clientIDs := []uint64{}
		clientIDs = append(clientIDs, uint64(i*4+1))
		clientIDs = append(clientIDs, uint64(i*4+2))
		clientIDs = append(clientIDs, uint64(i*4+3))
		clientIDs = append(clientIDs, uint64(i*4+4))
		messageSender.EXPECT().BroadcastPackage(clientIDs, gomock.Any(), gomock.Any()).Times(1)
	}

	jam := newApplyMgr(false)
	jam.applyChannel = make(chan uint64, 1024)

	go func() {
		for i := 0; i < groupCount*4; i++ {
			jam.applyChannel <- uint64(i + 1)
		}
		close(jam.applyChannel)
	}()
	jam.checkMatch()
}
