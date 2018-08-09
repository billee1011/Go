package gameconfig

// PlayerState 玩家状态
type PlayerState struct {
	PlayerID  uint64
	State     uint64
	GameID    uint64
	LevelID   uint64
	IPAddr    string
	GateAddr  string
	MatchAddr string
	RoomAddr  string
}
