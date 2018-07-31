package poker

import (
	"strconv"
)

type CardType int32

const (
	CardType_CT_NONE     CardType = 0
	CardType_CT_SINGLE   CardType = 1
	CardType_CT_PAIR     CardType = 2
	CardType_CT_SHUNZI   CardType = 3
	CardType_CT_PAIRS    CardType = 4
	CardType_CT_TRIPLE   CardType = 5
	CardType_CT_3AND1    CardType = 6
	CardType_CT_3AND2    CardType = 7
	CardType_CT_TRIPLES  CardType = 8
	CardType_CT_3SAND1S  CardType = 9
	CardType_CT_3SAND2S  CardType = 10
	CardType_CT_4SAND1S  CardType = 11
	CardType_CT_4SAND2S  CardType = 12
	CardType_CT_BOMB     CardType = 13
	CardType_CT_KINGBOMB CardType = 14
)

var CardType_name = map[int32]string{
	0:  "CT_NONE",
	1:  "CT_SINGLE",
	2:  "CT_PAIR",
	3:  "CT_SHUNZI",
	4:  "CT_PAIRS",
	5:  "CT_TRIPLE",
	6:  "CT_3AND1",
	7:  "CT_3AND2",
	8:  "CT_TRIPLES",
	9:  "CT_3SAND1S",
	10: "CT_3SAND2S",
	11: "CT_4SAND1S",
	12: "CT_4SAND2S",
	13: "CT_BOMB",
	14: "CT_KINGBOMB",
}

func (x CardType) String() string {
	s, ok := CardType_name[int32(x)]
	if ok {
		return s
	}
	return strconv.Itoa(int(x))
}
