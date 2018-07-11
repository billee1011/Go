package gutils

import (
	"steve/client_pb/room"
)

const (
	// SCXLGameID 四川血流
	SCXLGameID = 1
	// SCXZGameID 四川血战
	SCXZGameID = 2
	// DDZGameID 斗地主
	DDZGameID = 3
	// ERMJGameID 二人麻将
	ERMJGameID = 4
)

func GameIDServer2Client(sGameID int) (cGameID room.GameId) {
	switch sGameID {
	case SCXLGameID:
		cGameID = room.GameId_GAMEID_XUELIU
	case SCXZGameID:
		cGameID = room.GameId_GAMEID_XUEZHAN
	case DDZGameID:
		cGameID = room.GameId_GAMEID_DOUDIZHU
	case ERMJGameID:
		cGameID = room.GameId_GAMEID_ERRENMJ
	}
	return
}
