package scxl

import (
	"fmt"
	"steve/majong/interfaces"
	majongpb "steve/server_pb/majong"
	"testing"
)

//TestRoundSettle 校验单局结算
func TestRoundSettle(t *testing.T) {
	roundSettleParams := interfaces.RoundSettleParams{
		FlowerPigPlayers: []uint64{0},
		HuPlayers:        []uint64{2},
		NotTingPlayers:   []uint64{3},
		TingPlayersInfo: map[uint64]int{
			1: 6,
		},
		SettleInfos: []*majongpb.SettleInfo{
			&majongpb.SettleInfo{
				SettleType: majongpb.SettleType_settle_gang,
				Scores: map[uint64]int64{
					3: 11,
				},
				Id: 11,
			},
			&majongpb.SettleInfo{
				SettleType: majongpb.SettleType_settle_gang,
				Scores: map[uint64]int64{
					0: 11,
				},
				Id: 22,
			},
		},
	}
	roundSettle := new(RoundSettle)
	settleInfos, raxrebats := roundSettle.Settle(roundSettleParams)
	for _, settleInfo := range settleInfos {
		fmt.Println(settleInfo)

	}
	fmt.Println(raxrebats)
}
