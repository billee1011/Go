package constant

// GoldFuncType 金币增减调用函数类型
type GoldFuncType int32

const (
	// GFGAMESETTLE 游戏结算
	GFGAMESETTLE GoldFuncType = 0
	// GFGAMEPEIPAI 游戏配牌
	GFGAMEPEIPAI GoldFuncType = 1
)

// PropsFuncType 道具增减调用函数类型
type PropsFuncType int32

const (
	// PFGAMEUSE 游戏中使用
	PFGAMEUSE PropsFuncType = 0
)
