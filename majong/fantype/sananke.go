package fantype

// checkSanAnKe 三暗刻:胡牌时,含有 3 个暗刻
func checkSanAnKe(tc *typeCalculator) bool {
	return getPlayerMaxAnKeNum(tc.combines, 3)
}

// // checkSanAnKe 三暗刻:胡牌时,含有 3 个暗刻
// func checkSanAnKe(tc *typeCalculator) bool {
// 	keCombines := make([]combine, 0)
// 	for _, combine := range tc.combines {
// 		if len(combine.kes) >= 3 {
// 			keCombines = append(keCombines, combine)
// 		}
// 	}
// 	for _, keCombine := range keCombines {
// 		anKeCount := 0
// 		for _, ke := range keCombine.kes {
// 			if contains(tc.handCards, ke) {
// 				anKeCount++
// 			}
// 			if anKeCount >= 3 {
// 				return true
// 			}
// 		}

// 	}
// 	return false
// }

// func contains(cards []*majongpb.Card, check *majongpb.Card) bool {
// 	for _, card := range cards {
// 		if utils.ServerCard2Number(check) == utils.ServerCard2Number(card) {
// 			return true
// 		}
// 	}
// 	return false
// }
