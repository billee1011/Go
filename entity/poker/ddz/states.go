package ddz

import (
	"strconv"
)

// StateID 状态 ID
type StateID int32

const (
	StateID_state_init    StateID = 0
	StateID_state_deal    StateID = 1
	StateID_state_grab    StateID = 2
	StateID_state_double  StateID = 3
	StateID_state_playing StateID = 4
	StateID_state_settle  StateID = 5
	StateID_state_over    StateID = 6
)

var StateID_name = map[int32]string{
	0: "state_init",
	1: "state_deal",
	2: "state_grab",
	3: "state_double",
	4: "state_playing",
	5: "state_settle",
	6: "state_over",
}

func (x StateID) String() string {
	s, ok := StateID_name[int32(x)]
	if ok {
		return s
	}
	return strconv.Itoa(int(x))
}
