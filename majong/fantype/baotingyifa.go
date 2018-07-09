package fantype

// checkBaoTingYiFa 检测报听一发 报听后紧接着就胡牌（包括正常的点炮或自摸，以及补花后的点炮 或自摸） TODO
func checkBaoTingYiFa(tc *typeCalculator) bool {
	return tc.getPlayer().GetTingStateInfo().GetIsBaotingyifa()
}
