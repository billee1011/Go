package misctests

import (
	"steve/client_pb/common"
	"steve/client_pb/hall"
	"steve/client_pb/msgid"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

// Test_RealName 实名测试
func Test_RealName(t *testing.T) {
	player, err := utils.LoginNewPlayer()
	assert.Nil(t, err)
	assert.NotNil(t, player)

	rsp := hall.HallRealNameRsp{}

	err = player.GetClient().Request(utils.CreateMsgHead(msgid.MsgID_HALL_REAL_NAME_REQ), &hall.HallRealNameReq{
		Name:   proto.String("安佳玮"),
		IdCard: proto.String("410322199202152910"),
	}, global.DefaultWaitMessageTime, uint32(msgid.MsgID_HALL_REAL_NAME_RSP), &rsp)
	assert.Nil(t, err, "请求认证失败")
	assert.Equal(t, uint32(common.ErrCode_EC_SUCCESS), rsp.GetErrCode())
	assert.Equal(t, uint64(5000), rsp.GetCoinReward())
	assert.True(t, rsp.GetNewCoin() >= uint64(5000))

	// 获取玩家信息时已经认证
	getPlayerInfoRsp := hall.HallGetPlayerInfoRsp{}
	err = player.GetClient().Request(utils.CreateMsgHead(msgid.MsgID_HALL_GET_PLAYER_INFO_REQ),
		&hall.HallGetPlayerInfoReq{},
		global.DefaultWaitMessageTime,
		uint32(msgid.MsgID_HALL_GET_PLAYER_INFO_RSP),
		&getPlayerInfoRsp,
	)
	assert.Nil(t, err, "获取玩家信息失败")
	assert.Equal(t, uint32(common.ErrCode_EC_SUCCESS), getPlayerInfoRsp.GetErrCode())
	assert.Equal(t, uint32(1), getPlayerInfoRsp.GetRealnameStatus())
}
