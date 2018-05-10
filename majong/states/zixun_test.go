package states

import (
	"fmt"
	"steve/majong/interfaces"
	"steve/majong/utils"
	majongpb "steve/server_pb/majong"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/proto"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestZixunState_angang(t *testing.T) {
	mc := gomock.NewController(t)
	flow := interfaces.NewMockMajongFlow(mc)
	// playersID := []uint64{1}
	// ntf := &room.RoomAngangNtf{
	// 	Player: proto.Uint64(1),
	// 	Card: &room.Card{
	// 		Color: room.CardColor_ColorWan.Enum(),
	// 		Point: proto.Int32(1),
	// 	},
	// }
	// toClientMessage := interfaces.ToClientMessage{
	// 	MsgID: int(msgid.MsgID_room_angang_ntf),
	// 	Msg:   ntf,
	// }
	// flow.EXPECT().PushMessages(playersID, toClientMessage).DoAndReturn(
	// 	func(playerIDs []uint64, msgs ...interfaces.ToClientMessage) {

	// 	},
	// )
	flow.EXPECT().PushMessages(gomock.Any(), gomock.Any()).AnyTimes()
	flow.EXPECT().GetMajongContext().Return(
		&majongpb.MajongContext{
			GameId: 1,
			Players: []*majongpb.Player{
				&majongpb.Player{
					PalyerId:        1,
					HandCards:       []*majongpb.Card{&Card1W, &Card1W, &Card1W, &Card1W, &Card2W, &Card2W, &Card2W, &Card2W, &Card3W, &Card3W, &Card3W, &Card3W, &Card4W, &Card4W},
					PossibleActions: []majongpb.Action{majongpb.Action_action_angang, majongpb.Action_action_zimo},
				},
			},
			ActivePlayer: 1,
			WallCards:    []*majongpb.Card{&Card1T},
		},
	).AnyTimes()

	s := ZiXunState{}
	gangRequestEvent := &majongpb.GangRequestEvent{
		Head: &majongpb.RequestEventHead{
			PlayerId: 1,
		},
		Card: &Card1W,
	}

	requestEvent, err := proto.Marshal(gangRequestEvent)
	assert.Nil(t, err)
	context := flow.GetMajongContext()
	player := utils.GetPlayerByID(context.GetPlayers(), context.GetActivePlayer())
	beforeResults := ""
	beforeResults += fmt.Sprintln("before暗杠：")
	beforeResults += FmtPlayerInfo(player)
	logrus.Info(beforeResults)
	stateID, err := s.ProcessEvent(majongpb.EventID_event_gang_request, requestEvent, flow)
	assert.Nil(t, err)
	assert.Equal(t, majongpb.StateID_state_angang, stateID, "执行暗杠操作成功后，状态应该为暗杠状态")
	results := ""
	results += fmt.Sprintln("after暗杠：")
	results += FmtPlayerInfo(player)
	logrus.Info(results)
}

func TestZixunState_zimo(t *testing.T) {
	mc := gomock.NewController(t)
	flow := interfaces.NewMockMajongFlow(mc)
	flow.EXPECT().PushMessages(gomock.Any(), gomock.Any()).AnyTimes()
	flow.EXPECT().GetMajongContext().Return(
		&majongpb.MajongContext{
			GameId: 1,
			Players: []*majongpb.Player{
				&majongpb.Player{
					PalyerId:        1,
					HandCards:       []*majongpb.Card{&Card1W, &Card1W, &Card1W, &Card1W, &Card2W, &Card2W, &Card2W, &Card2W, &Card3W, &Card3W, &Card3W, &Card3W, &Card4W, &Card4W},
					PossibleActions: []majongpb.Action{majongpb.Action_action_angang, majongpb.Action_action_zimo},
					DingqueColor:    majongpb.CardColor_ColorTiao,
				},
			},
			ActivePlayer: 1,
			WallCards:    []*majongpb.Card{&Card1T},
		},
	).AnyTimes()

	s := ZiXunState{}
	huRequestEvent := &majongpb.HuRequestEvent{
		Head: &majongpb.RequestEventHead{
			PlayerId: 1,
		},
	}

	requestEvent, err := proto.Marshal(huRequestEvent)
	context := flow.GetMajongContext()
	player := utils.GetPlayerByID(context.GetPlayers(), context.GetActivePlayer())
	beforeResults := ""
	beforeResults += fmt.Sprintln("before自摸：")
	beforeResults += FmtPlayerInfo(player)
	logrus.Info(beforeResults)

	stateID, err := s.ProcessEvent(majongpb.EventID_event_hu_request, requestEvent, flow)
	assert.Nil(t, err)
	assert.Equal(t, majongpb.StateID_state_zimo, stateID, "执行自摸操作成功后，状态应该为自摸状态")
	results := ""
	results += fmt.Sprintln("after自摸：")
	results += FmtPlayerInfo(player)
	logrus.Info(results)
}

func TestZixunState_bugang(t *testing.T) {
	mc := gomock.NewController(t)
	flow := interfaces.NewMockMajongFlow(mc)
	// playersID := []uint64{1}
	// ntf := &room.RoomBugangNtf{
	// 	Player: proto.Uint64(1),
	// 	Card: &room.Card{
	// 		Color: room.CardColor_ColorWan.Enum(),
	// 		Point: proto.Int32(1),
	// 	},
	// }
	// toClientMessage := interfaces.ToClientMessage{
	// 	MsgID: int(msgid.MsgID_room_bugang_ntf),
	// 	Msg:   ntf,
	// }
	// flow.EXPECT().PushMessages(playersID, toClientMessage).DoAndReturn(
	// 	func(playerIDs []uint64, msgs ...interfaces.ToClientMessage) {},
	// )
	flow.EXPECT().PushMessages(gomock.Any(), gomock.Any()).AnyTimes()
	flow.EXPECT().GetMajongContext().Return(
		&majongpb.MajongContext{
			GameId: 1,
			Players: []*majongpb.Player{
				&majongpb.Player{
					PalyerId:        1,
					HandCards:       []*majongpb.Card{&Card1W, &Card2W, &Card2W, &Card2W, &Card2W, &Card3W, &Card3W, &Card3W, &Card3W, &Card4W, &Card4W},
					PossibleActions: []majongpb.Action{majongpb.Action_action_bugang, majongpb.Action_action_zimo},
					DingqueColor:    majongpb.CardColor_ColorTiao,
					PengCards: []*majongpb.PengCard{
						&majongpb.PengCard{
							Card:      &Card1W,
							SrcPlayer: 2,
						},
					},
				},
			},
			ActivePlayer: 1,
			WallCards:    []*majongpb.Card{&Card1T},
		},
	).AnyTimes()

	s := ZiXunState{}
	gangRequestEvent := &majongpb.BugangRequestEvent{
		Head: &majongpb.RequestEventHead{
			PlayerId: 1,
		},
		Cards: &Card1W,
	}
	requestEvent, err := proto.Marshal(gangRequestEvent)
	context := flow.GetMajongContext()
	player := utils.GetPlayerByID(context.GetPlayers(), context.GetActivePlayer())
	beforeResults := ""
	beforeResults += fmt.Sprintln("before补杠：")
	beforeResults += FmtPlayerInfo(player)
	logrus.Info(beforeResults)
	stateID, err := s.ProcessEvent(majongpb.EventID_event_gang_request, requestEvent, flow)
	assert.Nil(t, err)
	assert.Equal(t, majongpb.StateID_state_bugang, stateID, "执行补杠操作成功后，状态应该为补杠状态")
	results := ""
	results += fmt.Sprintln("after补杠：")
	results += FmtPlayerInfo(player)
	logrus.Info(results)
}

func TestZixunState_chupai(t *testing.T) {
	mc := gomock.NewController(t)
	flow := interfaces.NewMockMajongFlow(mc)
	// playersID := []uint64{1}
	// ntf := &room.RoomChupaiNtf{
	// 	Player: proto.Uint64(1),
	// 	Card: &room.Card{
	// 		Color: room.CardColor_ColorWan.Enum(),
	// 		Point: proto.Int32(1),
	// 	},
	// }
	// toClientMessage := interfaces.ToClientMessage{
	// 	MsgID: int(msgid.MsgID_room_chupai_ntf),
	// 	Msg:   ntf,
	// }
	// flow.EXPECT().PushMessages(playersID, toClientMessage).DoAndReturn(
	// 	func(playerIDs []uint64, msgs ...interfaces.ToClientMessage) {},
	// )
	flow.EXPECT().PushMessages(gomock.Any(), gomock.Any()).AnyTimes()
	flow.EXPECT().GetMajongContext().Return(
		&majongpb.MajongContext{
			GameId: 1,
			Players: []*majongpb.Player{
				&majongpb.Player{
					PalyerId:        1,
					HandCards:       []*majongpb.Card{&Card1W, &Card2W, &Card2W, &Card2W, &Card2W, &Card3W, &Card3W, &Card3W, &Card3W, &Card4W, &Card4W},
					PossibleActions: []majongpb.Action{majongpb.Action_action_bugang, majongpb.Action_action_zimo},
					DingqueColor:    majongpb.CardColor_ColorTiao,
					PengCards: []*majongpb.PengCard{
						&majongpb.PengCard{
							Card:      &Card1W,
							SrcPlayer: 2,
						},
					},
				},
			},
			ActivePlayer: 1,
			WallCards:    []*majongpb.Card{&Card1T},
		},
	).AnyTimes()

	s := ZiXunState{}
	bugangRequestEvent := &majongpb.ChupaiRequestEvent{
		Head: &majongpb.RequestEventHead{
			PlayerId: 1,
		},
		Cards: &Card1W,
	}
	requestEvent, err := proto.Marshal(bugangRequestEvent)
	context := flow.GetMajongContext()
	player := utils.GetPlayerByID(context.GetPlayers(), context.GetActivePlayer())
	beforeResults := ""
	beforeResults += fmt.Sprintln("before出牌：")
	beforeResults += FmtPlayerInfo(player)
	logrus.Info(beforeResults)
	// stateID, err := s.bugang(flow, bugangRequestEvent)
	stateID, err := s.ProcessEvent(majongpb.EventID_event_chupai_request, requestEvent, flow)
	assert.Nil(t, err)
	assert.Equal(t, majongpb.StateID_state_chupai, stateID, "执行出牌操作成功后，状态应该为出牌状态")
	results := ""
	results += fmt.Sprintln("after出牌：")
	results += FmtPlayerInfo(player)
	logrus.Info(results)
}
