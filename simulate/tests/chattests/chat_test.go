package chattests

import (
	"fmt"
	msgid "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

//TestChat 聊天测试
//游戏过程中庄家发送聊天信息
//期望：所有玩家都收到，庄家发送的聊天信息
func TestChat(t *testing.T) {
	// 开始游戏
	params := global.NewCommonStartGameParams()
	deskData, err := utils.StartGame(params)
	assert.NotNil(t, deskData)
	assert.Nil(t, err)

	// 庄家发送类型为打字的，信息为“大家好！”的聊天信息
	bankerSeat := params.BankerSeat
	//开局 0 自询
	assert.Nil(t, utils.WaitZixunNtf(deskData, bankerSeat))
	assert.Nil(t, sendChatReq(deskData, bankerSeat, 3, "大家好!"))

	// 所有人都要收到庄家发送来的聊天信息
	for _, player := range deskData.Players {
		assert.Nil(t, waitChatRsp(deskData, player.Seat, bankerSeat, 3, "大家好!"))
	}
}

// sendChatReq 发送聊天请求
func sendChatReq(deskData *utils.DeskData, seat int, chatType room.ChatType, chatInfo string) error {
	player := utils.GetDeskPlayerBySeat(seat, deskData)
	client := player.Player.GetClient()
	_, err := client.SendPackage(utils.CreateMsgHead(msgid.MsgID_ROOM_CHAT_REQ), &room.RoomDeskChatReq{
		ChatType: &chatType,
		ChatInfo: &chatInfo,
	})
	return err
}

// waitChatRsp 等待接收玩家聊天信息,sourceSeat 聊天发起人，seat接受聊天的人
func waitChatRsp(deskData *utils.DeskData, seat, sourceSeat int, chatType room.ChatType, chatInfo string) error {
	player := utils.GetDeskPlayerBySeat(seat, deskData)
	expector, _ := player.Expectors[msgid.MsgID_ROOM_CHAT_NTF]

	ntf := room.RoomDeskChatNtf{}
	if err := expector.Recv(global.DefaultWaitMessageTime, &ntf); err != nil {
		return err
	}
	sourcePlayer := utils.GetDeskPlayerBySeat(sourceSeat, deskData)
	sourcePlayerID := sourcePlayer.Player.GetID()
	if sourcePlayerID != ntf.GetPlayerId() {
		return fmt.Errorf("聊天发起人错误:%v", ntf.GetPlayerId())
	}
	if ntf.GetChatType() != chatType {
		return fmt.Errorf("聊天类型错误:%v", ntf.GetChatType())
	}
	if ntf.GetChatInfo() != chatInfo {
		return fmt.Errorf("聊天信息错误:%v", ntf.GetChatInfo())
	}
	return nil
}
