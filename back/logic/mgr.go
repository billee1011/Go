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

// SavePlayerGameInfo 向db更新player的信息
func SavePlayerGameInfo(summaryInfo gamelog.TGameSummary) {

}

// SetWinningPercentage 计算玩家的胜率
func SetWinningPercentage() {

}
