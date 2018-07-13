package abs

type Desk interface {
	GetUid() uint64
	GetGameId() int
	Start()
	Stop()
}