package utils

import (
	"fmt"
	"steve/client_pb/common"
	"steve/client_pb/hall"
	"steve/client_pb/msgid"
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/interfaces"
)

// SendRecoverGameReq 发送恢复游戏请求
func SendRecoverGameReq(seat int, deskData *DeskData) error {
	player := GetDeskPlayerBySeat(seat, deskData)
	client := player.Player.GetClient()
	_, err := client.SendPackage(CreateMsgHead(msgid.MsgID_ROOM_RESUME_GAME_REQ), &room.RoomResumeGameReq{})
	return err
}

// // SendNeedRecoverGameReq 发送恢复游戏请求
// func SendNeedRecoverGameReq(seat int, deskData *DeskData) error {
// 	player := GetDeskPlayerBySeat(seat, deskData)
// 	client := player.Player.GetClient()
// 	_, err := client.SendPackage(CreateMsgHead(msgid.MsgID_HALL_GET_PLAYER_STATE_REQ), &hall.HallGetPlayerStateReq{})
// 	return err
// }

// GetDeskPlayerState 获取牌桌玩家状态
func GetDeskPlayerState(player interfaces.ClientPlayer) (common.PlayerState, error) {
	client := player.GetClient()
	player.AddExpectors(msgid.MsgID_HALL_GET_PLAYER_STATE_RSP)
	_, err := client.SendPackage(CreateMsgHead(msgid.MsgID_HALL_GET_PLAYER_STATE_REQ), &hall.HallGetPlayerStateReq{})
	if err != nil {
		return common.PlayerState_PS_IDLE, fmt.Errorf("发送获取状态请求失败 %v", err)
	}
	expector := player.GetExpector(msgid.MsgID_HALL_GET_PLAYER_STATE_RSP)
	response := hall.HallGetPlayerStateRsp{}
	if err := expector.Recv(global.DefaultWaitMessageTime, &response); err != nil {
		return common.PlayerState_PS_IDLE, fmt.Errorf("没有收到获取状态响应 %v", err)
	}
	return response.GetPlayerState(), nil
}
