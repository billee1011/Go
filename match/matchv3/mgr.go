package matchv3

import (
	"container/list"
	"fmt"
	"steve/client_pb/common"
	"steve/client_pb/match"
	"steve/client_pb/msgid"
	"steve/common/data/player"
	"steve/gutils"
	"steve/match/web"
	"sync/atomic"
	"time"

	"github.com/Sirupsen/logrus"
)

// applyPlayer 压入申请通道的玩家信息
type applyPlayer struct {
	playerID uint64 // 玩家 ID
	winRate  uint8  // 具体游戏的胜率
	gold     uint64 // 金币数

	gameID int // 游戏ID
}

// continueApply 续局申请
type continueApply struct {
	playerID uint64 // 玩家 ID
	gameID   int    // 游戏 ID，如果玩家续局失败，则以该游戏 ID 重新申请匹配
	cancel   bool   // 是否退出
}

// playerOffline 玩家离线
type playerLogin struct {
	playerID uint64 // 玩家 ID
}

// gameLevelConfig 游戏场次配置数据，来自数据库
type gameLevelConfig struct {
	levelID     int    // 场次ID
	levelName   []byte // 场次名字
	minGold     uint64 // 金币要求下限
	maxGold     uint64 // 金币要求上限
	bottomScore int    // 底分
}

// gameConfig 游戏配置数据，来自数据库
type gameConfig struct {
	gameID          int                       // 游戏ID
	gameName        []byte                    // 游戏名字
	needPlayerCount uint8                     // 满桌所需玩家数量
	levelConfig     map[int32]gameLevelConfig // 所有的游戏场次
}

// gameInfo 单个游戏的匹配信息
type gameInfo struct {
	allLevelChan map[int32]chan applyPlayer // 本levelID戏所有场次的匹配申请通道,Key:场次ID, Value:该场次的匹配申请通道
	config       gameConfig                 // 游levelID配置数据
}

// matchManager 匹配管理器
type matchManager struct {
	bInitFinish bool               // 是否初始化完成
	allGame     map[int32]gameInfo // 所有游戏的匹配信息，Key:游戏ID, Value:该游戏的匹配信息
	deskStartID uint64             // 桌子ID起始值

	applyChannel    chan applyPlayer   // 申请通道
	continueChannel chan continueApply // 续局通道
	loginChannel    chan playerLogin   // 玩家登录通道
	maxDeskID       uint64             // 最大牌桌 ID

	desks      map[uint64]*desk  // 当前匹配中的牌桌
	playerDesk map[uint64]uint64 // 匹配中的玩家， playerID -> deskID
}

// matchMgr 匹配管理器
var matchMgr = &matchManager{
	bInitFinish:     false,
	applyChannel:    make(chan applyPlayer, 128),
	continueChannel: make(chan continueApply, 128),
	loginChannel:    make(chan playerLogin, 128),
	maxDeskID:       0,
	/* 	gameConfig: map[int]gameConfig{
		int(roomanager.GameId_GAMEID_XUELIU):   gameConfig{needPlayerCount: 4},
		int(roomanager.GameId_GAMEID_XUEZHAN):  gameConfig{needPlayerCount: 4},
		int(roomanager.GameId_GAMEID_DOUDIZHU): gameConfig{needPlayerCount: 3},
		int(roomanager.GameId_GAMEID_ERRENMJ):  gameConfig{needPlayerCount: 2},
	}, */
	desks:      make(map[uint64]*desk, 128),
	playerDesk: make(map[uint64]uint64, 1024),
}

//  获取指定游戏的满桌需要的玩家
func (manager *matchManager) getGameNeedPlayerCount(gameID int32) uint8 {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "getGameNeedPlayerCount",
		"gameID":    gameID,
	})

	logEntry.Debugln("进入函数")

	// 得到该游戏的信息
	gameInfo, exist := manager.allGame[gameID]

	// 该游戏不存在
	if !exist {
		return 0
	}

	logEntry.Debugf("离开函数,满桌人数为：%v\n", gameInfo.config.needPlayerCount)

	return gameInfo.config.needPlayerCount
}

// 运行
func init() {
	//go matchMgr.run()

	// 为所有游戏，所有场次开始匹配
	matchMgr.startAllMatch()
}

// startAllMatch 开始为所有的游戏，所有的场次开启协程，执行匹配过程
func (manager *matchManager) startAllMatch() {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "startAllMatch",
	})

	logEntry.Debugln("进入函数")

	// 未初始化完毕，直接返回
	if manager.bInitFinish == false {
		logEntry.Errorln("开始匹配时发现仍未初始化完毕，返回")
		return
	}

	// 所有的游戏
	for gameID, gameInfo := range manager.allGame {

		// 所有的场次
		for levelID, _ := range gameInfo.config.levelConfig {

			// 每个场次，新建协程进行匹配
			go manager.startLevelMatch(gameID, levelID)
		}
	}

	logEntry.Debugln("离开函数")
	return
}

// 生成唯一的桌子ID
func (manager *matchManager) generateDeskID() uint64 {
	return atomic.AddUint64(&manager.deskStartID, 1)
}

// 返回指定间隔秒数后的胜率浮动值
// interval为0时表示0秒时(即首次匹配)的胜率浮动值
// interval为1时表示1秒后的胜率浮动值
// interval为2时表示2秒后的胜率浮动值
func (manager *matchManager) getWinRateValue(interval time.Duration) uint8 {
	// todo
	return 7
}

// 为指定玩家执行首次匹配
func (manager *matchManager) firstMatch(gameID int32, levelID int32, allWinRate []list.List, reqPlayer applyPlayer) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "firstMatch",
		"gameID":    gameID,
		"levelID":   levelID,
		"reqPlayer": reqPlayer,
	})

	logEntry.Debugln("进入函数")

	// 胜率范围检测
	if reqPlayer.winRate < 0 || reqPlayer.winRate > 100 {
		logEntry.Errorln("玩家%v的胜率为%v，不再执行匹配", reqPlayer.playerID, reqPlayer.winRate)
		return
	}

	// 是否匹配成功
	bMatchSuc := false

	// 新建一个匹配玩家
	newMatchPlayer := matchPlayer{
		playerID: reqPlayer.playerID,
		robotLv:  0,
		seat:     -1,
		IP:       uint32(0),
	}

	///////////////////////////////////////////////////////// 先检测首次匹配范围是不是有桌子 ///////////////////////////////////////////////

	// 首次匹配的胜率值
	firstRateValue := manager.getWinRateValue(0 * time.Second)

	// 玩家胜率范围下限
	playerBeginRate := reqPlayer.winRate - firstRateValue
	if playerBeginRate < 0 {
		playerBeginRate = 0
	}

	// 胜率范围上限
	playerEndRate := reqPlayer.winRate + firstRateValue
	if playerEndRate > 100 {
		playerEndRate = 100
	}

	// 从低往高匹配，存在桌子即可加入，这些桌子的胜率都在玩家的匹配范围
	for i := playerBeginRate; i <= playerEndRate; i++ {

		// 遍历所有的桌子
		for iter := allWinRate[i].Front(); iter != nil; iter = iter.Next() {

			// 把玩家加入桌子
			if manager.addPlayerToDesk(&newMatchPlayer, iter.Value.(*matchDesk)) == false {
				logEntry.Errorln("玩家%v加入桌子失败，返回", reqPlayer.playerID)
				return
			}

			bMatchSuc = true
			break
		}
	}

	///////////////////////////////////  首次范围失败后，再检测那些不在首次范围的，但因胜率范围扩张造成现在可能在玩家匹配范围了  /////////////////////////////
	if !bMatchSuc {

		// 遍历从 (playerBeginRate - 0]
		for i := playerBeginRate - 1; i >= uint8(0); i-- {
			// 遍历该概率下的所有桌子
			for iter := allWinRate[i].Front(); iter != nil; iter = iter.Next() {

				desk := iter.Value.(*matchDesk)

				interval := time.Now().Sub(desk.createTime)

				// 距离桌子创建时间不足1秒的，不检测
				if interval < time.Second*1 {
					continue
				}

				// 该桌子的概率浮动值
				rateValue := manager.getWinRateValue(interval)

				// 该桌子的浮动值下限
				deskBeginRate := i - rateValue
				if deskBeginRate < 0 {
					deskBeginRate = 0
				}

				// 该桌子的浮动值上限
				deskEndRate := i + rateValue
				if deskEndRate > 100 {
					deskEndRate = 100
				}

				// 玩家的匹配范围和桌子的匹配范围无交集时说明不匹配
				if (playerBeginRate > deskEndRate) || (playerEndRate < deskBeginRate) {
					continue
				}

				// 把玩家加入桌子
				if manager.addPlayerToDesk(&newMatchPlayer, iter.Value.(*matchDesk)) == false {
					logEntry.Errorln("玩家%v加入桌子失败，返回", reqPlayer.playerID)
					return
				}

				bMatchSuc = true
				break
			}
		}

		// 遍历从 endRate - 100]
		if !bMatchSuc {

		}
	}

	//////////////////////////////////////////////////////////  所有的桌子都失败了，新建桌子  ////////////////////////////////////////////////
	// 创建桌子
	needPlayerCount := manager.getGameNeedPlayerCount(gameID)
	newDesk := createMatchDesk(gameID, levelID, needPlayerCount, reqPlayer.gold)
	if newDesk == nil {
		logEntry.Errorf("创建匹配桌子失败，返回")
		return
	}

	// 把该玩家压入桌子
	manager.addPlayerToDesk(&newMatchPlayer, newDesk)

	logEntry.Debugln("离开函数")
	return
}

// startLevelMatch 开始单个游戏单个场次的匹配
// gameID : 游戏ID
// levelID : 场次ID
func (manager *matchManager) startLevelMatch(gameID int32, levelID int32) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "singleLevelMatch",
		"gameID":    gameID,
		"levelID":   levelID,
	})

	logEntry.Debugln("进入函数")

	// 建立从胜率1% - 100%的所有匹配队列
	// allWinRate[0] 表胜率为0%的所有匹配桌子，allWinRate[1] 表胜率为1%的所有匹配桌子，.........，allWinRate[100] 表胜率为100%的所有匹配桌子
	allWinRate := make([]list.List, 101)

	// 1秒1次的定时器
	timer := time.NewTicker(time.Second * 1)

	// 本场次的匹配申请通道
	gameInfo, exist := manager.allGame[gameID]
	if !exist {
		logEntry.Errorln("游戏不存在！退出")
		return
	}

	reqMatchChan, exist := gameInfo.allLevelChan[levelID]
	if !exist {
		logEntry.Errorln("游戏场次不存在！退出")
		return
	}

	// 一直匹配
	for {
		select {
		case req := <-reqMatchChan: // 匹配申请
			{
				manager.firstMatch(gameID, levelID, allWinRate, req) // 首次匹配
			}
		case pl := <-manager.loginChannel: // 登录玩家
			{
				manager.onPlayerLogin(pl.playerID)
			}
		case <-timer.C: // 定时器触发
			{
				manager.handleRobotTick()
			}
		}
	}

	logEntry.Debugln("离开函数")
	return
}

// addContinueDesk 添加续局牌桌
func (manager *matchManager) addContinueDesk(players []deskPlayer, gameID int, fixBanker bool, bankerSeat int) {
	manager.maxDeskID++
	// 有玩家在匹配中，不创建
	for _, player := range players {
		if _, ok := manager.playerDesk[player.playerID]; ok {
			logrus.WithField("player_id", player.playerID).Infoln("添加续局牌桌时玩家已经在匹配中了")
			return
		}
	}
	deskID := manager.maxDeskID
	desk := createContinueDesk(gameID, deskID, players, fixBanker, bankerSeat)
	for _, player := range players {
		manager.playerDesk[player.playerID] = deskID
	}
	manager.desks[deskID] = desk
}

// dismissContinueDesk 解散续局牌桌
// emitPlayer 发起解散的玩家 ID，超时解散时为0
func (manager *matchManager) dismissContinueDesk(desk *desk, emitPlayer uint64) {
	logrus.WithFields(logrus.Fields{
		"func_name":    "mgr.dismissContinueDesk",
		"ready_player": desk.players,
	}).Debugln("解散续局牌桌")
	notify := match.MatchContinueDeskDimissNtf{}

	for _, deskPlayer := range desk.players {
		delete(manager.playerDesk, deskPlayer.playerID)
		if deskPlayer.playerID != emitPlayer {
			gutils.SendMessage(deskPlayer.playerID, msgid.MsgID_MATCH_CONTINUE_DESK_DIMISS_NTF, &notify)
		}
		// 更新状态为空闲状态
		player.SetPlayerPlayStates(deskPlayer.playerID, player.PlayStates{
			State:  int(common.PlayerState_PS_IDLE),
			GameID: int(desk.gameID),
		})
	}
	for playerID := range desk.continueWaitPlayers {
		delete(manager.playerDesk, playerID)
		if playerID != emitPlayer {
			gutils.SendMessage(playerID, msgid.MsgID_MATCH_CONTINUE_DESK_DIMISS_NTF, &notify)
		}
	}
	delete(manager.desks, desk.deskID)
}

// addPlayer 添加匹配玩家
func (manager *matchManager) addPlayer(playerID uint64, gameID int) {
	manager.applyChannel <- applyPlayer{
		playerID: playerID,
	}
	return
}

// 分发匹配请求
// playerID 	:	玩家ID
// gameID		：	请求匹配的游戏ID
// levelID		:   请求匹配的级别ID
// 返回string 	 ：	 返回的错误描述，成功时返回空
func (manager *matchManager) dispatchMatchReq(playerID uint64, gameID int32, levelID int32) string {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "dispatchMatchReq",
		"playerID":  playerID,
		"gameID":    gameID,
		"levelID":   levelID,
	})

	logEntry.Debugln("进入函数")

	// 得到该游戏的信息
	gameInfo, exist := manager.allGame[gameID]

	// 该游戏不存在
	if !exist {
		logrus.Errorln("请求匹配的游戏不存在")
		return fmt.Sprintf("请求匹配的游戏ID:%v不存在，请求的玩家ID:%v", gameID, playerID)
	}

	// 得到该场次的信息
	levelConfig, exist := gameInfo.config.levelConfig[levelID]

	// 该场次不存在
	if !exist {
		logrus.Errorln("请求匹配的游戏存在，但场次不存在")
		return fmt.Sprintf("请求匹配的游戏ID:%v存在，但场次ID:%v不存在，请求的玩家ID:%v", gameID, levelID, playerID)
	}

	// 玩家金币数
	playerGold := player.GetPlayerCoin(playerID)

	// 金币范围检测
	if playerGold < levelConfig.minGold {
		logrus.Errorln("玩家金币数小于游戏场次金币要求最小值，最小值：%v", levelConfig.minGold)
		return fmt.Sprintf("玩家金币数小于游戏场次金币要求最小值，请求匹配的游戏ID:%v，场次ID:%v，玩家ID:%v", gameID, levelID, playerID)
	}

	if playerGold > levelConfig.maxGold {
		logrus.Errorf("玩家金币数大于游戏场次金币要求最大值，最大值：%v", levelConfig.maxGold)
		return fmt.Sprintf("玩家金币数大于游戏场次金币要求最大值，请求匹配的游戏ID:%v，场次ID:%v，玩家ID:%v", gameID, levelID, playerID)
	}

	// 全部检测通过

	// 获取该场次的申请通道
	reqMatchChan, exist := gameInfo.allLevelChan[levelID]
	if !exist {
		logrus.Errorln("内部错误，请求匹配的游戏存在，场次存在，但找不到该场次的匹配申请通道")
		return fmt.Sprintf("请求匹配的游戏存在，场次存在，但找不到该场次的匹配申请通道，请求匹配的游戏ID:%v，场次ID:%v，玩家ID:%v", gameID, levelID, playerID)
	}

	// 计算该游戏的胜率，已经乘以100,比如：50表胜率为50%
	// 暂时写死，todo
	playerWinRate := 50

	// 压入通道
	reqMatchChan <- applyPlayer{
		playerID: playerID,
		winRate:  uint8(playerWinRate),
		gold:     playerGold,
	}

	logEntry.Debugln("离开函数")

	return ""
}

// addContinueApply 添加续局申请
func (manager *matchManager) addContinueApply(playerID uint64, cancel bool, gameID int) {
	manager.continueChannel <- continueApply{
		playerID: playerID,
		cancel:   cancel,
		gameID:   gameID,
	}
	return
}

// addLoginData 添加玩家登录信息
func (manager *matchManager) addLoginData(playerID uint64) {
	manager.loginChannel <- playerLogin{
		playerID: playerID,
	}
}

// run 执行匹配流程
func (manager *matchManager) run() {

	// 从DB读取游戏配置信息
	// todo

	// 机器人的定时器（1秒1次）
	robotTick := time.NewTicker(time.Second * 1)

	for {
		select {
		/* 		case ap := <-manager.applyChannel: // 普通匹配申请
		   			{
		   				manager.acceptApplyPlayer(ap.gameID, ap.playerID)
		   			}
		   		case cp := <-manager.continueChannel: // 续局匹配申请
		   			{
		   				manager.acceptContinuePlayer(cp.gameID, cp.playerID, cp.cancel)
		   			} */
		case pl := <-manager.loginChannel: // 登录玩家
			{
				manager.onPlayerLogin(pl.playerID)
			}
		case <-robotTick.C: // 机器人定时器
			{
				manager.handleRobotTick()
			}
		}
	}
}

// onPlayerLogin 玩家登录，取消玩家匹配
func (manager *matchManager) onPlayerLogin(playerID uint64) {
	entry := logrus.WithField("player_id", playerID)

	// 是否存在其桌子
	deskID, ok := manager.playerDesk[playerID]
	if !ok {
		return
	}

	// 得到该桌子
	desk, ok := manager.desks[deskID]
	if !ok {
		delete(manager.playerDesk, playerID)
		entry.Errorln("没有对应的牌桌")
		return
	}

	// 续局牌桌直接解散
	if desk.isContinue {
		entry.Debugln("玩家重新登录，解散续局牌桌")
		manager.dismissContinueDesk(desk, playerID)
		return
	}

	// 从桌子中删除
	//desk.removePlayer(playerID)

	// 删除 playerID -> deskID的映射
	delete(manager.playerDesk, playerID)

	entry.Debugln("玩家重新登录，移出普通匹配")

	// 桌子没人了，则解散桌子
	if len(desk.players) == 0 {
		delete(manager.desks, deskID)
	}
}

/* // acceptContinuePlayer 接收续局匹配玩家
func (manager *matchManager) acceptContinuePlayer(gameID int, playerID uint64, cancel bool) {
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "mgr.acceptContinuePlayer",
		"player_id": playerID,
		"cancel":    cancel,
	})
	deskID, ok := manager.playerDesk[playerID]
	if !ok && !cancel {
		manager.acceptApplyPlayer(gameID, playerID)
		return
	}
	entry = entry.WithField("desk_id", deskID)
	desk, ok := manager.desks[deskID]
	if !ok {
		delete(manager.playerDesk, playerID)
		entry.Errorln("牌桌不存在")
		manager.acceptApplyPlayer(gameID, playerID)
		return
	}
	// 非续局牌桌
	if !desk.isContinue {
		return
	}
	if cancel {
		entry.Debugf("玩家取消续局造成解散续局牌桌，玩家ID：%v", playerID)
		manager.dismissContinueDesk(desk, playerID)
		return
	}
	player, ok := desk.continueWaitPlayers[playerID]
	if !ok {
		return
	}
	entry.Debugln("接收续局玩家")
	delete(desk.continueWaitPlayers, playerID)
	manager.addDeskPlayer2Desk(&player, desk)
	return
} */

/* // acceptApplyPlayer 接收申请匹配玩家
func (manager *matchManager) acceptApplyPlayer(gameID int, playerID uint64) {
	deskID, ok := manager.playerDesk[playerID]
	logrus.WithFields(logrus.Fields{
		"func_name":   "mgr.acceptApplyPlayer",
		"player_id":   playerID,
		"game_id":     gameID,
		"old_desk_id": deskID,
	}).Debugln("接收申请匹配玩家")
	if ok {
		// 等待续局中
		if desk, exist := manager.desks[deskID]; exist && desk.isContinue {
			logrus.Debugf("普通匹配时发现玩家已经在续局牌桌中，解散续局牌桌，玩家ID:%v", playerID)
			manager.dismissContinueDesk(desk, playerID)
		} else {
			return // 匹配中
		}
	}
	// 加入到牌桌
	for _, desk := range manager.desks {
		if desk.gameID != gameID || desk.isContinue {
			continue
		}
		manager.addDeskPlayer2Desk(&deskPlayer{
			playerID: playerID,
		}, desk)
		return
	}
	manager.maxDeskID++
	desk := createDesk(gameID, manager.maxDeskID)
	manager.desks[desk.deskID] = desk
	manager.addDeskPlayer2Desk(&deskPlayer{
		playerID: playerID,
	}, desk)
} */

/* // addDeskPlayer2Desk 将玩家添加到牌桌
func (manager *matchManager) addDeskPlayer2Desk(deskPlayer *deskPlayer, desk *desk) {
	player.SetPlayerPlayStates(deskPlayer.playerID, player.PlayStates{
		State:  int(common.PlayerState_PS_MATCHING),
		GameID: int(desk.gameID),
	})
	desk.players = append(desk.players, *deskPlayer)
	manager.playerDesk[deskPlayer.playerID] = desk.deskID
	manager.removeOfflines(desk)
	config := manager.gameConfig[desk.gameID]
	if len(desk.players) >= config.needPlayerCount {
		manager.onDeskFinish(desk)
	}
} */

// 将玩家添加到牌桌
// pMatchPlayer : 匹配的玩家
//
func (manager *matchManager) addPlayerToDesk(pPlayer *matchPlayer, pDesk *matchDesk) (bSuc bool) {

	// 参数检测
	if pPlayer == nil || pDesk == nil {
		logrus.Error("参数错误，pPlayer == nil || pDesk == nil，返回")
		return false
	}

	logrus.WithFields(logrus.Fields{
		"player": pPlayer,
		"desk":   pDesk,
	})

	// 压入该玩家
	pDesk.players = append(pDesk.players, *pPlayer)

	logrus.Debugf("桌子%v压入了玩家%v\n", pDesk, pPlayer)

	// playerID与deskID的映射
	//manager.playerDesk[deskPlayer.playerID] = desk.deskID

	// 移除不在线的
	//manager.removeOfflines(desk)

	// 满桌需要的玩家数量
	needPlayerCount := manager.getGameNeedPlayerCount(pDesk.gameID)

	// 满员时的处理
	if uint8(len(pDesk.players)) >= needPlayerCount {
		manager.onDeskFull(pDesk)
	}

	return true
}

/* // fillRobots 填充机器人
func (manager *matchManager) fillRobots(desk *desk) {
	config := manager.gameConfig[desk.gameID]
	logrus.WithFields(logrus.Fields{
		"func_name":  "mgr.fillRobots",
		"desk":       desk,
		"need_count": config.needPlayerCount,
	}).Debugln("加入机器人")
	curPlayerCount := len(desk.players)
	for i := curPlayerCount; i < config.needPlayerCount; i++ {
		manager.addDeskPlayer2Desk(&deskPlayer{
			playerID: GetIdleRobot(1),
			robotLv:  1,
		}, desk)
	}
} */

// removeOfflines 移除 desk 中的离线玩家
func (manager *matchManager) removeOfflines(desk *desk) {
	newPlayers := make([]deskPlayer, 0, len(desk.players))
	for _, deskPlayer := range desk.players {
		// 机器人不移除
		if deskPlayer.robotLv != 0 {
			newPlayers = append(newPlayers, deskPlayer)
			continue
		}
		online := (player.GetPlayerGateAddr(deskPlayer.playerID) != "")
		if online {
			newPlayers = append(newPlayers, deskPlayer)
		} else {
			delete(manager.playerDesk, deskPlayer.playerID)
		}
	}
	desk.players = newPlayers
}

// onDeskFinish 牌桌匹配完成
func (manager *matchManager) onDeskFinish(desk *desk) {
	requestCreateDesk(desk)
	players := desk.players
	// 解除关联
	for _, player := range players {
		delete(manager.playerDesk, player.playerID)
	}
	// 移除 desk
	delete(manager.desks, desk.deskID)
}

// onDeskFull 桌子满员时的处理
func (manager *matchManager) onDeskFull(pDesk *matchDesk) {

	// 移除桌子

	// 通知room服创建桌子
	sendCreateDesk(pDesk)

	/* 	players := pDesk.players

	   	// 解除关联
	   	for _, player := range players {
	   		delete(manager.playerDesk, player.playerID)
	   	} */

	// 移除 desk
	//delete(manager.desks, pDesk.deskID)
}

// handleRobotTick 处理机器人 tick
func (manager *matchManager) handleRobotTick() {
	// 避免遍历时删除
	deskIDs := make([]uint64, 0, len(manager.desks))
	for deskID := range manager.desks {
		deskIDs = append(deskIDs, deskID)
	}
	for _, deskID := range deskIDs {
		desk := manager.desks[deskID]
		if !desk.isContinue && time.Now().Sub(desk.createTime) >= web.GetRobotJoinTime() {
			//manager.fillRobots(desk)
		}
	}
}

/* // checkContinueDesks 检查续局牌桌，超过 20s 解散
func (manager *matchManager) checkContinueDesks() {
	// 避免遍历时删除
	deskIDs := make([]uint64, 0, len(manager.desks))
	for deskID := range manager.desks {
		deskIDs = append(deskIDs, deskID)
	}
	for _, deskID := range deskIDs {
		desk := manager.desks[deskID]
		// 非续局牌桌
		if !desk.isContinue {
			continue
		}
		interval := time.Now().Sub(desk.createTime)
		// 超过解散时间
		if interval >= web.GetContinueDismissTime() {
			logrus.Debugf("续局牌桌超时，解散续局牌桌，桌子创建时间=%v，现在时间=%v,超时时间=%v", desk.createTime, time.Now(), web.GetContinueDismissTime())
			manager.dismissContinueDesk(desk, 0)
			continue
		}
		// 超过机器人续局时间
		if interval >= web.GetContinueRobotTime() {
			manager.robotContinue(desk)
		}
	}
} */

/* // robotContinue 机器人作续局决策
func (manager *matchManager) robotContinue(desk *desk) {
	robots := make([]uint64, 0, len(desk.continueWaitPlayers))

	for playerID := range desk.continueWaitPlayers {
		robots = append(robots, playerID)
	}

	for _, playerID := range robots {
		player := desk.continueWaitPlayers[playerID]
		if player.robotLv == 0 {
			continue
		}
		rate := web.GetRobotContinueRate(player.winner)
		continual := gutils.Probability(rate)
		manager.acceptContinuePlayer(desk.gameID, playerID, !continual)
	}
} */
