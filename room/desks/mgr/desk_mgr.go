package mgr

import (
	"errors"
	"fmt"
	"steve/common/data/player"
	"steve/room/interfaces"
	"steve/room/interfaces/facade"
	"steve/room/interfaces/global"
	"steve/structs/proto/gate_rpc"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

type deskMgr struct {
	deskMap       sync.Map // deskID -> desk
	playerDeskMap sync.Map // playerID -> deskID
	mu            sync.RWMutex

	deskCount int
}

var errPlayerAlreadyInDesk = errors.New("有玩家已经在牌桌上了")
var errDeskStartError = errors.New("牌桌启动失败")

func init() {
	mgr := new(deskMgr)
	global.SetDeskMgr(mgr)
	logrus.Debugln("初始化牌桌管理器")
}

// RunDesk 运转牌桌
func (dm *deskMgr) RunDesk(desk interfaces.Desk) error {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "deskMgr.RunDesk",
		"desk_uid":  desk.GetUID(),
	})

	playerIDs := facade.GetDeskPlayerIDs(desk)
	dm.bindPlayerRoomAddr(playerIDs)

	for _, playerID := range playerIDs {
		if _, ok := dm.playerDeskMap.Load(playerID); ok {
			logEntry.WithField("player_id", playerID).Errorln(errPlayerAlreadyInDesk)
			return errPlayerAlreadyInDesk
		}
	}
	deskUID := desk.GetUID()
	dm.deskMap.Store(deskUID, desk)
	dm.deskCount++
	for _, playerID := range playerIDs {
		dm.playerDeskMap.Store(playerID, deskUID)
	}

	if err := desk.Start(dm.deskFinish(desk)); err != nil {
		logEntry.WithError(err).Errorln(errDeskStartError)
		dm.finishDesk(deskUID, playerIDs)
		return errDeskStartError
	}
	return nil
}

// 绑定玩家所在 room 服
func (dm *deskMgr) bindPlayerRoomAddr(players []uint64) {
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "deskFactory.bindPlayerRoomAddr",
		"players":   players,
	})
	roomIP := viper.GetString("rpc_addr")
	roomPort := viper.GetInt("rpc_port")
	roomAddr := fmt.Sprintf("%s:%d", roomIP, roomPort)
	for _, playerID := range players {
		if err := player.SetPlayerRoomAddr(playerID, roomAddr); err != nil {
			entry.WithError(err).Errorln("绑定玩家所在 room 失败")
		}
	}
}

// 解除玩家的 room 服绑定
func (dm *deskMgr) unbindPlayerRoomAddr(players []uint64) {
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "deskFactory.unbindPlayerRoomAddr",
		"players":   players,
	})
	for _, playerID := range players {
		if err := player.SetPlayerRoomAddr(playerID, ""); err != nil {
			entry.WithError(err).Errorln("解除玩家的 room 服绑定失败")
		}
	}
}

func (dm *deskMgr) finishDesk(deskUID uint64, players []uint64) {
	logrus.WithFields(logrus.Fields{
		"func_name": "deskMgr.finishDesk",
		"desk_uid":  deskUID,
		"players":   players,
	}).Infoln("desk finished")

	dm.deskMap.Delete(deskUID)
	dm.deskCount--
	for _, playerID := range players {
		dm.playerDeskMap.Delete(playerID)
	}
	dm.unbindPlayerRoomAddr(players)
}

func (dm *deskMgr) deskFinish(desk interfaces.Desk) func() {
	deskUID := desk.GetUID()
	playerIDs := facade.GetDeskPlayerIDs(desk)
	return func() {
		dm.mu.Lock()
		defer dm.mu.Unlock()
		dm.finishDesk(deskUID, playerIDs)
	}
}

// HandlePlayerRequest 处理玩家请求
func (dm *deskMgr) HandlePlayerRequest(playerID uint64, head *steve_proto_gaterpc.Header, bodyData []byte) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name":  "deskMgr.HandlePlayerRequest",
		"player_id":  playerID,
		"message_id": head.GetMsgId(),
	})
	iDeskID, ok := dm.playerDeskMap.Load(playerID)
	if !ok {
		logEntry.Infoln("玩家不在牌桌上")
		return
	}
	deskID := iDeskID.(uint64)
	logEntry = logEntry.WithField("desk_id", deskID)

	iDesk, ok := dm.deskMap.Load(deskID)
	if !ok {
		logEntry.Infoln("牌桌可能已经结束")
		return
	}
	desk := iDesk.(interfaces.Desk)
	desk.PushRequest(playerID, head, bodyData)
}

// GetRunDeskByPlayerID
func (dm *deskMgr) GetRunDeskByPlayerID(playerID uint64) (desk interfaces.Desk, err error) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "deskMgr.GetRunDeskByPlayerID",
		"player_id": playerID,
	})

	iDeskID, ok := dm.playerDeskMap.Load(playerID)
	if !ok {
		logEntry.Infoln("玩家不在牌桌上")
		return nil, errors.New("玩家不在牌桌上")
	}
	deskID := iDeskID.(uint64)
	logEntry = logEntry.WithField("desk_id", deskID)

	iDesk, ok := dm.deskMap.Load(deskID)
	if !ok {
		logEntry.Infoln("牌桌可能已经结束")
		return nil, errors.New("牌桌可能已经结束")
	}

	desk = iDesk.(interfaces.Desk)
	return desk, nil
}

// RemoveDeskPlayerByPlayerID
func (dm *deskMgr) RemoveDeskPlayerByPlayerID(playerID uint64) {
	dm.mu.Lock()
	defer dm.mu.Unlock()
	dm.playerDeskMap.Delete(playerID)
}

func (dm *deskMgr) GetDeskCount() int {
	return dm.deskCount
}