package login

import (
	"steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/room/interfaces/global"
	"steve/structs/exchanger"
	"sync"

	"steve/structs/proto/gate_rpc"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

type player struct {
	playerID uint64
	coin     uint64
	clientID uint64
}

func (p *player) GetID() uint64 {
	return p.playerID
}
func (p *player) GetCoin() uint64 {
	return p.coin
}
func (p *player) GetClientID() uint64 {
	return p.clientID
}
func (p *player) SetCoin(coin uint64) {
	p.coin = coin
}

var maxPlayerID uint64
var maxPlayerIDMutex sync.Mutex

func allocPlayerID() uint64 {
	maxPlayerIDMutex.Lock()
	defer maxPlayerIDMutex.Unlock()
	maxPlayerID++
	return maxPlayerID
}

// HandleLogin 处理客户端登录消息
func HandleLogin(clientID uint64, header *steve_proto_gaterpc.Header, req room.RoomLoginReq) []exchanger.ResponseMsg {
	logentry := logrus.WithFields(logrus.Fields{
		"func_name": "HandleLogin",
		"client_id": clientID,
		"user_name": req.GetUserName(),
	})
	playerMgr := global.GetPlayerMgr()
	p := &player{
		playerID: allocPlayerID(),
		coin:     10000,
		clientID: clientID,
	}
	logentry = logentry.WithFields(logrus.Fields{
		"player_id": p.playerID,
		"coin":      p.coin,
	})
	logentry.Infoln("玩家登录")

	playerMgr.AddPlayer(p)
	return []exchanger.ResponseMsg{
		exchanger.ResponseMsg{
			MsgID: uint32(msgid.MsgID_ROOM_LOGIN_RSP),
			Body: &room.RoomLoginRsp{
				PlayerId: proto.Uint64(p.GetID()),
				Coin:     proto.Uint64(p.GetCoin()),
			},
		},
	}
}
