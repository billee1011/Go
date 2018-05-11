package settle

import (
	"fmt"
	"steve/majong/states"
	"steve/majong/utils"
	majongpb "steve/server_pb/majong"
	"testing"

	"github.com/stretchr/testify/assert"
)

//  0个花猪的情况,不进行查花猪
func TestCheckFlowerPigSettle(t *testing.T) {
	for i := 0; i < len(players); i++ {
		player := players[i]
		player.HandCards, _ = utils.CheckHuUtilCardsToHandCards(card1util)
		player.HuCards = []*majongpb.HuCard{}
		player.OutCards = []*majongpb.Card{
			&majongpb.Card{
				Color: majongpb.CardColor_ColorWan,
				Point: 1,
			},
		}
		player.DingqueColor = majongpb.CardColor_ColorWan
	}
	flowPigSettle := new(FlowerPigSettle)
	context := &majongpb.MajongContext{
		GameId:  1,
		Players: players,
	}
	settleInfos := flowPigSettle.SettleFlowerPig(context)
	assert.Equal(t, len(settleInfos), 0)
	fmt.Println(settleInfos)
}

//  1个花猪的情况,1个听，1个未听，1个胡
func TestCheckFlowerPigSettleB(t *testing.T) {
	player3.HuCards = append(player3.HuCards, &majongpb.HuCard{Card: &states.Card1B, SrcPlayer: uint64(0), Type: majongpb.HuType_hu_dianpao})
	flowPigSettle := new(FlowerPigSettle)
	context := &majongpb.MajongContext{
		GameId:  1,
		Players: players,
	}
	settleInfos := flowPigSettle.SettleFlowerPig(context)
	assert.Equal(t, len(settleInfos), 1)
	assert.Equal(t, settleInfos[0].Id, uint64(0))
	fmt.Println(settleInfos)
}
