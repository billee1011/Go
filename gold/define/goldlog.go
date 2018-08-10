package define



type GoldLog struct {
	TradeID string			// 交易ID
	PlayerID uint64			// 玩家ID
	Channel int64			// 渠道ID
	CurrencyType int16		// 货币类型：1.金币，2.元宝，3，房卡
	Amount int64			// 变化类型
	BeforeBalance int64		// 变化前余额
	AfterBalance  int64		// 变化后余额
	TradeTime  string		// 交易时间
	Status  int32			// 状态: 1.成功，2失败
	GameId int32			// 游戏ID
	Level int32				// 场次ID
	FuncId int32			// 行为ID或功能ID

}