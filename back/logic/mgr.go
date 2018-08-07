package logic

import (
	"steve/entity/gamelog"
)

// SaveDetailInfo 向db更新detail的信息
func SaveDetailInfo(detailInfo gamelog.TGameDetail) {
	insertDetailInfo(detailInfo)
}

// SaveSummaryInfo 向db储存Summary的信息
func SaveSummaryInfo(summaryInfo gamelog.TGameSummary) {
	insertSummaryInfo(summaryInfo)
}

// UpdatePlayerGameInfo 向db更新player的信息
func UpdatePlayerGameInfo(detailInfo gamelog.TGameDetail) {
	updatePlayerInfo(detailInfo)
}
