package matchv3

import (
	"container/list"
	"context"
	"fmt"
	"steve/external/goldclient"
	"steve/match/web"
	"steve/server_pb/gold"
	"steve/server_pb/user"
	"steve/structs"
	"sync/atomic"
	"time"

	"github.com/Sirupsen/logrus"
)

// reqMatchPlayer 压入申请匹配通道的玩家信息
type reqMatchPlayer struct {
	playerID uint64 // 玩家 ID
	winRate  int8   // 具体游戏的胜率
	gold     int64  // 金币数
	IP       string // IP地址
}

// playerOffline 玩家离线
type playerLogin struct {
	playerID uint64 // 玩家 ID
}

// gameLevelConfig 游戏场次配置数据
type gameLevelConfig struct {
	levelID            uint32 // 场次ID
	levelName          string // 场次名字
	bottomScore        uint32 // 底分
	minGold            int64  // 金币要求下限
	maxGold            int64  // 金币要求上限
	minNeedPlayerCount uint32 // 允许最低人数
	maxNeedPlayerCount uint32 // 允许最高人数
}

// gameConfig 游戏配置数据
type gameConfig struct {
	gameID      uint32                     // 游戏ID
	gameName    string                     // 游戏名字
	levelConfig map[uint32]gameLevelConfig // 所有的游戏场次
}

// levelGlobalInfo 单个游戏单个场次的全局信息
type levelGlobalInfo struct {
	// 游戏ID
	gameID uint32

	// 场次ID
	levelID uint32

	// 胜率1% - 100%的所有匹配桌子
	// allRateDesks[0] 表胜率为0%的所有匹配桌子，allRateDesks[1] 表胜率为1%的所有匹配桌子，.........，allRateDesks[100] 表胜率为100%的所有匹配桌子
	allRateDesks []list.List

	// 已成功匹配的玩家，Key:玩家ID，Value:桌子ID
	sucPlayers map[uint64]uint64

	// 已成功匹配的桌子，Key:桌子ID，Value:桌子信息
	sucDesks map[uint64]*sucDesk
}

// gameInfo 单个游戏的匹配信息
type gameInfo struct {
	allLevelChan map[uint32]chan reqMatchPlayer // 本levelID戏所有场次的匹配申请通道,Key:场次ID, Value:该场次的匹配申请通道
	config       gameConfig                     // 游levelID配置数据
}

// matchManager 匹配管理器
type matchManager struct {
	bInitFinish    bool                // 是否初始化完成
	allGame        map[uint32]gameInfo // 所有游戏的匹配信息，Key:游戏ID, Value:该游戏的匹配信息
	deskStartID    uint64              // 桌子ID起始值
	rateCompuValue []float32           // 胜率计算配置
	goldCompuValue []float32           // 金币计算配置

	applyChannel chan reqMatchPlayer // 申请通道
	//continueChannel chan continueApply // 续局通道
	loginChannel chan playerLogin // 玩家登录通道

	//maxDeskID    uint64           // 最大牌桌 ID

	//desks      map[uint64]*desk  // 当前匹配中的牌桌
	//playerDesk map[uint64]uint64 // 匹配中的玩家， playerID -> deskID
}

// matchMgr 匹配管理器
var matchMgr = &matchManager{
	bInitFinish:    false,
	allGame:        map[uint32]gameInfo{},
	deskStartID:    0,
	rateCompuValue: make([]float32, 0, web.GetMaxCompuValidTime()+1),
	goldCompuValue: make([]float32, 0, web.GetMaxCompuValidTime()+1),

	applyChannel: make(chan reqMatchPlayer, 128),
	//continueChannel: make(chan continueApply, 128),
	loginChannel: make(chan playerLogin, 128),
	//maxDeskID:    0,
	/* 	gameConfig: map[int]gameConfig{
		int(roomanager.GameId_GAMEID_XUELIU):   gameConfig{needPlayerCount: 4},
		int(roomanager.GameId_GAMEID_XUEZHAN):  gameConfig{needPlayerCount: 4},
		int(roomanager.GameId_GAMEID_DOUDIZHU): gameConfig{needPlayerCount: 3},
		int(roomanager.GameId_GAMEID_ERRENMJ):  gameConfig{needPlayerCount: 2},
	}, */
	//desks:      make(map[uint64]*desk, 128),
	//playerDesk: make(map[uint64]uint64, 1024),
}

//  获取指定游戏的满桌需要的玩家
func (manager *matchManager) getGameNeedPlayerCount(gameID uint32, levelID uint32) uint8 {
	logEntry := logrus.WithFields(logrus.Fields{
		"gameID":  gameID,
		"levelID": levelID,
	})

	logEntry.Debugln("进入函数")

	// 得到该游戏的信息
	gameInfo, exist := manager.allGame[gameID]
	// 该游戏不存在
	if !exist {
		logEntry.Errorln("游戏ID不存在")
		return 0
	}

	// 得到该场次的信息
	levelInfo, exist := gameInfo.config.levelConfig[levelID]
	// 该场次不存在
	if !exist {
		logEntry.Errorln("场次ID不存在")
		return 0
	}

	logEntry.Debugf("离开函数,最低满桌人数为：%v\n", levelInfo.minNeedPlayerCount)

	return uint8(levelInfo.minNeedPlayerCount)
}

// init 初始化并运行
func init() {

	// 初始化
	if !matchMgr.init() {
		return
	}

	// 运行
	matchMgr.start()
}

// init 初始化操作
func (manager *matchManager) init() bool {

	logrus.Debugln("进入函数")

	// 初始化胜率和金币差异
	if !manager.compuRateGold() {
		logrus.Errorln("初始化胜率和金币差异失败，返回")
		return false
	}

	// 获取游戏和场次配置
	if !manager.requestGameLevelConfig() {
		logrus.Errorln("获取游戏和场次配置失败，返回")
		return false
	}

	// 初始化结束
	manager.bInitFinish = true

	logrus.Debugln("离开函数")
	return true
}

// rateCompuformula 胜率计算公式
func (manager *matchManager) rateCompuformula(t uint8) float32 {
	return (0.05*float32(t) + web.GetWinRateCompuBase())
}

// goldCompuformula 金币计算公式
func (manager *matchManager) goldCompuformula(t uint8) float32 {
	return ((0.2*float32(t))*(0.2*float32(t)) + web.GetGoldCompuBase())
}

// compuRateGold 计算胜率浮动值和金币浮动值
func (manager *matchManager) compuRateGold() bool {
	logrus.Debugln("进入函数")

	// 计算胜率浮动值
	var i uint8 = 0
	for i = 0; i <= uint8(web.GetMaxCompuValidTime()); i++ {
		manager.rateCompuValue = append(manager.rateCompuValue, manager.rateCompuformula(i))
	}

	// 计算金币浮动值
	for i = 0; i <= uint8(web.GetMaxCompuValidTime()); i++ {
		manager.goldCompuValue = append(manager.goldCompuValue, manager.goldCompuformula(i))
	}

	logrus.Debugln("计算出来的胜率浮动值：", manager.rateCompuValue)
	logrus.Debugln("计算出来的金币浮动值：", manager.goldCompuValue)

	logrus.Debugln("离开函数")

	return true
}

// 向hall服请求游戏，场次的配置信息
func (manager *matchManager) requestGameLevelConfig() bool {
	logrus.Debugln("进入函数")

	exposer := structs.GetGlobalExposer()

	// 获取hall的connection
	hallConnection, err := exposer.RPCClient.GetConnectByServerName("hall")
	if err != nil || hallConnection == nil {
		logrus.WithError(err).Errorln("获得hall服的gRPC失败!!!")
		return false
	}

	hallClient := user.NewPlayerDataClient(hallConnection)

	// 调用room服的创建桌子
	rsp, err := hallClient.GetGameListInfo(context.Background(), &user.GetGameListInfoReq{})

	// 不成功时，报错
	if err != nil || rsp == nil {
		logrus.WithError(err).Errorln("从hall服获取游戏场次配置信息失败!!!")
		return false
	}

	// 返回的不是成功，报错
	if rsp.GetErrCode() != int32(user.ErrCode_EC_SUCCESS) {
		logrus.WithError(err).Errorln("从hall服获取游戏场次配置信息成功，但errCode显示失败")
		return false
	}

	// 游戏配置
	rspGameConfig := rsp.GetGameConfig()
	logrus.Debugf("hall服发送了%v个游戏配置信息", len(rspGameConfig))
	for i := 0; i < len(rspGameConfig); i++ {
		pGameConf := rspGameConfig[i]

		// 游戏需不存在
		_, exist := manager.allGame[pGameConf.GetGameId()]
		if exist {
			logrus.Errorln("游戏ID:%v存在重复", pGameConf.GetGameId())
			return false
		}

		// 新游戏配置信息
		newGameConf := gameConfig{
			gameID:      pGameConf.GetGameId(),
			gameName:    pGameConf.GetGameName(),
			levelConfig: map[uint32]gameLevelConfig{},
		}

		// 加入该游戏
		manager.allGame[pGameConf.GetGameId()] = gameInfo{
			allLevelChan: map[uint32]chan reqMatchPlayer{},
			config:       newGameConf}
	}

	// 场次配置
	rspLevelConfig := rsp.GetGameLevelConfig()
	logrus.Debugf("hall服发送了%v个场次配置信息", len(rspLevelConfig))
	for i := 0; i < len(rspLevelConfig); i++ {
		pLevelConf := rspLevelConfig[i]

		// 游戏需存在
		gInfo, exist := manager.allGame[pLevelConf.GetGameId()]
		if !exist {
			logrus.Errorln("游戏ID:%v中不存在", pLevelConf.GetGameId())
			return false
		}

		// 场次需不存在
		_, exist = gInfo.config.levelConfig[pLevelConf.GetLevelId()]
		if exist {
			logrus.Errorln("游戏ID:%v中场次ID:%v存在重复", pLevelConf.GetGameId(), pLevelConf.GetLevelId())
			return false
		}

		// 新场次配置信息
		newLevelConf := gameLevelConfig{
			levelID:            pLevelConf.GetLevelId(),
			levelName:          pLevelConf.GetLevelName(),
			bottomScore:        pLevelConf.GetBaseScores(),
			minGold:            int64(pLevelConf.GetLowScores()),
			maxGold:            int64(pLevelConf.GetHighScores()),
			minNeedPlayerCount: pLevelConf.GetMinPeople(),
			maxNeedPlayerCount: pLevelConf.GetMaxPeople(),
		}

		// 加入该场次
		gInfo.config.levelConfig[pLevelConf.GetLevelId()] = newLevelConf

		// 该场次的申请通道
		gInfo.allLevelChan[pLevelConf.GetLevelId()] = make(chan reqMatchPlayer, 1024)
	}

	logrus.Debugf("接收hall服的游戏配置信息结束，配置如下：%v", manager.allGame)

	logrus.Debugln("离开函数")
	return true
}

// start 开始为所有的游戏，所有的场次开启协程，执行匹配过程
func (manager *matchManager) start() {
	logrus.Debugln("进入函数")

	// 未初始化完毕，直接返回
	if manager.bInitFinish == false {
		logrus.Errorln("开始匹配时发现仍未初始化完毕，返回")
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

	logrus.Debugln("离开函数")
	return
}

// 生成唯一的桌子ID
func (manager *matchManager) generateDeskID() uint64 {

	atomic.AddUint64(&manager.deskStartID, 1)

	logrus.Debugln("产生桌子唯一ID:", manager.deskStartID)

	return manager.deskStartID
}

// 返回指定间隔秒数后的胜率差异浮动值，百分比
// interval为0时表示[0~1)秒时(也是首次匹配)的胜率浮动值
// interval为1时表示[1~2)秒时的胜率浮动值
// interval为2时表示[2~3)秒时的胜率浮动值
func (manager *matchManager) getWinRateValue(interval int64) int8 {
	logEntry := logrus.WithFields(logrus.Fields{
		"interval": interval,
	})

	if interval < 0 {
		logEntry.Errorf("参数错误，interval < 0 !!")
		return 100
	}

	// 有效时间内
	if interval >= 0 && uint32(interval) <= web.GetMaxCompuValidTime() {
		result := int8(manager.rateCompuValue[interval] * 100)
		logEntry.Debugln("返回的胜率浮动值:", result)
		return result
	}

	// 超过有效时间，认为是正无穷，即所有胜率
	logEntry.Debugln("返回的胜率浮动值:100")
	return 100
}

// 返回指定间隔秒数后的金币差异浮动值
// interval为0时表示0秒时(即首次匹配)的金币浮动值
// interval为1时表示1秒后的金币浮动值
// interval为2时表示2秒后的金币浮动值
func (manager *matchManager) getGoldValue(interval int64) float32 {
	logEntry := logrus.WithFields(logrus.Fields{
		"interval": interval,
	})

	if interval < 0 {
		logEntry.Errorf("参数错误，interval < 0 !!")
		return 10000
	}

	// 有效时间内
	if interval >= 0 && uint32(interval) <= web.GetMaxCompuValidTime() {
		result := manager.goldCompuValue[interval]
		logEntry.Debugln("返回的金币浮动值:", result)
		return result
	}

	// 超过有效时间，认为是正无穷，即所有金币数都可以
	logEntry.Debugln("返回的胜率浮动值:10000")
	return 10000
}

// 检测两个桌子玩家的IP是否有重复的
// true表有重复的，false表无重复的
func (manager *matchManager) checkDeskSameIP(pDesk1 *matchDesk, pDesk2 *matchDesk) bool {

	// 参数检测
	if pDesk1 == nil || pDesk2 == nil {
		logrus.Errorln("参数错误，pDesk1 == nil || pDesk2 == nil 返回")
		return false
	}

	for i := 0; i < len(pDesk1.players); i++ {
		for j := 0; i < len(pDesk2.players); j++ {
			// 存在IP相等的即返回
			if pDesk1.players[i].IP == pDesk2.players[j].IP {
				return true
			}
		}
	}

	return false
}

// 检测玩家和桌子的IP是否有重复
// true表有重复的，false表无重复的
func (manager *matchManager) checkPlayerSameIP(pPlayer *matchPlayer, pDesk *matchDesk) bool {

	// 参数检测
	if pPlayer == nil || pDesk == nil {
		logrus.Errorln("参数错误，pPlayer == nil || pDesk == nil 返回")
		return false
	}

	for i := 0; i < len(pDesk.players); i++ {
		// 存在IP相等的即返回
		if pPlayer.IP == pDesk.players[i].IP {
			return true
		}
	}

	return false
}

// 检测两个桌子玩家的上一局是否有同桌的
// true表有同桌的，false表无同桌的
func (manager *matchManager) checkDeskLastSameDesk(pDesk1 *matchDesk, pDesk2 *matchDesk, pGlobalInfo *levelGlobalInfo) bool {

	// 参数检测
	if pDesk1 == nil || pDesk2 == nil || pGlobalInfo == nil {
		logrus.Errorln("参数错误，pDesk1 == nil || pDesk2 == nil || pGlobalInfo == nil，返回")
		return false
	}

	for i := 0; i < len(pDesk1.players); i++ {
		// 某个玩家上一局有同桌即返回
		if manager.checkPlayerLastSameDesk(&pDesk1.players[i], pDesk2, pGlobalInfo) {
			return true
		}
	}

	return false
}

// 检测玩家和桌子内的玩家是否有上一局同桌的
// true表有同桌的，false表无同桌的
func (manager *matchManager) checkPlayerLastSameDesk(pPlayer *matchPlayer, pDesk *matchDesk, pGlobalInfo *levelGlobalInfo) bool {

	// 参数检测
	if pPlayer == nil || pDesk == nil || pGlobalInfo == nil {
		logrus.Errorln("参数错误，pPlayer == nil || pDesk == nil || pGlobalInfo == nil，返回")
		return false
	}

	// 自己是否存在上一局
	selfDeskID, selfExist := pGlobalInfo.sucPlayers[pPlayer.playerID]
	if !selfExist {
		return false
	}

	tNowTime := time.Now().Unix()

	for i := 0; i < len(pDesk.players); i++ {

		deskID, exist := pGlobalInfo.sucPlayers[pDesk.players[i].playerID]

		// 上一局的桌子ID相同，说明同桌
		if exist && selfDeskID == deskID {

			// 找到该桌子
			desk, exist := pGlobalInfo.sucDesks[deskID]
			if exist {
				// 接着检测成功时间，若未超过同桌限制时间，认为是上局同桌
				if tNowTime-desk.sucTime < web.GetSameDeskLimitTime() {
					return true
				}
			}
		}
	}

	return false
}

// 检测指定金币数和指定桌子是否匹配
// true表匹配，false表不匹配
func (manager *matchManager) checkGoldMatchDesk(goldNum int64, pDesk *matchDesk) bool {

	// 参数检测
	if pDesk == nil {
		logrus.Errorln("checkGoldMatchDesk() 参数错误，pDesk == nil，返回")
		return false
	}

	nowTime := time.Now().Unix()

	// 金币差异度(百分比)
	goldDiff := (float64(goldNum) - float64(pDesk.aveGold)) / float64(pDesk.aveGold)

	// 检测金币范围
	if (float32(goldDiff)) > manager.getGoldValue(nowTime-pDesk.createTime) {
		return false
	}

	return true
}

// 为指定玩家执行首次匹配
func (manager *matchManager) firstMatch(globalInfo *levelGlobalInfo, reqPlayer *reqMatchPlayer) {
	if globalInfo == nil || reqPlayer == nil {
		logrus.Errorln("参数错误，globalInfo == nil || reqPlayer == nil 返回")
		return
	}

	logEntry := logrus.WithFields(logrus.Fields{
		"reqPlayer": reqPlayer,
	})

	logEntry.Debugln("进入函数")

	// 胜率范围检测
	if reqPlayer.winRate < 0 || reqPlayer.winRate > 100 {
		logEntry.Errorf("错误，玩家%v的胜率为%v，不再执行匹配\n", reqPlayer.playerID, reqPlayer.winRate)
		return
	}

	// 找到的匹配的桌子
	var pFindDesk *matchDesk = nil

	// 新建一个匹配玩家
	newMatchPlayer := matchPlayer{
		playerID: reqPlayer.playerID,
		robotLv:  0,
		seat:     -1,
		IP:       IPStringToUInt32(reqPlayer.IP),
		gold:     reqPlayer.gold,
	}

	nowTime := time.Now().Unix()

	///////////////////////////////////////////////////////// 先检测首次匹配范围是不是有桌子 ///////////////////////////////////////////////

	// 首次匹配的胜率值
	firstRateValue := manager.getWinRateValue(0)

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

	logEntry.Debugf("玩家首次匹配胜率范围：(%v - %v)", playerBeginRate, playerEndRate)

	// 从低往高匹配，存在桌子即可加入，这些桌子的胜率都在玩家的匹配范围
	for i := playerBeginRate; i <= playerEndRate; i++ {

		// 遍历该概率下所有的桌子
		for iter := globalInfo.allRateDesks[i].Front(); iter != nil; iter = iter.Next() {

			desk := iter.Value.(*matchDesk)

			// 检测金币范围
			if !manager.checkGoldMatchDesk(reqPlayer.gold, desk) {
				continue
			}

			// 检测IP地址
			if manager.checkPlayerSameIP(&newMatchPlayer, desk) {
				continue
			}

			// 检测上局是否同桌
			if manager.checkPlayerLastSameDesk(&newMatchPlayer, desk, globalInfo) {
				continue
			}

			// 可以进桌子了，再找到创建时间最早的那个

			// 比较桌子创建时间，记录创建时间最早的
			if (pFindDesk == nil) || (pFindDesk != nil && desk.createTime < pFindDesk.createTime) {
				pFindDesk = desk
			}
		}
	}

	// 找到的话，则加入桌子，返回
	if pFindDesk != nil {
		// 把玩家加入桌子
		if manager.addPlayerToDesk(&newMatchPlayer, pFindDesk) == false {
			logEntry.Errorf("首次范围检测时，玩家%v加入桌子失败，返回\n", reqPlayer.playerID)
			return
		}

		// 成功入桌即返回
		logEntry.Debugf("首次范围检测时，胜率为%v的玩家%v匹配进桌子%v，正常返回", reqPlayer.winRate, reqPlayer.playerID, pFindDesk)
		return
	}

	///////////////////////////////////  首次范围失败后，再检测那些不在首次范围的，但因胜率范围扩张造成现在可能在玩家匹配范围了  /////////////////////////////

	// 剩下的需要检测的胜率值，也是allWinRate的下标值
	lastIndexs := make([]int8, 0, 101)

	// 从(playerBeginRate - 0]
	for i := playerBeginRate - 1; i >= int8(0); i-- {
		lastIndexs = append(lastIndexs, i)
	}

	// 从(playerEndRate - 100]
	for i := playerEndRate + 1; i <= 100; i++ {
		lastIndexs = append(lastIndexs, i)
	}

	// 遍历lastIndexs
	for i := 0; i < len(lastIndexs); i++ {
		index := int8(lastIndexs[i])

		// 遍历该概率下的所有桌子
		for iter := globalInfo.allRateDesks[index].Front(); iter != nil; iter = iter.Next() {

			desk := iter.Value.(*matchDesk)

			// 距离桌子创建时间的间隔
			interval := nowTime - desk.createTime

			// 不足一秒的不用比较，因为既然玩家上面已经比较过了，不足1秒的桌子肯定是不符合的
			if interval < 1 {
				continue
			}

			// 该桌子的胜率浮动值
			rateValue := manager.getWinRateValue(interval)

			// 该桌子的浮动值下限
			deskBeginRate := index - int8(rateValue)
			if deskBeginRate < 0 {
				deskBeginRate = 0
			}

			// 该桌子的浮动值上限
			deskEndRate := index + int8(rateValue)
			if deskEndRate > 100 {
				deskEndRate = 100
			}

			// 玩家的匹配范围和桌子的匹配范围无交集时说明不匹配
			if (playerBeginRate > deskEndRate) || (playerEndRate < deskBeginRate) {
				continue
			}

			// 检测金币范围
			if !manager.checkGoldMatchDesk(reqPlayer.gold, desk) {
				continue
			}

			// 检测IP地址
			if manager.checkPlayerSameIP(&newMatchPlayer, desk) {
				continue
			}

			// 检测上局是否同桌
			if manager.checkPlayerLastSameDesk(&newMatchPlayer, desk, globalInfo) {
				continue
			}

			// 可以进桌子了，再找到创建时间最早的那个

			// 比较桌子创建时间，记录创建时间最早的
			if (pFindDesk == nil) || (pFindDesk != nil && desk.createTime < pFindDesk.createTime) {
				pFindDesk = desk
			}
		}
	}

	// 找到的话，则加入桌子，返回
	if pFindDesk != nil {
		// 把玩家加入桌子
		if manager.addPlayerToDesk(&newMatchPlayer, pFindDesk) == false {
			logEntry.Errorf("二次范围检测时，玩家%v加入桌子失败，返回\n", reqPlayer.playerID)
			return
		}

		// 成功入桌即返回
		logEntry.Debugf("二次范围检测时，胜率为%v的玩家%v匹配进桌子%v，正常返回", reqPlayer.winRate, reqPlayer.playerID, pFindDesk)
		return
	}

	//////////////////////////////////////////////////////////  所有的桌子都失败了，新建桌子  ////////////////////////////////////////////////
	// 创建桌子
	needPlayerCount := manager.getGameNeedPlayerCount(globalInfo.gameID, globalInfo.levelID)
	// 桌子唯一ID
	deskID := manager.generateDeskID()
	newDesk := createMatchDesk(deskID, globalInfo.gameID, globalInfo.levelID, needPlayerCount, reqPlayer.gold)
	if newDesk == nil {
		logEntry.Errorf("创建匹配桌子失败，返回")
		return
	}

	// 把该玩家压入桌子
	if manager.addPlayerToDesk(&newMatchPlayer, newDesk) == false {
		logEntry.Errorf("新建桌子后，玩家%v加入桌子失败，返回\n", reqPlayer.playerID)
		return
	}

	logEntry.Debugf("胜率为%v的玩家%v匹配失败后，创建桌子并加入，正常返回", reqPlayer.winRate, reqPlayer.playerID)

	logEntry.Debugln("离开函数")
	return
}

// startLevelMatch 开始单个游戏单个场次的匹配
// gameID : 游戏ID
// levelID : 场次ID
func (manager *matchManager) startLevelMatch(gameID uint32, levelID uint32) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "singleLevelMatch",
		"gameID":    gameID,
		"levelID":   levelID,
	})

	logEntry.Debugln("进入函数")

	// 该场次的全局信息
	globalInfo := levelGlobalInfo{
		gameID:       gameID,
		levelID:      levelID,
		allRateDesks: make([]list.List, 101), // 胜率1% - 100%的所有匹配桌子
		sucPlayers:   map[uint64]uint64{},    // 已成功匹配的玩家，Key:玩家ID，Value:桌子ID
		sucDesks:     map[uint64]*sucDesk{},  // 已成功匹配的桌子，Key:桌子ID，Value:桌子信息
	}

	// 2秒1次的合并定时器
	mergeTimer := time.NewTicker(time.Second * 1)

	// 1秒1次的超时定时器
	timeoutTimer := time.NewTicker(time.Second * 1)

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
				manager.firstMatch(&globalInfo, &req) // 首次匹配
			}
		//case pl := <-manager.loginChannel: // 登录玩家
		//	{
		//		manager.onPlayerLogin(pl.playerID)
		//	}
		case <-mergeTimer.C: // 合并定时器
			{
				//manager.mergeDesks(&globalInfo)
			}
		case <-timeoutTimer.C: // 超时定时器
			{
				//manager.checkTimeout(&globalInfo)
			}
		}
	}

	logEntry.Debugln("离开函数")
	return
}

/* // addContinueDesk 添加续局牌桌
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
} */

// addPlayer 添加匹配玩家
func (manager *matchManager) addPlayer(playerID uint64, gameID int) {
	manager.applyChannel <- reqMatchPlayer{
		playerID: playerID,
	}
	return
}

// 分发匹配请求
// playerID 	:	玩家ID
// gameID		：	请求匹配的游戏ID
// levelID		:   请求匹配的级别ID
// 返回string 	 ：	 返回的错误描述，成功时返回空
func (manager *matchManager) dispatchMatchReq(playerID uint64, gameID uint32, levelID uint32) string {
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

	// 获取玩家金币数
	playerGold, err := goldclient.GetGold(playerID, int16(gold.GoldType_GOLD_COIN))
	if err != nil {
		logrus.Errorln("从gold服获取玩家金币失败")
		return fmt.Sprintf("从gold服获取玩家金币失败，游戏ID:%v，场次ID:%v，请求的玩家ID:%v", gameID, levelID, playerID)
	}

	// 金币范围检测
	if playerGold < levelConfig.minGold {
		logrus.Errorln("玩家金币数小于游戏场次金币要求最小值，最小值：%v", levelConfig.minGold)
		return fmt.Sprintf("玩家金币数小于游戏场次金币要求最小值，请求匹配的游戏ID:%v，场次ID:%v，玩家ID:%v", gameID, levelID, playerID)
	}

	if playerGold > levelConfig.maxGold {
		logrus.Errorf("玩家金币数大于游戏场次金币要求最大值，最大值：%v", levelConfig.maxGold)
		return fmt.Sprintf("玩家金币数大于游戏场次金币要求最大值，请求匹配的游戏ID:%v，场次ID:%v，玩家ID:%v", gameID, levelID, playerID)
	}

	// 获取该游戏的胜率
	// 计算该游戏的胜率，已经乘以100,比如：50表胜率为50%
	playerWinRate, err := requestPlayerWinRate(playerID, gameID)
	if err != nil {
		logrus.Errorln("从hall服获取玩家胜率失败")
		return fmt.Sprintf("从hall服获取玩家胜率失败，游戏ID:%v，场次ID:%v，请求的玩家ID:%v", gameID, levelID, playerID)
	}

	// 全部检测通过

	// 获取该场次的申请通道
	reqMatchChan, exist := gameInfo.allLevelChan[levelID]
	if !exist {
		logrus.Errorln("内部错误，请求匹配的游戏存在，场次存在，但找不到该场次的匹配申请通道")
		return fmt.Sprintf("请求匹配的游戏存在，场次存在，但找不到该场次的匹配申请通道，请求匹配的游戏ID:%v，场次ID:%v，玩家ID:%v", gameID, levelID, playerID)
	}

	clientIP, _ := requestPlayerIP(playerID)

	// 压入通道
	reqMatchChan <- reqMatchPlayer{
		playerID: playerID,            // playerID
		winRate:  int8(playerWinRate), // 胜率
		gold:     playerGold,          // 金币数
		IP:       clientIP,            // IP地址
	}

	logEntry.Debugln("离开函数")

	return ""
}

/* // addContinueApply 添加续局申请
func (manager *matchManager) addContinueApply(playerID uint64, cancel bool, gameID int) {
	manager.continueChannel <- continueApply{
		playerID: playerID,
		cancel:   cancel,
		gameID:   gameID,
	}
	return
} */

// addLoginData 添加玩家登录信息
func (manager *matchManager) addLoginData(playerID uint64) {
	manager.loginChannel <- playerLogin{
		playerID: playerID,
	}
}

/* // run 执行匹配流程
func (manager *matchManager) run() {

	// 从DB读取游戏配置信息
	// todo

	// 机器人的定时器（1秒1次）
	robotTick := time.NewTicker(time.Second * 1)

	for {
		select {
		case ap := <-manager.applyChannel: // 普通匹配申请
			{
				manager.acceptApplyPlayer(ap.gameID, ap.playerID)
			}
		case cp := <-manager.continueChannel: // 续局匹配申请
			{
				manager.acceptContinuePlayer(cp.gameID, cp.playerID, cp.cancel)
			}
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
} */

/* // onPlayerLogin 玩家登录，取消玩家匹配
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
} */

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

	// 原总金币数
	oldAllGold := pDesk.aveGold * int64(len(pDesk.players))

	// 压入该玩家
	pDesk.players = append(pDesk.players, *pPlayer)

	// 计算平均金币
	pDesk.aveGold = (oldAllGold + pPlayer.gold) / int64(len(pDesk.players))

	logrus.Debugf("桌子%v压入了玩家%v\n", pDesk, pPlayer)

	// playerID与deskID的映射
	//manager.playerDesk[deskPlayer.playerID] = desk.deskID

	// 移除不在线的
	//manager.removeOfflines(desk)

	// 满桌需要的玩家数量
	needPlayerCount := manager.getGameNeedPlayerCount(pDesk.gameID, pDesk.levelID)

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

/* // removeOfflines 移除 desk 中的离线玩家
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
*/
/* // onDeskFinish 牌桌匹配完成
func (manager *matchManager) onDeskFinish(desk *desk) {
	requestCreateDesk(desk)
	players := desk.players
	// 解除关联
	for _, player := range players {
		delete(manager.playerDesk, player.playerID)
	}
	// 移除 desk
	delete(manager.desks, desk.deskID)
} */

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

// deleteDesk 删除指定的桌子
func (manager *matchManager) deleteDesk(pDesk *matchDesk) {

}

/* // handleRobotTick 处理机器人 tick
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
} */

// mergeDesks 合并桌子
func (manager *matchManager) mergeDesks(globalInfo *levelGlobalInfo) {
	if globalInfo == nil {
		logrus.Errorln("mergeDesks()，globalInfo == nil，返回")
		return
	}

	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "mergeDesks",
		"gameID":    globalInfo.gameID,
		"levelID":   globalInfo.levelID,
	})

	logEntry.Debugln("进入函数")

	// 当前时间
	tNowTime := time.Now().Unix()

	// 所有的概率
	var index int8 = 0
	for ; index <= 100; index++ {
		// 该概率下所有的桌子
		for iter := globalInfo.allRateDesks[index].Front(); iter != nil; iter = iter.Next() {

			desk := iter.Value.(*matchDesk)

			// 距离桌子创建时间的间隔
			interval := tNowTime - desk.createTime

			// 不足1秒的，不检测，因为新建一个桌子时已检测过
			if interval < 1 {
				continue
			}

			// 前一秒的胜率浮动值
			lastRateValue := manager.getWinRateValue(interval - 1)

			// 这一秒的胜率浮动值
			nowRateValue := manager.getWinRateValue(interval)

			// 所有的需要检测的胜率值，也是allRateDesks的下标值
			checkIndexs := make([]int8, 0, 101)

			// 左段起始值(包含自身)
			leftStartIndex := index - int8(nowRateValue)
			if leftStartIndex < 0 {
				leftStartIndex = 0
			}
			// 左段结束值(不包含自身)
			leftEndIndex := index - int8(lastRateValue)
			if leftEndIndex < 0 {
				leftEndIndex = 0
			}
			// 从[leftStartIndex - leftEndIndex)
			for j := leftStartIndex; j < leftEndIndex; j++ {
				checkIndexs = append(checkIndexs, j)
			}

			// 右段起始值(不包含自身)
			rightStartIndex := index + int8(lastRateValue)
			if leftStartIndex > 100 {
				leftStartIndex = 100
			}
			// 右段结束值(包含自身)
			rightEndIndex := index + int8(nowRateValue)
			// 从[rightStartIndex - rightEndIndex)
			for j := rightStartIndex + 1; j <= rightEndIndex; j++ {
				checkIndexs = append(checkIndexs, j)
			}

			var pMergeDesk *matchDesk = nil

			// 和这些桌子尝试组合
			for k := 0; k < len(checkIndexs); k++ {
				merIndex := int8(checkIndexs[k])

				// 遍历该概率下的所有桌子
				for merIter := globalInfo.allRateDesks[merIndex].Front(); merIter != nil; merIter = merIter.Next() {

					merDesk := merIter.Value.(*matchDesk)

					// IP是否存在相同的
					if manager.checkDeskSameIP(desk, merDesk) {
						continue
					}

					// 上局是否存在同桌的
					if manager.checkDeskLastSameDesk(desk, merDesk, globalInfo) {
						continue
					}

					// 可以合并
					pMergeDesk = merDesk
					break
				}
			}

			// 有合并的桌子
			if pMergeDesk != nil {
				// 合并操作
				// 根据两个桌子的创建时间，把时间短的桌子拆了，玩家添加到时间长的桌子;若有剩余玩家，继续留在当前桌子;若没有剩余玩家，
				// todo
			}
		}
	}

	logEntry.Debugln("离开函数")
}

// checkTimeout 检测超时
func (manager *matchManager) checkTimeout(globalInfo *levelGlobalInfo) {
	if globalInfo == nil {
		logrus.Errorln("checkTimeout()，globalInfo == nil，返回")
		return
	}

	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "checkTimeout",
		"gameID":    globalInfo.gameID,
		"levelID":   globalInfo.levelID,
	})

	// 当前时间
	tNowTime := time.Now().Unix()

	// 所有的概率
	var index int8 = 0
	for ; index <= 100; index++ {
		// 该概率下所有的桌子
		for iter := globalInfo.allRateDesks[index].Front(); iter != nil; iter = iter.Next() {

			desk := iter.Value.(*matchDesk)

			if tNowTime-desk.createTime > 100 {

			}
		}
	}

	logEntry.Debugln("进入函数")
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
