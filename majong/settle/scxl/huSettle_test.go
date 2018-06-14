package scxl

import (
	"fmt"
	"steve/majong/interfaces"
	majongpb "steve/server_pb/majong"
	"testing"
)

//TestDianPaoHuSettle 校验自摸结算
func TestDianPaoHuSettle(t *testing.T) {
	// 自摸
	huParmas := interfaces.HuSettleParams{
		HuPlayers:     []uint64{1, 2},
		SrcPlayer:     uint64(0),
		AllPlayers:    []uint64{0, 1, 2, 3},
		CardTypeValue: int(2),
		SettleType:    majongpb.SettleType_settle_dianpao,
	}
	huSettle := new(HuSettle)
	settleInfos := huSettle.Settle(huParmas)
	fmt.Println(settleInfos)

}

//TestZiMoHuSettle 校验点炮结算
func TestZiMoHuSettle(t *testing.T) {
	// 点炮
	huParmas := interfaces.HuSettleParams{
		HuPlayers:  []uint64{1, 2},
		AllPlayers: []uint64{0, 1, 2, 3},
		//CardTypeValue: int(2),
		SettleType: majongpb.SettleType_settle_zimo,
	}
	huSettle := new(HuSettle)
	settleInfos := huSettle.Settle(huParmas)
	fmt.Println(settleInfos)
}
