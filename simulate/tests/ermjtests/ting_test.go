package ermjtest

import (
	"steve/client_pb/msgid"
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func TestTing(t *testing.T) {
	params := global.NewCommonStartGameParams()
	params.PlayerNum = 2
	params.BankerSeat = 0
	params.PeiPaiGame = "ermj"
	params.GameID = room.GameId_GAMEID_ERRENMJ
	params.IsDq = false
	params.IsHsz = false
	params.Cards = [][]uint32{
		{11, 11, 11, 51, 52, 12, 12, 12, 13, 13, 13, 14, 14},
		{53, 54, 15, 15, 15, 16, 16, 16, 17, 17, 17, 18, 18},
	}
	params.WallCards = []uint32{11, 55, 12, 56, 13, 14, 57, 58, 14, 19, 19, 19, 41, 41, 41}
	deskData, err := utils.StartGame(params)
	assert.NotNil(t, deskData)
	assert.Nil(t, err)
	utils.CheckZixunNtfWithTing(t, deskData, 0, false, true, true, true)
	//等補花結束
	// time.Sleep(time.Second * 2)
	player := utils.GetDeskPlayerBySeat(0, deskData)
	client := player.Player.GetClient()
	_, err = client.SendPackage(utils.CreateMsgHead(msgid.MsgID_ROOM_CHUPAI_REQ), &room.RoomChupaiReq{
		Card: proto.Uint32(14),
		TingAction: &room.TingAction{
			EnableTing: proto.Bool(true),
			TingType:   room.TingType_TT_TIAN_TING.Enum(),
		},
	})
	for _, s := range []int{0, 1} {
		p := utils.GetDeskPlayerBySeat(s, deskData)
		messageExpector := p.Expectors[msgid.MsgID_ROOM_CHUPAI_NTF]
		ntf := &room.RoomChupaiNtf{}
		assert.Nil(t, messageExpector.Recv(global.DefaultWaitMessageTime, ntf))
		assert.Equal(t, uint32(14), ntf.GetCard())
		assert.Equal(t, player.Player.GetID(), ntf.GetPlayer())
		assert.Equal(t, true, ntf.GetTingAction().GetEnableTing())
	}

	// utils.CheckChuPaiNotifyWithSeats(t, deskData, uint32(14), 0, []int{0, 1})
}
