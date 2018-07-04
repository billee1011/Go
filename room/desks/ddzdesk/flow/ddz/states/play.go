package states

import (
	"steve/room/desks/ddzdesk/flow/machine"
	"steve/server_pb/ddz"

	"github.com/Sirupsen/logrus"
	"steve/majong/global"
	"github.com/gogo/protobuf/proto"
	"steve/client_pb/room"
	"steve/client_pb/msgId"
	"github.com/pkg/errors"
)

type playState struct{}

func (s *playState) OnEnter(m machine.Machine) {
	logrus.WithField("context", getDDZContext(m)).Debugln("进入初始状态")
}

func (s *playState) OnExit(m machine.Machine) {
	logrus.WithField("context", getDDZContext(m)).Debugln("离开初始状态")
}

func (s *playState) OnEvent(m machine.Machine, event machine.Event) (int, error) {
	if event.EventID != int(ddz.EventID_event_chupai_request) {
		return int(ddz.StateID_state_playing), global.ErrInvalidEvent
	}

	message := &ddz.PlayCardRequestEvent{}
	err := proto.Unmarshal(event.EventData, message)
	if err != nil {
		return int(ddz.StateID_state_playing), global.ErrUnmarshalEvent
	}

	context := getDDZContext(m);
	playerId := message.GetHead().GetPlayerId()
	if context.CurrentPlayerId != playerId {
		sendToPlayer(m, playerId, msgid.MsgID_ROOM_DDZ_PLAY_CARD_RSP, &room.DDZPlayCardRsp{
			Result: &room.Result{ErrCode:proto.Uint32(1), ErrDesc: proto.String("未轮到本玩家出牌")},
		})
		return int(ddz.StateID_state_playing), global.ErrInvalidRequestPlayer
	}

	player := GetPlayerByID(context.GetPlayers(), playerId)
	outCards := message.GetCards()

	if(!ContainsAll(player.Cards, outCards)){
		sendToPlayer(m, playerId, msgid.MsgID_ROOM_DDZ_PLAY_CARD_RSP, &room.DDZPlayCardRsp{
			Result: &room.Result{ErrCode:proto.Uint32(2), ErrDesc: proto.String("手牌不存在")},
		})
		return int(ddz.StateID_state_playing), errors.New("手牌没有包含所有出的牌")
	}

	if len(player.Cards) == 0 {
		broadcast(m, msgid.MsgID_ROOM_DDZ_GAME_OVER_NTF, &room.DDZGameOverNtf{WinnerId:&playerId,ShowHandTime:proto.Uint32(4)})
		return int(ddz.StateID_state_over), nil
	} else {
		return int(ddz.StateID_state_playing), nil
	}
}
