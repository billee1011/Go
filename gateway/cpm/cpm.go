// Package cpm 为 connectPlayerMap ， 维护连接到玩家的映射关系
package cpm

import (
	"steve/gateway/global"
	"sync"
)

type connectPlayerMap struct {
	connectPlayer sync.Map
	playerConnect sync.Map
	mu            sync.Mutex
}

func (cpm *connectPlayerMap) GetConnectPlayer(clientID uint64) uint64 {
	v, ok := cpm.connectPlayer.Load(clientID)
	if !ok {
		return 0
	}
	return v.(uint64)
}

func (cpm *connectPlayerMap) GetPlayerConnect(playerID uint64) uint64 {
	v, ok := cpm.playerConnect.Load(playerID)
	if !ok {
		return 0
	}
	return v.(uint64)
}

func (cpm *connectPlayerMap) SaveConnectPlayer(clientID uint64, playerID uint64) {
	cpm.mu.Lock()
	cpm.connectPlayer.Store(clientID, playerID)
	cpm.playerConnect.Store(playerID, clientID)
	cpm.mu.Unlock()
}

func (cpm *connectPlayerMap) RemoveConnect(clientID uint64) {
	playerID := cpm.GetConnectPlayer(clientID)
	cpm.mu.Lock()
	cpm.connectPlayer.Delete(clientID)
	if playerID != 0 {
		cpm.playerConnect.Delete(playerID)
	}
	cpm.mu.Unlock()
}

func init() {
	global.SetConnectPlayerMap(&connectPlayerMap{})
}
