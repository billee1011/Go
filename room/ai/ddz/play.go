package ddz

import (
	"fmt"
	"steve/room/interfaces"
	"steve/server_pb/ddz"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

type playStateAI struct {
}

// GenerateAIEvent 生成 出牌AI 事件
// 无论是超时、托管还是机器人，胡过了自动胡，没胡过的其他操作都默认弃， 并且产生相应的事件
func (playAI *playStateAI) GenerateAIEvent(params interfaces.AIEventGenerateParams) (result interfaces.AIEventGenerateResult, err error) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "play.go:GenerateAIEvent()"})

	// 产生的事件结果
	result, err = interfaces.AIEventGenerateResult{
		Events: []interfaces.AIEvent{},
	}, nil

	ddzContext := params.DDZContext

	// 当前玩家
	var curPlayer *ddz.Player
	for _, player := range ddzContext.GetPlayers() {
		if player.GetPalyerId() == params.PlayerID {
			curPlayer = player
		}
	}

	// 找不到玩家
	if curPlayer == nil {
		logEntry.Errorf("找不到玩家%d", params.PlayerID)
		return result, fmt.Errorf("找不到玩家%d", params.PlayerID)
	}

	// 不该自己打牌
	if ddzContext.GetCurrentPlayerId() != curPlayer.GetPalyerId() {
		logEntry.Errorf("此时不应该调用玩家%d的出牌AI", params.PlayerID)
		return result, fmt.Errorf("此时不应该调用玩家%d的出牌AI", params.PlayerID)
	}

	// 没有牌型时说明是主动打牌
	if ddzContext.GetCurCardType() == ddz.CardType_CT_NONE {

		// 主动产生
		if event := playAI.getActivePlayCardEvent(ddzContext, curPlayer); event != nil {
			result.Events = append(result.Events, *event)
		}
	} else {
		// 被动产生
		if event := playAI.getPassivePlayCardEvent(ddzContext, curPlayer); event != nil {
			result.Events = append(result.Events, *event)
		}
	}

	return
}

// Play 生成出牌请求事件(被动出牌)
func (playAI *playStateAI) getPassivePlayCardEvent(ddzContext *ddz.DDZContext, player *ddz.Player) *interfaces.AIEvent {
	// 最终打出去的牌
	resultCards := []uint32{}

	// 最终打出去的牌型
	resultCardType := ddz.CardType_CT_NONE

	// 玩家手中的牌
	handCards := player.GetHandCards()

	// 转换为poke
	handPokes := ToDDZCards(handCards)

	// 按照排序权重进行排序
	//DdzPokerSort(handPokes)

	// 当前牌型
	curCardType := ddzContext.GetCurCardType()

	// 上家出的牌，转换为poke
	curOutPokes := ToDDZCards(ddzContext.GetCurOutCards())

	// 是否有压制的牌
	bSuc := false

	// 压制牌的数组
	sendPukes := []Poker{}

	switch curCardType {
	case ddz.CardType_CT_SINGLE: // 单牌
		bSuc, sendPukes = GetMinBiggerSingle(handPokes, curOutPokes)
		break
	case ddz.CardType_CT_PAIR: // 对子
		bSuc, sendPukes = GetMinBiggerPair(handPokes, curOutPokes)
		break
	case ddz.CardType_CT_SHUNZI: // 顺子
		bSuc, sendPukes = GetMinBiggerShunzi(handPokes, curOutPokes)
		break
	case ddz.CardType_CT_PAIRS: // 连对
		bSuc, sendPukes = GetMinBiggerPairs(handPokes, curOutPokes)
		break
	case ddz.CardType_CT_TRIPLE: // 三张
		bSuc, sendPukes = GetMinBiggerTriple(handPokes, curOutPokes)
		break
	case ddz.CardType_CT_3AND1: // 三带一
		bSuc, sendPukes = GetMinBigger3And1(handPokes, curOutPokes)
		break
	case ddz.CardType_CT_3AND2: // 三带二
		bSuc, sendPukes = GetMinBigger3And2(handPokes, curOutPokes)
		break
	case ddz.CardType_CT_TRIPLES: // 飞机
		bSuc, sendPukes = GetMinBiggerTriples(handPokes, curOutPokes)
		break
	case ddz.CardType_CT_3SAND1S: // 飞机带翅膀1，例：JJJQQQKKK + 856
		bSuc, sendPukes = GetMinBigger3sAnd1s(handPokes, curOutPokes)
		break
	case ddz.CardType_CT_3SAND2S: // 飞机带翅膀2，例：JJJQQQKKK + 885566
		bSuc, sendPukes = GetMinBiggerTriples(handPokes, curOutPokes)
		break
	case ddz.CardType_CT_4SAND1S: // 四带两个单张
		bSuc, sendPukes = GetMinBiggerTriples(handPokes, curOutPokes)
		break
	case ddz.CardType_CT_4SAND2S: // 四带两个对子
		bSuc, sendPukes = GetMinBiggerTriples(handPokes, curOutPokes)
		break
	case ddz.CardType_CT_BOMB: // 炸弹
		bSuc, sendPukes = GetMinBiggerBoom(handPokes, curOutPokes)
		break
	case ddz.CardType_CT_KINGBOMB: // 火箭
		bSuc, sendPukes = GetMinBiggerKingBoom(handPokes, curOutPokes)
		break
	}

	// 有压制的牌，则出的牌和上家牌型一致
	if bSuc {
		resultCardType = curCardType
	} else {

		// 无压制的牌，且当前牌型是炸弹，则判断自己有无火箭
		if !bSuc && curCardType == ddz.CardType_CT_BOMB {
			bSuc, sendPukes = GetKingBoom(handPokes)

			if bSuc {
				resultCardType = ddz.CardType_CT_KINGBOMB
			}
		}

		// 无压制的牌，且当前牌型不是炸弹，也不是火箭，则判断自己有无炸弹，无炸弹时再检测火箭
		if !bSuc && curCardType != ddz.CardType_CT_BOMB && curCardType != ddz.CardType_CT_KINGBOMB {

			// 优先检测炸弹
			bSuc, sendPukes = GetBoom(handPokes)
			if bSuc {
				resultCardType = ddz.CardType_CT_BOMB // 用炸弹来压
			} else {
				// 无炸弹时检测有无火箭
				bSuc, sendPukes = GetKingBoom(handPokes)
				if bSuc {
					resultCardType = ddz.CardType_CT_KINGBOMB // 用火箭来压
				}
			}
		}
	}

	// 下面是回复消息

	// 有压制的牌，转换数组
	if bSuc {
		resultCards = toInts(sendPukes)
	}

	request := ddz.PlayCardRequestEvent{
		Head: &ddz.RequestEventHead{
			PlayerId: player.GetPalyerId()},
		Cards:    resultCards,    // 打出去的牌
		CardType: resultCardType, // 打出去的牌型
	}

	data, _ := proto.Marshal(&request)

	event := interfaces.AIEvent{
		ID:      int32(ddz.EventID_event_chupai_request),
		Context: data,
	}

	return &event
}

// Play 生成出牌请求事件(主动出牌)
func (playAI *playStateAI) getActivePlayCardEvent(ddzContext *ddz.DDZContext, player *ddz.Player) *interfaces.AIEvent {

	// 玩家手中的牌
	handCards := player.GetHandCards()

	// 转换为poke
	handPokes := ToDDZCards(handCards)

	// 按照排序权重进行排序
	DdzPokerSort(handPokes)

	// 最终打出去的牌（打最小的那个牌）
	resultCards := []uint32{handPokes[0].toInt()}

	// 最终打出去的牌型（单张）
	resultCardType := ddz.CardType_CT_SINGLE

	// 下面是回复消息
	request := ddz.PlayCardRequestEvent{
		Head: &ddz.RequestEventHead{
			PlayerId: player.GetPalyerId()},
		Cards:    resultCards,    // 打出去的牌
		CardType: resultCardType, // 打出去的牌型
	}

	data, _ := proto.Marshal(&request)

	event := interfaces.AIEvent{
		ID:      int32(ddz.EventID_event_chupai_request),
		Context: data,
	}

	return &event
}

// GetPokeCount 统计各个牌的个数
// @inparam 	pokes ： 需统计的牌
// @outparam	map[uint32]uint32 :	key:牌的无花色权重，value:牌的个数
func GetPokeCount(pokes []Poker) map[uint32]uint32 {

	counts := make(map[uint32]uint32)

	for _, poke := range pokes {
		pointWeight := poke.PointWeight

		count, exists := counts[pointWeight]

		if !exists {
			counts[pointWeight] = 1
		} else {
			counts[pointWeight] = count + 1
		}
	}

	return counts
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
	DdzPokerSort(allPokes)

	// 统计各个牌的个数
	counts := GetPokeCount(allPokes)

	// 炸弹的无花色权重
	var boomPointWeight uint32 = 0
	for pointWeight, count := range counts {
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
	DdzPokerSort(allPokes)

	// 统计各个牌的个数
	counts := GetPokeCount(allPokes)

	// 找到的符合条件的牌的无花色权重
	var findPointWeight uint32 = 0

	// 由于前面已经排序，所以找到个数>=1，且无花色权重比speciPoke大的牌即可
	for pointWeight, count := range counts {
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
	DdzPokerSort(allPokes)

	//Map<无花色权重点数, 牌的个数>
	counts := GetPokeCount(allPokes)

	// 找到的符合条件的牌的无花色权重
	var findPointWeight uint32 = 0

	// 由于前面已经排序，所以找到个数>=2的，且无花色权重比speciPoke大的牌即可
	for pointWeight, count := range counts {
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

	// 顺子检测
	for i := 0; i < lenSpecialPoke-1; i++ {
		if speciPoke[i+1].PointWeight-speciPoke[i].PointWeight != 1 {
			logEntry.Errorln("参数错误2")
			return false, nil
		}
	}

	// 先排序，从小到大
	DdzPokerSort(allPokes)

	//Map<无花色权重点数, 牌的个数>
	counts := GetPokeCount(allPokes)

	// A的无花色权重
	pointWeightA := 14

	// 符合条件的最小顺子的起点牌的无花色权重
	var resultStartPokePoint uint32 = 0

	// 规则：若3-4-5-6-7的顺子，则从4-5-6-7-8开始判断，一直判断到10-11-12-13-14
	for startPokePoint := speciPoke[0].PointWeight + 1; startPokePoint <= uint32(pointWeightA-lenSpecialPoke+1); startPokePoint++ {

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

	// 连对检测
	for i := 0; i < lenSpecialPoke; i += 2 {
		if speciPoke[i+1].PointWeight != speciPoke[i].PointWeight {
			logEntry.Errorln("参数错误2")
			return false, nil
		}
	}

	// 先排序，从小到大
	DdzPokerSort(allPokes)

	//Map<无花色权重点数, 牌的个数>
	counts := GetPokeCount(allPokes)

	// A的无花色权重
	pointWeightA := 14

	// 符合条件的最小连对的起点牌的无花色权重
	var resultStartPokePoint uint32 = 0

	// 办法：若3-4-5-6-7的顺子，则从4-5-6-7-8开始判断，一直判断到10-11-12-13-14
	for startPokePoint := speciPoke[0].PointWeight + 1; startPokePoint <= uint32(pointWeightA-(lenSpecialPoke/2)+1); startPokePoint++ {
		count1, _ := counts[startPokePoint]
		// 连对的起点牌个数要>=2
		if count1 >= 2 {
			// 后面的每张牌个数都>=2
			for secondPokePoint := startPokePoint + 1; secondPokePoint <= uint32(uint32(lenSpecialPoke)+startPokePoint-2); secondPokePoint++ {
				count, _ := counts[secondPokePoint]

				// 牌数不足则跳出
				if count < 2 {
					break
				}
			}

			// 全部检测通过，说明存在最小连对了，且startPokePoint就是起点
			resultStartPokePoint = startPokePoint
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
	DdzPokerSort(allPokes)

	//Map<无花色权重点数, 牌的个数>
	counts := GetPokeCount(allPokes)

	// 找到的符合条件的牌的无花色权重
	var findPointWeight uint32 = 0

	// 由于前面已经排序，所以找到个数>=3，且无花色权重比speciPoke大的牌即可
	for pointWeight, count := range counts {
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
	DdzPokerSort(allPokes)

	//Map<无花色权重点数, 牌的个数>
	counts := GetPokeCount(allPokes)

	// 找到的符合条件的牌的无花色权重
	var findPointWeight uint32 = 0

	// 由于前面已经排序，所以找到个数>=3，且无花色权重比speciPoke大的牌即可
	for pointWeight, count := range counts {
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
			DdzPokerSort(lastPoke)

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
	DdzPokerSort(allPokes)

	//Map<无花色权重点数, 牌的个数>
	counts := GetPokeCount(allPokes)

	// 找到的符合条件的牌的无花色权重
	var findPointWeight uint32 = 0

	// 由于前面已经排序，所以找到个数>=3，且无花色权重比speciPoke大的牌即可
	for pointWeight, count := range counts {
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
			DdzPokerSort(lastPoke)

			//Map<无花色权重点数, 牌的个数>
			lastCounts := GetPokeCount(lastPoke)

			var lastFindPointWeight uint32 = 0

			// 找到个数>=2的即可
			for pointWeight, count := range lastCounts {
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

	// 飞机检测
	for i := 0; i < lenSpiciPoke; i += 3 {
		if speciPoke[i].PointWeight != speciPoke[i+1].PointWeight ||
			speciPoke[i].PointWeight != speciPoke[i+2].PointWeight {
			logEntry.Errorln("参数错误2")
			return false, nil
		}
	}

	// 先排序，从小到大
	DdzPokerSort(allPokes)

	//Map<无花色权重点数, 牌的个数>
	counts := GetPokeCount(allPokes)

	// A的无花色权重
	pointWeightA := 14

	// 符合条件的最小飞机的起点牌的无花色权重
	var resultStartPokePoint uint32 = 0

	// 办法：若333444555的飞机，则从666777888开始判断，一直判断到QQQKKKAAA
	for startPokePoint := speciPoke[lenSpiciPoke-1].PointWeight + 1; startPokePoint <= uint32(pointWeightA-(lenSpiciPoke/3)+1); startPokePoint++ {

		// 飞机的开始牌>=3
		count1, _ := counts[startPokePoint]
		if count1 >= 3 {

			// 后面的每张牌都>=3
			for secondPokePoint := startPokePoint + 1; secondPokePoint <= uint32(uint32(lenSpiciPoke)+startPokePoint-2); secondPokePoint++ {
				count2, _ := counts[secondPokePoint]

				// 牌数不足则跳出
				if count2 < 3 {
					break
				}
			}

			// 全部检测通过，说明存在最小连对了，且startPokePoint就是起点
			resultStartPokePoint = startPokePoint
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

// GetMinBigger3sAnd1s 从allPokes中获取比指定飞机带翅膀speciPoke大的最小的飞机带翅膀
// 例如：887777666554 中找到比333444大的牌，应该返回666777;
// 例如：998887766554 中找到比333444大的牌，应该返回空;
func GetMinBigger3sAnd1s(allPokes []Poker, speciPoke []Poker) (bool, []Poker) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "play.go:GetMinBigger3sAnd1s",
		"allPokes":  allPokes,
		"speciPoke": speciPoke,
	})

	lenSpiciPoke := len(speciPoke)

	// 参数检测
	if len(allPokes) == 0 || lenSpiciPoke < 6 || lenSpiciPoke%3 != 0 {
		logEntry.Errorln("参数错误1")
		return false, nil
	}

	// 飞机检测
	for i := 0; i < lenSpiciPoke; i += 3 {
		if speciPoke[i].PointWeight != speciPoke[i+1].PointWeight ||
			speciPoke[i].PointWeight != speciPoke[i+2].PointWeight {
			logEntry.Errorln("参数错误2")
			return false, nil
		}
	}

	// 先排序，从小到大
	DdzPokerSort(allPokes)

	//Map<无花色权重点数, 牌的个数>
	counts := GetPokeCount(allPokes)

	// A的无花色权重
	pointWeightA := 14

	// 符合条件的最小飞机的起点牌的无花色权重
	var resultStartPokePoint uint32 = 0

	// 办法：若333444555的飞机，则从666777888开始判断，一直判断到QQQKKKAAA
	for startPokePoint := speciPoke[lenSpiciPoke-1].PointWeight + 1; startPokePoint <= uint32(pointWeightA-(lenSpiciPoke/3)+1); startPokePoint++ {

		// 飞机的开始牌>=3
		count1, _ := counts[startPokePoint]
		if count1 >= 3 {

			// 后面的每张牌都>=3
			for secondPokePoint := startPokePoint + 1; secondPokePoint <= uint32(uint32(lenSpiciPoke)+startPokePoint-2); secondPokePoint++ {
				count2, _ := counts[secondPokePoint]

				// 牌数不足则跳出
				if count2 < 3 {
					break
				}
			}

			// 全部检测通过，说明存在最小连对了，且startPokePoint就是起点
			resultStartPokePoint = startPokePoint
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
	DdzPokerSort(allPokes)

	//Map<无花色权重点数, 牌的个数>
	counts := GetPokeCount(allPokes)

	// 找到的符合条件的牌的无花色权重
	var findPointWeight uint32 = 0

	// 由于前面已经排序，所以找到个数>=4的，且无花色权重比speciPoke大的牌即可
	for pointWeight, count := range counts {
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
	if len(allPokes) == 0 || len(speciPoke) != 2 || !Contains(speciPoke, redJoker) || !Contains(speciPoke, blackJoker) {
		logEntry.Errorln("参数错误")
		return false, nil
	}

	// 没有牌能压得过火箭
	return false, nil
}
