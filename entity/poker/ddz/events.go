package ddz

import (
	"steve/entity/poker"
	"strconv"
)

// EventID 事件 ID
type EventID int32

const (
	EventID_event_invalid         EventID = 0
	EventID_event_start_game      EventID = 1
	EventID_event_deal_finish     EventID = 2
	EventID_event_grab_request    EventID = 3
	EventID_event_double_request  EventID = 4
	EventID_event_chupai_request  EventID = 5
	EventID_event_showhand_finish EventID = 6
	EventID_event_resume_request  EventID = 7
)

var EventID_name = map[int32]string{
	0: "event_invalid",
	1: "event_start_game",
	2: "event_deal_finish",
	3: "event_grab_request",
	4: "event_double_request",
	5: "event_chupai_request",
	6: "event_showhand_finish",
	7: "event_resume_request",
}

func (x EventID) String() string {
	s, ok := StateID_name[int32(x)]
	if ok {
		return s
	}
	return strconv.Itoa(int(x))
}

// RequestEventHead 请求事件通用头
type RequestEventHead struct {
	PlayerId uint64 `protobuf:"varint,1,opt,name=player_id,json=playerId" json:"player_id,omitempty"`
}

func (m *RequestEventHead) GetPlayerId() uint64 {
	if m != nil {
		return m.PlayerId
	}
	return 0
}

// GrabRequestEvent 叫/抢地主请求事件
type GrabRequestEvent struct {
	Head *RequestEventHead `protobuf:"bytes,1,opt,name=head" json:"head,omitempty"`
	Grab bool              `protobuf:"varint,2,opt,name=grab" json:"grab,omitempty"`
}

func (m *GrabRequestEvent) GetHead() *RequestEventHead {
	if m != nil {
		return m.Head
	}
	return nil
}

func (m *GrabRequestEvent) GetGrab() bool {
	if m != nil {
		return m.Grab
	}
	return false
}

// DoubleRequestEvent 加倍请求事件
type DoubleRequestEvent struct {
	Head     *RequestEventHead `protobuf:"bytes,1,opt,name=head" json:"head,omitempty"`
	IsDouble bool              `protobuf:"varint,2,opt,name=is_double,json=isDouble" json:"is_double,omitempty"`
}

func (m *DoubleRequestEvent) GetHead() *RequestEventHead {
	if m != nil {
		return m.Head
	}
	return nil
}

func (m *DoubleRequestEvent) GetIsDouble() bool {
	if m != nil {
		return m.IsDouble
	}
	return false
}

// PlayCardRequestEvent 出牌请求事件
type PlayCardRequestEvent struct {
	Head     *RequestEventHead `protobuf:"bytes,1,opt,name=head" json:"head,omitempty"`
	Cards    []uint32          `protobuf:"varint,2,rep,packed,name=cards" json:"cards,omitempty"`
	CardType poker.CardType    `protobuf:"varint,3,opt,name=card_type,json=cardType,enum=ddz.CardType" json:"card_type,omitempty"`
}

func (m *PlayCardRequestEvent) GetHead() *RequestEventHead {
	if m != nil {
		return m.Head
	}
	return nil
}

func (m *PlayCardRequestEvent) GetCards() []uint32 {
	if m != nil {
		return m.Cards
	}
	return nil
}

func (m *PlayCardRequestEvent) GetCardType() poker.CardType {
	if m != nil {
		return m.CardType
	}
	return poker.CardType_CT_NONE
}

// ResumeRequestEvent 恢复对局请求事件
type ResumeRequestEvent struct {
	Head *RequestEventHead `protobuf:"bytes,1,opt,name=head" json:"head,omitempty"`
}

func (m *ResumeRequestEvent) GetHead() *RequestEventHead {
	if m != nil {
		return m.Head
	}
	return nil
}
