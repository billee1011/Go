package fantype

// checkBaoTing 检测报听 主动选择听牌，自动摸牌打牌后胡牌 TODO
func checkBaoTing(tc *typeCalculator) bool {
	return tc.getPlayer().GetTingStateInfo().GetIsTing()
}
