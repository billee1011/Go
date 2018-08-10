package fixed

type LogType int32

const (
	LOG_TYPE_GAME_PERSON_NUM LogType = 1 //定时上报每种游戏玩家数量 (5分钟一次)

	//实时上包类型
	LOG_TYPE_REG       LogType = 4  //新增注册 每次发生注册时上报
	LOG_TYPE_ACT       LogType = 5  //活跃人数 玩家每次登录时上报
	LOG_TYPE_GAM       LogType = 6  //游戏人数 玩家完成牌局时上报
	LOG_TYPE_GOLD_ADD  LogType = 7  //金豆总产出 玩家获得金豆时上报
	LOG_TYPE_GOLD_REMV LogType = 8  //金豆总消耗 玩家消耗金豆时上报
	LOG_TYPE_YB_ADD    LogType = 9  //元宝总产出 玩家获得元宝时上报
	LOG_TYPE_YB_REMV   LogType = 10 //元宝总消耗 玩家消耗元宝时上报
	LOG_TYPE_CARD_ADD  LogType = 11 //房卡总产出 玩家获得房卡时上报
	LOG_TYPE_CARD_REMV LogType = 12 //房卡总消耗 玩家消耗房卡时上报
)
