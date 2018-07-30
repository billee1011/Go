package majong

// CardColor 牌花色
type CardColor int32

const (
	CardColor_ColorWan  CardColor = 0
	CardColor_ColorTong CardColor = 1
	CardColor_ColorTiao CardColor = 2
	CardColor_ColorZi   CardColor = 3
	CardColor_ColorHua  CardColor = 4
)

// TingType 听牌类型
type TingType int32

const (
	TingType_TT_NORMAL_TING TingType = 0
	TingType_TT_TIAN_TING   TingType = 1
)

// Card 卡牌结构
type Card struct {
	Color CardColor `protobuf:"varint,1,opt,name=color,enum=majong.CardColor" json:"color,omitempty"`
	Point int32     `protobuf:"varint,2,opt,name=point" json:"point,omitempty"`
}

func (m *Card) GetColor() CardColor {
	if m != nil {
		return m.Color
	}
	return CardColor_ColorWan
}

func (m *Card) GetPoint() int32 {
	if m != nil {
		return m.Point
	}
	return 0
}

// TingAction 听牌动作
type TingAction struct {
	EnableTing bool     `protobuf:"varint,1,opt,name=enable_ting,json=enableTing" json:"enable_ting,omitempty"`
	TingType   TingType `protobuf:"varint,2,opt,name=ting_type,json=tingType,enum=majong.TingType" json:"ting_type,omitempty"`
}

func (m *TingAction) GetEnableTing() bool {
	if m != nil {
		return m.EnableTing
	}
	return false
}

func (m *TingAction) GetTingType() TingType {
	if m != nil {
		return m.TingType
	}
	return TingType_TT_NORMAL_TING
}
