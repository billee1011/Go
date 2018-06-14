package scxl

import (
	"fmt"
	"steve/majong/interfaces"
	majongpb "steve/server_pb/majong"

	"testing"
)

//TestAnGangSettle 校验暗杠结算
func TestAnGangSettle(t *testing.T) {
	gangSettleParams := interfaces.GangSettleParams{
		GangPlayer: uint64(0),
		AllPlayers: []uint64{0, 1, 2, 3},
		GangType:   majongpb.GangType_gang_angang,
	}
	gangSettle := new(GangSettle)
	settleInfos := gangSettle.Settle(gangSettleParams)
	fmt.Println(settleInfos)
}

//TestBuGangSettle 校验补杠结算
func TestBuGangSettle(t *testing.T) {
	gangSettleParams := interfaces.GangSettleParams{
		GangPlayer: uint64(0),
		AllPlayers: []uint64{0, 1, 2, 3},
		GangType:   majongpb.GangType_gang_bugang,
	}
	gangSettle := new(GangSettle)
	settleInfos := gangSettle.Settle(gangSettleParams)
	fmt.Println(settleInfos)
}

//TestMingGangSettle 校验明杠结算
func TestMingGangSettle(t *testing.T) {
	gangSettleParams := interfaces.GangSettleParams{
		GangPlayer: uint64(0),
		SrcPlayer:  uint64(1),
		AllPlayers: []uint64{0, 1, 2, 3},
		GangType:   majongpb.GangType_gang_minggang,
	}
	gangSettle := new(GangSettle)
	settleInfos := gangSettle.Settle(gangSettleParams)
	fmt.Println(settleInfos)
}
