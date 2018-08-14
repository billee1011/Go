package majong

import "strconv"

// StateID 状态 ID
type StateID int32

const (
	StateID_state_init               StateID = 0
	StateID_state_xipai              StateID = 1
	StateID_state_fapai              StateID = 2
	StateID_state_huansanzhang       StateID = 3
	StateID_state_dingque            StateID = 4
	StateID_state_chupai             StateID = 5
	StateID_state_angang             StateID = 6
	StateID_state_zimo               StateID = 7
	StateID_state_peng               StateID = 8
	StateID_state_gang               StateID = 9
	StateID_state_hu                 StateID = 10
	StateID_state_mopai              StateID = 11
	StateID_state_zixun              StateID = 12
	StateID_state_bugang             StateID = 13
	StateID_state_waitqiangganghu    StateID = 14
	StateID_state_qiangganghu        StateID = 15
	StateID_state_chupaiwenxun       StateID = 16
	StateID_state_gameover           StateID = 17
	StateID_state_gang_settle        StateID = 18
	StateID_state_zimo_settle        StateID = 19
	StateID_state_hu_settle          StateID = 20
	StateID_state_qiangganghu_settle StateID = 21
	StateID_state_gamestart_buhua    StateID = 22
	StateID_state_xingpai_buhua      StateID = 23
	StateID_state_chi                StateID = 24
)

var StateID_name = map[int32]string{
	0:  "state_init",
	1:  "state_xipai",
	2:  "state_fapai",
	3:  "state_huansanzhang",
	4:  "state_dingque",
	5:  "state_chupai",
	6:  "state_angang",
	7:  "state_zimo",
	8:  "state_peng",
	9:  "state_gang",
	10: "state_hu",
	11: "state_mopai",
	12: "state_zixun",
	13: "state_bugang",
	14: "state_waitqiangganghu",
	15: "state_qiangganghu",
	16: "state_chupaiwenxun",
	17: "state_gameover",
	18: "state_gang_settle",
	19: "state_zimo_settle",
	20: "state_hu_settle",
	21: "state_qiangganghu_settle",
	22: "state_gamestart_buhua",
	23: "state_xingpai_buhua",
	24: "state_chi",
}

var StateID_value = map[string]int32{
	"state_init":               0,
	"state_xipai":              1,
	"state_fapai":              2,
	"state_huansanzhang":       3,
	"state_dingque":            4,
	"state_chupai":             5,
	"state_angang":             6,
	"state_zimo":               7,
	"state_peng":               8,
	"state_gang":               9,
	"state_hu":                 10,
	"state_mopai":              11,
	"state_zixun":              12,
	"state_bugang":             13,
	"state_waitqiangganghu":    14,
	"state_qiangganghu":        15,
	"state_chupaiwenxun":       16,
	"state_gameover":           17,
	"state_gang_settle":        18,
	"state_zimo_settle":        19,
	"state_hu_settle":          20,
	"state_qiangganghu_settle": 21,
	"state_gamestart_buhua":    22,
	"state_xingpai_buhua":      23,
	"state_chi":                24,
}

func (x StateID) String() string {
	s, ok := StateID_name[int32(x)]
	if ok {
		return s
	}
	return strconv.Itoa(int(x))
}
