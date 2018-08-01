// purpose	：给定拥有的牌，获取指定牌，指定牌型的压制牌
//			: 给定牌（6665443），指定对子（33），返回找到的压制牌，成功即true,44
// 			：给定牌（6665443），指定对子（77），返回找到的压制牌，失败即false,nil
// author	: yzf
// date		：2018-7-18

package states

import (
	"sort"

	"github.com/Sirupsen/logrus"
)

// GetPokeCount 统计各个牌的个数
// @inparam 	pokes ： 需统计的牌
// @outparam	map[uint32]uint32 :	key:牌的无花色权重，value:牌的个数
// @outparam	[]int :	所有牌的无花色权重，已去重，已排序
func GetPokeCount(pokes []Poker) (map[uint32]uint32, []int) {

	counts := make(map[uint32]uint32)
	keys := []int{}

	for _, poke := range pokes {
		pointWeight := poke.PointWeight

		count, exists := counts[pointWeight]

		if !exists {
			counts[pointWeight] = 1
			keys = append(keys, int(pointWeight))
		} else {
			counts[pointWeight] = count + 1
		}
	}

	// 排序
	sort.Ints(keys)

	return counts, keys
}

// GetBoom 若有炸弹，返回炸弹;没有则返回false
func GetBoom(allPokes []Poker) (bool, []Poker) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "play.go:GetBoom",
		"allPokes":  allPokes,
	})

	// 参数检测
	if len(allPokes) == 0 {
		logEntry.Errorln("参数错误")
		return false, []Poker{}
	}

	// 先排序，从小到大
	DDZPokerSort(allPokes)

	// 统计各个牌的个数
	counts, keys := GetPokeCount(allPokes)

	// 炸弹的无花色权重
	var boomPointWeight uint32 = 0
	for i := 0; i < len(keys); i++ {
		pointWeight := (uint32)(keys[i])
		count := counts[pointWeight]
		if count == 4 {
			boomPointWeight = pointWeight
			break
		}
	}

	// 炸弹的poke
	boomPokes := []Poker{}

	// 有则找到炸弹的这些牌
	if boomPointWeight != 0 {
		for i := 0; i < len(allPokes); i++ {
			if allPokes[i].PointWeight == boomPointWeight {
				boomPokes = append(boomPokes, allPokes[i])
			}
		}
	}

	if len(boomPokes) == 0 {
		return false, nil
	}

	return true, boomPokes
}

// GetKingBoom 若有炸弹，返回炸弹;没有则返回false
func GetKingBoom(allPokes []Poker) (bool, []Poker) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "play.go:GetKingBoom",
		"allPokes":  allPokes,
	})

	// 参数检测
	if len(allPokes) == 0 {
		logEntry.Errorln("参数错误")
		return false, []Poker{}
	}

	// 火箭的pokes
	boomPokes := []Poker{}

	for i := 0; i < len(allPokes); i++ {
		if allPokes[i].Suit == uint32(0x00) && (allPokes[i].Point == uint32(0x0E) || allPokes[i].Point == uint32(0x0F)) {
			boomPokes = append(boomPokes, allPokes[i])
		}
	}

	// 不是火箭则清空
	if len(boomPokes) != 2 {
		boomPokes = []Poker{}
	}

	if len(boomPokes) == 0 {
		return false, nil
	}

	return true, boomPokes
}

// GetMinBiggerSingle 从allPokes中获取比指定单牌speciPoke大的最小的单牌
// 例如：87777655544 中找到比3大的牌，应该返回6;
// 例如：77776665544 中找到比3大的牌，应该返回空;
func GetMinBiggerSingle(allPokes []Poker, speciPoke []Poker) (bool, []Poker) {

	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "play.go:GetMinBiggerSingle",
		"allPokes":  allPokes,
		"speciPoke": speciPoke,
	})

	// 参数检测
	if len(allPokes) == 0 || len(speciPoke) != 1 {
		logEntry.Errorln("参数错误")
		return false, []Poker{}
	}

	// 先排序，从小到大
	DDZPokerSort(allPokes)

	// 统计各个牌的个数
	counts, keys := GetPokeCount(allPokes)

	// 找到的符合条件的牌的无花色权重
	var findPointWeight uint32 = 0

	// 由于前面已经排序，所以找到个数>=1，且无花色权重比speciPoke大的牌即可
	for i := 0; i < len(keys); i++ {
		pointWeight := (uint32)(keys[i])
		count := counts[pointWeight]
		if (count >= 1) && (pointWeight > speciPoke[0].PointWeight) {
			findPointWeight = pointWeight
			break
		}
	}

	// 找到了符合条件的牌
	if findPointWeight != 0 {
		// 找到这张牌
		for _, poke := range allPokes {
			if poke.PointWeight == findPointWeight {
				return true, []Poker{poke}
			}
		}
	}

	return false, []Poker{}
}

// GetMinBiggerPair 从allPokes中获取比指定对子speciPoke大的最小的对子
// 例如：77776655544 中找到比33大的牌，应该返回44;
// 例如：77776665444 中找到比33大的牌，应该返回空;
func GetMinBiggerPair(allPokes []Poker, speciPoke []Poker) (bool, []Poker) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "play.go:GetMinBiggerPair",
		"allPokes":  allPokes,
		"speciPoke": speciPoke,
	})

	// 参数检测
	if len(allPokes) == 0 || len(speciPoke) != 2 || (speciPoke[0].PointWeight != speciPoke[1].PointWeight) {
		logEntry.Errorln("参数错误")
		return false, nil
	}

	// 先排序，从小到大
	DDZPokerSort(allPokes)

	//Map<无花色权重点数, 牌的个数>
	counts, keys := GetPokeCount(allPokes)

	// 找到的符合条件的牌的无花色权重
	var findPointWeight uint32 = 0

	// 由于前面已经排序，所以找到个数>=2的，且无花色权重比speciPoke大的牌即可
	for i := 0; i < len(keys); i++ {
		pointWeight := (uint32)(keys[i])
		count := counts[pointWeight]
		if (count >= 2) && (pointWeight > speciPoke[0].PointWeight) {
			findPointWeight = pointWeight
			break
		}
	}

	resultPoke := []Poker{}

	// 找到了符合条件的牌
	if findPointWeight != 0 {
		// 找到这张牌及后面的牌
		for _, poke := range allPokes {
			if poke.PointWeight == findPointWeight {
				resultPoke = append(resultPoke, poke)

				if len(resultPoke) == 2 {
					break
				}
			}
		}
	}

	if len(resultPoke) == 2 {
		return true, resultPoke
	}

	return false, nil
}

// GetMinBiggerShunzi 从allPokes中获取比指定顺子speciPoke大的最小的顺子
// 例如：877776544 中找到比34567大的牌，应该返回45678;
// 例如：J109987775444 中找到比33大的牌，应该返回空;
func GetMinBiggerShunzi(allPokes []Poker, speciPoke []Poker) (bool, []Poker) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "play.go:GetMinBiggerShunzi",
		"allPokes":  allPokes,
		"speciPoke": speciPoke,
	})

	// 指定牌的长度
	lenSpecialPoke := len(speciPoke)

	// 参数检测
	if len(allPokes) == 0 || lenSpecialPoke < 5 {
		logEntry.Errorln("参数错误1")
		return false, nil
	}

	// 是否是顺子
	bSuc, _ := IsShunZi(speciPoke)
	if !bSuc {
		logEntry.Errorln("参数错误2，传入的不是顺子")
		return false, nil
	}

	// 顺子最小牌的无花色权重
	minPointWeight := GetMinCard(speciPoke).PointWeight

	// 先排序，从小到大
	DDZPokerSort(allPokes)

	//Map<无花色权重点数, 牌的个数>
	counts, _ := GetPokeCount(allPokes)

	// A的无花色权重
	pointWeightA := 14

	// 符合条件的最小顺子的起点牌的无花色权重
	var resultStartPokePoint uint32 = 0

	// 规则：若3-4-5-6-7的顺子，则从4-5-6-7-8开始判断，一直判断到10-11-12-13-14
	for startPokePoint := minPointWeight + 1; startPokePoint <= uint32(pointWeightA-lenSpecialPoke+1); startPokePoint++ {

		_, exist := counts[startPokePoint]

		// 起点牌必须存在
		if !exist {
			break
		}

		bAllExist := true

		// 后面的每张牌都要存在
		for secondPokePoint := startPokePoint + 1; secondPokePoint <= uint32(uint32(lenSpecialPoke)+startPokePoint-1); secondPokePoint++ {
			_, exist = counts[secondPokePoint]
			// 有一个不存在就失败，跳出
			if !exist {
				bAllExist = false
				break
			}
		}

		// 全部检测通过，说明存在最小顺子了，且startPokePoint就是起点
		if bAllExist {
			resultStartPokePoint = startPokePoint
			break
		}
	}

	// 没找到就返回吧
	if resultStartPokePoint == 0 {
		return false, nil
	}

	resultPoke := []Poker{}

	// 已经排序了，就从低往高遍历，找到需要的牌
	for i := 0; i < len(allPokes); i++ {
		if allPokes[i].PointWeight == resultStartPokePoint {
			resultPoke = append(resultPoke, allPokes[i])

			// 下次压入的是下一张牌
			resultStartPokePoint++

			// 牌数压够了，就返回吧
			if len(resultPoke) == lenSpecialPoke {
				return true, resultPoke
			}
		}
	}

	return false, nil
}

// GetMinBiggerPairs 从allPokes中获取比指定连对speciPoke大的最小的连对
// 例如：9988777766554 中找到比33445566大的牌，应该返回55667788;
// 例如：J99887665544 中找到比33445566大的牌，应该返回空;
func GetMinBiggerPairs(allPokes []Poker, speciPoke []Poker) (bool, []Poker) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "play.go:GetMinBiggerPairs",
		"allPokes":  allPokes,
		"speciPoke": speciPoke,
	})

	// 指定牌的长度
	lenSpecialPoke := len(speciPoke)

	// 参数检测
	if len(allPokes) == 0 || lenSpecialPoke < 6 || lenSpecialPoke%2 != 0 {
		logEntry.Errorln("参数错误1")
		return false, nil
	}

	// 是否是连对
	bSuc, _ := IsPairs(speciPoke)
	if !bSuc {
		logEntry.Errorln("参数错误2，传入的不是连对")
		return false, nil
	}

	// 连对最小牌的无花色权重
	minPointWeight := GetMinCard(speciPoke).PointWeight

	// 先排序，从小到大
	DDZPokerSort(allPokes)

	//Map<无花色权重点数, 牌的个数>
	counts, _ := GetPokeCount(allPokes)

	// A的无花色权重
	pointWeightA := 14

	// 符合条件的最小连对的起点牌的无花色权重
	var resultStartPokePoint uint32 = 0

	// 办法：若33-44-55-66的连对，则从44-55-66-77开始判断，一直判断到JJ-QQ-KK-AA
	for startPokePoint := minPointWeight + 1; startPokePoint <= uint32(pointWeightA-(lenSpecialPoke/2)+1); startPokePoint++ {
		count1, _ := counts[startPokePoint]
		// 连对的起点牌个数要>=2
		if count1 >= 2 {
			bAllExist := true

			// 后面的每张牌个数都>=2
			for secondPokePoint := startPokePoint + 1; secondPokePoint <= uint32(uint32(lenSpecialPoke/2)+startPokePoint-1); secondPokePoint++ {
				count, _ := counts[secondPokePoint]

				// 牌数不足则跳出
				if count < 2 {
					bAllExist = false
					break
				}
			}

			// 全部检测通过，说明存在最小连对了，且startPokePoint就是起点
			if bAllExist {
				resultStartPokePoint = startPokePoint
				break
			}
		}

	}

	// 没找到就返回吧
	if resultStartPokePoint == 0 {
		return false, nil
	}

	resultPoke := []Poker{}

	// 已经排序了，就从低往高遍历，找到需要的牌
	pushCount := 0
	for i := 0; i < len(allPokes); i++ {
		if allPokes[i].PointWeight == resultStartPokePoint {
			resultPoke = append(resultPoke, allPokes[i])
			pushCount++

			// 压入两张才下一个
			if pushCount == 2 {
				pushCount = 0
				// 这样下次压入的就是下一张牌了
				resultStartPokePoint++
			}

			// 牌数压够了，就返回吧
			if len(resultPoke) == len(speciPoke) {
				return true, resultPoke
			}
		}
	}

	return false, nil
}

// GetMinBiggerTriple 从allPokes中获取比指定三张speciPoke大的最小的三张
// 例如：77776655544 中找到比444大的牌，应该返回777;
// 例如：7766544 中找到比444大的牌，应该返回空;
func GetMinBiggerTriple(allPokes []Poker, speciPoke []Poker) (bool, []Poker) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "play.go:GetMinBiggerTriple",
		"allPokes":  allPokes,
		"speciPoke": speciPoke,
	})

	// 参数检测
	if len(allPokes) == 0 || len(speciPoke) != 3 || !IsAllSamePoint(speciPoke) {
		logEntry.Errorln("参数错误")
		return false, nil
	}

	// 先排序，从小到大
	DDZPokerSort(allPokes)

	//Map<无花色权重点数, 牌的个数>
	counts, keys := GetPokeCount(allPokes)

	// 找到的符合条件的牌的无花色权重
	var findPointWeight uint32 = 0

	// 由于前面已经排序，所以找到个数>=3，且无花色权重比speciPoke大的牌即可
	for i := 0; i < len(keys); i++ {
		pointWeight := (uint32)(keys[i])
		count := counts[pointWeight]
		if (count >= 3) && (pointWeight > speciPoke[0].PointWeight) {
			findPointWeight = pointWeight
			break
		}
	}

	resultPoke := []Poker{}

	// 找到了符合条件的牌
	if findPointWeight != 0 {
		// 找到这张牌及后面的牌
		for _, poke := range allPokes {
			if poke.PointWeight == findPointWeight {
				resultPoke = append(resultPoke, poke)
				if len(resultPoke) == 3 {
					break
				}
			}
		}
	}

	if len(resultPoke) == 3 {
		return true, resultPoke
	}

	return false, nil
}

// GetMinBigger3And1 从allPokes中获取比指定3带1speciPoke大的最小的3带1
// 例如：77776655544 中找到比4443大的牌，应该返回5554;
// 例如：7766544 中找到比4443大的牌，应该返回空;
func GetMinBigger3And1(allPokes []Poker, speciPoke []Poker) (bool, []Poker) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "play.go:GetMinBigger3And1",
		"allPokes":  allPokes,
		"speciPoke": speciPoke,
	})

	// 参数检测
	if len(allPokes) == 0 || len(speciPoke) != 4 {
		logEntry.Errorln("参数错误")
		return false, nil
	}

	// 参数检测
	samePokes := GetMaxSamePointCards(speciPoke)
	if len(samePokes) != 3 {
		logEntry.Errorln("参数错误2")
		return false, nil
	}

	// 先排序，从小到大
	DDZPokerSort(allPokes)

	//Map<无花色权重点数, 牌的个数>
	counts, keys := GetPokeCount(allPokes)

	// 找到的符合条件的牌的无花色权重
	var findPointWeight uint32 = 0

	// 由于前面已经排序，所以找到个数>=3，且无花色权重比speciPoke大的牌即可
	for i := 0; i < len(keys); i++ {
		pointWeight := (uint32)(keys[i])
		count := counts[pointWeight]
		if (count >= 3) && (pointWeight > speciPoke[0].PointWeight) {
			findPointWeight = pointWeight
			break
		}
	}

	resultPoke := []Poker{}

	// 找到了符合条件的主体牌
	if findPointWeight != 0 {
		// 压入主题牌
		for _, poke := range allPokes {
			if poke.PointWeight == findPointWeight {
				resultPoke = append(resultPoke, poke)
				if len(resultPoke) == 3 {
					break
				}
			}
		}

		// 移除这三个主体牌，再查找单张牌
		lastPoke := RemoveAll(allPokes, resultPoke)

		// 符合条件规则：剩下至少一张牌，且存在和主题牌不同的牌
		if len(lastPoke) >= 1 {

			// 重新排序
			DDZPokerSort(lastPoke)

			// 从小往大找
			for i := 0; i < len(lastPoke); i++ {
				if lastPoke[i].PointWeight != findPointWeight {
					// 压入单牌
					resultPoke = append(resultPoke, lastPoke[i])
					break
				}
			}
		}

		// 若处理完，结果仍为3张，说明没有找到合适的单牌，清空
		if len(resultPoke) == 3 {
			lastPoke = []Poker{}
		}
	}

	if len(resultPoke) == 4 {
		return true, resultPoke
	}

	return false, nil
}

// GetMinBigger3And2 从allPokes中获取比指定3带2speciPoke大的最小的3带2
// 例如：7777665554 中找到比44433大的牌，应该返回55566;
// 例如：765554 中找到比44433大的牌，应该返回空;
func GetMinBigger3And2(allPokes []Poker, speciPoke []Poker) (bool, []Poker) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "play.go:GetMinBigger3And2",
		"allPokes":  allPokes,
		"speciPoke": speciPoke,
	})

	// 参数检测
	if len(allPokes) == 0 || len(speciPoke) != 5 {
		logEntry.Errorln("参数错误")
		return false, nil
	}

	// 参数检测-主体为三张一样的
	samePokes := GetMaxSamePointCards(speciPoke)
	if len(samePokes) != 3 {
		logEntry.Errorln("参数错误2")
		return false, nil
	}

	// 参数检测-带的牌是2张一样的，且和主题牌不一样
	lastPokes := RemoveAll(speciPoke, samePokes)
	if !IsAllSamePoint(lastPokes) || lastPokes[0].PointWeight == samePokes[0].PointWeight {
		logEntry.Errorln("参数错误3")
		return false, nil
	}

	// 先排序，从小到大
	DDZPokerSort(allPokes)

	//Map<无花色权重点数, 牌的个数>
	counts, keys := GetPokeCount(allPokes)

	// 找到的符合条件的牌的无花色权重
	var findPointWeight uint32 = 0

	// 由于前面已经排序，所以找到个数>=3，且无花色权重比speciPoke大的牌即可
	for i := 0; i < len(keys); i++ {
		pointWeight := (uint32)(keys[i])
		count := counts[pointWeight]
		if (count >= 3) && (pointWeight > speciPoke[0].PointWeight) {
			findPointWeight = pointWeight
			break
		}
	}

	resultPoke := []Poker{}

	// 找到了符合条件的主体牌
	if findPointWeight != 0 {

		// 压入主体牌
		for _, poke := range allPokes {
			if poke.PointWeight == findPointWeight {
				resultPoke = append(resultPoke, poke)
				if len(resultPoke) == 3 {
					break
				}
			}
		}

		// 移除这三个主体牌，再查找对子
		lastPoke := RemoveAll(allPokes, resultPoke)

		// 符合条件规则：剩下的牌中存在对子
		if len(lastPoke) >= 2 {

			// 重新排序
			DDZPokerSort(lastPoke)

			//Map<无花色权重点数, 牌的个数>
			lastCounts, keys := GetPokeCount(lastPoke)

			var lastFindPointWeight uint32 = 0

			// 找到个数>=2的即可
			for i := 0; i < len(keys); i++ {
				pointWeight := (uint32)(keys[i])
				count := lastCounts[pointWeight]
				if count >= 2 {
					lastFindPointWeight = pointWeight
					break
				}
			}

			if lastFindPointWeight > 0 {
				// 压入对子
				for _, poke := range lastPoke {
					if poke.PointWeight == lastFindPointWeight {
						resultPoke = append(resultPoke, poke)

						// 对子压入后，满5张，跳出
						if len(resultPoke) == 5 {
							break
						}
					}
				}
			}
		}

		// 若处理完，结果仍为3张，说明没有找到合适的对子，清空
		if len(resultPoke) == 3 {
			lastPoke = []Poker{}
		}
	}

	if len(resultPoke) == 5 {
		return true, resultPoke
	}

	return false, nil
}

// GetMinBiggerTriples 从allPokes中获取比指定飞机speciPoke大的最小的飞机
// 例如：887777666554 中找到比333444大的牌，应该返回666777;
// 例如：998887766554 中找到比333444大的牌，应该返回空;
func GetMinBiggerTriples(allPokes []Poker, speciPoke []Poker) (bool, []Poker) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "play.go:GetMinBiggerTriples",
		"allPokes":  allPokes,
		"speciPoke": speciPoke,
	})

	lenSpiciPoke := len(speciPoke)

	// 参数检测
	if len(allPokes) == 0 || lenSpiciPoke < 6 || lenSpiciPoke%3 != 0 {
		logEntry.Errorln("参数错误1")
		return false, nil
	}

	// 是否是飞机
	bSuc, maxSamePoke := IsTriples(speciPoke)
	if !bSuc {
		logEntry.Errorln("参数错误2，传入的不是飞机")
		return false, nil
	}

	// 先排序，从小到大
	DDZPokerSort(allPokes)

	//Map<无花色权重点数, 牌的个数>
	counts, _ := GetPokeCount(allPokes)

	// A的无花色权重
	pointWeightA := 14

	// 符合条件的最小飞机的起点牌的无花色权重
	var resultStartPokePoint uint32 = 0

	// 办法：若333444555的飞机，则从666777888开始判断，一直判断到QQQKKKAAA
	for startPokePoint := maxSamePoke.PointWeight + 1; startPokePoint <= uint32(pointWeightA-(lenSpiciPoke/3)+1); startPokePoint++ {

		// 飞机的开始牌>=3
		count1, _ := counts[startPokePoint]
		if count1 >= 3 {

			bAllExist := true

			// 后面的每张牌都>=3
			for secondPokePoint := startPokePoint + 1; secondPokePoint <= uint32(uint32(lenSpiciPoke/3)+startPokePoint-1); secondPokePoint++ {
				count2, _ := counts[secondPokePoint]

				// 牌数不足则跳出
				if count2 < 3 {
					bAllExist = false
					break
				}
			}

			// 全部检测通过，说明存在最小飞机了，且startPokePoint就是起点
			if bAllExist {
				resultStartPokePoint = startPokePoint
				break
			}
		}
	}

	// 没找到就返回吧
	if resultStartPokePoint == 0 {
		return false, nil
	}

	resultPoke := []Poker{}

	// 已经排序了，就从低往高遍历，找到需要的牌
	pushCount := 0
	for i := 0; i < len(allPokes); i++ {
		if allPokes[i].PointWeight == resultStartPokePoint {
			resultPoke = append(resultPoke, allPokes[i])
			pushCount++

			// 压入三张才下一个
			if pushCount == 3 {
				pushCount = 0
				// 这样下次压入的就是下一张牌了
				resultStartPokePoint++
			}

			// 牌数压够了，就返回吧
			if len(resultPoke) == len(speciPoke) {
				return true, resultPoke
			}
		}
	}

	return false, nil
}

// GetMinBigger3sAnd1s 从allPokes中获取比指定飞机带单张speciPoke大的最小的飞机带单张
// 例如：88777666554 中找到比33344468大的牌，应该返回66677745，其中666777为主体牌，48为带牌，带牌时可拆;
// 例如：99888776664 中找到比33344468大的牌，应该返回空;
func GetMinBigger3sAnd1s(allPokes []Poker, speciPoke []Poker) (bool, []Poker) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "play.go:GetMinBigger3sAnd1s",
		"allPokes":  allPokes,
		"speciPoke": speciPoke,
	})

	lenSpiciPoke := len(speciPoke)

	// 参数检测
	if (len(allPokes) == 0) || (lenSpiciPoke < 8) || (lenSpiciPoke%4 != 0) {
		logEntry.Errorln("参数错误")
		return false, nil
	}

	// 剩余的指定牌
	lastPokes := make([]Poker, len(speciPoke))
	copy(lastPokes, speciPoke)

	// 是否是飞机带单张
	bSuc, maxSamePoker := IsTriplesAndSingles(speciPoke)
	if !bSuc {
		logEntry.Errorln("参数错误2，传入的不是飞机带单张")
		return false, nil
	}

	// 最大相同的牌的无花色权重
	maxSamePokePoint := maxSamePoker.PointWeight

	// 先对手牌排序，从小到大
	DDZPokerSort(allPokes)

	//Map<无花色权重点数, 牌的个数>
	counts, _ := GetPokeCount(allPokes)

	// A的无花色权重
	pointWeightA := 14

	// 符合条件的最小飞机的起点牌的无花色权重
	var resultStartPokePoint uint32 = 0

	// 办法：若333444555789的飞机带单张，则从666777888开始判断，一直判断到QQQKKKAAA
	for startPokePoint := maxSamePokePoint + 1; startPokePoint <= uint32(pointWeightA-(lenSpiciPoke/4)+1); startPokePoint++ {

		// 飞机的开始牌>=3
		count, _ := counts[startPokePoint]
		if count >= 3 {

			// 标记后续的是否都存在
			bSecondSuc := true

			// 后面的每张牌都>=3
			for secondPokePoint := startPokePoint + 1; secondPokePoint <= uint32(uint32(lenSpiciPoke/4)+startPokePoint-1); secondPokePoint++ {
				count, _ = counts[secondPokePoint]

				// 牌数不足则跳出
				if count < 3 {
					bSecondSuc = false
					break
				}
			}

			// 后续的全部检测通过，说明存在最小连对了，且startPokePoint就是起点
			if bSecondSuc {
				resultStartPokePoint = startPokePoint
				break
			}
		}
	}

	// 没找到就返回吧
	if resultStartPokePoint == 0 {
		return false, nil
	}

	// 最终发出的牌
	resultPoke := []Poker{}

	// 已经排序了，就从低往高遍历，找到需要的主体牌
	pushCount := 0
	for i := 0; i < len(allPokes); i++ {
		if allPokes[i].PointWeight == resultStartPokePoint {
			resultPoke = append(resultPoke, allPokes[i])
			pushCount++

			// 压入3张才下一个
			if pushCount == 3 {
				// 这样下次压入的就是下一张牌了
				resultStartPokePoint++
				pushCount = 0
			}

			// 主体牌数压够了，就跳出
			if len(resultPoke) == (len(speciPoke) / 4 * 3) {
				break
			}
		}
	}

	// 此时resultPoke的长度应为 lenSpiciPoke/4 * 3
	if len(resultPoke) != (len(speciPoke) / 4 * 3) {
		logEntry.Errorf("主体牌压入完毕后，牌的张数错误，期望：%d，实际：%d", len(speciPoke)/4*3, len(resultPoke))
		return false, nil
	}

	// 移除主体牌后剩下的牌
	lastPokes = RemoveAll(allPokes, resultPoke)

	// 寻找需要的单张，张数为：len(speciPoke)/4，且和主体牌不同，可以拆牌

	// 重新排序
	DDZPokerSort(lastPokes)

	// 从小往大找
	singleCount := 0
	for i := 0; i < len(lastPokes); i++ {

		// 直接压入，点数相同也没关系，比如前面666,后面可以带6
		//if ContainsPoint(resultPoke, lastPokes[i].Point) == false {
		resultPoke = append(resultPoke, lastPokes[i])
		singleCount++

		// 压入够了就跳出lastPokes
		if singleCount == len(speciPoke)/4 {
			break
		}
		//}
	}

	// 最终的牌数不同，则说明没找到足够的单张，失败
	if len(resultPoke) != len(speciPoke) {
		return false, nil
	}

	return true, resultPoke
}

// GetMinBigger3sAnd2s 从allPokes中获取比指定飞机带对子speciPoke大的最小的飞机带对子
// 例如：AAAKKKQQQJJ1010998 中找到比888777666554433大的牌，应该返回AAAKKKQQQJJ101099，其中AAAKKKQQQ为主体牌，J101099为带牌，带牌时可拆;
// 例如：KKKQQQJJJ101099876 中找到比888777666554433大的牌，应该返回空;
func GetMinBigger3sAnd2s(allPokes []Poker, speciPoke []Poker) (bool, []Poker) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "play.go:GetMinBigger3sAnd2s",
		"allPokes":  allPokes,
		"speciPoke": speciPoke,
	})

	lenSpiciPoke := len(speciPoke)

	// 参数检测
	if (len(allPokes) == 0) || (lenSpiciPoke < 10) || (lenSpiciPoke%5 != 0) {
		logEntry.Errorln("参数错误")
		return false, nil
	}

	// 剩余的指定牌
	lastPokes := make([]Poker, len(speciPoke))
	copy(lastPokes, speciPoke)

	// 是否是飞机带对子
	bSuc, maxSamePoker := IsTriplesAndPairs(speciPoke)
	if !bSuc {
		logEntry.Errorln("参数错误2，传入的不是飞机带对子")
		return false, nil
	}

	// 最大相同的牌的无花色权重
	maxSamePokePoint := maxSamePoker.PointWeight

	// 先把手牌排序，从小到大
	DDZPokerSort(allPokes)

	//Map<无花色权重点数, 牌的个数>
	counts, _ := GetPokeCount(allPokes)

	// A的无花色权重
	pointWeightA := 14

	// 符合条件的最小飞机的起点牌的无花色权重
	var resultStartPokePoint uint32 = 0

	// 办法：若666777888554433的飞机带对子，则从999101010JJJ开始判断，一直判断到QQQKKKAAA
	for startPokePoint := maxSamePokePoint + 1; startPokePoint <= uint32(pointWeightA-(lenSpiciPoke/5)+1); startPokePoint++ {

		// 飞机的开始牌>=3
		count, _ := counts[startPokePoint]
		if count >= 3 {

			// 标记后续的是否都存在
			bSecondSuc := true

			// 后面的每张牌都>=3
			for secondPokePoint := startPokePoint + 1; secondPokePoint <= uint32(uint32(lenSpiciPoke/5)+startPokePoint-1); secondPokePoint++ {
				count, _ = counts[secondPokePoint]

				// 牌数不足则跳出
				if count < 3 {
					bSecondSuc = false
					break
				}
			}

			// 后续的全部检测通过，说明存在最小连对了，且startPokePoint就是起点
			if bSecondSuc {
				resultStartPokePoint = startPokePoint
				break
			}
		}
	}

	// 没找到就返回吧
	if resultStartPokePoint == 0 {
		return false, nil
	}

	// 最终发出的牌
	resultPoke := []Poker{}

	// 已经排序了，就从低往高遍历，找到需要的主体牌
	pushCount := 0
	for i := 0; i < len(allPokes); i++ {
		if allPokes[i].PointWeight == resultStartPokePoint {
			resultPoke = append(resultPoke, allPokes[i])
			pushCount++

			// 压入3张才下一个
			if pushCount == 3 {
				// 这样下次压入的就是下一张牌了
				resultStartPokePoint++
				pushCount = 0
			}

			// 主体牌数压够了，就跳出
			if len(resultPoke) == (len(speciPoke) / 5 * 3) {
				break
			}
		}
	}

	// 此时resultPoke的长度应为 lenSpiciPoke/5 * 3
	if len(resultPoke) != (len(speciPoke) / 5 * 3) {
		logEntry.Errorf("主体牌压入完毕后，牌的张数错误，期望：%d，实际：%d", len(speciPoke)/5*3, len(resultPoke))
		return false, nil
	}

	// 移除主体牌后剩下的牌
	lastPokes = RemoveAll(allPokes, resultPoke)

	// 寻找需要的对子，张数为：len(speciPoke)/5 * 2，且和主体牌不同，可以拆牌

	// 重新排序
	DDZPokerSort(lastPokes)

	//Map<无花色权重点数, 牌的个数>
	counts, keys := GetPokeCount(lastPokes)

	// 存在的对子的无花色权重
	pairPointWeight := []uint32{}

	for i := 0; i < len(keys); i++ {
		pointWeight := (uint32)(keys[i])
		count := counts[pointWeight]

		// 4个牌时认为是两个对子
		if count >= 4 {
			pairPointWeight = append(pairPointWeight, pointWeight)
			pairPointWeight = append(pairPointWeight, pointWeight)
		} else if count >= 2 {
			pairPointWeight = append(pairPointWeight, pointWeight)
		}
	}

	// 对子是否足够
	if len(pairPointWeight) < len(speciPoke)/5 {
		return false, nil
	}

	// 压入所有的对子
	for i := 0; i < len(pairPointWeight); i++ {

		// 压入牌的数量
		pushCount := 0

		// 压入的牌
		pushPoke := []Poker{}

		// 遍历所有的剩余牌
		for j := 0; j < len(lastPokes); j++ {
			// 是前面找到的对子牌，则压入
			if pairPointWeight[i] == lastPokes[j].PointWeight {

				// 压入
				resultPoke = append(resultPoke, lastPokes[j])
				pushCount++

				pushPoke = append(pushPoke, lastPokes[j])

				// 压入两张就跳出，因为pairPointWeight[i]这张牌压入结束
				if pushCount == 2 {
					break
				}
			}
		}

		// 压入足够的牌后就跳出，防止压入过多的对子
		if len(resultPoke) == len(speciPoke) {
			break
		}

		// 压入结束后，删除这两张牌
		lastPokes = RemoveAll(lastPokes, pushPoke)
	}

	// 最终的牌数不同，则说明没找到足够的对子，失败
	if len(resultPoke) != len(speciPoke) {
		return false, nil
	}

	return true, resultPoke
}

// GetMinBigger4sAnd1s 从allPokes中获取比指定4带2单张speciPoke大的最小的4带2单张
// 例如：87777666554 中找到比333378大的牌，应该返回777745，其中7777为主体牌，45为带牌，只比较主体牌，带牌时可拆，带的单牌需要和主题牌不同;
// 例如：87776665544 中找到比333378大的牌，应该返回空;
func GetMinBigger4sAnd1s(allPokes []Poker, speciPoke []Poker) (bool, []Poker) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "play.go:GetMinBigger4sAnd1s",
		"allPokes":  allPokes,
		"speciPoke": speciPoke,
	})

	lenSpiciPoke := len(speciPoke)

	// 参数检测
	if (len(allPokes) == 0) || (lenSpiciPoke != 6) {
		logEntry.Errorln("参数错误")
		return false, nil
	}

	// 是否是4带2单张
	bSuc, maxSamePoker := IsBombAndSingles(speciPoke)
	if !bSuc {
		logEntry.Errorln("参数错误2，传入的不是4带2单张")
		return false, nil
	}

	// 先排序，从小到大
	DDZPokerSort(allPokes)

	//Map<无花色权重点数, 牌的个数>
	counts, keys := GetPokeCount(allPokes)

	// 找到的符合条件的牌的无花色权重
	var findPointWeight uint32 = 0

	// 由于前面已经排序，所以找到个数>=4，且无花色权重比speciPoke中主体牌大的牌即可
	for i := 0; i < len(keys); i++ {
		pointWeight := (uint32)(keys[i])
		count := counts[pointWeight]
		if (count >= 4) && (pointWeight > maxSamePoker.PointWeight) {
			findPointWeight = pointWeight
			break
		}
	}

	resultPoke := []Poker{}

	// 找到了符合条件的主体牌
	if findPointWeight != 0 {
		// 压入主题牌
		for _, poke := range allPokes {
			if poke.PointWeight == findPointWeight {
				resultPoke = append(resultPoke, poke)
				if len(resultPoke) == 4 {
					break
				}
			}
		}

		// 移除主体牌，再查找2个单张
		lastPoke := RemoveAll(allPokes, resultPoke)

		// 符合条件规则：剩下至少一张牌，且是和主题牌不同的牌
		if len(lastPoke) >= 2 {

			// 剩下的牌重新排序
			DDZPokerSort(lastPoke)

			// 从小往大找
			for i := 0; i < len(lastPoke); i++ {
				// 和主体牌不同即可压入
				if lastPoke[i].PointWeight != findPointWeight {
					resultPoke = append(resultPoke, lastPoke[i])

					if len(resultPoke) == 6 {
						break
					}
				}
			}

			// 现在应该是6张牌了，不是则置空，并报错
			if len(resultPoke) != 6 {
				resultPoke = []Poker{}
				logEntry.Errorln("添加两个单牌后发现总牌数不是6张")
			}
		}

		// 若处理完，结果不是6张，说明没有找到合适的2张单牌，清空
		if len(resultPoke) != 6 {
			resultPoke = []Poker{}
		}
	}

	if len(resultPoke) == 6 {
		return true, resultPoke
	}

	return false, nil
}

// GetMinBigger4sAnd2s 从allPokes中获取比指定4带2对speciPoke大的最小的4带2对
// 例如：87777666554 中找到比33338899大的牌，应该返回77775566，其中7777为主体牌，5566为带牌，只比较主体牌，带对子时可拆，带的对子需要和主题牌不同;
// 例如：J9877666654 中找到比33338899大的牌，应该返回空;
func GetMinBigger4sAnd2s(allPokes []Poker, speciPoke []Poker) (bool, []Poker) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "play.go:GetMinBigger4sAnd2s",
		"allPokes":  allPokes,
		"speciPoke": speciPoke,
	})

	lenSpiciPoke := len(speciPoke)

	// 参数检测
	if (len(allPokes) == 0) || (lenSpiciPoke != 8) {
		logEntry.Errorln("参数错误")
		return false, nil
	}

	// 是否是4带2对
	bSuc, maxSamePoker := IsBombAndPairs(speciPoke)
	if !bSuc {
		logEntry.Errorln("参数错误2，传入的不是4带2对")
		return false, nil
	}

	// 先排序，从小到大
	DDZPokerSort(allPokes)

	//Map<无花色权重点数, 牌的个数>
	counts, keys := GetPokeCount(allPokes)

	// 找到的符合条件的牌的无花色权重
	var findPointWeight uint32 = 0

	// 由于前面已经排序，所以找到个数>=4，且无花色权重比speciPoke中主体牌大的牌即可
	for i := 0; i < len(keys); i++ {
		pointWeight := (uint32)(keys[i])
		count := counts[pointWeight]
		if (count >= 4) && (pointWeight > maxSamePoker.PointWeight) {
			findPointWeight = pointWeight
			break
		}
	}

	resultPoke := []Poker{}

	// 找到了符合条件的主体牌
	if findPointWeight != 0 {

		// 先压入主题牌
		for _, poke := range allPokes {
			if poke.PointWeight == findPointWeight {
				resultPoke = append(resultPoke, poke)
				if len(resultPoke) == 4 {
					break
				}
			}
		}

		// 移除主体牌，再查找2个对子
		lastPoke := RemoveAll(allPokes, resultPoke)

		// 符合条件规则：剩下至少4张牌，且是和主题牌不同的牌
		if len(lastPoke) >= 4 {

			// 剩下的牌重新排序
			DDZPokerSort(lastPoke)

			//Map<无花色权重点数, 牌的个数>
			lastCounts, keys := GetPokeCount(lastPoke)

			// 存在的对子的无花色权重
			pairPointWeight := []uint32{}

			for i := 0; i < len(keys); i++ {
				pointWeight := (uint32)(keys[i])
				count := lastCounts[pointWeight]

				// 4个牌时认为是两个对子
				if count >= 4 {
					pairPointWeight = append(pairPointWeight, pointWeight)
					pairPointWeight = append(pairPointWeight, pointWeight)
				} else if count >= 2 {
					pairPointWeight = append(pairPointWeight, pointWeight)
				}
			}

			// 对子是否足够
			if len(pairPointWeight) < 2 {
				return false, nil
			}

			// 压入所有的对子
			for i := 0; i < len(pairPointWeight); i++ {

				// 压入牌的数量
				pushCount := 0

				// 压入的牌
				pushPoke := []Poker{}

				// 遍历所有的剩余牌
				for j := 0; j < len(lastPoke); j++ {
					// 是前面找到的对子牌，则压入
					if pairPointWeight[i] == lastPoke[j].PointWeight {

						// 压入
						resultPoke = append(resultPoke, lastPoke[j])
						pushCount++

						pushPoke = append(pushPoke, lastPoke[j])

						// 压入两张就跳出，因为pairPointWeight[i]这种牌压入结束
						if pushCount == 2 {
							break
						}
					}
				}

				// 够8张就退出，防止压入过多对子
				if len(resultPoke) == 8 {
					break
				}

				// 压入结束后，删除这两张牌
				lastPoke = RemoveAll(lastPoke, pushPoke)
			}

			// 现在应该是8张牌了，不是则置空，并报错
			if len(resultPoke) != 8 {
				resultPoke = []Poker{}
				logEntry.Errorln("添加两个对子后发现总牌数不是8张")
			}
		}

		// 若处理完，结果不是8张，说明没有找到合适的2个对子，清空
		if len(resultPoke) != 8 {
			resultPoke = []Poker{}
		}
	}

	if len(resultPoke) == 8 {
		return true, resultPoke
	}

	return false, nil
}

// GetMinBiggerBoom 从allPokes中获取比指定炸弹speciPoke大的最小的炸弹
// 例如：777766555 中找到比4444大的牌，应该返回7777;
// 例如：755543333 中找到比4444大的牌，应该返回空;
func GetMinBiggerBoom(allPokes []Poker, speciPoke []Poker) (bool, []Poker) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "play.go:GetMinBiggerBoom",
		"allPokes":  allPokes,
		"speciPoke": speciPoke,
	})

	// 参数检测
	if len(allPokes) == 0 || len(speciPoke) != 4 || !IsAllSamePoint(speciPoke) {
		logEntry.Errorln("参数错误")
		return false, nil
	}

	// 先排序，从小到大
	DDZPokerSort(allPokes)

	//Map<无花色权重点数, 牌的个数>
	counts, keys := GetPokeCount(allPokes)

	// 找到的符合条件的牌的无花色权重
	var findPointWeight uint32 = 0

	// 由于前面已经排序，所以找到个数>=4的，且无花色权重比speciPoke大的牌即可
	for i := 0; i < len(keys); i++ {
		pointWeight := (uint32)(keys[i])
		count := counts[pointWeight]
		if (count >= 4) && (pointWeight > speciPoke[0].PointWeight) {
			findPointWeight = pointWeight
			break
		}
	}

	resultPoke := []Poker{}

	// 找到了符合条件的牌
	if findPointWeight != 0 {
		// 找到这张牌及后面的牌
		for _, poke := range allPokes {
			if poke.PointWeight == findPointWeight {
				resultPoke = append(resultPoke, poke)

				if len(resultPoke) == 4 {
					break
				}
			}
		}
	}

	if len(resultPoke) == 4 {
		return true, resultPoke
	}

	return false, nil
}

// GetMinBiggerKingBoom 从allPokes中获取比指定炸弹speciPoke大的最小的炸弹
// 没有牌能压得过火箭，直接返回
func GetMinBiggerKingBoom(allPokes []Poker, speciPoke []Poker) (bool, []Poker) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "play.go:GetMinBiggerKingBoom",
		"allPokes":  allPokes,
		"speciPoke": speciPoke,
	})

	// 参数检测
	if len(allPokes) == 0 || len(speciPoke) != 2 || !Contains(speciPoke, RedJoker) || !Contains(speciPoke, BlackJoker) {
		logEntry.Errorln("参数错误")
		return false, nil
	}

	// 没有牌能压得过火箭
	return false, nil
}
