package states

import (
	"fmt"
	"steve/majong/utils"
	majongpb "steve/server_pb/majong"

	"github.com/Sirupsen/logrus"
)

//CheckFlowerPigSettle 查花猪结算 playerAll-所有玩家，flowerPigPlayers-花猪玩家，noPigPlayers-不是花猪的未听玩家
func CheckFlowerPigSettle(playerAll []*majongpb.Player) ([]*majongpb.SettleInfo, error) {
	// 所有花猪玩家
	flowerPigPlayers := utils.GetFlowerPigPlayers(playerAll)
	flowerPigSum := len(flowerPigPlayers)
	// 查花猪信息
	settleInfos := make([]*majongpb.SettleInfo, 0)
	if flowerPigSum > 0 {
		//di-麻将底分，total-赢玩家结算的总分，
		var di, total int64
		di = 1
		total = 16 * di
		// 胡过玩家
		huEdPlayers := utils.GetHuEdPlayers(playerAll)
		//未听玩家,不包括花猪，因为花猪玩家之间会直接抵消
		noPigPlayers, err := utils.GetBustedHandPlayers(playerAll, false)
		if err != nil {
			return []*majongpb.SettleInfo{}, fmt.Errorf("查花猪结算-未听玩家：%v", err)
		}
		// 听玩家，听牌玩家ID及最大倍数
		tingPlayersMap, err := utils.GetTingPlayerIDAndMultiple(playerAll)
		if err != nil {
			return []*majongpb.SettleInfo{}, fmt.Errorf("查花猪结算-听玩家：%v", err)
		}
		// 对每个花猪玩家，进行查花猪结算，包括了查大叫
		for i := 0; i < flowerPigSum; i++ {
			// loseTotal-当前输家总共输掉的分数
			var loseTotal int64
			// 花猪玩家
			flowerPig := flowerPigPlayers[i]
			// 玩家ID对应的输赢分，除了当前花猪玩家要叠加扣的分数，其他都直接覆盖，因为玩家不能能有1个以上的，状态，如不可能拥有胡，和未听，2种
			settleInfoMap := make(map[uint64]int64)
			//胡过玩家结算处理
			for j := 0; j < len(huEdPlayers); j++ {
				settleInfoMap[huEdPlayers[j].PalyerId] = total
				loseTotal = loseTotal - total
			}
			// 不是花猪的未听玩家结算处理
			for n := 0; n < len(noPigPlayers); n++ {
				settleInfoMap[noPigPlayers[n].PalyerId] = total
				loseTotal = loseTotal - total
			}
			// 听玩家结算处理
			for playerID, multiple := range tingPlayersMap {
				tingPlayer := utils.GetPlayerByID(playerAll, playerID)
				if tingPlayer == nil {
					return []*majongpb.SettleInfo{}, fmt.Errorf("查花猪结算-听牌玩家ID不存在： %v ", playerID)
				}
				// 16*di + multiple*di = (16+multiple)*di
				total2 := total + (multiple * di)
				settleInfoMap[tingPlayer.PalyerId] = total2
				loseTotal = loseTotal - total2
			}
			// 查花猪玩家结算信息，例子：id:2 为花猪玩家-[id:2 scores:<key:1 value:2 > scores:<key:2 value:-6 > scores:<key:3 value:2 > scores:<key:4 value:2 > ]
			if len(settleInfoMap) > 0 {
				settleInfoMap[flowerPig.PalyerId] = loseTotal
				flowerSettleInfo := &majongpb.SettleInfo{
					Id:     flowerPig.PalyerId,
					Scores: settleInfoMap,
				}
				settleInfos = append(settleInfos, flowerSettleInfo)
			}
		}
		logrus.WithFields(logrus.Fields{
			"flowerPigSum": flowerPigSum,
			"settleInfos":  settleInfos,
		}).Info("-------查花猪结算")
	}
	return settleInfos, nil
}

//CheckYellSettle  查大叫结算 playerAll-所有玩家，noPigPlayers-不是花猪的未听玩家
func CheckYellSettle(playerAll []*majongpb.Player) ([]*majongpb.SettleInfo, error) {
	// 所有未听玩家，不包括花猪，因为查花猪包括了查大叫，所以未听玩家，中是花猪的，都不用再进行查大叫
	noPigPlayers, err := utils.GetBustedHandPlayers(playerAll, false)
	if err != nil {
		return []*majongpb.SettleInfo{}, fmt.Errorf("查大叫结算-未听玩家：%v", err)
	}
	noPigSum := len(noPigPlayers)
	// 查大叫结算信息
	settleInfos := make([]*majongpb.SettleInfo, 0)
	if noPigSum > 0 {
		//di-麻将底分
		var di int64
		di = 1
		// 听玩家,听牌玩家ID及最大倍数
		tingPlayersMap, err := utils.GetTingPlayerIDAndMultiple(playerAll)
		if err != nil {
			return []*majongpb.SettleInfo{}, err
		}
		// 对每个未听玩家，进行查大叫结算
		for i := 0; i < noPigSum; i++ {
			//total-赢玩家结算的总分，loseTotal-当前输家总共输掉的分数
			var total, loseTotal int64
			// 未听玩家
			noPigBustedHandPlayer := noPigPlayers[i]
			// 玩家ID对应的输赢分，除了当前未听玩家要叠加扣的分数，其他都直接覆盖，
			settleInfoMap := make(map[uint64]int64)
			// 听玩家结算处理
			for playerID, multiple := range tingPlayersMap {
				tingPlayer := utils.GetPlayerByID(playerAll, playerID)
				if tingPlayer == nil {
					return []*majongpb.SettleInfo{}, fmt.Errorf("查大叫结算-听牌玩家ID不存在： %v ", playerID)
				}
				total = multiple * di
				loseTotal = loseTotal - total
				settleInfoMap[tingPlayer.PalyerId] = total
			}
			// 结算信息记录
			if len(settleInfoMap) > 0 {
				settleInfoMap[noPigBustedHandPlayer.PalyerId] = loseTotal
				yellSettleInfo := &majongpb.SettleInfo{
					Id:     noPigBustedHandPlayer.PalyerId,
					Scores: settleInfoMap,
				}
				settleInfos = append(settleInfos, yellSettleInfo)
			}
		}
		logrus.WithFields(logrus.Fields{
			"noPigSum":    noPigSum,
			"settleInfos": settleInfos,
		}).Info("-------查大叫结算")
	}
	return settleInfos, nil
}

//CallDivertSettle 呼叫转移 huType 胡类型，palyerAll所有玩家,winPalyers点炮玩家，loserPlayer被点炮玩家,当前的杠钱转移给其他人杠后炮一响，谁胡，杠钱就转移给谁
func CallDivertSettle(huType majongpb.HuType, palyerAll, winPalyers []*majongpb.Player, loserPlayer *majongpb.Player) ([]*majongpb.SettleInfo, error) {
	log := logrus.WithFields(logrus.Fields{
		"huType": huType,
	})
	// 呼叫转移结算信息
	settleInfos := make([]*majongpb.SettleInfo, 0)
	// 判断是否是杠后炮
	if huType == majongpb.HuType_hu_ganghoupao {
		// 玩家ID对应的输赢分
		settleInfoMap := make(map[uint64]int64)
		// di底分
		var di int64 = 1

		// 判断输家是否有杠牌，防止数组越界
		gangSum := len(loserPlayer.GangCards)
		if gangSum == 0 {
			return []*majongpb.SettleInfo{}, fmt.Errorf("呼叫转移事件错误-被点炮玩家没有杠")
		}
		// 获取输家最后一个杠，即被转移的杠
		gang := loserPlayer.GangCards[gangSum-1]

		// 暗杠-每个人给2倍底，补杠-每个人给1倍底，明杠-点杠人给2倍底
		gangScoreMap := map[majongpb.GangType]int64{
			majongpb.GangType_gang_angang:   2,
			majongpb.GangType_gang_bugang:   1,
			majongpb.GangType_gang_minggang: 2,
		}
		// 获取杠类型分书数
		score, isExist := gangScoreMap[gang.Type]
		if !isExist {
			return []*majongpb.SettleInfo{}, fmt.Errorf("呼叫转移事件错误-杠类型不存在：%v", gang.Type)
		}

		// 赢家人数
		winSum := len(winPalyers)
		//输家所得的总杠分
		total := di * score
		// 不是明杠类型，要乘上所有玩家数量
		if gang.Type != majongpb.GangType_gang_minggang {
			total = total * int64(len(palyerAll))
		}

		// 日志
		log.WithFields(logrus.Fields{"winSum": winSum, "gang": gang, "gangtotal": total})

		// 只有一个赢家，谁胡，杠钱，不管什么杠，都给胡的人
		if winSum == 1 {
			win := winPalyers[0]
			settleInfoMap[win.PalyerId] = total
		} else if winSum > 1 { // 多个赢家的情况
			isDivideEqually := true // 是需要平分杠钱
			switch gang.Type {
			//杠类型为明杠
			case majongpb.GangType_gang_minggang:
				// 赢家数组中包含有明杠玩家的情况，杠钱都给点杠人
				winGangPlayer := utils.GetPlayerByID(winPalyers, gang.SrcPlayer)
				if winGangPlayer != nil {
					settleInfoMap[winGangPlayer.PalyerId] = settleInfoMap[winGangPlayer.PalyerId] + total
					// 钱都转了，不用平分
					isDivideEqually = false

					// 日志
					log.WithFields(logrus.Fields{"winGangPlayerID": winGangPlayer.PalyerId})
				}
			//杠类型为补杠
			case majongpb.GangType_gang_bugang:
				// 在赢家人数为2的情况下，补杠分数，无法平分，多余的分数要给第一个胡家
				if winSum == 2 {
					// 剩余分数
					surplusTotal := total % int64(winSum)
					// 获取第一个胡的赢家
					firstHuPlayer := utils.GetFirstHuPlayerByID(palyerAll, winPalyers, loserPlayer.PalyerId)
					if firstHuPlayer == nil {
						return []*majongpb.SettleInfo{}, fmt.Errorf("呼叫转移事件错误-获取第一胡家失败：%v ", winPalyers)
					}
					settleInfoMap[firstHuPlayer.PalyerId] = settleInfoMap[firstHuPlayer.PalyerId] + surplusTotal

					// 日志
					log.WithFields(logrus.Fields{"firstHuPlayerID": firstHuPlayer.PalyerId, "surplusTotal": surplusTotal})
				}
			}
			if isDivideEqually {
				// 平分
				equallyTotal := total / int64(winSum)
				for i := 0; i < winSum; i++ {
					win := winPalyers[i]
					settleInfoMap[win.PalyerId] = settleInfoMap[win.PalyerId] + equallyTotal
				}

				// 日志
				log.WithFields(logrus.Fields{"equallyTotal": equallyTotal})
			}
		} else {
			return []*majongpb.SettleInfo{}, fmt.Errorf("呼叫转移事件错误-赢家人数为0")
		}

		// 输家总共扣的分
		loserTotal := -total
		settleInfoMap[loserPlayer.PalyerId] = loserTotal
		// 呼叫转移结算信息
		if len(settleInfoMap) > 0 {
			yellSettleInfo := &majongpb.SettleInfo{
				Id:     loserPlayer.PalyerId,
				Scores: settleInfoMap,
			}
			settleInfos = append(settleInfos, yellSettleInfo)
		}

		// 日志
		log.WithFields(logrus.Fields{"settleInfos": settleInfos})
	}
	log.Info("----------呼叫转移结算")
	return settleInfos, nil
}

//TaxRebateSettle 退税结算 TODO 呼叫转移的杠，不用退税
func TaxRebateSettle(playerAll []*majongpb.Player) ([]*majongpb.SettleInfo, error) {
	// 所有未听玩家，包括花猪
	bustedHandPlayers, err := utils.GetBustedHandPlayers(playerAll, true)
	if err != nil {
		return []*majongpb.SettleInfo{}, fmt.Errorf("退税结算-未听玩家：%v", err)
	}
	bustedHandPlayerSum := len(bustedHandPlayers)
	// 退税结算信息
	settleInfos := make([]*majongpb.SettleInfo, 0)
	if bustedHandPlayerSum > 0 {
		//di-麻将底分，total-与玩家结算的总分，loseTotal-输家总共输掉的分数
		var di int64
		di = 1
		for i := 0; i < bustedHandPlayerSum; i++ {
			//total-与玩家结算的总分，loseTotal-输家总共输掉的分数
			var total, loseTotal int64
			// 未听玩家
			bHandPlayer := bustedHandPlayers[i]
			// 玩家ID对应的输赢分
			settleInfoMap := make(map[uint64]int64)
			gangCards := bHandPlayer.GetGangCards()
			for j := 0; j < len(gangCards); j++ {
				gang := gangCards[j]
				// 赢钱玩家
				winPlayers := make([]*majongpb.Player, 0)
				switch gang.Type {
				case majongpb.GangType_gang_angang: //暗杠类型
					total = di * 2
					winPlayers = playerAll
				case majongpb.GangType_gang_bugang: //补杠类型
					total = di * 1
					winPlayers = playerAll
				case majongpb.GangType_gang_minggang: //明杠类型
					total = di * 2
					// 明杠玩家
					mingGangPlayer := utils.GetPlayerByID(playerAll, gang.GetSrcPlayer())
					if mingGangPlayer == nil {
						return []*majongpb.SettleInfo{}, fmt.Errorf("退税结算-该玩家ID不存在：%v", gang.SrcPlayer)
					}
					winPlayers = append(winPlayers, mingGangPlayer)
				default:
					return []*majongpb.SettleInfo{}, fmt.Errorf("退税结算-玩家ID（%v）-该杠类型不存在：%v", bHandPlayer.PalyerId, gang.Type)
				}
				//还杠钱
				for n := 0; n < len(winPlayers); n++ {
					win := winPlayers[n]
					if win.GetPalyerId() != bHandPlayer.GetPalyerId() {
						// 当玩家退税俩个暗杠，退的都是一样的人，不能覆盖，要叠加
						settleInfoMap[win.PalyerId] = settleInfoMap[win.PalyerId] + total
						loseTotal = loseTotal - total
					}
				}
			}
			if len(settleInfoMap) > 0 {
				settleInfoMap[bHandPlayer.PalyerId] = loseTotal
				yellSettleInfo := &majongpb.SettleInfo{
					Id:     bHandPlayer.PalyerId,
					Scores: settleInfoMap,
				}
				settleInfos = append(settleInfos, yellSettleInfo)
			}
		}
		logrus.WithFields(logrus.Fields{
			"bustedHandPlayerSum": bustedHandPlayerSum,
			"settleInfos":         settleInfos,
		}).Info("-------退税结算")
	}
	return settleInfos, nil
}
