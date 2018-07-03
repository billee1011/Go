package ddzdesk

import (
	"fmt"
	"steve/room/desks/deskbase"
	"steve/room/interfaces"
)

// CreateDDZDesk 创建斗地主房间
func CreateDDZDesk(players []uint64, gameID int, opt interfaces.CreateDeskOptions, alloc interfaces.DeskIDAllocator) (result interfaces.CreateDeskResult, err error) {
	id, err := alloc.AllocDeskID()
	if err != nil {
		err = fmt.Errorf("分配牌桌 ID 失败")
		return
	}
	return interfaces.CreateDeskResult{
		Desk: &desk{
			DeskBase: deskbase.NewDeskBase(id, gameID, players),
		},
	}, nil
}
