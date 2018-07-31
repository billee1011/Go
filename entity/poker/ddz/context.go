package ddz

import (
	"steve/entity/poker"
	"strconv"
)

type DDZStage int32

const (
	DDZStage_DDZ_STAGE_NONE    DDZStage = 0
	DDZStage_DDZ_STAGE_DEAL    DDZStage = 1
	DDZStage_DDZ_STAGE_CALL    DDZStage = 2
	DDZStage_DDZ_STAGE_GRAB    DDZStage = 3
	DDZStage_DDZ_STAGE_DOUBLE  DDZStage = 4
	DDZStage_DDZ_STAGE_PLAYING DDZStage = 5
	DDZStage_DDZ_STAGE_OVER    DDZStage = 6
)

var DDZStage_name = map[int32]string{
	0: "DDZ_STAGE_NONE",
	1: "DDZ_STAGE_DEAL",
	2: "DDZ_STAGE_CALL",
	3: "DDZ_STAGE_GRAB",
	4: "DDZ_STAGE_DOUBLE",
	5: "DDZ_STAGE_PLAYING",
	6: "DDZ_STAGE_OVER",
}

func (x DDZStage) String() string {
	s, ok := DDZStage_name[int32(x)]
	if ok {
		return s
	}
	return strconv.Itoa(int(x))
}

// Player 玩家数据
type Player struct {
	PlayerId    uint64   `protobuf:"varint,1,opt,name=player_id,json=playerId" json:"player_id,omitempty"`
	Lord        bool     `protobuf:"varint,2,opt,name=lord" json:"lord,omitempty"`
	Grab        bool     `protobuf:"varint,3,opt,name=grab" json:"grab,omitempty"`
	IsDouble    bool     `protobuf:"varint,4,opt,name=is_double,json=isDouble" json:"is_double,omitempty"`
	HandCards   []uint32 `protobuf:"varint,5,rep,packed,name=hand_cards,json=handCards" json:"hand_cards,omitempty"`
	OutCards    []uint32 `protobuf:"varint,6,rep,packed,name=out_cards,json=outCards" json:"out_cards,omitempty"`
	AllOutCards []uint32 `protobuf:"varint,7,rep,packed,name=all_out_cards,json=allOutCards" json:"all_out_cards,omitempty"`
	Win         bool     `protobuf:"varint,8,opt,name=win" json:"win,omitempty"`
}

func (m *Player) GetPlayerId() uint64 {
	if m != nil {
		return m.PlayerId
	}
	return 0
}

func (m *Player) GetLord() bool {
	if m != nil {
		return m.Lord
	}
	return false
}

func (m *Player) GetGrab() bool {
	if m != nil {
		return m.Grab
	}
	return false
}

func (m *Player) GetIsDouble() bool {
	if m != nil {
		return m.IsDouble
	}
	return false
}

func (m *Player) GetHandCards() []uint32 {
	if m != nil {
		return m.HandCards
	}
	return nil
}

func (m *Player) GetOutCards() []uint32 {
	if m != nil {
		return m.OutCards
	}
	return nil
}

func (m *Player) GetAllOutCards() []uint32 {
	if m != nil {
		return m.AllOutCards
	}
	return nil
}

func (m *Player) GetWin() bool {
	if m != nil {
		return m.Win
	}
	return false
}

// DDZContext 麻将现场
type DDZContext struct {
	GameId          int32     `protobuf:"varint,1,opt,name=game_id,json=gameId" json:"game_id,omitempty"`
	CurState        StateID   `protobuf:"varint,2,opt,name=cur_state,json=curState,enum=ddz.StateID" json:"cur_state,omitempty"`
	CurStage        DDZStage  `protobuf:"varint,3,opt,name=cur_stage,json=curStage,enum=ddz.DDZStage" json:"cur_stage,omitempty"`
	Players         []*Player `protobuf:"bytes,4,rep,name=players" json:"players,omitempty"`
	Dipai           []uint32  `protobuf:"varint,5,rep,packed,name=dipai" json:"dipai,omitempty"`
	CurrentPlayerId uint64    `protobuf:"varint,6,opt,name=current_player_id,json=currentPlayerId" json:"current_player_id,omitempty"`
	// 叫地主阶段
	CallPlayerId      uint64 `protobuf:"varint,7,opt,name=call_player_id,json=callPlayerId" json:"call_player_id,omitempty"`
	FirstGrabPlayerId uint64 `protobuf:"varint,8,opt,name=first_grab_player_id,json=firstGrabPlayerId" json:"first_grab_player_id,omitempty"`
	LastGrabPlayerId  uint64 `protobuf:"varint,9,opt,name=last_grab_player_id,json=lastGrabPlayerId" json:"last_grab_player_id,omitempty"`
	GrabbedCount      uint32 `protobuf:"varint,10,opt,name=grabbed_count,json=grabbedCount" json:"grabbed_count,omitempty"`
	LordPlayerId      uint64 `protobuf:"varint,11,opt,name=lord_player_id,json=lordPlayerId" json:"lord_player_id,omitempty"`
	AllAbandonCount   uint32 `protobuf:"varint,12,opt,name=all_abandon_count,json=allAbandonCount" json:"all_abandon_count,omitempty"`
	TotalGrab         uint32 `protobuf:"varint,13,opt,name=total_grab,json=totalGrab" json:"total_grab,omitempty"`
	// 加倍阶段
	DoubledPlayers []uint64 `protobuf:"varint,14,rep,packed,name=doubled_players,json=doubledPlayers" json:"doubled_players,omitempty"`
	TotalDouble    uint32   `protobuf:"varint,15,opt,name=total_double,json=totalDouble" json:"total_double,omitempty"`
	// 出牌阶段
	CurOutCards   []uint32       `protobuf:"varint,16,rep,packed,name=cur_out_cards,json=curOutCards" json:"cur_out_cards,omitempty"`
	CurCardType   poker.CardType `protobuf:"varint,17,opt,name=cur_card_type,json=curCardType,enum=ddz.CardType" json:"cur_card_type,omitempty"`
	CardTypePivot uint32         `protobuf:"varint,18,opt,name=card_type_pivot,json=cardTypePivot" json:"card_type_pivot,omitempty"`
	PassCount     uint32         `protobuf:"varint,19,opt,name=pass_count,json=passCount" json:"pass_count,omitempty"`
	TotalBomb     uint32         `protobuf:"varint,20,opt,name=total_bomb,json=totalBomb" json:"total_bomb,omitempty"`
	// 结算阶段
	WinnerId   uint64 `protobuf:"varint,21,opt,name=winner_id,json=winnerId" json:"winner_id,omitempty"`
	Spring     bool   `protobuf:"varint,22,opt,name=spring" json:"spring,omitempty"`
	AntiSpring bool   `protobuf:"varint,23,opt,name=anti_spring,json=antiSpring" json:"anti_spring,omitempty"`
	// 恢复对局
	GrabbedPlayers []uint64 `protobuf:"varint,24,rep,packed,name=grabbed_players,json=grabbedPlayers" json:"grabbed_players,omitempty"`
	// 倒计时
	CountDownPlayers []uint64 `protobuf:"varint,25,rep,packed,name=count_down_players,json=countDownPlayers" json:"count_down_players,omitempty"`
	StartTime        []byte   `protobuf:"bytes,26,opt,name=startTime,proto3" json:"startTime,omitempty"`
	Duration         uint32   `protobuf:"varint,27,opt,name=duration" json:"duration,omitempty"`
	WallCards        []uint32 `protobuf:"varint,28,rep,packed,name=wall_cards,json=wallCards" json:"wall_cards,omitempty"`
}

func (m *DDZContext) GetGameId() int32 {
	if m != nil {
		return m.GameId
	}
	return 0
}

func (m *DDZContext) GetCurState() StateID {
	if m != nil {
		return m.CurState
	}
	return StateID_state_init
}

func (m *DDZContext) GetCurStage() DDZStage {
	if m != nil {
		return m.CurStage
	}
	return DDZStage_DDZ_STAGE_NONE
}

func (m *DDZContext) GetPlayers() []*Player {
	if m != nil {
		return m.Players
	}
	return nil
}

func (m *DDZContext) GetDipai() []uint32 {
	if m != nil {
		return m.Dipai
	}
	return nil
}

func (m *DDZContext) GetCurrentPlayerId() uint64 {
	if m != nil {
		return m.CurrentPlayerId
	}
	return 0
}

func (m *DDZContext) GetCallPlayerId() uint64 {
	if m != nil {
		return m.CallPlayerId
	}
	return 0
}

func (m *DDZContext) GetFirstGrabPlayerId() uint64 {
	if m != nil {
		return m.FirstGrabPlayerId
	}
	return 0
}

func (m *DDZContext) GetLastGrabPlayerId() uint64 {
	if m != nil {
		return m.LastGrabPlayerId
	}
	return 0
}

func (m *DDZContext) GetGrabbedCount() uint32 {
	if m != nil {
		return m.GrabbedCount
	}
	return 0
}

func (m *DDZContext) GetLordPlayerId() uint64 {
	if m != nil {
		return m.LordPlayerId
	}
	return 0
}

func (m *DDZContext) GetAllAbandonCount() uint32 {
	if m != nil {
		return m.AllAbandonCount
	}
	return 0
}

func (m *DDZContext) GetTotalGrab() uint32 {
	if m != nil {
		return m.TotalGrab
	}
	return 0
}

func (m *DDZContext) GetDoubledPlayers() []uint64 {
	if m != nil {
		return m.DoubledPlayers
	}
	return nil
}

func (m *DDZContext) GetTotalDouble() uint32 {
	if m != nil {
		return m.TotalDouble
	}
	return 0
}

func (m *DDZContext) GetCurOutCards() []uint32 {
	if m != nil {
		return m.CurOutCards
	}
	return nil
}

func (m *DDZContext) GetCurCardType() poker.CardType {
	if m != nil {
		return m.CurCardType
	}
	return poker.CardType_CT_NONE
}

func (m *DDZContext) GetCardTypePivot() uint32 {
	if m != nil {
		return m.CardTypePivot
	}
	return 0
}

func (m *DDZContext) GetPassCount() uint32 {
	if m != nil {
		return m.PassCount
	}
	return 0
}

func (m *DDZContext) GetTotalBomb() uint32 {
	if m != nil {
		return m.TotalBomb
	}
	return 0
}

func (m *DDZContext) GetWinnerId() uint64 {
	if m != nil {
		return m.WinnerId
	}
	return 0
}

func (m *DDZContext) GetSpring() bool {
	if m != nil {
		return m.Spring
	}
	return false
}

func (m *DDZContext) GetAntiSpring() bool {
	if m != nil {
		return m.AntiSpring
	}
	return false
}

func (m *DDZContext) GetGrabbedPlayers() []uint64 {
	if m != nil {
		return m.GrabbedPlayers
	}
	return nil
}

func (m *DDZContext) GetCountDownPlayers() []uint64 {
	if m != nil {
		return m.CountDownPlayers
	}
	return nil
}

func (m *DDZContext) GetStartTime() []byte {
	if m != nil {
		return m.StartTime
	}
	return nil
}

func (m *DDZContext) GetDuration() uint32 {
	if m != nil {
		return m.Duration
	}
	return 0
}

func (m *DDZContext) GetWallCards() []uint32 {
	if m != nil {
		return m.WallCards
	}
	return nil
}
