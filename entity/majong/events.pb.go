package majong

import "strconv"

// EventID 事件 ID
type EventID int32

const (
	EventID_event_invalid              EventID = 0
	EventID_event_start_game           EventID = 1
	EventID_event_xipai_finish         EventID = 2
	EventID_event_fapai_finish         EventID = 3
	EventID_event_huansanzhang_request EventID = 4
	EventID_event_dingque_request      EventID = 5
	EventID_event_peng_request         EventID = 6
	EventID_event_gang_request         EventID = 7
	EventID_event_hu_request           EventID = 8
	EventID_event_qi_request           EventID = 9
	EventID_event_chupai_finish        EventID = 10
	EventID_event_angang_finish        EventID = 11
	EventID_event_zimo_finish          EventID = 12
	EventID_event_chupai_request       EventID = 13
	EventID_event_gang_finish          EventID = 14
	EventID_event_hu_finish            EventID = 15
	EventID_event_mopai_finish         EventID = 16
	// event_zimo_request = 17;            // 自摸请求事件
	// event_angang_request = 18;          // 暗杠请求事件
	// event_bugang_request = 19;          // 补杠请求事件
	EventID_event_bugang_finish EventID = 20
	// event_qiangganghu_request = 21;     // 抢杠胡请求事件
	// event_qiqiangganghu_request = 22;   // 放弃抢杠胡请求事件
	EventID_event_qiangganghu_finish     EventID = 23
	EventID_event_wenxun_overtime        EventID = 24
	EventID_event_peng_finish            EventID = 25
	EventID_event_timer                  EventID = 26
	EventID_event_huansanzhang_finish    EventID = 27
	EventID_event_cartoon_finish_request EventID = 28
	EventID_event_settle_finish          EventID = 29
	EventID_event_gamestart_buhua_finish EventID = 30
	EventID_event_xingpai_buhua_finish   EventID = 31
	EventID_event_chi_finish             EventID = 32
	EventID_event_chi_request            EventID = 33
)

var EventID_name = map[int32]string{
	0:  "event_invalid",
	1:  "event_start_game",
	2:  "event_xipai_finish",
	3:  "event_fapai_finish",
	4:  "event_huansanzhang_request",
	5:  "event_dingque_request",
	6:  "event_peng_request",
	7:  "event_gang_request",
	8:  "event_hu_request",
	9:  "event_qi_request",
	10: "event_chupai_finish",
	11: "event_angang_finish",
	12: "event_zimo_finish",
	13: "event_chupai_request",
	14: "event_gang_finish",
	15: "event_hu_finish",
	16: "event_mopai_finish",
	20: "event_bugang_finish",
	23: "event_qiangganghu_finish",
	24: "event_wenxun_overtime",
	25: "event_peng_finish",
	26: "event_timer",
	27: "event_huansanzhang_finish",
	28: "event_cartoon_finish_request",
	29: "event_settle_finish",
	30: "event_gamestart_buhua_finish",
	31: "event_xingpai_buhua_finish",
	32: "event_chi_finish",
	33: "event_chi_request",
}

var EventID_value = map[string]int32{
	"event_invalid":                0,
	"event_start_game":             1,
	"event_xipai_finish":           2,
	"event_fapai_finish":           3,
	"event_huansanzhang_request":   4,
	"event_dingque_request":        5,
	"event_peng_request":           6,
	"event_gang_request":           7,
	"event_hu_request":             8,
	"event_qi_request":             9,
	"event_chupai_finish":          10,
	"event_angang_finish":          11,
	"event_zimo_finish":            12,
	"event_chupai_request":         13,
	"event_gang_finish":            14,
	"event_hu_finish":              15,
	"event_mopai_finish":           16,
	"event_bugang_finish":          20,
	"event_qiangganghu_finish":     23,
	"event_wenxun_overtime":        24,
	"event_peng_finish":            25,
	"event_timer":                  26,
	"event_huansanzhang_finish":    27,
	"event_cartoon_finish_request": 28,
	"event_settle_finish":          29,
	"event_gamestart_buhua_finish": 30,
	"event_xingpai_buhua_finish":   31,
	"event_chi_finish":             32,
	"event_chi_request":            33,
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

// StartGameEvent 开始游戏事件
type StartGameEvent struct {
}

// XipaiFinishEvent 洗牌完成事件
type XipaiFinishEvent struct {
}

// FapaiFinishEvent 发牌完成事件
type FapaiFinishEvent struct {
}

// HuansanzhangRequestEvent 换三张请求事件
type HuansanzhangRequestEvent struct {
	Head  *RequestEventHead `protobuf:"bytes,1,opt,name=head" json:"head,omitempty"`
	Cards []*Card           `protobuf:"bytes,2,rep,name=cards" json:"cards,omitempty"`
	Sure  bool              `protobuf:"varint,3,opt,name=sure" json:"sure,omitempty"`
}

func (m *HuansanzhangRequestEvent) GetHead() *RequestEventHead {
	if m != nil {
		return m.Head
	}
	return nil
}

func (m *HuansanzhangRequestEvent) GetCards() []*Card {
	if m != nil {
		return m.Cards
	}
	return nil
}

func (m *HuansanzhangRequestEvent) GetSure() bool {
	if m != nil {
		return m.Sure
	}
	return false
}

// DingqueRequestEvent 定缺请求事件
type DingqueRequestEvent struct {
	Head  *RequestEventHead `protobuf:"bytes,1,opt,name=head" json:"head,omitempty"`
	Color CardColor         `protobuf:"varint,3,opt,name=color,enum=majong.CardColor" json:"color,omitempty"`
}

func (m *DingqueRequestEvent) GetHead() *RequestEventHead {
	if m != nil {
		return m.Head
	}
	return nil
}

func (m *DingqueRequestEvent) GetColor() CardColor {
	if m != nil {
		return m.Color
	}
	return CardColor_ColorWan
}

// BugangRequestEvent 补杠请求事件
type BugangRequestEvent struct {
	Head  *RequestEventHead `protobuf:"bytes,1,opt,name=head" json:"head,omitempty"`
	Cards *Card             `protobuf:"bytes,2,opt,name=cards" json:"cards,omitempty"`
}

func (m *BugangRequestEvent) GetHead() *RequestEventHead {
	if m != nil {
		return m.Head
	}
	return nil
}

func (m *BugangRequestEvent) GetCards() *Card {
	if m != nil {
		return m.Cards
	}
	return nil
}

// ChupaiRequestEvent 出牌请求事件
type ChupaiRequestEvent struct {
	Head       *RequestEventHead `protobuf:"bytes,1,opt,name=head" json:"head,omitempty"`
	Cards      *Card             `protobuf:"bytes,2,opt,name=cards" json:"cards,omitempty"`
	TingAction *TingAction       `protobuf:"bytes,3,opt,name=ting_action,json=tingAction" json:"ting_action,omitempty"`
}

func (m *ChupaiRequestEvent) GetHead() *RequestEventHead {
	if m != nil {
		return m.Head
	}
	return nil
}

func (m *ChupaiRequestEvent) GetCards() *Card {
	if m != nil {
		return m.Cards
	}
	return nil
}

func (m *ChupaiRequestEvent) GetTingAction() *TingAction {
	if m != nil {
		return m.TingAction
	}
	return nil
}

// PengRequestEvent 碰请求事件
type PengRequestEvent struct {
	Head *RequestEventHead `protobuf:"bytes,1,opt,name=head" json:"head,omitempty"`
}

func (m *PengRequestEvent) GetHead() *RequestEventHead {
	if m != nil {
		return m.Head
	}
	return nil
}

// GangRequestEvent 明杠请求事件
type GangRequestEvent struct {
	Head *RequestEventHead `protobuf:"bytes,1,opt,name=head" json:"head,omitempty"`
	Card *Card             `protobuf:"bytes,2,opt,name=card" json:"card,omitempty"`
}

func (m *GangRequestEvent) GetHead() *RequestEventHead {
	if m != nil {
		return m.Head
	}
	return nil
}

func (m *GangRequestEvent) GetCard() *Card {
	if m != nil {
		return m.Card
	}
	return nil
}

// HuRequestEvent 胡请求事件
type HuRequestEvent struct {
	Head *RequestEventHead `protobuf:"bytes,1,opt,name=head" json:"head,omitempty"`
}

func (m *HuRequestEvent) GetHead() *RequestEventHead {
	if m != nil {
		return m.Head
	}
	return nil
}

// QiRequestEvent 弃请求事件
type QiRequestEvent struct {
	Head *RequestEventHead `protobuf:"bytes,1,opt,name=head" json:"head,omitempty"`
}

func (m *QiRequestEvent) GetHead() *RequestEventHead {
	if m != nil {
		return m.Head
	}
	return nil
}

type ChiRequestEvent struct {
	Head  *RequestEventHead `protobuf:"bytes,1,opt,name=head" json:"head,omitempty"`
	Cards []*Card           `protobuf:"bytes,2,rep,name=cards" json:"cards,omitempty"`
}

func (m *ChiRequestEvent) GetHead() *RequestEventHead {
	if m != nil {
		return m.Head
	}
	return nil
}

func (m *ChiRequestEvent) GetCards() []*Card {
	if m != nil {
		return m.Cards
	}
	return nil
}

// CartoonFinishRequestEvent 动画完成时间
type CartoonFinishRequestEvent struct {
	CartoonType int32  `protobuf:"varint,1,opt,name=cartoon_type,json=cartoonType" json:"cartoon_type,omitempty"`
	PlayerId    uint64 `protobuf:"varint,2,opt,name=player_id,json=playerId" json:"player_id,omitempty"`
}

func (m *CartoonFinishRequestEvent) GetCartoonType() int32 {
	if m != nil {
		return m.CartoonType
	}
	return 0
}

func (m *CartoonFinishRequestEvent) GetPlayerId() uint64 {
	if m != nil {
		return m.PlayerId
	}
	return 0
}

// SettleFinishEvent 结算完成事件
type SettleFinishEvent struct {
	PlayerId []uint64 `protobuf:"varint,1,rep,packed,name=player_id,json=playerId" json:"player_id,omitempty"`
}

func (m *SettleFinishEvent) GetPlayerId() []uint64 {
	if m != nil {
		return m.PlayerId
	}
	return nil
}
