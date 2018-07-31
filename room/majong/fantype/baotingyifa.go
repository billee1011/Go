package fantype

// checkBaoTingYiFa 检测报听一发 报听后紧接着就胡牌（包括正常的点炮或自摸，以及补花后的点炮 或自摸） TODO
func checkBaoTingYiFa(tc *typeCalculator) bool {
	if tc.callCheckFunc(baotingFuncID) {
		player := tc.getPlayer()
		chupaiCount := player.GetChupaiCount()
		actionCount := len(player.GetGangCards())
		baoTingCount := player.GetTingStateInfo().GetBaotingyifaCount()
		return baoTingCount+1 == chupaiCount+int32(actionCount)
	}
	return false
}
