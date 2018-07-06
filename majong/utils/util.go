package utils

import (
	"fmt"
	"steve/client_pb/room"
	"steve/gutils"
	"steve/majong/interfaces"
	majongpb "steve/server_pb/majong"

	"github.com/golang/protobuf/proto"
)

//GetPlayerByID 根据玩家id获取玩家
func GetPlayerByID(players []*majongpb.Player, id uint64) *majongpb.Player {
	for _, player := range players {
		if player.PalyerId == id {
			return player
		}
	}
	return nil
}

//GetNextPlayerByID 根据玩家id获取下个玩家
func GetNextPlayerByID(players []*majongpb.Player, id uint64) *majongpb.Player {
	for k, player := range players {
		if player.PalyerId == id {
			index := (k + 1) % len(players)
			return players[index]
		}
	}
	return nil
}

//CardsToUtilCards 用来辅助查ting胡的工具,将Card转为适合查胡的工具
func CardsToUtilCards(cards []*majongpb.Card) []Card {
	cardsCard := make([]Card, 0)
	for _, v := range cards {
		cardInt, _ := CardToInt(*v)
		cardsCard = append(cardsCard, Card(*cardInt))
	}
	return cardsCard
}

//HuCardsToUtilCards 用来辅助查ting胡的工具,将Card转为适合查胡的工具
func HuCardsToUtilCards(cards []*majongpb.HuCard) []Card {
	cardsCard := make([]Card, 0)
	for _, v := range cards {

		cardInt, _ := CardToInt(*v.Card)
		cardsCard = append(cardsCard, Card(*cardInt))
	}
	return cardsCard
}

//CardToInt 将Card转换成牌值（int32）
func CardToInt(card majongpb.Card) (*int32, error) {
	var color int32
	switch card.GetColor() {
	case majongpb.CardColor_ColorWan:
		color = 1
	case majongpb.CardColor_ColorTiao:
		color = 2
	case majongpb.CardColor_ColorTong:
		color = 3
	default:
		return &color, fmt.Errorf("cant trans card ")
	}
	tValue := card.Point
	value := color*10 + tValue
	return &value, nil
}

//DeleteIntCardFromLast 从int32类型的牌组中，找到对应的目标牌，并且移除
func DeleteIntCardFromLast(cards []int32, targetCard int32) ([]int32, bool) {
	index := -1
	l := len(cards)
	if l == 0 {
		return cards, false
	}
	for i := l - 1; i >= 0; i-- {
		if targetCard == cards[i] {
			index = i
			break
		}
	}
	if index != -1 {
		cards = append(cards[:index], cards[index+1:]...)
	}
	return cards, index != -1
}

//CardEqual 判断两张牌是否一样
func CardEqual(card1 *majongpb.Card, card2 *majongpb.Card) bool {
	return card1.GetColor() == card2.GetColor() && card1.GetPoint() == card2.GetPoint()
}

//DeleteCardFromLast 从majongpb.Card类型的牌组中，找到对应的目标牌，并且移除
func DeleteCardFromLast(cards []*majongpb.Card, targetCard *majongpb.Card) ([]*majongpb.Card, bool) {
	index := -1
	l := len(cards)
	if l == 0 {
		return cards, false
	}
	for i := l - 1; i >= 0; i-- {
		if CardEqual(targetCard, cards[i]) {
			index = i
			break
		}
	}
	if index != -1 {
		cards = append(cards[:index], cards[index+1:]...)
	}
	return cards, index != -1
}

// RemoveCards 从玩家的手牌中移除指定数量的某张牌
func RemoveCards(cards []*majongpb.Card, card *majongpb.Card, count int) ([]*majongpb.Card, bool) {
	newCards := []*majongpb.Card{}
	removeCount := 0
	for index, c := range cards {
		if CardEqual(c, card) {
			removeCount++
			if removeCount == count {
				newCards = append(newCards, cards[index+1:]...)
				break
			}
		} else {
			newCards = append(newCards, c)
		}
	}
	if removeCount != count {
		return cards, false
	}
	return newCards, true
}

//IntToCard int32类型转majongpb.Card类型
func IntToCard(cardValue int32) (*majongpb.Card, error) {
	colorValue := cardValue / 10
	value := cardValue % 10
	var color majongpb.CardColor
	switch colorValue {
	case 1:
		color = majongpb.CardColor_ColorWan
	case 2:
		color = majongpb.CardColor_ColorTiao
	case 3:
		color = majongpb.CardColor_ColorTong
	case 4:
		color = majongpb.CardColor_ColorFeng
	case 5:
		color = majongpb.CardColor_ColorHua
	default:
		return nil, fmt.Errorf("cant trans card %d", cardValue)
	}
	return &majongpb.Card{
		Color: color,
		Point: value,
	}, nil
}

//IntToRoomCard int32类型转room.Card类型
func IntToRoomCard(cardValue int32) (*room.Card, error) {
	colorValue := cardValue / 10
	value := cardValue % 10
	var color room.CardColor
	switch colorValue {
	case 1:
		color = room.CardColor_CC_WAN
	case 2:
		color = room.CardColor_CC_TIAO
	case 3:
		color = room.CardColor_CC_TONG
	default:
		return nil, fmt.Errorf("cant trans card %d", cardValue)
	}
	return &room.Card{
		Color: color.Enum(),
		Point: proto.Int32(value),
	}, nil
}

//CardToRoomCard majongpb.card类型转room.Card类型
func CardToRoomCard(card *majongpb.Card) (*room.Card, error) {
	var color room.CardColor
	if card.Color.String() == room.CardColor_CC_WAN.String() {
		color = room.CardColor_CC_WAN
	}
	if card.Color.String() == room.CardColor_CC_TIAO.String() {
		color = room.CardColor_CC_TIAO
	}
	if card.Color.String() == room.CardColor_CC_TONG.String() {
		color = room.CardColor_CC_TONG
	}

	return &room.Card{
		Color: color.Enum(),
		Point: proto.Int32(card.Point),
	}, nil
}

// ServerCard2Number 服务器的 Card 转换成数字
func ServerCard2Number(card *majongpb.Card) int {
	var color int
	if card.Color == majongpb.CardColor_ColorWan {
		color = 1
	} else if card.Color == majongpb.CardColor_ColorTiao {
		color = 2
	} else if card.Color == majongpb.CardColor_ColorTong {
		color = 3
	} else if card.Color == majongpb.CardColor_ColorFeng {
		color = 4
	}
	value := color*10 + int(card.Point)
	return value
}

// ServerCards2Numbers 服务器的 Card 数组转 int 数组
func ServerCards2Numbers(cards []*majongpb.Card) []int {
	result := []int{}
	for _, c := range cards {
		result = append(result, ServerCard2Number(c))
	}
	return result
}

// ServerCard2Uint32 服务器的 Card 转换成数字
func ServerCard2Uint32(card *majongpb.Card) uint32 {
	return uint32(ServerCard2Number(card))
}

// ServerCards2Uint32 服务器的 Card 数组转 int 数组
func ServerCards2Uint32(cards []*majongpb.Card) []uint32 {
	result := []uint32{}
	for _, c := range cards {
		result = append(result, ServerCard2Uint32(c))
	}
	return result
}

// CardsToRoomCards 将Card转换为room package中的Card
func CardsToRoomCards(cards []*majongpb.Card) []*room.Card {
	var rCards []*room.Card
	for i := 0; i < len(cards); i++ {
		rCards = append(rCards, &room.Card{
			Color: room.CardColor(cards[i].Color).Enum(),
			Point: &cards[i].Point,
		})
	}
	return rCards
}

// ContainCard 验证card是否存在于玩家手牌中，存在返回true,否则返回false
func ContainCard(cards []*majongpb.Card, card *majongpb.Card) bool {
	for i := 0; i < len(cards); i++ {
		if CardEqual(cards[i], card) {
			return true
		}
	}
	return false
}

//IntToUtilCard uint32类型的数组强转成类型
func IntToUtilCard(cards []int32) []Card {
	cardsCard := make([]Card, 0, 0)
	for _, v := range cards {

		utilCard := Card(v)
		cardsCard = append(cardsCard, utilCard)
	}
	return cardsCard
}

//ContainHuCards 判断当前可以胡的牌中是否包含已经胡过的所有牌
func ContainHuCards(targetHuCards []Card, HuCards []Card) bool {
	sameHuCards := 0
	for _, tagetCard := range targetHuCards {
		for _, Card := range HuCards {
			if tagetCard == Card {
				sameHuCards++
			}
		}
	}
	if len(HuCards) == sameHuCards {
		return true
	}
	return false
}

// CheckHuResult 查胡结果
type CheckHuResult struct {
	Can      bool
	Combines Combines // 推倒胡组合
}

// CheckHu 用来辅助胡牌查胡工具 cards玩家的所有牌，huCard点炮的牌（自摸时huCard为0）
// needCombines 是否需要返回所有组合
func CheckHu(cards []*majongpb.Card, huCard uint32, needCombines bool) CheckHuResult {
	result := CheckHuResult{}
	cardsCard := CardsToUtilCards(cards)
	if huCard > 0 {
		cardsCard = append(cardsCard, Card(huCard))
	}
	laizi := make(map[Card]bool)
	flag, combines := FastCheckHuV2(cardsCard, laizi, needCombines) // 检测玩家能否推倒胡
	canQidui := FastCheckQiDuiHu(cardsCard)
	result.Can = result.Can || flag || canQidui
	result.Combines = combines
	return result
}

//CheckHuUtilCardsToHandCards 将推到胡工具的util.Card转为玩家手牌的类型
func CheckHuUtilCardsToHandCards(cards []Card) ([]*majongpb.Card, error) {
	handCards := make([]*majongpb.Card, 0)
	for i := 0; i < len(cards); i++ {
		handCard, err := IntToCard(int32(cards[i]))
		if err != nil {
			return []*majongpb.Card{}, err
		}
		handCards = append(handCards, handCard)
	}
	return handCards, nil
}

//SeekCardSum 相同的牌的数量
func SeekCardSum(cards []*majongpb.Card, targetCard *majongpb.Card) int {
	count := 0
	for i := 0; i < len(cards); i++ {
		if CardEqual(cards[i], targetCard) {
			count++
		}
	}
	return count
}

//GetTingPlayerIDAndMultiple 获取所有听玩家,和返回每个听玩家最大倍数
func GetTingPlayerIDAndMultiple(players []*majongpb.Player, laizi map[Card]bool) (map[uint64]int64, error) {
	tingPlayers := make(map[uint64]int64, 0)
	for i := 0; i < len(players); i++ {
		// 胡过的不算
		if len(players[i].HuCards) > 0 {
			continue
		}
		// 查能不能听，能听，返回返回最大番型，及ID
		isTing, multiple, err := IsCanTingAndGetMultiple(players[i], laizi)
		if err != nil {
			return nil, err
		}
		if isTing {
			tingPlayers[players[i].GetPalyerId()] = multiple
		}
	}
	return tingPlayers, nil
}

//IsCanTingAndGetMultiple 判断玩家是否能听,和返回能听玩家的最大倍数 TODO
//未上听者需赔上听者最大可能番数（杠后炮、杠上开花、抢杠胡、海底捞、海底炮不参与）的牌型钱。注：查大叫时，
//若上听者牌型中有根，则根也要未上听者包给上听者。
func IsCanTingAndGetMultiple(player *majongpb.Player, laizi map[Card]bool) (bool, int64, error) {
	var max int64
	handCardSum := len(player.HandCards)
	//只差1张牌就能胡，并且玩家手牌不存在花牌
	if handCardSum%3 == 1 && !gutils.CheckHasDingQueCard(player.HandCards, player.DingqueColor) {
		cardCombines, err := GetTingCards(player.HandCards, laizi)
		if err != nil {
			return false, 0, err
		}
		handCards := player.GetHandCards()
		for card := range cardCombines {
			pbCard, _ := IntToCard(int32(card))
			handCards = append(handCards, pbCard)
			// TODO 获取最大番型
			mult := int64(2)
			if max < mult {
				max = mult
			}
			handCards = player.GetHandCards()
		}
	}
	return max > 0, max, nil
}

//GetTingCards 获取玩家能胡的牌,必须是缺一张
func GetTingCards(handCards []*majongpb.Card, laizi map[Card]bool) (CardCombines, error) {
	result := CardCombines{}

	if len(handCards)%3 != 1 {
		return result, fmt.Errorf("获取玩家能胡的牌,必须是缺一张")
	}
	cardsCard := CardsToUtilCards(handCards)
	// 推倒胡
	result = FastCheckTingV2(cardsCard, laizi)
	// 七对
	cardAll := []Card{11, 12, 13, 14, 15, 16, 17, 18, 19, 21, 22, 23, 24, 25, 26, 27, 28, 29, 31, 32, 33, 34, 35, 36, 37, 38, 39}
	qiCards := FastCheckQiDuiTing(cardsCard, cardAll)
	for _, card := range qiCards {
		if _, ok := result[card]; ok {
			continue
		}
		result[card] = Combines{}
	}
	return result, nil
}

//MergeAndNoRepeat 合并去重复UtilCard
func MergeAndNoRepeat(srcCards1 []Card, srcCards2 []Card) []Card {
	newCards := make([]Card, 0)
	newCards = append(newCards, srcCards1...)
	for _, card2 := range srcCards2 {
		flag := true
		for _, card1 := range srcCards1 {
			if card2 == card1 {
				flag = false
				break
			}
		}
		if flag {
			newCards = append(newCards, card2)
		}
	}
	return newCards
}

//GetFirstHuPlayerByID 获取第一个胡的玩家,winPlayers源数组， loserPlayerID输家ID
func GetFirstHuPlayerByID(playerAll, winPlayers []*majongpb.Player, loserPlayerID uint64) *majongpb.Player {
	// 获取输家的下家
	nextPlayer := GetNextPlayerByID(playerAll, loserPlayerID)
	for nextPlayer != nil {
		// 判断赢家里面是否有输家的下家
		for i := 0; i < len(winPlayers); i++ {
			if winPlayers[i].PalyerId == nextPlayer.PalyerId {
				return winPlayers[i]
			}
		}
		// 获取输家的下家的下家
		nextPlayer = GetNextPlayerByID(playerAll, nextPlayer.PalyerId)
	}
	return nil
}

// //GetPlayCardCheckTing 出牌查听，获取可以出那些牌，和出了这张牌，可以胡那些牌，返回map[Card][]Card
// func GetPlayCardCheckTing(handCards []*majongpb.Card, laizi map[Card]bool) map[Card][]Card {
// 	tingInfo := make(map[Card][]Card)
// 	// 不能少一张
// 	if len(handCards)%3 != 2 {
// 		return tingInfo
// 	}
// 	// 手牌转查胡的工具牌
// 	cardsCard := CardsToUtilCards(handCards)
// 	// 推倒胡查胡，打那张牌可以胡那些牌
// 	tingCombines := FastCheckTingInfoV2(cardsCard, laizi)
// 	for card, cardCombines := range tingCombines {
// 		tingcards := []Card{}
// 		for card := range cardCombines {
// 			tingcards = append(tingcards, card)
// 		}
// 		tingInfo[card] = tingcards
// 	}

// 	// 1-9所有牌
// 	cardAll := []Card{11, 12, 13, 14, 15, 16, 17, 18, 19, 21, 22, 23, 24, 25, 26, 27, 28, 29, 31, 32, 33, 34, 35, 36, 37, 38, 39}
// 	// 七对查胡，打那张牌可以胡那些牌
// 	qiStrategy := FastCheckQiDuiTingInfo(cardsCard, cardAll)
// 	// 存在相同的playCard,去重复
// 	for playCard, huCard := range tingInfo {
// 		tInfo, exite := qiStrategy[playCard]
// 		if exite {
// 			tingInfo[playCard] = MergeAndNoRepeat(tInfo, huCard)
// 		}
// 	}
// 	// 存在不相同的playCard,合并,把推倒胡中不存在的听，加进去
// 	for playCard, huCards := range qiStrategy {
// 		_, exite := tingInfo[playCard]
// 		if !exite {
// 			tingInfo[playCard] = huCards
// 		}
// 	}
// 	return tingInfo
// }

//GetPlayCardCheckTing 出牌查听，获取可以出那些牌，和出了这张牌，可以胡那些牌
// 返回可胡的牌与对应的组合
func GetPlayCardCheckTing(handCards []*majongpb.Card, laizi map[Card]bool) map[Card]CardCombines {
	result := make(map[Card]CardCombines)
	// 不能少一张
	if len(handCards)%3 != 2 {
		return result
	}
	// 手牌转查胡的工具牌
	cardsCard := CardsToUtilCards(handCards)
	// 推倒胡查胡，打那张牌可以胡那些牌
	result = FastCheckTingInfoV2(cardsCard, laizi)
	// 1-9所有牌
	cardAll := []Card{11, 12, 13, 14, 15, 16, 17, 18, 19, 21, 22, 23, 24, 25, 26, 27, 28, 29, 31, 32, 33, 34, 35, 36, 37, 38, 39}
	// 七对查胡，打那张牌可以胡那些牌
	qiStrategy := FastCheckQiDuiTingInfo(cardsCard, cardAll)

	for card, huCards := range qiStrategy {
		cardCombines := result[card]
		if cardCombines == nil {
			result[card] = make(CardCombines)
			cardCombines = result[card]
		}
		for _, huCard := range huCards {
			if cardCombines[huCard] == nil {
				cardCombines[huCard] = Combines{}
			}
		}
	}
	return result
}

// TransPengCard 碰牌转Card
func TransPengCard(pengCards []*majongpb.PengCard) []*majongpb.Card {
	cards := make([]*majongpb.Card, 0)
	for _, pengCard := range pengCards {
		cards = append(cards, pengCard.Card)
	}
	return cards
}

// TransGangCard 杠牌转Card
func TransGangCard(gangCards []*majongpb.GangCard) []*majongpb.Card {
	cards := make([]*majongpb.Card, 0)
	for _, gangCard := range gangCards {
		cards = append(cards, gangCard.Card)
	}
	return cards
}

// GetAllMopaiCount 获取所有人的摸牌数总和
func GetAllMopaiCount(mjContext *majongpb.MajongContext) int {
	count := 0
	for _, player := range mjContext.GetPlayers() {
		count += int(player.GetMopaiCount())
	}
	return count
}

// HasAvailableWallCards 判断是否有墙牌可摸
func HasAvailableWallCards(flow interfaces.MajongFlow) bool {
	context := flow.GetMajongContext()
	if len(context.WallCards) == 0 {
		return false
	}
	// 由配牌控制是否gameover,配牌长度为0走正常gameover,配牌长度不为0走配牌长度流局
	length := context.GetOption().GetWallcardsLength()
	if GetAllMopaiCount(context) == int(length)-53 {
		return false
	}
	return true
}

// CardsToInt card 转换
func CardsToInt(cards []*majongpb.Card) ([]int32, error) {
	return gutils.ServerCards2Int32(cards), nil
}
