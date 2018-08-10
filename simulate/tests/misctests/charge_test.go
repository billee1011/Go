package misctests

import (
	"fmt"
	"steve/client_pb/common"
	"steve/client_pb/hall"
	"steve/client_pb/msgid"
	"steve/simulate/facade"
	"steve/simulate/global"
	"steve/simulate/interfaces"
	"steve/simulate/utils"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func getChargeInfo(player interfaces.ClientPlayer) (chargeInfoRsp *hall.GetChargeInfoRsp, err error) {
	request := &hall.GetChargeInfoReq{
		Platform: common.Platform_Android.Enum(),
	}
	response := hall.GetChargeInfoRsp{}
	if err := player.GetClient().Request(utils.CreateMsgHead(msgid.MsgID_GET_CHARGE_INFO_REQ), request,
		global.DefaultWaitMessageTime, uint32(msgid.MsgID_GET_CHARGE_INFO_RSP), &response); err != nil {
		return nil, fmt.Errorf("请求获取充值信息失败:%v", err)
	}
	return &response, err
}

func Test_GetChargeInfo(t *testing.T) {
	player, err := utils.LoginNewPlayer()
	assert.Nil(t, err)
	assert.NotNil(t, player)

	response, err := getChargeInfo(player)
	assert.Nil(t, err)
	assert.Equal(t, response.GetResult().GetErrCode(), common.ErrCode_EC_SUCCESS)
	assert.NotZero(t, len(response.GetItems()))
	assert.Zero(t, response.GetTodayCharge())
	assert.NotZero(t, response.GetDayMaxCharge())
}

func Test_Charge(t *testing.T) {
	player, err := utils.LoginNewPlayer()
	assert.Nil(t, err)
	assert.NotNil(t, player)
	oldcoin := player.GetCoin()

	chargeInfoRsp, err := getChargeInfo(player)
	assert.Nil(t, err)
	assert.Equal(t, chargeInfoRsp.GetResult().GetErrCode(), common.ErrCode_EC_SUCCESS)
	item := chargeInfoRsp.GetItems()[0]

	chargeReq := &hall.ChargeReq{
		ItemId:   proto.Uint64(item.GetItemId()),
		Cost:     proto.Uint64(item.GetPrice()),
		Platform: common.Platform_Android.Enum(),
	}
	chargeRsp := &hall.ChargeRsp{}
	assert.Nil(t, facade.Request(player.GetClient(), msgid.MsgID_CHARGE_REQ, chargeReq, global.DefaultWaitMessageTime,
		msgid.MsgID_CHARGE_RSP, chargeRsp))
	// 校验充值到账
	assert.Equal(t, chargeRsp.GetResult().GetErrCode(), common.ErrCode_EC_SUCCESS, chargeRsp.GetResult().GetErrDesc())
	assert.Equal(t, item.GetCoin()+item.GetPresentCoin(), chargeRsp.GetObtainedCoin())
	assert.Equal(t, oldcoin+chargeRsp.GetObtainedCoin(), chargeRsp.GetNewCoin())
}
