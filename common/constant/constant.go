package constant

// GoldFuncType 金币修改调用函数类型
type GoldFuncType int32

const (
	// GFGAMESETTLE 游戏结算
	GFGAMESETTLE GoldFuncType = 0
	// GFGAMEPEIPAI 游戏配牌
	GFGAMEPEIPAI GoldFuncType = 1
)
