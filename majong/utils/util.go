package utils

import (
	"fmt"
	"steve/client_pb/room"
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

//CheckHasDingQueCard 检查牌里面是否含有定缺的牌
func CheckHasDingQueCard(cards []*majongpb.Card, color majongpb.CardColor) bool {
	for _, card := range cards {
		if card.Color == color {
			return true
		}
	}
	return false
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

//CardsToInt 将Card转换成牌值（int32）
func CardsToInt(cards []*majongpb.Card) ([]int32, error) {
	var cardsInt []int32
	var color int32
	for _, card := range cards {
		switch card.GetColor() {
		case majongpb.CardColor_ColorWan:
			color = 1
		case majongpb.CardColor_ColorTiao:
			color = 2
		case majongpb.CardColor_ColorTong:
			color = 3
		default:
			return nil, fmt.Errorf("cant trans card ")
		}
		tValue := card.Point
		value := color*10 + tValue
		cardsInt = append(cardsInt, value)
	}
	return cardsInt, nil
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
		color = room.CardColor_ColorWan
	case 2:
		color = room.CardColor_ColorTiao
	case 3:
		color = room.CardColor_ColorTong
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
	if card.Color.String() == room.CardColor_ColorWan.String() {
		color = room.CardColor_ColorWan
	}
	if card.Color.String() == room.CardColor_ColorTiao.String() {
		color = room.CardColor_ColorTiao
	}
	if card.Color.String() == room.CardColor_ColorTong.String() {
		color = room.CardColor_ColorTong
	}

	return &room.Card{
		Color: color.Enum(),
		Point: proto.Int32(card.Point),
	}, nil
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

//CheckHu 用来辅助胡牌查胡工具 cards玩家的所有牌，huCard点炮的牌（自摸时huCard为0）
func CheckHu(cards []*majongpb.Card, huCard uint32) bool {
	cardsCard := CardsToUtilCards(cards)
	if huCard > 0 {
		cardsCard = append(cardsCard, Card(huCard))
	}
	// flag, _ := util.FastCheckHuV1(cardsCard) // 检测玩家能否推倒胡
	laizi := make(map[Card]bool)
	flag := FastCheckHuV2(cardsCard, laizi) // 检测玩家能否推倒胡
	if flag != true {
		flag = FastCheckQiDuiHu(cardsCard) // 检测玩家能否七对胡
	}
	return flag
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

//GetHuEdPlayers 获取胡过牌玩家
func GetHuEdPlayers(players []*majongpb.Player) []*majongpb.Player {
	huEdPlayers := make([]*majongpb.Player, 0)
	for i := 0; i < len(players); i++ {
		if len(players[i].HuCards) > 0 {
			huEdPlayers = append(huEdPlayers, players[i])
		}
	}
	return huEdPlayers
}

//GetBustedHandPlayers 获取未听玩家,包括花猪,isIncludeFlower-是否包含花猪，true 包含，false 不包含
func GetBustedHandPlayers(players []*majongpb.Player, isIncludeFlower bool) ([]*majongpb.Player, error) {
	bustedHandPlayers := make([]*majongpb.Player, 0)
	for i := 0; i < len(players); i++ {
		// 胡过的不算
		if len(players[i].HuCards) > 0 {
			continue
		}
		//查听
		isTing, _, err := IsCanTingAndGetMultiple(players[i])
		if err != nil {
			return []*majongpb.Player{}, err
		}
		// 不能听
		if !isTing && (isIncludeFlower || !IsFlowerPig(players[i])) {
			bustedHandPlayers = append(bustedHandPlayers, players[i])
		}
	}
	return bustedHandPlayers, nil
}

//GetFlowerPigPlayers 获取花猪玩家
func GetFlowerPigPlayers(players []*majongpb.Player) []*majongpb.Player {
	flowerPigPlayers := make([]*majongpb.Player, 0)
	for i := 0; i < len(players); i++ {
		if IsFlowerPig(players[i]) {
			flowerPigPlayers = append(flowerPigPlayers, players[i])
		}
	}
	return flowerPigPlayers
}

//IsOutNoDingQueColorCard 玩家properties中的key，代表玩家是否出过非定缺颜色的牌
const IsOutNoDingQueColorCard = "isoutnodingquecolorcard"

//IsFlowerPig 判断玩家是否是花猪 TODO
func IsFlowerPig(bustedHandPlayer *majongpb.Player) bool {
	//在出牌逻辑设置玩家一旦出过非定缺颜色的牌，[]byte{1}, 玩家是否出过非定缺的牌 TODO
	if len(bustedHandPlayer.Properties[IsOutNoDingQueColorCard]) != 0 {
		// 玩家手牌中是否存在花牌
		return CheckHasDingQueCard(bustedHandPlayer.HandCards, bustedHandPlayer.DingqueColor)
	}
	return false
}

//GetTingPlayerIDAndMultiple 获取所有听玩家,和返回每个听玩家最大倍数
func GetTingPlayerIDAndMultiple(players []*majongpb.Player) (map[uint64]int64, error) {
	tingPlayers := make(map[uint64]int64, 0)
	for i := 0; i < len(players); i++ {
		// 胡过的不算
		if len(players[i].HuCards) > 0 {
			continue
		}
		// 查能不能听，能听，返回返回最大番型，及ID
		isTing, multiple, err := IsCanTingAndGetMultiple(players[i])
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
func IsCanTingAndGetMultiple(player *majongpb.Player) (bool, int64, error) {
	var max int64
	handCardSum := len(player.HandCards)
	//只差1张牌就能胡，并且玩家手牌不存在花牌
	if handCardSum%3 == 1 && !CheckHasDingQueCard(player.HandCards, player.DingqueColor) {
		tingCards, err := GetTingCards(player.HandCards)
		if err != nil {
			return false, 0, err
		}
		handCards := player.GetHandCards()
		for i := 0; i < len(tingCards); i++ {
			handCards = append(handCards, tingCards[i])
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
func GetTingCards(handCards []*majongpb.Card) ([]*majongpb.Card, error) {
	if len(handCards)%3 != 1 {
		return []*majongpb.Card{}, fmt.Errorf("获取玩家能胡的牌,必须是缺一张")
	}
	cardsCard := CardsToUtilCards(handCards)
	laizi := make(map[Card]bool)
	// 推倒胡
	huCards := FastCheckTingV2(cardsCard, laizi)
	// 七对
	cardAll := []Card{11, 12, 13, 14, 15, 16, 17, 18, 19, 21, 22, 23, 24, 25, 26, 27, 28, 29, 31, 32, 33, 34, 35, 36, 37, 38, 39}
	qiCards := FastCheckQiDuiTing(cardsCard, cardAll)
	// 合并去重复
	tingCards := MergeAndNoRepeat(huCards, qiCards)
	newTingCards, err := CheckHuUtilCardsToHandCards(tingCards)
	return newTingCards, err
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

//GetDingQueCardSum 获取定缺牌数量
func GetDingQueCardSum(handCards []*majongpb.Card, dingQueColor majongpb.CardColor) int {
	sum := 0
	for _, card := range handCards {
		if card.Color == dingQueColor {
			sum++
		}
	}
	return sum
}

//GetPlayCardCheckTing 出牌查听，获取可以出那些牌，和出了这张牌，可以胡那些牌，返回map[Card][]Card
func GetPlayCardCheckTing(handCards []*majongpb.Card) map[Card][]Card {
	tingInfo := make(map[Card][]Card)
	// 不能少一张
	if len(handCards)%3 != 2 {
		return tingInfo
	}
	// 手牌转查胡的工具牌
	cardsCard := CardsToUtilCards(handCards)
	laizi := make(map[Card]bool)
	// 推倒胡查胡，打那张牌可以胡那些牌
	tingInfo = FastCheckTingInfoV2(cardsCard, laizi)
	// 1-9所有牌
	cardAll := []Card{11, 12, 13, 14, 15, 16, 17, 18, 19, 21, 22, 23, 24, 25, 26, 27, 28, 29, 31, 32, 33, 34, 35, 36, 37, 38, 39}
	// 七对查胡，打那张牌可以胡那些牌
	qiStrategy := FastCheckQiDuiTingInfo(cardsCard, cardAll)
	// 存在相同的playCard,去重复
	for playCard, huCard := range tingInfo {
		tInfo, exite := qiStrategy[playCard]
		if exite {
			tingInfo[playCard] = MergeAndNoRepeat(tInfo, huCard)
		}
	}
	// 存在不相同的playCard,合并,把推倒胡中不存在的听，加进去
	for playCard, huCards := range qiStrategy {
		_, exite := tingInfo[playCard]
		if !exite {
			tingInfo[playCard] = huCards
		}
	}
	return qiStrategy
}

//GetPlayCardHint 出牌提示，出牌这张牌，提示胡的牌和胡的牌的倍数，返回map[int32]map[int32]int64, error
func GetPlayCardHint(palyer *majongpb.Player) (map[int32]map[int32]int64, error) {
	// map:palyCard-map:[tingCard-multiple]
	tingMultiple := make(map[int32]map[int32]int64)
	// 获取手牌定缺牌数量
	sum := GetDingQueCardSum(palyer.HandCards, palyer.DingqueColor)
	// 手中少于2张定缺牌才能进行查听
	if sum < 2 {
		// 获取出牌提示
		tingInfo := GetPlayCardCheckTing(palyer.HandCards)
		// 手牌数量
		handCardSum := len(palyer.HandCards)
		handCard := make([]*majongpb.Card, handCardSum)
		// 打那张牌可以胡那些牌，和胡这些牌的倍数
		for playCard, tingCards := range tingInfo {
			// util.card转麻将牌
			playCard2, err := IntToCard(int32(playCard))
			if err != nil {
				return tingMultiple, err
			}
			// 复制手牌
			handCard = append(handCard[:0],palyer.HandCards...)
			// 删除出牌
			newHanCard, isSucceed := DeleteCardFromLast(handCard, playCard2)
			if !isSucceed {
				return tingMultiple, fmt.Errorf("获取出牌提示：删除牌失败：")
			}
			// 能胡牌
			for _, tingCard := range tingCards {
				// util.card转麻将牌
				tingCard2, err := IntToCard(int32(tingCard))
				if err != nil {
					return tingMultiple, err
				}
				// 听的定缺牌不用查倍数
				if tingCard2.Color != palyer.DingqueColor {
					// 添加能胡的牌
					newHanCard = append(newHanCard, tingCard2)
					// 查询能胡的最大倍数 TODO
					multiple := int64(1)
					huMutipleMap := map[int32]int64{int32(tingCard):multiple}
					tingMultiple[int32(playCard)] = huMutipleMap
					// 删除能胡的牌
					newHanCard = newHanCard[:len(newHanCard)-1]
				}
			}
		}
	}
	return tingMultiple, nil
}

//GetHuHint 胡牌提示倍数，缺一张，返回map[int32]int64, error
func GetHuHint(palyer *majongpb.Player) (map[int32]int64, error) {
	// map:tingCard-multiple
	tingMultiple := make(map[int32]int64)
	// 手中没有定缺牌
	if !CheckHasDingQueCard(palyer.HandCards, palyer.DingqueColor) {
		// 获取出牌提示
		tingInfo, err := GetTingCards(palyer.HandCards)
		if err != nil {
			return tingMultiple, err
		}
		// 手牌数量
		handCardSum := len(palyer.HandCards)
		handCard := make([]*majongpb.Card, handCardSum)
		// 复制手牌
		copy(handCard, palyer.HandCards)
		// 可以胡那些牌
		for _, tingCard := range tingInfo {
			// 听的定缺牌不用查倍数
			if tingCard.Color != palyer.DingqueColor {
				// 添加能胡的牌
				newHandCard := append(handCard, tingCard)
				// 查询能胡的最大倍数 TODO
				// 倍数
				multiple := int64(1)
				// 麻将牌转Int32
				cardInt, err := CardToInt(*tingCard)
				if err != nil {
					return tingMultiple, err
				}
				tingMultiple[*cardInt] = multiple
				// 删除能胡的牌
				newHandCard = newHandCard[:len(newHandCard)-1]
			}
		}
	}
	return tingMultiple, nil
}
