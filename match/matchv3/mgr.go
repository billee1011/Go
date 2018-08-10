package matchv3

import (
	"container/list"
	"context"
	"fmt"
	"steve/external/goldclient"
	"steve/external/hallclient"
	"steve/external/robotclient"
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
	isMatch  bool   // 申请匹配为true，取消匹配为false
	playerID uint64 // 玩家 ID
	winRate  int32  // 具体游戏的胜率
	gold     int64  // 金币数
	IP       uint32 // IP地址
}

// playerOffline 玩家离线
type playerLogin struct {
	playerID uint64 // 玩家 ID
}

// gameLevelConfig 游戏场次配置数据
type gameLevelConfig struct {
	levelID     uint32 // 场次ID
	levelName   string // 场次名字
	bottomScore uint32 // 底分
	minGold     int64  // 金币要求下限
	maxGold     int64  // 金币要求上限
}

// gameConfig 游戏配置数据
type gameConfig struct {
	gameID             uint32                     // 游戏ID
	gameName           string                     // 游戏名字
	minNeedPlayerCount uint32                     // 允许最低人数
	maxNeedPlayerCount uint32                     // 允许最高人数
	levelConfig        map[uint32]gameLevelConfig // 所有的游戏场次
}

// levelGlobalInfo 单个游戏单个场次的全局信息
type levelGlobalInfo struct {
	// 游戏ID
	gameID uint32

	// 场次ID
	levelID uint32

	// 满桌所需人数
	needPlayerCount uint8

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
	allLevelChan map[uint32]chan reqMatchPlayer // 本levelID戏所有场次的匹配/取消匹配通道,Key:场次ID, Value:该场次的匹配/取消匹配通道
	config       gameConfig                     // 游levelID配置数据
}

// matchManager 匹配管理器
type matchManager struct {
	bInitFinish    bool                // 是否初始化完成
	allGame        map[uint32]gameInfo // 所有游戏的匹配信息，Key:游戏ID, Value:该游戏的匹配信息
	deskStartID    uint64              // 桌子ID起始值
	rateCompuValue []float32           // 胜率计算配置
	goldCompuValue []float32           // 金币计算配置
}

// matchMgr 匹配管理器
var matchMgr = &matchManager{
	bInitFinish:    false,
	allGame:        map[uint32]gameInfo{},
	deskStartID:    0,
	rateCompuValue: make([]float32, 0, web.GetMaxCompuValidTime()+1),
	goldCompuValue: make([]float32, 0, web.GetMaxCompuValidTime()+1),
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

	/* 	// 得到该场次的信息
	   	levelInfo, exist := gameInfo.config.levelConfig[levelID]
	   	// 该场次不存在
	   	if !exist {
	   		logEntry.Errorln("场次ID不存在")
	   		return 0
	   	} */

	logEntry.Debugf("离开函数,最低满桌人数为：%v", gameInfo.config.minNeedPlayerCount)

	return uint8(gameInfo.config.minNeedPlayerCount)
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
			gameID:             pGameConf.GetGameId(),
			gameName:           pGameConf.GetGameName(),
			minNeedPlayerCount: pGameConf.GetMinPeople(),
			maxNeedPlayerCount: pGameConf.GetMaxPeople(),
			levelConfig:        map[uint32]gameLevelConfig{},
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
			levelID:     pLevelConf.GetLevelId(),
			levelName:   pLevelConf.GetLevelName(),
			bottomScore: pLevelConf.GetBaseScores(),
			minGold:     int64(pLevelConf.GetLowScores()),
			maxGold:     int64(pLevelConf.GetHighScores()),
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
func (manager *matchManager) getWinRateValue(interval int64) int32 {
	logEntry := logrus.WithFields(logrus.Fields{
		"interval": interval,
	})

	if interval < 0 {
		logEntry.Errorf("参数错误，interval < 0 !!")
		return 100
	}

	// 有效时间内
	if interval >= 0 && uint32(interval) <= web.GetMaxCompuValidTime() {
		result := int32(manager.rateCompuValue[interval] * 100)
		//logEntry.Debugln("返回的胜率浮动值:", result)
		return result
	}

	// 超过有效时间，认为是正无穷，即所有胜率
	//logEntry.Debugln("返回的胜率浮动值:100")
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
	logEntry.Debugln("返回的金币浮动值:10000")
	return 10000
}

// 检测两个桌子玩家的IP是否有重复的，只比较真实玩家
// true表有重复的，false表无重复的
func (manager *matchManager) checkDeskSameIP(pDesk1 *matchDesk, pDesk2 *matchDesk) bool {

	// 参数检测
	if pDesk1 == nil || pDesk2 == nil {
		logrus.Errorln("参数错误，pDesk1 == nil || pDesk2 == nil 返回")
		return false
	}

	// 若不限制相同IP，返回false
	if !web.GetLimitSameIP() {
		return false
	}

	for i := 0; i < len(pDesk1.players); i++ {
		for j := 0; i < len(pDesk2.players); j++ {
			// 都是真实玩家，存在IP相等的即返回
			if (pDesk1.players[i].robotLv == 0) && (pDesk2.players[i].robotLv == 0) && (pDesk1.players[i].IP == pDesk2.players[j].IP) {
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

	// 若不限制相同IP，返回false
	if !web.GetLimitSameIP() {
		return false
	}

	// 自己是机器人，不检测
	if pPlayer.robotLv != 0 {
		return false
	}

	for i := 0; i < len(pDesk.players); i++ {
		// 对方是真实玩家，且存在IP相等的即返回
		if (pDesk.players[i].robotLv == 0) && (pPlayer.IP == pDesk.players[i].IP) {
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

	// 若不限制上局同桌，返回false
	if !web.GetLimitLastSameDesk() {
		return false
	}

	for i := 0; i < len(pDesk1.players); i++ {
		// 只检测真实玩家
		if pDesk1.players[i].robotLv == 0 {
			// 与pDesk2的某个玩家上一局有同桌即返回
			if manager.checkPlayerLastSameDesk(&pDesk1.players[i], pDesk2, pGlobalInfo) {
				return true
			}
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

	// 若不限制上局同桌，返回false
	if !web.GetLimitLastSameDesk() {
		return false
	}

	// 自己是机器人，不检测
	if pPlayer.robotLv != 0 {
		return false
	}

	// 自己不存在上一局，不检测
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
				// 接着检测距离上次同桌的时间间隔，若未超过同桌限制时间，认为是上局同桌
				if tNowTime-desk.sucTime < web.GetSameDeskLimitTime() {
					return true
				}
			}
		}
	}

	return false
}

// 获取指定金币数经过间隔为interval秒后的金币匹配范围
// goldNum  : 金币数
// inteval  : 间隔的秒数
// 返回值
// minGold  : 金币最小值
// maxGold  : 金币最大值
func (manager *matchManager) getGoldRange(goldNum int64, inteval int64) (minGold int64, maxGold int64) {

	// 差异百分比
	goldValue := manager.getGoldValue(inteval)

	// 金币匹配范围最小值
	minGold = int64(float64(goldNum) * float64(1-goldValue))
	if minGold < 0 {
		minGold = 0
	}

	//金币匹配范围最大值
	maxGold = int64(float64(goldNum) * float64(1+goldValue))

	return
}

// 获取指定的胜率经过间隔为interval秒后的胜率匹配范围
// rate     : 指定的胜率
// inteval  : 间隔的秒数
// 返回值
// minRate  : 胜率起始值
// maxRate  : 金币最大值
func (manager *matchManager) getWinRateRange(rate int32, inteval int64) (minRate int32, maxRate int32) {

	// 胜率浮动值
	rateValue := manager.getWinRateValue(inteval)

	// 胜率匹配范围最小值
	minRate = rate - rateValue
	if minRate < 0 {
		minRate = 0
	}

	//桌子金币匹配范围最大值
	maxRate = rate + rateValue
	if maxRate > 100 {
		maxRate = 100
	}

	return
}

// 为指定玩家执行首次匹配
func (manager *matchManager) firstMatch(globalInfo *levelGlobalInfo, reqPlayer *reqMatchPlayer) {
	if globalInfo == nil || reqPlayer == nil {
		logrus.Errorln("firstMatch(), 参数错误，globalInfo == nil || reqPlayer == nil 返回")
		return
	}

	logEntry := logrus.WithFields(logrus.Fields{
		"reqPlayer": reqPlayer,
	})

	logEntry.Debugln("进入函数")

	// 胜率范围检测
	if reqPlayer.winRate < 0 || reqPlayer.winRate > 100 {
		logEntry.Errorf("数据错误，玩家%v的胜率为%v，不再执行匹配", reqPlayer.playerID, reqPlayer.winRate)
		return
	}

	// 找到的匹配桌子
	var pFindIter *list.Element = nil

	// 找到的匹配桌子所在index
	var pFindIndex int32 = -1

	// 新建一个匹配玩家
	newMatchPlayer := matchPlayer{
		playerID: reqPlayer.playerID,
		robotLv:  0,
		seat:     -1,
		IP:       reqPlayer.IP,
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

			desk := *(iter.Value.(**matchDesk))

			// 检测金币范围
			minGold, maxGold := manager.getGoldRange(desk.aveGold, nowTime-desk.createTime)
			if reqPlayer.gold < minGold || reqPlayer.gold > maxGold {
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
			if pFindIter == nil {
				pFindIter = iter
				pFindIndex = i
			} else {
				pFindDesk := *(pFindIter.Value.(**matchDesk))
				if desk.createTime < pFindDesk.createTime {
					pFindIter = iter
					pFindIndex = i
				}
			}
		}
	}

	// 找到的话，则加入桌子，返回
	if pFindIter != nil && pFindIndex != -1 {
		pFindDesk := *(pFindIter.Value.(**matchDesk))

		// 把玩家加入桌子，若桌子已满，则从列表中移除
		if manager.addPlayerToDesk(&newMatchPlayer, pFindDesk, globalInfo) {
			logEntry.Debugf("首次范围检测时，胜率为%v的玩家%v匹配进桌子%v，桌子满员，已删除", reqPlayer.winRate, newMatchPlayer, pFindDesk)
			globalInfo.allRateDesks[pFindIndex].Remove(pFindIter)
		} else {
			logEntry.Debugf("首次范围检测时，胜率为%v的玩家%v匹配进桌子%v，桌子未满员，继续匹配", reqPlayer.winRate, newMatchPlayer, pFindDesk)
		}

		// 成功入桌即返回
		return
	}

	///////////////////////////////////  首次范围失败后，再检测那些不在首次范围的，但因胜率范围扩张造成现在可能在玩家匹配范围了  /////////////////////////////

	// 剩下的需要检测的胜率值，也是allWinRate的下标值
	lastIndexs := make([]int32, 0, 101)

	// 从(playerBeginRate - 0]
	for i := playerBeginRate - 1; i >= 0; i-- {
		lastIndexs = append(lastIndexs, i)
	}

	// 从(playerEndRate - 100]
	for i := playerEndRate + 1; i <= 100; i++ {
		lastIndexs = append(lastIndexs, i)
	}

	// 遍历lastIndexs
	for i := 0; i < len(lastIndexs); i++ {
		index := lastIndexs[i]

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
			deskBeginRate := index - rateValue
			if deskBeginRate < 0 {
				deskBeginRate = 0
			}

			// 该桌子的浮动值上限
			deskEndRate := index + rateValue
			if deskEndRate > 100 {
				deskEndRate = 100
			}

			// 玩家的匹配范围和桌子的匹配范围无交集时说明不匹配
			if (playerBeginRate > deskEndRate) || (playerEndRate < deskBeginRate) {
				continue
			}

			// 检测金币范围
			minGold, maxGold := manager.getGoldRange(desk.aveGold, nowTime-desk.createTime)
			if reqPlayer.gold < minGold || reqPlayer.gold > maxGold {
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
			if (pFindIter == nil) || (pFindIter != nil && desk.createTime < pFindIter.Value.(*matchDesk).createTime) {
				pFindIter = iter
				pFindIndex = index
			}
		}
	}

	// 找到的话，则加入桌子，返回
	if pFindIter != nil && pFindIndex != -1 {
		pFindDesk := pFindIter.Value.(*matchDesk)

		// 把玩家加入桌子，若桌子已满，则从列表中移除
		if manager.addPlayerToDesk(&newMatchPlayer, pFindDesk, globalInfo) {

			logEntry.Debugf("二次范围检测时，胜率为%v的玩家%v匹配进桌子%v，桌子满员，已删除", reqPlayer.winRate, newMatchPlayer, pFindDesk)
			globalInfo.allRateDesks[pFindIndex].Remove(pFindIter)
		} else {
			logEntry.Debugf("二次范围检测时，胜率为%v的玩家%v匹配进桌子%v，桌子未满员，继续匹配", reqPlayer.winRate, newMatchPlayer, pFindDesk)
		}

		// 成功入桌即返回
		return
	}

	//////////////////////////////////////////////////////////  所有的桌子都失败了，新建桌子  ////////////////////////////////////////////////
	// 创建桌子
	// 桌子唯一ID
	deskID := manager.generateDeskID()
	newDesk := createMatchDesk(deskID, globalInfo.gameID, globalInfo.levelID, globalInfo.needPlayerCount, reqPlayer.gold)
	if newDesk == nil {
		logEntry.Errorf("创建匹配桌子失败，返回")
		return
	}

	// 把玩家加入桌子，若桌子已满，则从列表中移除
	if manager.addPlayerToDesk(&newMatchPlayer, newDesk, globalInfo) {
		logEntry.Debugf("胜率为%v的玩家%v匹配失败后，创建桌子并加入，桌子满员，不再加入列表", reqPlayer.winRate, newMatchPlayer)
	} else {
		logEntry.Debugf("胜率为%v的玩家%v匹配失败后，创建桌子并加入，桌子未满员，加入列表，继续匹配", reqPlayer.winRate, newMatchPlayer)

		globalInfo.allRateDesks[reqPlayer.winRate].PushBack(&newDesk)
	}

	logEntry.Debugln("离开函数")
	return
}

// 取消玩家匹配
func (manager *matchManager) cancelMatch(globalInfo *levelGlobalInfo, reqPlayer *reqMatchPlayer) {
	if globalInfo == nil || reqPlayer == nil {
		logrus.Errorln("cancelMatch() 参数错误，globalInfo == nil || reqPlayer == nil 返回")
		return
	}

	logEntry := logrus.WithFields(logrus.Fields{
		"reqPlayer": reqPlayer,
	})

	logEntry.Debugln("进入函数")

	// 胜率范围检测
	if reqPlayer.winRate < 0 || reqPlayer.winRate > 100 {
		logEntry.Errorf("数据错误，玩家%v的胜率为%v，不再执行取消匹配", reqPlayer.playerID, reqPlayer.winRate)
		return
	}

	///////////////////////////////////////////////////////// 以玩家胜率为中心，向左右依次搜索，直到所有的桌子 ///////////////////////////////////////////////
	serIndexs := make([]int32, 0, 100)

	// 先压入自身胜率的index
	serIndexs = append(serIndexs, reqPlayer.winRate)

	var i int32 = 0
	var leftIndex int32 = 0
	var rightIndex int32 = 0

	for i = 1; i <= 100; i++ {

		// 左侧,只压入有效的
		leftIndex = reqPlayer.winRate - i
		if leftIndex >= 0 {
			serIndexs = append(serIndexs, leftIndex)
		}

		// 右侧,只压入有效的
		rightIndex = reqPlayer.winRate + i
		if rightIndex <= 100 {
			serIndexs = append(serIndexs, rightIndex)
		}
	}

	bRemovePlayer := false

	var index int32 = 0
	for i := 0; i < len(serIndexs); i++ {

		index = serIndexs[i]

		var next *list.Element = nil

		// 该概率下所有的桌子
		for iter := globalInfo.allRateDesks[index].Front(); iter != nil; iter = next {

			next = iter.Next()

			desk := *(iter.Value.(**matchDesk))

			// 该概率下所有的玩家
			for i := 0; i < len(desk.players); i++ {

				// 找到该玩家
				if desk.players[i].playerID == reqPlayer.playerID {
					logEntry.Debugf("取消匹配时，在胜率为%v的桌子列表中找到了玩家，游戏ID:%v，级别ID:%v，从桌子删除前信息为:%v", index, globalInfo.gameID, globalInfo.levelID, desk)
					bRemovePlayer = true
					break
				}
			}

			// 找到桌子后，删掉该玩家
			if bRemovePlayer {
				tempPlayers := make([]matchPlayer, 0, desk.needPlayerCount)
				for i := 0; i < len(desk.players); i++ {
					// 不是该玩家则压入
					if desk.players[i].playerID != reqPlayer.playerID {
						tempPlayers = append(tempPlayers, desk.players[i])
					}
				}
				desk.players = tempPlayers
				logEntry.Debugf("取消匹配时，在胜率为%v的桌子列表中找到了玩家，游戏ID:%v，级别ID:%v，从桌子删除后信息为:%v", index, globalInfo.gameID, globalInfo.levelID, desk)

				// 剩余是否存在真实玩家
				bExistTruePlayer := false
				for i := 0; i < len(desk.players); i++ {
					if desk.players[i].robotLv == 0 {
						bExistTruePlayer = true
						break
					}
				}

				// 没真实玩家了，就删除桌子
				if !bExistTruePlayer {
					globalInfo.allRateDesks[index].Remove(iter)
					logEntry.Debugf("由于删除该玩家后，桌子里不再有真实玩家，删除桌子")
				}

				// 找到即跳出
				break
			}
		}

		// 找到即跳出
		if bRemovePlayer {
			break
		}
	}

	// 已删除玩家，重置玩家状态
	if bRemovePlayer {
		// 设置为空闲状态
		bSuc, err := hallclient.UpdatePlayerState(reqPlayer.playerID, user.PlayerState_PS_MATCHING, user.PlayerState_PS_IDIE, 0, 0)
		if err != nil || !bSuc {
			logEntry.WithError(err).Errorln("内部错误，通知hall服设置玩家状态为空闲状态时失败")
			return
		}

		// 更新玩家所在的服务器类型和地址，地址置空
		bSuc, err = hallclient.UpdatePlayeServerAddr(reqPlayer.playerID, user.ServerType_ST_MATCH, "")
		if err != nil || !bSuc {
			logEntry.WithError(err).Errorln("内部错误，通知hall服设置玩家的服务器类型及地址时失败")
			return
		}
	} else {
		// 没找到该玩家，报错
		logEntry.Errorf("玩家取消匹配时在匹配桌子中找不到该玩家，游戏ID:%v，级别ID:%v", globalInfo.gameID, globalInfo.levelID)
	}
}

// startLevelMatch 开始单个游戏单个场次的匹配
// gameID : 游戏ID
// levelID : 场次ID
func (manager *matchManager) startLevelMatch(gameID uint32, levelID uint32) {
	logEntry := logrus.WithFields(logrus.Fields{
		"gameID":  gameID,
		"levelID": levelID,
	})

	logEntry.Debugln("进入函数")

	// 该场次的全局信息
	globalInfo := levelGlobalInfo{
		gameID:          gameID,
		levelID:         levelID,
		needPlayerCount: manager.getGameNeedPlayerCount(gameID, levelID),
		allRateDesks:    make([]list.List, 101), // 胜率1% - 100%的所有匹配桌子
		sucPlayers:      map[uint64]uint64{},    // 已成功匹配的玩家，Key:玩家ID，Value:桌子ID
		sucDesks:        map[uint64]*sucDesk{},  // 已成功匹配的桌子，Key:桌子ID，Value:桌子信息
	}

	// 1秒1次的合并定时器
	mergeTimer := time.NewTicker(time.Second * 1)

	// 1秒1次的超时定时器(检测匹配桌子是否超时，添加机器人)
	deskTimer := time.NewTicker(time.Second * 1)

	// 60秒1次的超时定时器(检测之前成功的玩家是否超时)
	sucTimer := time.NewTicker(time.Second * 60)

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
		case req := <-reqMatchChan: // 匹配/取消匹配申请
			{
				// 分辨匹配还是取消
				if req.isMatch {
					manager.firstMatch(&globalInfo, &req) // 首次匹配
				} else {
					manager.cancelMatch(&globalInfo, &req) // 取消匹配
				}
			}
		case <-mergeTimer.C: // 合并定时器
			{
				manager.mergeDesks(&globalInfo)
			}
		case <-deskTimer.C: // 桌子超时定时器
			{
				manager.checkDeskTimeout(&globalInfo)
			}
		case <-sucTimer.C: // 匹配成功超时定时器
			{
				manager.checkSucTimeout(&globalInfo)
			}
		}
	}

	logEntry.Debugln("离开函数")
	return
}

// 分发匹配请求
// playerID 	:	玩家ID
// gameID		：	请求匹配的游戏ID
// levelID		:   请求匹配的级别ID
// IP			: 	客户端IP地址
// 返回string 	 ：	 返回的错误描述，成功时返回空
func (manager *matchManager) dispatchMatchReq(playerID uint64, gameID uint32, levelID uint32, clientIP uint32) string {
	logEntry := logrus.WithFields(logrus.Fields{
		"playerID": playerID,
		"gameID":   gameID,
		"levelID":  levelID,
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
		logrus.WithError(err).Errorln("从gold服获取玩家金币失败")
		return fmt.Sprintf("从gold服获取玩家金币失败，游戏ID:%v，场次ID:%v，请求的玩家ID:%v", gameID, levelID, playerID)
	}

	// 金币范围检测
	if playerGold < levelConfig.minGold {
		logrus.Errorln("玩家金币数小于游戏场次金币要求最小值，最小值：%v，玩家金币:%v", levelConfig.minGold, playerGold)
		return fmt.Sprintf("玩家金币数小于游戏场次金币要求最小值，请求匹配的游戏ID:%v，场次ID:%v，玩家ID:%v", gameID, levelID, playerID)
	}

	if playerGold > levelConfig.maxGold {
		logrus.Errorf("玩家金币数大于游戏场次金币要求最大值，最大值：%v，玩家金币:%v", levelConfig.maxGold, playerGold)
		return fmt.Sprintf("玩家金币数大于游戏场次金币要求最大值，请求匹配的游戏ID:%v，场次ID:%v，玩家ID:%v", gameID, levelID, playerID)
	}

	// 获取该玩家的游戏信息
	playerGameInfo, err := hallclient.GetPlayerGameInfo(playerID, gameID)
	if err != nil || playerGameInfo == nil {
		logrus.WithError(err).Errorln("从hall服获取玩家游戏信息失败")
		return fmt.Sprintf("从hall服获取玩家游戏信息失败，游戏ID:%v，场次ID:%v，玩家ID:%v", gameID, levelID, playerID)
	}
	logEntry.Debugln("从hall服获取的玩家游戏信息：", playerGameInfo)

	// 计算胜率
	var playerWinRate int32 = 0

	// 场数不足时，采用默认胜率
	if playerGameInfo.GetTotalBurea() < web.GetMinGameTimes() {
		playerWinRate = web.GetDefaultWinRate()
	} else {
		playerWinRate = int32((float64(playerGameInfo.GetWinningBurea()) / float64(playerGameInfo.GetTotalBurea())) * 100)
	}

	// 全部检测通过

	// 获取该场次的申请通道
	reqMatchChan, exist := gameInfo.allLevelChan[levelID]
	if !exist {
		logrus.Errorln("内部错误，请求匹配的游戏存在，场次存在，但找不到该场次的匹配申请通道")
		return fmt.Sprintf("请求匹配的游戏存在，场次存在，但找不到该场次的匹配申请通道，请求匹配的游戏ID:%v，场次ID:%v，玩家ID:%v", gameID, levelID, playerID)
	}

	req := reqMatchPlayer{
		isMatch:  true,          // 申请匹配
		playerID: playerID,      // playerID
		winRate:  playerWinRate, // 胜率
		gold:     playerGold,    // 金币数
		IP:       clientIP,      // IP地址
	}

	// 压入通道
	reqMatchChan <- req

	logEntry.Debugln("已经把匹配请求压入对应的通道，压入的请求的:", req)

	logEntry.Debugln("离开函数")

	return ""
}

// 分发取消匹配请求
// playerID 	:	玩家ID
// gameID		：	请求取消匹配的游戏ID
// levelID		:   请求取消匹配的级别ID
// 返回string 	 ：	 返回的错误描述，成功时返回空
func (manager *matchManager) dispatchCancelMatchReq(playerID uint64, gameID uint32, levelID uint32) string {
	logEntry := logrus.WithFields(logrus.Fields{
		"playerID": playerID,
		"gameID":   gameID,
		"levelID":  levelID,
	})

	logEntry.Debugln("进入函数")

	// 得到该游戏的信息
	gameInfo, exist := manager.allGame[gameID]

	// 该游戏不存在
	if !exist {
		logrus.Errorln("内部错误，取消匹配时，发现正在匹配的游戏不存在")
		return fmt.Sprintf("取消匹配时，发现正在匹配的游戏ID:%v不存在，请求的玩家ID:%v", gameID, playerID)
	}

	// 获取该场次的申请通道
	reqMatchChan, exist := gameInfo.allLevelChan[levelID]
	if !exist {
		logrus.Errorln("内部错误，取消匹配时，发现游戏存在，场次存在，但找不到该场次的申请通道")
		return fmt.Sprintf("请求匹配的游戏存在，场次存在，但找不到该场次的匹配申请通道，请求匹配的游戏ID:%v，场次ID:%v，玩家ID:%v", gameID, levelID, playerID)
	}

	// 获取该玩家的游戏信息
	pPlayerGameInfo, err := hallclient.GetPlayerGameInfo(playerID, gameID)
	if err != nil || pPlayerGameInfo == nil {
		logrus.WithError(err).Errorln("从hall服获取玩家游戏信息失败")
		return fmt.Sprintf("从hall服获取玩家游戏信息失败，游戏ID:%v，场次ID:%v，玩家ID:%v", gameID, levelID, playerID)
	}

	// 计算胜率
	var playerWinRate int32 = 0

	// 场数不足时，采用默认胜率
	if pPlayerGameInfo.GetTotalBurea() < web.GetMinGameTimes() {
		playerWinRate = web.GetDefaultWinRate()
	} else {
		playerWinRate = int32((float64(pPlayerGameInfo.GetWinningBurea()) / float64(pPlayerGameInfo.GetTotalBurea())) * 100)
	}

	// 压入通道
	reqMatchChan <- reqMatchPlayer{
		isMatch:  false,         // 取消匹配
		playerID: playerID,      // playerID
		winRate:  playerWinRate, // 胜率
		gold:     0,             // 金币数，暂时不需要
		IP:       0,             // IP地址，暂时不需要
	}

	logEntry.Debugln("离开函数")

	return ""
}

// 将玩家添加到牌桌
// pMatchPlayer : 要加入的玩家
// pDesk 		: 要加入的桌子
// globalInfo	: 该场次的全局信息
// 返回 true表桌子已满且已发送给room，可以删除，false表桌子未满，继续匹配
func (manager *matchManager) addPlayerToDesk(pPlayer *matchPlayer, pDesk *matchDesk, globalInfo *levelGlobalInfo) bool {

	// 参数检测
	if pPlayer == nil || pDesk == nil {
		logrus.Error("严重错误，pPlayer == nil || pDesk == nil，返回")
		return false
	}

	logrus.WithFields(logrus.Fields{
		"player": pPlayer,
		"desk":   pDesk,
	})

	// 压入该玩家
	pDesk.players = append(pDesk.players, *pPlayer)

	// 总金币
	var allGold int64 = 0

	// 把不是指定玩家的其他玩家加进来
	for i := 0; i < len(pDesk.players); i++ {
		allGold += pDesk.players[i].gold
	}

	// 重新计算平均金币
	pDesk.aveGold = allGold / int64(len(pDesk.players))

	logrus.Debugf("桌子%v压入了玩家%v，当前人数:%v", pDesk, pPlayer, len(pDesk.players))

	// playerID与deskID的映射
	//manager.playerDesk[deskPlayer.playerID] = desk.deskID

	// 移除不在线的
	//manager.removeOfflines(desk)

	// 满员时的处理
	if uint8(len(pDesk.players)) >= pDesk.needPlayerCount {
		return manager.onDeskFull(pDesk, globalInfo)
	}

	return false
}

// onDeskFull 桌子满员时的处理
// 返回true : 表已经发送给room服创建桌子，桌子可以删除
// 返回false: 表某玩家离线或其他原因被移除，重新等待匹配
func (manager *matchManager) onDeskFull(pDesk *matchDesk, globalInfo *levelGlobalInfo) bool {

	// 参数检测
	if pDesk == nil {
		logrus.Errorln("内部错误，onDeskFull()，参数 pDesk == nil 为空，返回")
		return false
	}

	logEntry := logrus.WithFields(logrus.Fields{
		"pDesk": pDesk,
	})

	logEntry.Debugln("进入函数")

	tempPlayers := make([]matchPlayer, 0, pDesk.needPlayerCount)

	// 重新检测玩家状态，目的是移除不在线的
	for i := 0; i < len(pDesk.players); i++ {
		player := pDesk.players[i]

		// 机器人直接压入
		if player.robotLv != 0 {
			tempPlayers = append(tempPlayers, player)
			continue
		}

		// 获取玩家当前状态
		rsp, err := hallclient.GetPlayerState(player.playerID)
		if err != nil || rsp == nil {
			logEntry.Errorf("内部错误，从hall服获取玩家状态出错,玩家:%v", player)

			// 暂时不移除玩家
			tempPlayers = append(tempPlayers, player)
			continue
		}

		// 游戏状态不符合，则移除该玩家
		if rsp.GetState() != user.PlayerState_PS_MATCHING || rsp.GetGameId() != pDesk.gameID || rsp.GetLevelId() != pDesk.levelID {
			// 检测是否是离线状态，其他状态则报错
			if rsp.GetGateAddr() == "" {
				logEntry.Warningf("桌子满员时发现，玩家%v最新状态错误，已离线，最新state:%v， gameID:%v，levelID:%v，该玩家被移除出桌子", player, rsp.GetState(), rsp.GetGameId(), rsp.GetLevelId())
			} else {
				logEntry.Errorf("桌子满员时发现，玩家%v最新状态错误，最新state:%v， gameID:%v，levelID:%v，该玩家被移除出桌子", player, rsp.GetState(), rsp.GetGameId(), rsp.GetLevelId())
			}

			continue
		}

		tempPlayers = append(tempPlayers, player)
	}

	pDesk.players = tempPlayers

	// 人数不足了，则继续等待匹配
	if uint8(len(pDesk.players)) != pDesk.needPlayerCount {
		logEntry.Debugln("桌子满员后，由于有玩家被移出桌子，导致不满员，重新匹配")
		return false
	}

	// 通知room服创建桌子
	sendCreateDesk(*pDesk, globalInfo)

	logEntry.Debugln("离开函数")

	return true
}

// deleteDesk 删除指定的桌子
func (manager *matchManager) deleteDesk(pDesk *matchDesk) {

}

// mergeDesks 合并桌子
func (manager *matchManager) mergeDesks(globalInfo *levelGlobalInfo) {
	if globalInfo == nil {
		logrus.Errorln("mergeDesks()，globalInfo == nil，返回")
		return
	}

	logEntry := logrus.WithFields(logrus.Fields{
		"gameID":  globalInfo.gameID,
		"levelID": globalInfo.levelID,
	})

	//logEntry.Debugln("进入桌子合并函数")

	// 当前时间
	tNowTime := time.Now().Unix()

	// 所有的概率
	var index int32 = 0
	for ; index <= 100; index++ {

		var iterNext *list.Element = nil

		// 该概率下所有的桌子
		for iter := globalInfo.allRateDesks[index].Front(); iter != nil; iter = iterNext {

			iterNext = iter.Next()

			desk := *(iter.Value.(**matchDesk))

			// 距离桌子创建时间的间隔
			interval := tNowTime - desk.createTime

			// 不足1秒的，不检测，因为新建一个桌子时已检测过，不存在可合并的
			if interval < 1 {
				continue
			}

			// 前一秒的胜率浮动值
			lastRateValue := manager.getWinRateValue(interval - 1)

			// 这一秒的胜率浮动值
			nowRateValue := manager.getWinRateValue(interval)

			// 所有的需要检测的胜率值，也是allRateDesks的下标值
			checkIndexs := make([]int32, 0, 101)

			// 左段起始值(包含自身)
			leftStartIndex := index - nowRateValue
			if leftStartIndex < 0 {
				leftStartIndex = 0
			}
			// 左段结束值(不包含自身)
			leftEndIndex := index - lastRateValue
			if leftEndIndex < 0 {
				leftEndIndex = 0
			}
			// 从[leftStartIndex - leftEndIndex)
			for j := leftStartIndex; j < leftEndIndex; j++ {
				checkIndexs = append(checkIndexs, j)
			}

			// 右段起始值(不包含自身)
			rightStartIndex := index + lastRateValue
			if rightStartIndex > 100 {
				rightStartIndex = 100
			}

			// 右段结束值(包含自身)
			rightEndIndex := index + nowRateValue
			if rightEndIndex > 100 {
				rightEndIndex = 100
			}

			// 从[rightStartIndex - rightEndIndex)
			for j := rightStartIndex + 1; j <= rightEndIndex; j++ {
				checkIndexs = append(checkIndexs, j)
			}

			// 两段不应有重叠
			if (leftStartIndex > leftEndIndex) || (rightStartIndex > rightEndIndex) || (leftEndIndex > rightStartIndex) {
				logEntry.Errorf("左段或右段数据错误，跳过该桌子，左段起值：%v,左段终值：%v,右段起值：%v,右段终值：%v ", leftStartIndex, leftEndIndex, rightStartIndex, rightEndIndex)
				continue
			}

			logEntry.Debugf("左段起值：%v,左段终值：%v,右段起值：%v,右段终值：%v ", leftStartIndex, leftEndIndex, rightStartIndex, rightEndIndex)

			// 可合并桌子所在的信息
			var pList2 *list.List = nil
			var iter2 *list.Element = nil

			// 和这些桌子尝试组合
			for k := 0; k < len(checkIndexs); k++ {
				merIndex := checkIndexs[k]

				// 遍历该概率下的所有桌子
				for merIter := globalInfo.allRateDesks[merIndex].Front(); merIter != nil; merIter = merIter.Next() {

					merDesk := *(merIter.Value.(**matchDesk))

					// 检测金币范围
					minGold, maxGold := manager.getGoldRange(merDesk.aveGold, tNowTime-merDesk.createTime)
					if desk.aveGold < minGold || desk.aveGold > maxGold {
						continue
					}

					// IP是否存在相同的
					if manager.checkDeskSameIP(desk, merDesk) {
						continue
					}

					// 上局是否存在同桌的
					if manager.checkDeskLastSameDesk(desk, merDesk, globalInfo) {
						continue
					}

					// 找到可合并的即跳出
					pList2 = &globalInfo.allRateDesks[merIndex]
					iter2 = merIter
					break
				}

				// 找到可合并的即跳出
				if iter2 != nil {
					break
				}
			}

			// 有合并的桌子
			if iter2 != nil {
				iter1 := iter
				pDesk1 := iter1.Value.(*matchDesk)
				pList1 := &globalInfo.allRateDesks[index]

				pDesk2 := iter2.Value.(*matchDesk)

				logEntry.Debugln("找到了可以合并的，桌子1:%v，桌子2:%v", pDesk1, pDesk2)

				// desk1需作为时间最早的桌子，desk2需作为被拆的桌子
				if pDesk1.createTime > pDesk2.createTime {
					iter1, iter2 = iter2, iter1
					pDesk1, pDesk2 = pDesk2, pDesk1
					pList1, pList2 = pList2, pList1
				}

				// 把desk2的玩家移入到desk1

				// 临时玩家，和桌子2的玩家一一对应
				tempPlayers := make([]matchPlayer, 0, len(pDesk2.players))
				for i := 0; i < len(pDesk2.players); i++ {
					tempPlayers = append(tempPlayers, pDesk2.players[i])
				}

				for i := 0; i < len(tempPlayers); i++ {

					// 先从desk2中删除这个玩家
					manager.removePlayerFromDesk(&pDesk2.players[i], pDesk2)

					// desk1桌子满，则从列表中删除desk1桌子，跳出
					if manager.addPlayerToDesk(&tempPlayers[i], pDesk1, globalInfo) {
						logEntry.Debugln("由于拆入的玩家已满桌，删除桌子1")
						pList1.Remove(iter1)
						break
					}
				}

				// 是否存在真实玩家
				bExistTruePlayer := false
				for i := 0; i < len(pDesk2.players); i++ {
					if pDesk2.players[i].robotLv == 0 {
						bExistTruePlayer = true
					}
				}

				// 桌子2没真实玩家了，就删除桌子
				if !bExistTruePlayer {
					pList2.Remove(iter2)
				}
			}
		}
	}

	//logEntry.Debugln("离开桌子合并函数")
}

// 从桌子中删除玩家
func (manager *matchManager) removePlayerFromDesk(pPlayer *matchPlayer, pDesk *matchDesk) {
	// 参数检测
	if pPlayer == nil || pDesk == nil {
		logrus.Error("严重错误，removePlayerFromDesk(), pPlayer == nil || pDesk == nil，返回")
		return
	}

	logrus.WithFields(logrus.Fields{
		"player": pPlayer,
		"desk":   pDesk,
	})

	// 总金币
	var allGold int64 = 0

	// 新玩家
	tempPlayers := make([]matchPlayer, 0, pDesk.needPlayerCount)

	// 把不是指定玩家的其他玩家加进来
	for i := 0; i < len(pDesk.players); i++ {
		if pPlayer.playerID != pDesk.players[i].playerID {
			tempPlayers = append(tempPlayers, pDesk.players[i])
			allGold += pDesk.players[i].gold
		}
	}

	// 新玩家
	pDesk.players = tempPlayers

	// 重新计算平均金币
	pDesk.aveGold = allGold / int64(len(pDesk.players))
}

// checkDeskTimeout 检测桌子是否超时
func (manager *matchManager) checkDeskTimeout(globalInfo *levelGlobalInfo) {
	if globalInfo == nil {
		logrus.Errorln("checkTimeout()，参数错误，globalInfo == nil，返回")
		return
	}

	logEntry := logrus.WithFields(logrus.Fields{
		"gameID":  globalInfo.gameID,
		"levelID": globalInfo.levelID,
	})

	// 当前时间
	tNowTime := time.Now().Unix()

	// 机器人加入时间
	joinTime := int64(web.GetRobotJoinTime().Seconds())

	//logEntry.Debugf("进入桌子超时检测函数，当前时间：%v", tNowTime)

	// 所有的概率
	var index int32 = 0
	for ; index <= 100; index++ {

		var next *list.Element
		// 该概率下所有的桌子进入桌子超时检测函数
		for iter := globalInfo.allRateDesks[index].Front(); iter != nil; iter = next {

			// 提前保存下一个
			next = iter.Next()

			desk := *(iter.Value.(**matchDesk))

			// 间隔秒数
			interval := tNowTime - desk.createTime

			// logEntry.Debugf("开始检测桌子:%v是否超时，桌子已创建时间:%v秒", desk, interval)

			// 超过时间，则开始加入机器人
			if interval >= joinTime {

				// 胜率范围
				beginRate := index - manager.getWinRateValue(interval)
				if beginRate < 0 {
					beginRate = 0
				}

				endRate := int32(index) + int32(manager.getWinRateValue(interval))
				if endRate > 100 {
					endRate = 100
				}

				// logEntry.Debugf("桌子的平均金币:%v", desk.aveGold)

				// 金币范围
				minGold, maxGold := manager.getGoldRange(desk.aveGold, interval)

				reqRobot := robotclient.LeisureRobotReqInfo{
					CoinHigh:    maxGold,
					CoinLow:     minGold,
					WinRateHigh: endRate,
					WinRateLow:  beginRate,
					GameID:      desk.gameID,
					LevelID:     desk.levelID,
				}

				// logEntry.Debugf("请求的机器人参数:%v", reqRobot)

				// 从hall服获取一个空闲的机器人
				robotPlayerID, robotGold, robotRate, err := robotclient.GetLeisureRobotInfoByInfo(reqRobot)
				if err != nil {
					logEntry.WithError(err).Error("从hall服获取机器人失败,继续下一个桌子")
					continue
				}

				// 新建一个匹配玩家(机器人)
				newMatchPlayer := matchPlayer{
					playerID: robotPlayerID,
					robotLv:  1, // todo
					seat:     -1,
					IP:       0,
					gold:     robotGold,
				}

				// 更新机器人状态
				// 设置为匹配状态，后面匹配过程中出错删除时再标记为空闲状态，匹配成功时不需处理(room服会标记为游戏状态)
				bSuc, err := hallclient.UpdatePlayerState(robotPlayerID, user.PlayerState_PS_IDIE, user.PlayerState_PS_MATCHING, globalInfo.gameID, globalInfo.levelID)
				if err != nil || !bSuc {
					logEntry.WithError(err).Errorf("内部错误，通知hall服设置机器人状态为匹配状态时失败，游戏ID:%v，场次ID:%v，机器人玩家ID:%v", globalInfo.gameID, globalInfo.levelID, robotPlayerID)
				}

				// 更新机器人所在的服务器类型和地址
				bSuc, err = hallclient.UpdatePlayeServerAddr(robotPlayerID, user.ServerType_ST_MATCH, GetServerAddr())
				if err != nil || !bSuc {
					logEntry.WithError(err).Errorf("内部错误，通知hall服设置机器人的服务器类型及地址时失败，游戏ID:%v，场次ID:%v，机器人玩家ID:%v", globalInfo.gameID, globalInfo.levelID, robotPlayerID)
				}

				logEntry.Debugf("从hall服获取机器人成功，playerID:%v，金币数:%v，胜率:%v ", robotPlayerID, robotGold, robotRate)

				// 把机器人加入桌子
				if manager.addPlayerToDesk(&newMatchPlayer, desk, globalInfo) {

					logEntry.Debugf("请求到的机器人%v加入了桌子%v，桌子满员，桌子已删除", newMatchPlayer, desk)

					// 桌子已满，则删除
					globalInfo.allRateDesks[index].Remove(iter)
				} else {
					logEntry.Debugf("请求到的机器人%v加入了桌子%v，桌子未满员，继续匹配", newMatchPlayer, desk)
				}
			}
		}
	}

	//logEntry.Debugln("离开桌子超时检测函数")
}

// checkSucTimeout 检测之前匹配成功的是否超时
func (manager *matchManager) checkSucTimeout(globalInfo *levelGlobalInfo) {

	// 参数检测
	if globalInfo == nil {
		logrus.Errorln("checkSucTimeout()，参数错误，globalInfo == nil，返回")
		return
	}

	// logEntry := logrus.WithFields(logrus.Fields{
	// 	"gameID":  globalInfo.gameID,
	// 	"levelID": globalInfo.levelID,
	// })

	// 当前时间
	tNowTime := time.Now().Unix()

	// logEntry.Debugf("进入匹配成功超时检测函数，当前时间：%v", tNowTime)

	// 新建，然后再替换
	newSucDesks := map[uint64]*sucDesk{}
	newSucPlayers := map[uint64]uint64{}

	// 先遍历桌子，只记录未超时的
	for key, desk := range globalInfo.sucDesks {
		if tNowTime-desk.sucTime < web.GetSameDeskLimitTime() {
			newSucDesks[key] = desk
		}
	}

	// 替换桌子
	globalInfo.sucDesks = newSucDesks

	// 再遍历玩家，只记录桌子存在的
	for playerID, deskID := range globalInfo.sucPlayers {
		_, exist := globalInfo.sucDesks[playerID]
		if exist {
			newSucPlayers[playerID] = deskID
		}
	}

	// 替换玩家
	globalInfo.sucPlayers = newSucPlayers

	// logEntry.Debugln("离开匹配成功超时检测函数")
}
