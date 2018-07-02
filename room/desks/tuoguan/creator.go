package tuoguan

import "steve/room/interfaces"

// CreateTuoguanManager 创建托管管理器
func CreateTuoguanManager() interfaces.TuoGuanMgr {
	return &tuoGuanMgr{
		players:      make(map[uint64]*tuoGuanPlayer),
		maxOverTimer: 2,
	}
}
