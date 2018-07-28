package user

import (
	"context"
	"fmt"
	"steve/common/data/player"
	"steve/entity/db"
	"steve/hall/data"
	"steve/server_pb/user"
	"time"

	"github.com/Sirupsen/logrus"
)

type playerDataService struct{}

// Default 默认服务
var Default user.PlayerDataServer = new(playerDataService)

// GetPlayerByAccount 根据账号获取玩家 ID
func (pds *playerDataService) GetPlayerByAccount(ctx context.Context, req *user.GetPlayerByAccountReq) (rsp *user.GetPlayerByAccountRsp, err error) {
	rsp, err = &user.GetPlayerByAccountRsp{
		ErrCode: int32(user.ErrCode_EC_FAIL),
	}, nil
	accID := req.GetAccountId()
	exist, playerID, err := data.GetPlayerIDByAccountID(accID)
	if exist && err == nil {
		rsp.PlayerId, rsp.ErrCode = playerID, int32(user.ErrCode_EC_SUCCESS)
		return
	}
	var err2 error
	playerID, err2 = createPlayer(accID)
	if err2 != nil {
		logrus.WithField("account_id", accID).Errorln(err2)
		return
	}
	rsp.PlayerId, rsp.ErrCode = playerID, int32(user.ErrCode_EC_SUCCESS)
	return
}

// createPlayer 创建玩家
func createPlayer(accID uint64) (uint64, error) {
	playerID := data.AllocPlayerID()

	if playerID == 0 {
		return 0, fmt.Errorf("分配玩家 ID 失败")
	}
	// TODO: 使用正式的金币服
	player.SetPlayerCoin(playerID, 10*10000)

	if err := data.InitPlayerData(db.TPlayer{
		Accountid:    int64(accID),
		Playerid:     int64(playerID),
		Type:         1,
		Channelid:    0,                                 // TODO ，渠道 ID
		Nickname:     fmt.Sprintf("player%d", playerID), // TODO,昵称
		Gender:       1,
		Avatar:       "", // TODO , 头像
		Provinceid:   0,  // TODO， 省ID
		Cityid:       0,  // TODO 市ID
		Name:         "", // TODO: 真实姓名
		Phone:        "", // TODO: 电话
		Idcard:       "", // TODO 身份证
		Iswhitelist:  0,
		Zipcode:      0,
		Shippingaddr: "",
		Status:       1,
		Remark:       "",
		Createtime:   time.Now(),
		Createby:     "",
		Updatetime:   time.Now(),
		Updateby:     "",
	}); err != nil {
		return 0, fmt.Errorf("初始化玩家(%d)数据失败: %v", playerID, err)
	}
	return playerID, nil
}
