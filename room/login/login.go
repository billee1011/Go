package login

import (
	"fmt"
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
func (p *player) SetClientID(clientID uint64) {
	logrus.WithFields(logrus.Fields{
		"client_id":     clientID,
		"old_client_id": p.clientID,
	}).Debugln("设置客户端ID")
	p.clientID = clientID
}

var maxPlayerID uint64
var maxPlayerIDMutex sync.Mutex

func allocPlayerID() uint64 {
	maxPlayerIDMutex.Lock()
	defer maxPlayerIDMutex.Unlock()
	maxPlayerID++
	return maxPlayerID
}

func loginPlayer(clientID uint64, playerID uint64) *player {
	playerMgr := global.GetPlayerMgr()
	pm := playerMgr.GetPlayer(playerID)
	if pm == nil {
		p := &player{
			playerID: playerID,
			coin:     10000,
			clientID: clientID,
		}
		playerMgr.AddPlayer(p)
		return p
	}
	pm.SetClientID(clientID)
	return &player{
		playerID: pm.GetID(),
		coin:     pm.GetCoin(),
		clientID: pm.GetClientID(),
	}
}

// HandleLogin 处理客户端登录消息
func HandleLogin(clientID uint64, header *steve_proto_gaterpc.Header, req room.RoomLoginReq) []exchanger.ResponseMsg {
	logentry := logrus.WithFields(logrus.Fields{
		"func_name": "HandleLogin",
		"client_id": clientID,
		"user_name": req.GetUserName(),
	})
	playerID := allocPlayerID()

	p := loginPlayer(clientID, playerID)
	logentry = logentry.WithFields(logrus.Fields{
		"player_id": p.playerID,
		"coin":      p.coin,
	})
	logentry.Infoln("玩家登录")
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

// HandleVisitorLogin 处理游客登录
func HandleVisitorLogin(clientID uint64, header *steve_proto_gaterpc.Header, req room.RoomVisitorLoginReq) []exchanger.ResponseMsg {
	logentry := logrus.WithFields(logrus.Fields{
		"func_name":   "HandleVisitorLogin",
		"client_id":   clientID,
		"device_info": req.GetDeviceInfo(),
	})
	exsit := EqualUser(req.GetDeviceInfo())
	var playerID uint64
	if !exsit {
		playerID = allocPlayerID()
		YoukeInfos[req.GetDeviceInfo().GetUuid()] = &YoukeInfo{
			DeviceInfo: req.GetDeviceInfo(),
			PlayerID:   playerID,
		}
	} else {
		youkeInfo := YoukeInfos[req.GetDeviceInfo().GetUuid()]
		playerID = youkeInfo.PlayerID
	}
	p := loginPlayer(clientID, playerID)
	userName := fmt.Sprintf("youke%v", playerID)

	logentry = logentry.WithFields(logrus.Fields{
		"player_id": p.playerID,
		"coin":      p.coin,
		"user_name": userName,
	})
	logentry.Infoln("游客玩家登录")

	return []exchanger.ResponseMsg{
		exchanger.ResponseMsg{
			MsgID: uint32(msgid.MsgID_ROOM_VISITOR_LOGIN_RSP),
			Body: &room.RoomVisitorLoginRsp{
				ErrCode:  room.RoomError_Success.Enum(),
				UserName: proto.String(userName),
				PlayerId: proto.Uint64(p.GetID()),
				Coin:     proto.Uint64(p.GetCoin()),
			},
		},
	}
}

// YoukeInfos 储存游客登录的信息
var YoukeInfos map[string]*YoukeInfo

// YoukeInfo 游客信息
type YoukeInfo struct {
	DeviceInfo *room.DeviceInfo
	PlayerID   uint64
}

func init() {
	YoukeInfos = make(map[string]*YoukeInfo)
}

// EqualUser 判断游客是否第一次登录
func EqualUser(dvInfo *room.DeviceInfo) bool {
	info, ok := YoukeInfos[dvInfo.GetUuid()]
	if !ok {
		return false
	}
	if info.DeviceInfo.GetDeviceType() == dvInfo.GetDeviceType() &&
		info.DeviceInfo.GetUuid() == dvInfo.GetUuid() {
		return true
	}
	return false
}
