package prop

const (
	invalidPropType = int32(0)
	gold            = int32(1)
	prop            = int32(2)
)

const (
	invalidProp = int32(0)
	rose        = int32(1)
	beer        = int32(2)
	bomb        = int32(3)
	grabChicken = int32(4)
	eggGun      = int32(5)
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
