package user

import (
	"context"
	"fmt"
	"steve/entity/cache"
	"steve/entity/db"
	"steve/hall/data"
	"steve/server_pb/user"
	"strconv"
	"time"

	"github.com/Sirupsen/logrus"
)

// PlayerDataService 实现 user.PlayerServer
type PlayerDataService struct{}

var _ user.PlayerDataServer = new(PlayerDataService)

// GetPlayerByAccount 根据账号获取玩家 ID
func (pds *PlayerDataService) GetPlayerByAccount(ctx context.Context, req *user.GetPlayerByAccountReq) (rsp *user.GetPlayerByAccountRsp, err error) {
	logrus.Debugln("GetPlayerByAccount req", *req)

	// 默认返回消息
	rsp, err = &user.GetPlayerByAccountRsp{
		ErrCode: int32(user.ErrCode_EC_FAIL),
	}, nil

	// 请求参数
	accID := req.GetAccountId()

	// 逻辑处理
	exist, playerID, err := data.GetPlayerIDByAccountID(accID)

	// 返回消息
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

	// 返回消息
	rsp.PlayerId, rsp.ErrCode = playerID, int32(user.ErrCode_EC_SUCCESS)

	return
}

// GetPlayerInfo 获取玩家基本信息
func (pds *PlayerDataService) GetPlayerInfo(ctx context.Context, req *user.GetPlayerInfoReq) (rsp *user.GetPlayerInfoRsp, err error) {
	logrus.Debugln("GetPlayerInfo req", *req)

	// 默认返回消息
	rsp, err = &user.GetPlayerInfoRsp{
		ErrCode: int32(user.ErrCode_EC_FAIL),
	}, nil

	// 请求参数
	playerID := req.GetPlayerId()

	// 逻辑处理
	info, err := data.GetPlayerInfo(playerID)

	// 返回消息
	if err == nil {
		rsp.PlayerId, rsp.ErrCode = playerID, int32(user.ErrCode_EC_SUCCESS)
		rsp.NickName, rsp.Avatar = info[cache.NickNameField], info[cache.AvatarField]
		rsp.Name, rsp.Phone = info[cache.NameField], info[cache.PhoneField]
		value, _ := strconv.ParseInt(info[cache.GenderField], 10, 64)
		rsp.IpAddr, rsp.Gender = string("127.0.0.1"), uint32(value)
	}

	return
}

// UpdatePlayerInfo 设置玩家信息
func (pds *PlayerDataService) UpdatePlayerInfo(ctx context.Context, req *user.UpdatePlayerInfoReq) (rsp *user.UpdatePlayerInfoRsp, err error) {
	logrus.Debugln("SetPlayerInfo req", *req)

	// 默认返回消息
	rsp, err = &user.UpdatePlayerInfoRsp{
		ErrCode: int32(user.ErrCode_EC_FAIL),
		Result:  false,
	}, nil

	// 请求参数
	playerID := req.GetPlayerId()
	nickName := req.GetNickName() // 昵称
	avatar := req.GetAvatar()     // 头像
	name := req.GetName()         // 姓名
	phone := req.GetPhone()       // 电话
	gender := req.GetGender()     // 性别

	// 校验入参
	correct := validatePlayerInfoArgs()
	if !correct {
		rsp.ErrCode = int32(user.ErrCode_EC_Args)
		return
	}

	// 逻辑处理
	exist, result, err := data.UpdatePlayerInfo(playerID, nickName, avatar, name, phone, gender)

	// 返回消息
	if exist {
		rsp.Result = result
		rsp.ErrCode = int32(user.ErrCode_EC_SUCCESS)
	}

	return
}

// GetPlayerState 获取玩家状态
func (pds *PlayerDataService) GetPlayerState(ctx context.Context, req *user.GetPlayerStateReq) (rsp *user.GetPlayerStateRsp, err error) {
	logrus.Debugln("GetPlayerState req", *req)

	// 默认返回
	rsp, err = &user.GetPlayerStateRsp{
		ErrCode: int32(user.ErrCode_EC_FAIL),
		State:   user.PlayerState_PS_IDIE,
	}, nil

	// 逻辑处理
	state, _, err := data.GetPlayerState(req.GetPlayerId())

	if err == nil {
		rsp.State, rsp.ErrCode = user.PlayerState(state), int32(user.ErrCode_EC_SUCCESS)
	}
	return
}

// GetPlayerGameInfo 获取玩家游戏信息
func (pds *PlayerDataService) GetPlayerGameInfo(ctx context.Context, req *user.GetPlayerGameInfoReq) (rsp *user.GetPlayerGameInfoRsp, err error) {
	logrus.Debugln("GetPlayerState req", *req)

	// 请求参数
	playerID := req.GetPlayerId()
	gameID := req.GetGameId()

	// 默认返回消息
	rsp, err = &user.GetPlayerGameInfoRsp{
		ErrCode: int32(user.ErrCode_EC_FAIL),
		GameId:  gameID,
	}, nil

	// 逻辑处理
	exists, info, err := data.GetPlayerGameInfo(playerID, gameID)

	// 返回消息
	if !exists {
		rsp.ErrCode = int32(user.ErrCode_EC_EMPTY)
	}
	if err == nil {
		rsp.WinningRate, rsp.ErrCode = uint32(info.Winningrate), int32(user.ErrCode_EC_SUCCESS)
	}

	return
}

// UpdatePlayerState 设置玩家状态
func (pds *PlayerDataService) UpdatePlayerState(ctx context.Context, req *user.UpdatePlayerStateReq) (rsp *user.UpdatePlayerStateRsp, err error) {
	logrus.Debugln("SetPlayerState req", *req)

	// 默认返回消息
	rsp, err = &user.UpdatePlayerStateRsp{
		ErrCode: int32(user.ErrCode_EC_FAIL),
		Result:  false,
	}, nil

	// 请求参数
	playerID := req.GetPlayerId()
	oldState := uint32(req.GetOldState())
	newState := uint32(req.GetNewState())
	serverType := uint32(req.GetServerType()) // 服务端类型
	serverAddr := req.GetServerAddr()         // 服务端地址

	// 校验入参
	correct := validateSateArgs(oldState, newState, serverType, serverAddr)
	if !correct {
		rsp.ErrCode = int32(user.ErrCode_EC_Args)
		return
	}

	// 逻辑处理
	result, err := data.UpdatePlayerState(playerID, oldState, newState, serverType, serverAddr)

	// 返回消息
	if result && err == nil {
		rsp.Result, rsp.ErrCode = true, int32(user.ErrCode_EC_SUCCESS)
	}

	return
}

// GetGameListInfo 获取玩家游戏列表信息
func (pds *PlayerDataService) GetGameListInfo(ctx context.Context, req *user.GetGameListInfoReq) (rsp *user.GetGameListInfoRsp, err error) {
	logrus.Debugln("GetGameListInfo req", *req)

	// 默认返回消息
	rsp, err = &user.GetGameListInfoRsp{
		ErrCode:         int32(user.ErrCode_EC_FAIL),
		GameConfig:      []*user.GameConfig{},
		GameLevelConfig: []*user.GameLevelConfig{},
	}, nil

	// 逻辑处理
	rsp.GameConfig, rsp.GameLevelConfig, err = data.GetGameInfoList()

	// 返回消息
	if err == nil {
		rsp.ErrCode = int32(user.ErrCode_EC_SUCCESS)
	}
	return
}

// createPlayer 创建玩家
func createPlayer(accID uint64) (uint64, error) {
	playerID := data.AllocPlayerID()

	if playerID == 0 {
		return 0, fmt.Errorf("分配玩家 ID 失败")
	}
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
	if err := data.InitPlayerCoin(db.TPlayerCurrency{
		Playerid:       int64(playerID),
		Coins:          10000,
		Ingots:         0,
		Keycards:       0,
		Obtainingots:   0,
		Obtainkeycards: 0,
		Costingots:     0,
		Costkeycards:   0,
		Remark:         "",
		Createtime:     time.Now(),
		Createby:       "",
		Updatetime:     time.Now(),
		Updateby:       "",
	}); err != nil {
		return playerID, fmt.Errorf("初始化玩家(%d)金币数据失败: %v", playerID, err)
	}
	if err := data.InitPlayerState(int64(playerID)); err != nil {
		return playerID, fmt.Errorf("初始化玩家(%d)状态失败: %v", playerID, err)
	}
	return playerID, nil
}

// validatePlayerInfoArgs 校验更新玩家个人资料入参
func validatePlayerInfoArgs() bool {
	return true
}

// validateSateArgs 校验更新玩家状态入参
func validateSateArgs(oldState, newState, serverType uint32, serverAddr string) bool {
	userState := map[user.PlayerState]bool{
		user.PlayerState_PS_IDIE:     true,
		user.PlayerState_PS_MATCHING: true,
		user.PlayerState_PS_GAMEING:  true,
	}
	if !userState[user.PlayerState(oldState)] || !userState[user.PlayerState(newState)] {
		logrus.Warningln("player_state is incorrect, oldState:%d,newState:%d", oldState, newState)
		return false
	}

	userServerType := map[user.ServerType]bool{
		user.ServerType_ST_GATE:  true,
		user.ServerType_ST_MATCH: true,
		user.ServerType_ST_ROOM:  true,
	}

	if !userServerType[user.ServerType(serverType)] {
		logrus.Warningln("server_type is incorrect, server_type:%d", serverType)
		return false
	}

	if len(serverAddr) == 0 {
		logrus.Warningln("server_addr is empty, server_addr:%d", serverAddr)
		return false
	}
	return true
}
