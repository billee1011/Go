package interfaces

import (
	"steve/server_pb/majong"
	"time"
	"steve/server_pb/ddz"
)

// AutoEventGenerateParams 生成自动事件的参数
type AutoEventGenerateParams struct {
	Desk          Desk
	MajongContext *majong.MajongContext
	DDZContext    *ddz.DDZContext
	PlayerIds     []uint64
	StartTime     time.Time
	Duration      uint32
	RobotLv       map[uint64]int
}

// AutoEventGenerateResult 自动事件生成结果
type AutoEventGenerateResult struct {
	Events []Event
}

// DeskAutoEventGenerator 牌桌自动事件产生器
type DeskAutoEventGenerator interface {
	GenerateV2(params *AutoEventGenerateParams) AutoEventGenerateResult
	RegisterAI(gameID int, stateID int32, AI CommonAI) // 注册 AI
}
