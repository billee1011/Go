package models

import "steve/room2/entity"

type EventModel struct {
	event chan entity.DeskEvent         // 牌桌事件通道
}

func (model EventModel) GetName() string {
	return Event
}

func (model EventModel) Start() {

}

func (model EventModel) Stop() {

}
