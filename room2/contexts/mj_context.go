package contexts

import (
	"steve/entity/majong"
	"time"
)

// deskContext 牌桌现场
type MjContext struct {
	MjContext   majong.MajongContext // 牌局现场
	StateNumber int                  // 状态序号
	StateTime   time.Time            // 状态时间
}

func (c *MjContext) SetStateNumber(state int){
	c.StateNumber = state
	println("更新statenumber-------->",c.StateNumber)
}
