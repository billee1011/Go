package prop

const (
	// InvalidPropType 非法道具值
	InvalidPropType = int32(0)
	// Gold 货币
	Gold = int32(1)
	// Props 道具
	Props = int32(2)
)

const (
	// InvalidProp 非法值
	InvalidProp = int32(0)
	// Rose 玫瑰花
	Rose = int32(1)
	// Beer 啤酒
	Beer = int32(2)
	// Bomb 炸弹
	Bomb = int32(3)
	// GrabChicken 抓鸡
	GrabChicken = int32(4)
	// EggGun 鸡蛋机枪
	EggGun = int32(5)
)

// Prop 道具
type Prop struct {
	PropID int32
	Count  int64
}

// PlayerProps 玩家拥有的道具
type PlayerProps struct {
	PlayerID uint64
	Props    []Prop
}
