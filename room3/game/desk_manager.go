package game

import (
	"steve/room3/game/utils"
	"steve/structs/exchanger"
	"steve/structs/proto/gate_rpc"
	"sync"
)

var DefaultDeskManager = NewDeskManager()

func NewDeskManager() *DeskManager {
	return &DeskManager{
		desks: make(map[uint64]*Desk),
	}
}

type DeskManager struct {
	desks map[uint64]*Desk // deskID - *Desk
	mutex sync.Mutex
}

func (dm *DeskManager) CreateDesk(playerIDs []uint64, gameID int, option *DeskOption) (*Desk, error) {

	// 1. 分配桌子ID
	// 2. 分配位置号，创建Player，并加入PlayerManager
	// 3. 设置结算器，建立deskevent chan

	d := NewDesk(utils.DefaultIDAlloc.AllocDeskID(), gameID, option)

	// 创建桌子玩家
	var seatID uint32
	for _, playerID := range playerIDs {
		player := NewPlayer(playerID, seatID, d)

		d.players = append(d.players, player)
		DefaultPlayManager.AddPlayer(player)

		seatID++
	}

	dm.desks[d.deskID] = d

	return d, nil
}

func (dm *DeskManager) RunDesk(deskID uint64) {
	d := dm.desks[deskID]

	if err := d.Start(); err != nil {
		d.Stop()
	}
}

func (dm *DeskManager) HandleDeskRequest(playerID uint64, head *steve_proto_gaterpc.Header, bodyData []byte) (rspMsg []exchanger.ResponseMsg) {
	player := DefaultPlayManager.GetPlayer(playerID)
	player.desk.HandlePlayerRequest(playerID, head, bodyData)
	return []exchanger.ResponseMsg{}
}

func (dm *DeskManager) HandleCancelTuoGuanRequest(playerID uint64, head *steve_proto_gaterpc.Header, bodyData []byte) (rspMsg []exchanger.ResponseMsg) {
	return []exchanger.ResponseMsg{}
}

func (dm *DeskManager) HandleExitRequest(playerID uint64, head *steve_proto_gaterpc.Header, bodyData []byte) (rspMsg []exchanger.ResponseMsg) {
	return []exchanger.ResponseMsg{}
}
