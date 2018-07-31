package fantype

// checkTianTing 检测天听
// 庄家打出第一张牌前报听称为天听，发完牌后闲家摸牌便报听也称为天听；
// 若发完牌后有玩家补花，补花之后报听也算天听；
// 若庄家在发完牌后有暗杠，则庄家不算天听，但算报听；
func checkTianTing(tc *typeCalculator) bool {
	return tc.getPlayer().GetTingStateInfo().GetIsTianting()
}
