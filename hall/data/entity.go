package data

// gameDetail 游戏详情
type gameDetail struct {
	gameID      uint64 `xorm:"gameID"`
	gname       string `xorm:"name"`
	gtype       uint64 `xorm:"type"`
	glevelID    uint64 `xorm:"levelID"`
	gbaseScores uint64 `xorm:"baseScores"`
	glowScores  uint64 `xorm:"lowScores"`
	ghighScores uint64 `xorm:"highScores"`
	gstatus     uint64 `xorm:"status"`
}
