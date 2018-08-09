package user

import (
	"context"
	"fmt"
	"math/rand"
	"steve/entity/cache"
	"steve/entity/db"
	"steve/hall/data"
	"steve/server_pb/user"
	"time"

	"github.com/Sirupsen/logrus"
)

// PlayerDataService 实现 user.PlayerServer
type PlayerDataService struct{}

var _ user.PlayerDataServer = new(PlayerDataService)

// GetPlayerByAccount 根据账号获取玩家 ID
func (pds *PlayerDataService) GetPlayerByAccount(ctx context.Context, req *user.GetPlayerByAccountReq) (rsp *user.GetPlayerByAccountRsp, err error) {
	logrus.Debugf("GetPlayerByAccount req: (%v)", *req)

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
	logrus.Debugf("GetPlayerInfo req : (%v)", *req)

	// 默认返回消息
	rsp, err = &user.GetPlayerInfoRsp{
		ErrCode: int32(user.ErrCode_EC_FAIL),
	}, nil

	// 请求参数
	playerID := req.GetPlayerId()

	// 逻辑处理
	fields := []string{cache.NickName, cache.Gender, cache.Avatar, cache.ChannelID, cache.ProvinceID, cache.CityID}
	player, err := data.GetPlayerInfo(playerID, fields...)

	// 返回消息
	if err == nil {
		rsp.ErrCode = int32(user.ErrCode_EC_SUCCESS)
		rsp.PlayerId, rsp.Gender = playerID, uint32(player.Gender)
		rsp.NickName, rsp.Avatar = player.Nickname, player.Avatar
		rsp.ChannelId, rsp.ProvinceId, rsp.CityId = uint32(player.Channelid), uint32(player.Provinceid), uint32(player.Cityid)
	}

	return
}

// UpdatePlayerInfo 设置玩家信息
func (pds *PlayerDataService) UpdatePlayerInfo(ctx context.Context, req *user.UpdatePlayerInfoReq) (rsp *user.UpdatePlayerInfoRsp, err error) {
	logrus.Debugf("SetPlayerInfo req: (%v)", *req)

	// 默认返回消息
	rsp, err = &user.UpdatePlayerInfoRsp{
		ErrCode: int32(user.ErrCode_EC_FAIL),
		Result:  false,
	}, nil

	// 请求参数
	playerID := req.GetPlayerId()
	nickName := req.GetNickName() // 昵称
	avatar := req.GetAvatar()     // 头像
	gender := req.GetGender()     // 性别

	// 校验入参
	correct := validatePlayerInfoArgs()
	if !correct {
		rsp.ErrCode = int32(user.ErrCode_EC_Args)
		return
	}

	// 逻辑处理
	fields := []string{cache.NickName, cache.Avatar, cache.Gender}
	dbPlayer := db.TPlayer{
		Playerid: int64(playerID),
		Nickname: nickName,
		Gender:   int(gender),
		Avatar:   avatar,
	}
	err = data.SetPlayerFields(playerID, fields, &dbPlayer)

	// 返回消息
	if err == nil {
		rsp.Result = true
		rsp.ErrCode = int32(user.ErrCode_EC_SUCCESS)
	}

	return
}

// GetPlayerState 获取玩家状态
func (pds *PlayerDataService) GetPlayerState(ctx context.Context, req *user.GetPlayerStateReq) (rsp *user.GetPlayerStateRsp, err error) {
	logrus.Debugf("GetPlayerState req: (%v)", *req)

	// 默认返回
	rsp, err = &user.GetPlayerStateRsp{
		ErrCode: int32(user.ErrCode_EC_FAIL),
		State:   user.PlayerState_PS_IDIE,
	}, nil

	// 逻辑处理
	pState, err := data.GetPlayerState(req.GetPlayerId(), []string{cache.GameID, cache.LevelID, cache.GameState, cache.IPAddr, cache.GateAddr, cache.MatchAddr, cache.RoomAddr}...)

	if err == nil {
		rsp.GameId, rsp.LevelId = uint32(pState.GameID), uint32(pState.LevelID)
		rsp.State, rsp.IpAddr = user.PlayerState(pState.State), pState.IPAddr
		rsp.GateAddr, rsp.MatchAddr, rsp.RoomAddr = pState.GateAddr, pState.MatchAddr, pState.RoomAddr
		rsp.ErrCode = int32(user.ErrCode_EC_SUCCESS)
	}
	logrus.Debugln("get player state rsp", rsp)
	return
}

// GetPlayerGameInfo 获取玩家游戏信息
func (pds *PlayerDataService) GetPlayerGameInfo(ctx context.Context, req *user.GetPlayerGameInfoReq) (rsp *user.GetPlayerGameInfoRsp, err error) {
	logrus.Debugf("GetPlayerState req :(%v)", *req)

	// 请求参数
	playerID := req.GetPlayerId()
	gameID := req.GetGameId()

	// 默认返回消息
	rsp, err = &user.GetPlayerGameInfoRsp{
		ErrCode: int32(user.ErrCode_EC_FAIL),
		GameId:  gameID,
	}, nil

	// 逻辑处理
	fields := []string{cache.WinningRate, cache.WinningBurea, cache.TotalBurea, cache.MaxWinningStream, cache.MaxMultiple}
	exist, info, err := data.GetPlayerGameInfo(playerID, gameID, fields...)

	// 返回消息
	if !exist {
		rsp.ErrCode = int32(user.ErrCode_EC_SUCCESS)
		return
	}
	if err == nil {
		rsp.ErrCode = int32(user.ErrCode_EC_SUCCESS)
		rsp.WinningRate, rsp.WinningBurea, rsp.TotalBurea = uint32(info.Winningrate), uint32(info.Winningburea), uint32(info.Totalbureau)
		rsp.MaxWinningStream, rsp.MaxMultiple = uint32(info.Maxwinningstream), uint32(info.Maxmultiple)
	}

	return
}

// UpdatePlayerState 设置玩家状态
func (pds *PlayerDataService) UpdatePlayerState(ctx context.Context, req *user.UpdatePlayerStateReq) (rsp *user.UpdatePlayerRsp, err error) {
	logrus.Debugf("UpdatePlayerState req :(%v)", *req)

	// 默认返回消息
	rsp, err = &user.UpdatePlayerRsp{
		ErrCode: int32(user.ErrCode_EC_FAIL),
		Result:  false,
	}, nil

	// 请求参数
	playerID := req.GetPlayerId()
	oldState := uint32(req.GetOldState())
	newState := uint32(req.GetNewState())
	gameID := uint32(req.GetGameId())
	levelID := uint32(req.GetLevelId())
	// 校验入参
	correct := validateUserSate(oldState, newState)
	if !correct {
		rsp.ErrCode = int32(user.ErrCode_EC_Args)
		return
	}

	// 逻辑处理
	result, err := data.UpdatePlayerState(playerID, oldState, newState, gameID, levelID)

	// 返回消息
	if result && err == nil {
		rsp.Result, rsp.ErrCode = true, int32(user.ErrCode_EC_SUCCESS)
	}

	return
}

// UpdatePlayerGateInfo 更新玩家网关信息
func (pds *PlayerDataService) UpdatePlayerGateInfo(ctx context.Context, req *user.UpdatePlayerGateInfoReq) (rsp *user.UpdatePlayerRsp, err error) {
	logrus.Debugf("UpdatePlayerGateInfo req :(%v)", *req)

	// 默认返回消息
	rsp, err = &user.UpdatePlayerRsp{
		ErrCode: int32(user.ErrCode_EC_FAIL),
		Result:  false,
	}, nil

	// 请求参数
	playerID := req.GetPlayerId() // 玩家ID
	ipAddr := req.GetIpAddr()     // 客户端IP地址
	gateAddr := req.GetGateAddr() // 网关服地址

	// 校验
	if gateAddr == "" {
		rsp.ErrCode = int32(user.ErrCode_EC_Args)
		return
	}

	// 逻辑处理
	result, err := data.UpdatePlayerGateInfo(playerID, ipAddr, gateAddr)

	// 返回消息
	if result && err == nil {
		rsp.Result, rsp.ErrCode = true, int32(user.ErrCode_EC_SUCCESS)
	}
	return
}

// UpdatePlayerServerAddr 更新玩家服务端地址
func (pds *PlayerDataService) UpdatePlayerServerAddr(ctx context.Context, req *user.UpdatePlayerServerAddrReq) (rsp *user.UpdatePlayerRsp, err error) {
	logrus.Debugf("UpdatePlayerGateInfo req :(%v)", *req)

	// 默认返回消息
	rsp, err = &user.UpdatePlayerRsp{
		ErrCode: int32(user.ErrCode_EC_FAIL),
		Result:  false,
	}, nil

	// 请求参数
	playerID := req.GetPlayerId()
	serverType := req.GetServerType()
	serverAddr := req.GetServerAddr()

	// 校验
	correct := validateServerType(serverType, serverAddr)
	if !correct {
		rsp.ErrCode = int32(user.ErrCode_EC_Args)
		return
	}

	//逻辑处理
	result, err := data.UpdatePlayerServerAddr(playerID, uint32(serverType), serverAddr)

	// 返回消息
	if result && err == nil {
		rsp.Result, rsp.ErrCode = true, int32(user.ErrCode_EC_SUCCESS)
	}

	return
}

// GetGameListInfo 获取玩家游戏列表信息
func (pds *PlayerDataService) GetGameListInfo(ctx context.Context, req *user.GetGameListInfoReq) (rsp *user.GetGameListInfoRsp, err error) {
	logrus.Debugf("GetGameListInfo req :(%v)", *req)

	// 默认返回消息
	rsp, err = &user.GetGameListInfoRsp{
		ErrCode:         int32(user.ErrCode_EC_FAIL),
		GameConfig:      []*user.GameConfig{},
		GameLevelConfig: []*user.GameLevelConfig{},
	}, nil

	// 逻辑处理
	gameConfig, gameLevelConfig, err := data.GetGameInfoList()
	// 返回消息
	if err == nil {
		rsp.GameConfig, rsp.GameLevelConfig = DBGameConfig2Server(gameConfig), DBGamelevelConfig2Sercer(gameLevelConfig)
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
	showUID := data.AllocShowUID()
	if err := data.InitPlayerData(db.TPlayer{
		Accountid:    int64(accID),
		Playerid:     int64(playerID),
		Showuid:      showUID,
		Type:         1,
		Channelid:    0,                                // TODO ，渠道 ID
		Nickname:     fmt.Sprintf("player%d", showUID), // TODO,昵称
		Gender:       1,
		Avatar:       getRandomAvator(), // TODO , 头像
		Provinceid:   0,                 // TODO， 省ID
		Cityid:       0,                 // TODO 市ID
		Name:         "",                // TODO: 真实姓名
		Phone:        "",                // TODO: 电话
		Idcard:       "",                // TODO 身份证
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

func getRandomAvator() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%d", r.Intn(2)+1)
}

// validatePlayerInfoArgs 校验更新玩家个人资料入参
func validatePlayerInfoArgs() bool {
	return true
}

// validateUserSate 校验更新玩家状态入参
func validateUserSate(oldState, newState uint32) bool {
	userState := map[user.PlayerState]bool{
		user.PlayerState_PS_IDIE:     true,
		user.PlayerState_PS_MATCHING: true,
		user.PlayerState_PS_GAMEING:  true,
	}
	if !userState[user.PlayerState(oldState)] || !userState[user.PlayerState(newState)] {
		logrus.Warningln("player_state is incorrect, oldState:%d,newState:%d", oldState, newState)
		return false
	}

	return true
}

// validateServerType 校验更新玩家服务端地址
func validateServerType(serverType user.ServerType, serverAddr string) bool {

	userServerType := map[user.ServerType]bool{
		user.ServerType_ST_GATE:  true,
		user.ServerType_ST_MATCH: true,
		user.ServerType_ST_ROOM:  true,
	}

	if !userServerType[user.ServerType(serverType)] {
		logrus.Warningln("server_type is incorrect, server_type:%d", serverType)
		return false
	}

	return true
}
