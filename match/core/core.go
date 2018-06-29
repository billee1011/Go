package core

import (
	"container/list"
	"context"
	"errors"
	"fmt"
	msgid "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/server_pb/room_mgr"
	"steve/structs"
	"steve/structs/exchanger"
	"steve/structs/proto/gate_rpc"
	"steve/structs/service"
	"time"

	"github.com/Sirupsen/logrus"
)

const (
	playersOneDesk int = 4 // 一个桌子需要的人数
)

type matchCore struct {
	e            *structs.Exposer
	matchManager *matchManager // 匹配管理器
}

// 等待匹配的玩家
type matchPlayer struct {
	playerID  uint64 // 玩家唯一ID
	startTime int64  // 匹配开始时间
	//playerGrade   uint8        // 玩家段位
}

// 匹配管理器
type matchManager struct {
	matchCore   *matchCore // matchCore
	waitPlayers *list.List // 单个队列的等待匹配玩家信息
	//sendPlayers set      // 已发送给room创建桌子的players
	//allWaitPlayers *map  // 所有匹配中的玩家信息
}

// 加入一个新的匹配玩家
func (mm *matchManager) addPlayer(pMatchPlayer *matchPlayer) {
	// 日志
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "matchManager::addPlayer",
	})

	if pMatchPlayer == nil {
		logEntry.Errorln("matchManager::addPlayer() pMatchPlayer == nil")
		return
	}

	logEntry.WithField("player", pMatchPlayer).Debugln("matchManager::addPlayer()")

	//_, found := mm.waitPlayers.Get(pMatchPlayer.playerID)

	// 已经存在则报错
	//if found {
	//	logEntry.WithField("playerID", pMatchPlayer.playerID).Errorln("该匹配玩家已经存在于匹配列表中")
	//	return
	//}

	// 插入
	//mm.waitPlayers.Put(pMatchPlayer.playerID, pMatchPlayer)
	mm.waitPlayers.PushBack(pMatchPlayer)

	logEntry.Errorln("mm.waitPlayers.Size() = ", mm.waitPlayers.Len())

	// 个数满足就匹配一次
	if mm.waitPlayers.Len() > playersOneDesk {
		mm.match()
	}
}

// 遍历检测所有的待匹配玩家，超时的匹配机器人
func (mm *matchManager) CheckTimeout() {
	// todo
}

// 具体的匹配操作
func (mm *matchManager) match() {

	// 日志
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "matchManager::match",
	})

	logEntry.Debugln("matchManager::match()")

	// 人数不足时退出
	// todo 检测具体类型的等待匹配人数matchPlayer
	if mm.waitPlayers.Len() < playersOneDesk {
		return
	}

	// todo 检测玩家段位

	// 目前办法：够四个就请求room服开启一个桌子

	// 临时数组，存储玩家的playerID(一个桌子的人数）
	deskPlayers := make([]uint64, 0, playersOneDesk)

	var nextNode *list.Element

	for iter := mm.waitPlayers.Front(); iter != nil; iter = nextNode {

		nextNode = iter.Next()

		tempPlayer := (iter.Value).(*matchPlayer)

		// 加入到临时数组
		deskPlayers = append(deskPlayers, tempPlayer.playerID)

		// 删除iter(暂时，以后应更改为缓存在本地，等room服创建成功后再删除)
		mm.waitPlayers.Remove(iter)

		// 每满一个桌子，就通知room创建
		if len(deskPlayers) >= playersOneDesk {
			createErr := mm.matchCore.NofityRoomCreateDesk(deskPlayers)

			// 创建失败
			if createErr != nil {
				logEntry.WithError(createErr).Errorln("matchManager::match() NofityRoomCreateDesk() error, playerIDs = ", deskPlayers)
				return
			}

			// 创建成功

			// 把数组里面的playerID从匹配列表删除
			//for _, value := range deskPlayers {
			//	mm.waitPlayers.Remove(value)
			//}

			// 清空临时数组
			deskPlayers = deskPlayers[playersOneDesk:]
		}
	}
}

// NewService 创建服务
func NewService() service.Service {
	return new(matchCore)
}

func (c *matchCore) Init(e *structs.Exposer, param ...string) error {
	entry := logrus.WithField("name", "matchCore.Init")

	c.e = e

	// 创建匹配管理器
	c.matchManager = new(matchManager)
	c.matchManager.matchCore = c
	c.matchManager.waitPlayers = list.New()

	// 注册消息处理
	if err := c.registerHandles(e.Exchanger); err != nil {
		entry.WithError(err).Error("注册消息处理器失败")
		return err
	}
	return nil
}

func (c *matchCore) Start() error {
	return nil
}

// 注册消息处理
func (c *matchCore) registerHandles(e exchanger.Exchanger) error {
	registe := func(id msgid.MsgID, handler interface{}) {
		err := e.RegisterHandle(uint32(id), handler)
		if err != nil {
			panic(err)
		}
	}

	registe(msgid.MsgID_MATCH_REQ, c.handleMatch) // 匹配请求消息

	return nil
}

// 匹配请求的处理(来自网关服)
func (c *matchCore) handleMatch(playerID uint64, header *steve_proto_gaterpc.Header, req room.RoomJoinDeskReq) (ret []exchanger.ResponseMsg) {

	// 日志
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "matchCore::handleMatch()",
	})

	// 单个回应消息体
	response := &room.RoomJoinDeskRsp{
		ErrCode: room.RoomError_SUCCESS.Enum(),
	}

	// 结果
	ret = []exchanger.ResponseMsg{{
		MsgID: uint32(msgid.MsgID_MATCH_RSP),
		Body:  response,
	}}

	// 新建一个匹配玩家
	newMatchPlayer := &matchPlayer{
		playerID:  playerID,
		startTime: time.Now().Unix(),
	}

	logEntry.WithField("newMatchPlayer", newMatchPlayer).Debugln("matchCore::handleMatch()加入新的匹配玩家")

	// 添加
	c.matchManager.addPlayer(newMatchPlayer)

	// 执行一次匹配
	c.matchManager.match()

	return
}

// 通知room服创建desk
func (c *matchCore) NofityRoomCreateDesk(playersID []uint64) error {
	// 日志
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "matchCore::NofityRoomCreateDesk",
	})
	logEntry.Debugln("matchCore::NofityRoomCreateDesk()")

	// 获取一个room服
	roomConnection, roomErr := c.e.RPCClient.GetConnectByServerName("room")
	if roomErr != nil {
		logEntry.WithError(roomErr).Errorln("获取room服失败")
	}

	if roomConnection == nil {
		logEntry.Errorln("获取room服失败，room_connection == nil")
		return errors.New("获取room服失败，matchCore::NofityRoomCreateDesk() room_connection == nil")
	}

	// 建立一个新的连接
	roomMgrClient := roommgr.NewRoomMgrClient(roomConnection)
	deskResp, deskErr := roomMgrClient.CreateDesk(context.Background(), &roommgr.CreateDeskRequest{
		PlayerId: playersID,
	})

	if deskErr != nil {
		logEntry.WithError(deskErr).Errorln("调用room服的CreateDesk()失败")
		return fmt.Errorf("call room::CreateDesk() failed: %v", deskErr)
	}

	// 打印创建桌子的返回消息
	fmt.Println("收到room服 CreateDesk()返回消息 : ", deskResp.GetErrCode())
	return nil
}
