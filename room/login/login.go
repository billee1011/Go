package login

import (
	"fmt"
	"steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/room/interfaces"
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
	userName string
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

func (p *player) GetUserName() string {
	return p.userName
}

var maxPlayerID uint64
var maxPlayerIDMutex sync.Mutex

func allocPlayerID() uint64 {
	maxPlayerIDMutex.Lock()
	defer maxPlayerIDMutex.Unlock()
	maxPlayerID++
	return maxPlayerID
}

func loginPlayerByUserName(clientID uint64, userName string) interfaces.Player {
	playerMgr := global.GetPlayerMgr()
	pm := playerMgr.GetPlayerByUserName(userName)
	if pm == nil {
		p := &player{
			playerID: allocPlayerID(),
			coin:     10000000, //血战天胡地胡赢分数太高，容易认输进入游戏结束
			clientID: clientID,
			userName: userName,
		}
		playerMgr.AddPlayer(p)
		return p
	}
	playerMgr.UpdatePlayerClientID(pm.GetID(), clientID)
	return pm
}

// HandleLogin 处理客户端登录消息
func HandleLogin(clientID uint64, header *steve_proto_gaterpc.Header, req room.RoomLoginReq) (ret []exchanger.ResponseMsg) {
	logentry := logrus.WithFields(logrus.Fields{
		"func_name": "HandleLogin",
		"client_id": clientID,
		"user_name": req.GetUserName(),
	})
	rsp := &room.RoomLoginRsp{}
	ret = []exchanger.ResponseMsg{
		exchanger.ResponseMsg{
			MsgID: uint32(msgId.MsgID_ROOM_LOGIN_RSP),
			Body:  rsp,
		},
	}
	userName := req.GetUserName()
	if userName == "" {
		rsp.ErrCode = room.RoomError_EMPTY_USER_NAME.Enum()
		return
	}

	p := loginPlayerByUserName(clientID, userName)
	logentry = logentry.WithFields(logrus.Fields{
		"player_id": p.GetID(),
		"coin":      p.GetCoin(),
	})
	rsp.Coin = proto.Uint64(p.GetCoin())
	rsp.PlayerId = proto.Uint64(p.GetID())
	logentry.Infoln("玩家登录")
	return
}

func loginNewVisitor(clientID uint64, deviceInfo *room.DeviceInfo) interfaces.Player {
	playerMgr := global.GetPlayerMgr()
	playerID := allocPlayerID()
	p := &player{
		playerID: playerID,
		coin:     10000,
		clientID: clientID,
		userName: fmt.Sprintf("youke_%v", playerID),
	}
	playerMgr.AddPlayer(p)
	saveYoukeInfo(deviceInfo, playerID)
	return p
}

func loginExistVisitor(clientID uint64, youkeInfo *YoukeInfo) interfaces.Player {
	playerMgr := global.GetPlayerMgr()
	pm := playerMgr.GetPlayer(youkeInfo.PlayerID)
	if pm == nil {
		logrus.WithFields(logrus.Fields{
			"func_name":  "loginExistVisitor",
			"youke_info": youkeInfo,
		}).Errorln("游客不存在")
		return nil
	}
	playerMgr.UpdatePlayerClientID(pm.GetID(), clientID)
	return pm
}

// HandleVisitorLogin 处理游客登录
func HandleVisitorLogin(clientID uint64, header *steve_proto_gaterpc.Header, req room.RoomVisitorLoginReq) (ret []exchanger.ResponseMsg) {
	deviceInfo := req.GetDeviceInfo()
	logentry := logrus.WithFields(logrus.Fields{
		"func_name":   "HandleVisitorLogin",
		"client_id":   clientID,
		"device_info": deviceInfo,
	})
	rsp := &room.RoomVisitorLoginRsp{}
	ret = []exchanger.ResponseMsg{
		exchanger.ResponseMsg{
			MsgID: uint32(msgId.MsgID_ROOM_VISITOR_LOGIN_RSP),
			Body:  rsp,
		},
	}
	youkeInfo := getYoukeInfo(deviceInfo)
	var player interfaces.Player
	if youkeInfo == nil {
		player = loginNewVisitor(clientID, deviceInfo)
	} else {
		player = loginExistVisitor(clientID, youkeInfo)
	}
	if player == nil {
		rsp.ErrCode = room.RoomError_FAILED.Enum()
		return
	}
	logentry = logentry.WithFields(logrus.Fields{
		"player_id": player.GetID(),
		"coin":      player.GetCoin(),
		"user_name": player.GetUserName(),
	})
	logentry.Infoln("游客登录")
	rsp.Coin = proto.Uint64(player.GetCoin())
	rsp.ErrCode = room.RoomError_SUCCESS.Enum()
	rsp.PlayerId = proto.Uint64(player.GetID())
	rsp.UserName = proto.String(player.GetUserName())
	return
}

// youkeInfos 储存游客登录的信息
// var youkeInfos = map[string]*YoukeInfo{}

var youkeInfos sync.Map

// YoukeInfo 游客信息
type YoukeInfo struct {
	PlayerID uint64
}

// getYoukeInfo 获取设备游客信息
func getYoukeInfo(dvInfo *room.DeviceInfo) *YoukeInfo {
	iinfo, ok := youkeInfos.Load(dvInfo.GetUuid())
	if !ok {
		return nil
	}
	return iinfo.(*YoukeInfo)
}

func saveYoukeInfo(dvInfo *room.DeviceInfo, playerID uint64) {
	youkeInfos.Store(dvInfo.GetUuid(), &YoukeInfo{
		PlayerID: playerID,
	})
}
