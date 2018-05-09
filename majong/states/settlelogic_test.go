package states

import (
	"fmt"
	"steve/majong/utils"
	majongpb "steve/server_pb/majong"
	"testing"

	"github.com/stretchr/testify/assert"
)

var players = make([]*majongpb.Player, 4)

// 初始化玩家,默认玩家1是花猪，玩家2是未听玩家，玩家3是听玩家，玩家4是胡玩家
func init() {
	prop := map[string][]byte{utils.IsOutNoDingQueColorCard: []byte{1}}
	// 1W
	card1W := &majongpb.Card{Color: majongpb.CardColor_ColorWan, Point: 1}
	mG := &majongpb.GangCard{
		Card:      card1W,
		Type:      majongpb.GangType_gang_minggang,
		SrcPlayer: uint64(2),
	}
	// 明杠2
	gangCards1 := []*majongpb.GangCard{mG}

	// 1条
	card1T := &majongpb.Card{Color: majongpb.CardColor_ColorTiao, Point: 1}
	aG := &majongpb.GangCard{
		Card: card1T,
		Type: majongpb.GangType_gang_angang,
	}
	//暗杠
	gangCards2 := []*majongpb.GangCard{aG}

	// 2筒
	card2B := &majongpb.Card{Color: majongpb.CardColor_ColorTong, Point: 2}
	bG := &majongpb.GangCard{
		Card: card2B,
		Type: majongpb.GangType_gang_bugang,
	}
	//补杠
	gangCards4 := []*majongpb.GangCard{bG}

	// 胡牌
	huCards := make([]*majongpb.HuCard, 0)
	// 玩家1手牌
	card1 := []utils.Card{11, 11, 12, 12, 13, 13, 14, 14, 15, 15, 16, 17, 18}
	mjCard1, _ := utils.CheckHuCardsToHandCards(card1)
	player1 := &majongpb.Player{
		PalyerId:     uint64(1),
		HandCards:    mjCard1,
		DingqueColor: majongpb.CardColor_ColorWan,
		HuCards:      huCards,
		GangCards:    gangCards1,
		Properties:   prop,
	}
	players[0] = player1

	// 玩家2手牌
	card2 := []utils.Card{21, 21, 22, 22, 23, 23, 24, 24, 25, 25, 19, 26, 29}
	mjCard2, _ := utils.CheckHuCardsToHandCards(card2)
	player2 := &majongpb.Player{
		PalyerId:     uint64(2),
		HandCards:    mjCard2,
		DingqueColor: majongpb.CardColor_ColorTong,
		HuCards:      huCards,
		GangCards:    gangCards2,
		Properties:   prop,
	}
	players[1] = player2

	// 玩家3手牌
	card3 := []utils.Card{11, 11, 12, 12, 13, 13, 14, 14, 15, 15, 16, 16, 17}
	mjCard3, _ := utils.CheckHuCardsToHandCards(card3)
	player3 := &majongpb.Player{
		PalyerId:     uint64(3),
		HandCards:    mjCard3,
		DingqueColor: majongpb.CardColor_ColorTong,
		HuCards:      huCards,
		Properties:   prop,
	}
	players[2] = player3

	// 1筒
	card1B := &majongpb.Card{Color: majongpb.CardColor_ColorTong, Point: 1}
	//玩家4胡的牌
	huCard := &majongpb.HuCard{
		SrcPlayer: uint64(1),
		Card:      card1B,
		Type:      majongpb.HuType_hu_dianpao,
	}
	huCards = append(huCards, huCard)
	// 玩家4手牌
	card4 := []utils.Card{31, 31, 32, 32, 33, 33, 34, 34, 35, 35, 36, 37, 38}
	mjCard4, _ := utils.CheckHuCardsToHandCards(card4)
	player4 := &majongpb.Player{
		PalyerId:     uint64(4),
		HandCards:    mjCard4,
		DingqueColor: majongpb.CardColor_ColorTiao,
		HuCards:      huCards,
		GangCards:    gangCards4,
		Properties:   prop,
	}
	players[3] = player4
}

//  0个花猪的情况,不进行花猪
func TestCheckFlowerPigSettle(t *testing.T) {
	card := []utils.Card{21, 21, 22, 22, 23, 23, 24, 24, 25, 25, 19, 26, 29}
	mjCard, _ := utils.CheckHuCardsToHandCards(card)

	players[0].DingqueColor = majongpb.CardColor_ColorTiao
	for i := 0; i < len(players); i++ {
		player := players[i]
		player.HandCards = mjCard
		player.HuCards = []*majongpb.HuCard{}
		player.DingqueColor = majongpb.CardColor_ColorTong
	}
	settleInfos, err := CheckFlowerPigSettle(players)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(settleInfos)
}

//  4个花猪的情况,不进行花猪
func TestCheckFlowerPigSettleA(t *testing.T) {
	for i := 0; i < len(players); i++ {
		player := players[i]
		player.DingqueColor = player.HandCards[0].Color
		player.HuCards = []*majongpb.HuCard{}
	}
	settleInfos, err := CheckFlowerPigSettle(players)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(settleInfos)
}

//  3个花猪的情况,一个胡玩家
func TestCheckFlowerPigSettleB(t *testing.T) {
	for i := 0; i < len(players); i++ {
		player := players[i]
		if len(player.HuCards) == 0 {
			player.DingqueColor = player.HandCards[0].Color
		}
	}
	settleInfos, err := CheckFlowerPigSettle(players)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(settleInfos)
}

//  3个花猪的情况,一个听玩家
func TestCheckFlowerPigSettleC(t *testing.T) {
	for i := 0; i < len(players); i++ {
		player := players[i]
		if len(player.HuCards) > 0 {
			player.HuCards = []*majongpb.HuCard{}
		} else {
			player.DingqueColor = player.HandCards[0].Color
		}
	}
	settleInfos, err := CheckFlowerPigSettle(players)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(settleInfos)
}

//  3个花猪的情况,一个未听玩家
func TestCheckFlowerPigSettleD(t *testing.T) {
	for i := 0; i < len(players); i++ {
		player := players[i]
		player.HuCards = []*majongpb.HuCard{}
		if player.PalyerId != 2 {
			player.DingqueColor = player.HandCards[0].Color
		}
	}
	settleInfos, err := CheckFlowerPigSettle(players)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(settleInfos)
}

//  2个花猪的情况,一个未听玩家,一个听玩家
func TestCheckFlowerPigSettlE(t *testing.T) {
	for i := 0; i < len(players); i++ {
		player := players[i]
		player.HuCards = []*majongpb.HuCard{}
		if player.PalyerId != 2 && player.PalyerId != 4 {
			player.DingqueColor = player.HandCards[0].Color
		}
	}
	settleInfos, err := CheckFlowerPigSettle(players)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(settleInfos)
}

//  2个花猪的情况,一个未听玩家,一个胡玩家
func TestCheckFlowerPigSettlF(t *testing.T) {
	for i := 0; i < len(players); i++ {
		player := players[i]
		if player.PalyerId != 2 && player.PalyerId != 4 {
			player.DingqueColor = player.HandCards[0].Color
		}
		if player.PalyerId != 4 {
			player.HuCards = []*majongpb.HuCard{}
		}
	}
	settleInfos, err := CheckFlowerPigSettle(players)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(settleInfos)
}

//  2个花猪的情况,一个听玩家,一个胡玩家
func TestCheckFlowerPigSettlG(t *testing.T) {
	for i := 0; i < len(players); i++ {
		player := players[i]
		if player.PalyerId != 3 && player.PalyerId != 4 {
			player.DingqueColor = player.HandCards[0].Color
		}
		if player.PalyerId != 4 {
			player.HuCards = []*majongpb.HuCard{}
		}
	}
	settleInfos, err := CheckFlowerPigSettle(players)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(settleInfos)
}

// 1个花猪的情况,一个未听玩家,一个听玩家,一个胡玩家
func TestCheckFlowerPigSettlH(t *testing.T) {
	settleInfos, err := CheckFlowerPigSettle(players)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(settleInfos)
}

// 4个未听玩家,不进行查大叫
func TestCheckYellSettle1A(t *testing.T) {
	card := []utils.Card{21, 21, 22, 22, 23, 23, 24, 24, 25, 25, 19, 26, 29}
	mjCard, _ := utils.CheckHuCardsToHandCards(card)

	players[0].DingqueColor = majongpb.CardColor_ColorTiao
	for i := 0; i < len(players); i++ {
		player := players[i]
		player.HandCards = mjCard
		player.HuCards = []*majongpb.HuCard{}
		player.DingqueColor = majongpb.CardColor_ColorTong
	}
	settleInfos, err := CheckYellSettle(players)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(settleInfos)
}

// 3个未听玩家,1个胡玩家,不进行查大叫
func TestCheckYellSettle1B(t *testing.T) {
	card := []utils.Card{21, 21, 22, 22, 23, 23, 24, 24, 25, 25, 19, 26, 29}
	mjCard, _ := utils.CheckHuCardsToHandCards(card)

	players[0].DingqueColor = majongpb.CardColor_ColorTiao
	for i := 0; i < len(players); i++ {
		player := players[i]
		if player.PalyerId == 4 {
			continue
		}
		player.HandCards = mjCard
		player.HuCards = []*majongpb.HuCard{}
		player.DingqueColor = majongpb.CardColor_ColorTong
	}
	settleInfos, err := CheckYellSettle(players)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(settleInfos)
}

// 3个未听玩家,1个听玩家
func TestCheckYellSettle1C(t *testing.T) {
	card := []utils.Card{21, 21, 22, 22, 23, 23, 24, 24, 25, 25, 19, 26, 29}
	mjCard, _ := utils.CheckHuCardsToHandCards(card)
	players[0].DingqueColor = majongpb.CardColor_ColorTiao
	for i := 0; i < len(players); i++ {
		player := players[i]
		if player.PalyerId == 3 {
			continue
		}
		player.HandCards = mjCard
		player.HuCards = []*majongpb.HuCard{}
		player.DingqueColor = majongpb.CardColor_ColorTong
	}
	settleInfos, err := CheckYellSettle(players)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(settleInfos)
}

// 2个未听玩家,1个听玩家,1个胡玩家
func TestCheckYellSettle1D(t *testing.T) {
	card := []utils.Card{21, 21, 22, 22, 23, 23, 24, 24, 25, 25, 19, 26, 29}
	mjCard, _ := utils.CheckHuCardsToHandCards(card)
	players[0].DingqueColor = majongpb.CardColor_ColorTiao
	for i := 0; i < len(players); i++ {
		player := players[i]
		if player.PalyerId == 3 || player.PalyerId == 4 {
			continue
		}
		player.HandCards = mjCard
		player.HuCards = []*majongpb.HuCard{}
		player.DingqueColor = majongpb.CardColor_ColorTong
	}
	settleInfos, err := CheckYellSettle(players)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(settleInfos)
}

// 1个未听玩家,2个听玩家,1个胡玩家
func TestCheckYellSettle1E(t *testing.T) { // 1个未听玩家,1个花猪,1个听玩家,1个胡玩家
	card := []utils.Card{21, 21, 22, 22, 23, 23, 24, 24, 25, 25, 19, 26, 29}
	mjCard, _ := utils.CheckHuCardsToHandCards(card)
	players[0].DingqueColor = majongpb.CardColor_ColorTiao
	for i := 0; i < len(players); i++ {
		player := players[i]
		if player.PalyerId == 3 || player.PalyerId == 4 {
			continue
		}
		if player.PalyerId == 1 {
			player.DingqueColor = majongpb.CardColor_ColorTong
			continue
		}
		player.HandCards = mjCard
		player.HuCards = []*majongpb.HuCard{}
		player.DingqueColor = majongpb.CardColor_ColorTong
	}
	settleInfos, err := CheckYellSettle(players)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(settleInfos)
}

// 1个未听玩家,1个花猪,1个听玩家,1个胡玩家
func TestFlowerAndYell(t *testing.T) {
	settleInfos2, err := CheckFlowerPigSettle(players)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(settleInfos2)
	settleInfos, err := CheckYellSettle(players)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(settleInfos)
}

// 未听玩家暗杠，未听玩家补杠，花猪明杠
func TestTaxRebateSettleA(t *testing.T) {
	card2 := []utils.Card{21, 21, 22, 22, 23, 23, 24, 24, 25, 25, 19, 26, 29}
	mjCard2, _ := utils.CheckHuCardsToHandCards(card2)
	players[3].HuCards = players[3].HuCards[:0]
	players[3].HandCards = mjCard2
	players[3].DingqueColor = majongpb.CardColor_ColorTong
	settleInfos, err := TaxRebateSettle(players)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(settleInfos)
}

// 花猪玩家明杠，暗杠，补杠,未听玩家没有杠
func TestTaxRebateSettleB(t *testing.T) {
	for i := 0; i < len(players); i++ {
		players[i].GangCards = players[i].GangCards[:0]
	}
	// 1W
	mG := &majongpb.GangCard{
		Card:      &majongpb.Card{Color: majongpb.CardColor_ColorWan, Point: 1},
		Type:      majongpb.GangType_gang_minggang,
		SrcPlayer: uint64(2),
	}

	// 1条
	aG := &majongpb.GangCard{
		Card: &majongpb.Card{Color: majongpb.CardColor_ColorTiao, Point: 1},
		Type: majongpb.GangType_gang_angang,
	}

	// 2筒
	bG := &majongpb.GangCard{
		Card: &majongpb.Card{Color: majongpb.CardColor_ColorTong, Point: 2},
		Type: majongpb.GangType_gang_bugang,
	}
	players[0].GangCards = append(players[0].GangCards, mG, aG, bG)
	settleInfos, err := TaxRebateSettle(players)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(settleInfos)
}

// 花猪玩家和未听玩家,胡家，听家都没有杠
func TestTaxRebateSettleC(t *testing.T) {
	for i := 0; i < len(players); i++ {
		players[i].GangCards = players[i].GangCards[:0]
	}
	settleInfos, err := TaxRebateSettle(players)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(settleInfos)
}

// 花猪玩家和未听玩家，胡家没有杠，听家有杠
func TestTaxRebateSettleD(t *testing.T) {
	for i := 0; i < len(players); i++ {
		players[i].GangCards = players[i].GangCards[:0]
	}
	// 1W
	mG := &majongpb.GangCard{
		Card:      &majongpb.Card{Color: majongpb.CardColor_ColorWan, Point: 1},
		Type:      majongpb.GangType_gang_minggang,
		SrcPlayer: uint64(2),
	}

	// 1条
	aG := &majongpb.GangCard{
		Card: &majongpb.Card{Color: majongpb.CardColor_ColorTiao, Point: 1},
		Type: majongpb.GangType_gang_angang,
	}

	// 2筒
	bG := &majongpb.GangCard{
		Card: &majongpb.Card{Color: majongpb.CardColor_ColorTong, Point: 2},
		Type: majongpb.GangType_gang_bugang,
	}
	players[2].GangCards = append(players[2].GangCards, mG, aG, bG)
	settleInfos, err := TaxRebateSettle(players)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(settleInfos)
}

// 花猪玩家和未听玩家，听家没有杠，胡家有杠
func TestTaxRebateSettleE(t *testing.T) {
	for i := 0; i < len(players); i++ {
		players[i].GangCards = players[i].GangCards[:0]
	}
	// 1W
	mG := &majongpb.GangCard{
		Card:      &majongpb.Card{Color: majongpb.CardColor_ColorWan, Point: 1},
		Type:      majongpb.GangType_gang_minggang,
		SrcPlayer: uint64(2),
	}

	// 1条
	aG := &majongpb.GangCard{
		Card: &majongpb.Card{Color: majongpb.CardColor_ColorTiao, Point: 1},
		Type: majongpb.GangType_gang_angang,
	}

	// 2筒
	bG := &majongpb.GangCard{
		Card: &majongpb.Card{Color: majongpb.CardColor_ColorTong, Point: 2},
		Type: majongpb.GangType_gang_bugang,
	}
	players[3].GangCards = append(players[3].GangCards, mG, aG, bG)
	settleInfos, err := TaxRebateSettle(players)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(settleInfos)
}

// 未听玩家明杠，暗杠，补杠,花猪玩家没有杠
func TestTaxRebateSettleF(t *testing.T) {
	for i := 0; i < len(players); i++ {
		players[i].GangCards = players[i].GangCards[:0]
	}
	// 1W
	mG := &majongpb.GangCard{
		Card:      &majongpb.Card{Color: majongpb.CardColor_ColorWan, Point: 1},
		Type:      majongpb.GangType_gang_minggang,
		SrcPlayer: uint64(1),
	}

	// 1条
	aG := &majongpb.GangCard{
		Card: &majongpb.Card{Color: majongpb.CardColor_ColorTiao, Point: 1},
		Type: majongpb.GangType_gang_angang,
	}

	// 2筒
	bG := &majongpb.GangCard{
		Card: &majongpb.Card{Color: majongpb.CardColor_ColorTong, Point: 2},
		Type: majongpb.GangType_gang_bugang,
	}
	players[1].GangCards = append(players[1].GangCards, mG, aG, bG)
	settleInfos, err := TaxRebateSettle(players)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(settleInfos)
}

//呼叫转移,明杠，一人胡
func TestCallDivertSettle(t *testing.T) {
	// 1W
	mG := &majongpb.GangCard{
		Card:      &majongpb.Card{Color: majongpb.CardColor_ColorWan, Point: 1},
		Type:      majongpb.GangType_gang_minggang,
		SrcPlayer: uint64(2),
	}
	players[1].GangCards = []*majongpb.GangCard{mG}
	winPlayers := []*majongpb.Player{players[0]}
	settleInfos, err := CallDivertSettle(majongpb.HuType_hu_ganghoupao, players, winPlayers, players[1])
	assert.Nil(t, err)

	settleInfoMap := make(map[uint64]int64)
	settleInfoMap[players[0].PalyerId] = 2
	settleInfoMap[players[1].PalyerId] = -2
	yellSettleInfo := &majongpb.SettleInfo{
		Id:     players[1].PalyerId,
		Scores: settleInfoMap,
	}
	settleInfos2 := []*majongpb.SettleInfo{yellSettleInfo}
	assert.Equal(t, len(settleInfos), len(settleInfos2))
	for k, settleInfo := range settleInfos {
		assert.Equal(t, settleInfo.Id, settleInfos2[k].Id)
		for ID, Score := range settleInfos[k].Scores {
			assert.Equal(t, Score, settleInfos2[k].Scores[ID])
		}
	}
	fmt.Println(settleInfos2)
	fmt.Println(settleInfos)
}

//呼叫转移,明杠，2人胡,包含点明人玩家
func TestCallDivertSettleB(t *testing.T) {
	// 1W
	mG := &majongpb.GangCard{
		Card:      &majongpb.Card{Color: majongpb.CardColor_ColorWan, Point: 1},
		Type:      majongpb.GangType_gang_minggang,
		SrcPlayer: uint64(1),
	}
	players[1].GangCards = []*majongpb.GangCard{mG}
	winPlayers := []*majongpb.Player{players[0], players[2]}
	settleInfos, err := CallDivertSettle(majongpb.HuType_hu_ganghoupao, players, winPlayers, players[1])
	assert.Nil(t, err)

	settleInfoMap := make(map[uint64]int64)
	settleInfoMap[players[0].PalyerId] = 2
	settleInfoMap[players[1].PalyerId] = -2
	yellSettleInfo := &majongpb.SettleInfo{
		Id:     players[1].PalyerId,
		Scores: settleInfoMap,
	}
	settleInfos2 := append(settleInfos[:0], yellSettleInfo)
	assert.Equal(t, len(settleInfos), len(settleInfos2))
	for k, settleInfo := range settleInfos {
		assert.Equal(t, settleInfo.Id, settleInfos2[k].Id)
		for ID, Score := range settleInfos[k].Scores {
			assert.Equal(t, Score, settleInfos2[k].Scores[ID])
		}
	}
	fmt.Println(settleInfos2)
	fmt.Println(settleInfos)
}

//呼叫转移,明杠，2人胡,不包含点明人玩家
func TestCallDivertSettleC(t *testing.T) {
	// 1W
	mG := &majongpb.GangCard{
		Card:      &majongpb.Card{Color: majongpb.CardColor_ColorWan, Point: 1},
		Type:      majongpb.GangType_gang_minggang,
		SrcPlayer: uint64(1),
	}
	players[1].GangCards = []*majongpb.GangCard{mG}
	winPlayers := []*majongpb.Player{players[3], players[2]}
	settleInfos, err := CallDivertSettle(majongpb.HuType_hu_ganghoupao, players, winPlayers, players[1])
	assert.Nil(t, err)

	settleInfoMap := make(map[uint64]int64)
	settleInfoMap[players[2].PalyerId] = 1
	settleInfoMap[players[3].PalyerId] = 1
	settleInfoMap[players[1].PalyerId] = -2
	yellSettleInfo := &majongpb.SettleInfo{
		Id:     players[1].PalyerId,
		Scores: settleInfoMap,
	}
	settleInfos2 := append(settleInfos[:0], yellSettleInfo)
	assert.Equal(t, len(settleInfos), len(settleInfos2))
	for k, settleInfo := range settleInfos {
		assert.Equal(t, settleInfo.Id, settleInfos2[k].Id)
		for ID, Score := range settleInfos[k].Scores {
			assert.Equal(t, Score, settleInfos2[k].Scores[ID])
		}
	}
	fmt.Println(settleInfos2)
	fmt.Println(settleInfos)
}

//呼叫转移,明杠，3人胡,包含点明人玩家
func TestCallDivertSettleD(t *testing.T) {
	// 1W
	mG := &majongpb.GangCard{
		Card:      &majongpb.Card{Color: majongpb.CardColor_ColorWan, Point: 1},
		Type:      majongpb.GangType_gang_minggang,
		SrcPlayer: uint64(1),
	}
	players[1].GangCards = []*majongpb.GangCard{mG}
	winPlayers := []*majongpb.Player{players[0], players[2],players[3]}
	settleInfos, err := CallDivertSettle(majongpb.HuType_hu_ganghoupao, players, winPlayers, players[1])
	assert.Nil(t, err)

	settleInfoMap := make(map[uint64]int64)
	settleInfoMap[players[0].PalyerId] = 2
	settleInfoMap[players[1].PalyerId] = -2
	yellSettleInfo := &majongpb.SettleInfo{
		Id:     players[1].PalyerId,
		Scores: settleInfoMap,
	}
	settleInfos2 := append(settleInfos[:0], yellSettleInfo)
	assert.Equal(t, len(settleInfos), len(settleInfos2))
	for k, settleInfo := range settleInfos {
		assert.Equal(t, settleInfo.Id, settleInfos2[k].Id)
		for ID, Score := range settleInfos[k].Scores {
			assert.Equal(t, Score, settleInfos2[k].Scores[ID])
		}
	}
	fmt.Println(settleInfos2)
	fmt.Println(settleInfos)
}

//呼叫转移,补杠，1人胡
func TestCallDivertSettleE(t *testing.T) {
	// 1W
	mG := &majongpb.GangCard{
		Card:      &majongpb.Card{Color: majongpb.CardColor_ColorWan, Point: 1},
		Type:      majongpb.GangType_gang_bugang,
	}
	players[1].GangCards = []*majongpb.GangCard{mG}
	winPlayers := []*majongpb.Player{players[0]}
	settleInfos, err := CallDivertSettle(majongpb.HuType_hu_ganghoupao, players, winPlayers, players[1])
	assert.Nil(t, err)

	settleInfoMap := make(map[uint64]int64)
	settleInfoMap[players[0].PalyerId] = 3
	settleInfoMap[players[1].PalyerId] = -3
	yellSettleInfo := &majongpb.SettleInfo{
		Id:     players[1].PalyerId,
		Scores: settleInfoMap,
	}
	settleInfos2 := append(settleInfos[:0], yellSettleInfo)
	assert.Equal(t, len(settleInfos), len(settleInfos2))
	for k, settleInfo := range settleInfos {
		assert.Equal(t, settleInfo.Id, settleInfos2[k].Id)
		for ID, Score := range settleInfos[k].Scores {
			assert.Equal(t, Score, settleInfos2[k].Scores[ID])
		}
	}
	fmt.Println(settleInfos2)
	fmt.Println(settleInfos)
}

//呼叫转移,补杠，2人胡,第一个胡玩家多一分
func TestCallDivertSettleF(t *testing.T) {
	// 1W
	mG := &majongpb.GangCard{
		Card:      &majongpb.Card{Color: majongpb.CardColor_ColorWan, Point: 1},
		Type:      majongpb.GangType_gang_bugang,
	}
	players[1].GangCards = []*majongpb.GangCard{mG}
	winPlayers := []*majongpb.Player{players[0],players[2]}
	settleInfos, err := CallDivertSettle(majongpb.HuType_hu_ganghoupao, players, winPlayers, players[1])
	assert.Nil(t, err)

	settleInfoMap := make(map[uint64]int64)
	settleInfoMap[players[2].PalyerId] = 2
	settleInfoMap[players[0].PalyerId] = 1
	settleInfoMap[players[1].PalyerId] = -3
	yellSettleInfo := &majongpb.SettleInfo{
		Id:     players[1].PalyerId,
		Scores: settleInfoMap,
	}
	settleInfos2 := append(settleInfos[:0], yellSettleInfo)
	assert.Equal(t, len(settleInfos), len(settleInfos2))
	for k, settleInfo := range settleInfos {
		assert.Equal(t, settleInfo.Id, settleInfos2[k].Id)
		for ID, Score := range settleInfos[k].Scores {
			assert.Equal(t, Score, settleInfos2[k].Scores[ID])
		}
	}
	fmt.Println(settleInfos2)
	fmt.Println(settleInfos)
}

//呼叫转移,补杠，3人胡,第一个胡玩家多一分
func TestCallDivertSettleG(t *testing.T) {
	// 1W
	mG := &majongpb.GangCard{
		Card:      &majongpb.Card{Color: majongpb.CardColor_ColorWan, Point: 1},
		Type:      majongpb.GangType_gang_bugang,
	}
	players[1].GangCards = []*majongpb.GangCard{mG}
	winPlayers := []*majongpb.Player{players[0],players[2],players[3]}
	settleInfos, err := CallDivertSettle(majongpb.HuType_hu_ganghoupao, players, winPlayers, players[1])
	assert.Nil(t, err)

	settleInfoMap := make(map[uint64]int64)
	settleInfoMap[players[3].PalyerId] = 1
	settleInfoMap[players[2].PalyerId] = 1
	settleInfoMap[players[0].PalyerId] = 1
	settleInfoMap[players[1].PalyerId] = -3
	yellSettleInfo := &majongpb.SettleInfo{
		Id:     players[1].PalyerId,
		Scores: settleInfoMap,
	}
	settleInfos2 := append(settleInfos[:0], yellSettleInfo)
	assert.Equal(t, len(settleInfos), len(settleInfos2))
	for k, settleInfo := range settleInfos {
		assert.Equal(t, settleInfo.Id, settleInfos2[k].Id)
		for ID, Score := range settleInfos[k].Scores {
			assert.Equal(t, Score, settleInfos2[k].Scores[ID])
		}
	}
	fmt.Println(settleInfos2)
	fmt.Println(settleInfos)
}


//呼叫转移,暗杠，1人胡
func TestCallDivertSettleH(t *testing.T) {
	// 1W
	mG := &majongpb.GangCard{
		Card:      &majongpb.Card{Color: majongpb.CardColor_ColorWan, Point: 1},
		Type:      majongpb.GangType_gang_bugang,
		SrcPlayer: uint64(1),
	}
	players[1].GangCards = []*majongpb.GangCard{mG}
	winPlayers := []*majongpb.Player{players[0]}
	settleInfos, err := CallDivertSettle(majongpb.HuType_hu_ganghoupao, players, winPlayers, players[1])
	assert.Nil(t, err)

	settleInfoMap := make(map[uint64]int64)
	settleInfoMap[players[0].PalyerId] = 6
	settleInfoMap[players[1].PalyerId] = -6
	yellSettleInfo := &majongpb.SettleInfo{
		Id:     players[1].PalyerId,
		Scores: settleInfoMap,
	}
	settleInfos2 := append(settleInfos[:0], yellSettleInfo)
	assert.Equal(t, len(settleInfos), len(settleInfos2))
	for k, settleInfo := range settleInfos {
		assert.Equal(t, settleInfo.Id, settleInfos2[k].Id)
		for ID, Score := range settleInfos[k].Scores {
			assert.Equal(t, Score, settleInfos2[k].Scores[ID])
		}
	}
	fmt.Println(settleInfos2)
	fmt.Println(settleInfos)
}


//呼叫转移,暗杠，2人胡
func TestCallDivertSettleI(t *testing.T) {
	// 1W
	mG := &majongpb.GangCard{
		Card:      &majongpb.Card{Color: majongpb.CardColor_ColorWan, Point: 1},
		Type:      majongpb.GangType_gang_bugang,
		SrcPlayer: uint64(1),
	}
	players[1].GangCards = []*majongpb.GangCard{mG}
	winPlayers := []*majongpb.Player{players[0],players[2]}
	settleInfos, err := CallDivertSettle(majongpb.HuType_hu_ganghoupao, players, winPlayers, players[1])
	assert.Nil(t, err)

	settleInfoMap := make(map[uint64]int64)
	settleInfoMap[players[2].PalyerId] = 3
	settleInfoMap[players[0].PalyerId] = 3
	settleInfoMap[players[1].PalyerId] = -6
	yellSettleInfo := &majongpb.SettleInfo{
		Id:     players[1].PalyerId,
		Scores: settleInfoMap,
	}
	settleInfos2 := append(settleInfos[:0], yellSettleInfo)
	assert.Equal(t, len(settleInfos), len(settleInfos2))
	for k, settleInfo := range settleInfos {
		assert.Equal(t, settleInfo.Id, settleInfos2[k].Id)
		for ID, Score := range settleInfos[k].Scores {
			assert.Equal(t, Score, settleInfos2[k].Scores[ID])
		}
	}
	fmt.Println(settleInfos2)
	fmt.Println(settleInfos)
}

//呼叫转移,暗杠，3人胡
func TestCallDivertSettleJ(t *testing.T) {
	// 1W
	mG := &majongpb.GangCard{
		Card:      &majongpb.Card{Color: majongpb.CardColor_ColorWan, Point: 1},
		Type:      majongpb.GangType_gang_bugang,
		SrcPlayer: uint64(1),
	}
	players[1].GangCards = []*majongpb.GangCard{mG}
	winPlayers := []*majongpb.Player{players[0],players[2],players[3]}
	settleInfos, err := CallDivertSettle(majongpb.HuType_hu_ganghoupao, players, winPlayers, players[1])
	assert.Nil(t, err)

	settleInfoMap := make(map[uint64]int64)
	settleInfoMap[players[3].PalyerId] = 2
	settleInfoMap[players[2].PalyerId] = 2
	settleInfoMap[players[0].PalyerId] = 2
	settleInfoMap[players[1].PalyerId] = -6
	yellSettleInfo := &majongpb.SettleInfo{
		Id:     players[1].PalyerId,
		Scores: settleInfoMap,
	}
	settleInfos2 := append(settleInfos[:0], yellSettleInfo)
	assert.Equal(t, len(settleInfos), len(settleInfos2))
	for k, settleInfo := range settleInfos {
		assert.Equal(t, settleInfo.Id, settleInfos2[k].Id)
		for ID, Score := range settleInfos[k].Scores {
			assert.Equal(t, Score, settleInfos2[k].Scores[ID])
		}
	}
	fmt.Println(settleInfos2)
	fmt.Println(settleInfos)
}