package utils

import (
	"errors"
	 "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/simulate/global"

	"github.com/golang/protobuf/proto"
)

// SendChupaiReq 发送出牌请求
func SendChupaiReq(deskData *DeskData, seat int, card uint32) error {
	player := GetDeskPlayerBySeat(seat, deskData)
	client := player.Player.GetClient()
	_, err := client.SendPackage(CreateMsgHead(msgid.MsgID_ROOM_CHUPAI_REQ), &room.RoomChupaiReq{
		Card: proto.Uint32(card),
	})
	return err
}

// WaitChupaiWenxunNtf 等待出牌问询通知
func WaitChupaiWenxunNtf(desk *DeskData, seat int, canPeng bool, canDianpao bool, canGang bool) error {
	player := GetDeskPlayerBySeat(seat, desk)
	expector, _ := player.Expectors[msgid.MsgID_ROOM_CHUPAIWENXUN_NTF]

	ntf := room.RoomChupaiWenxunNtf{}
	if err := expector.Recv(global.DefaultWaitMessageTime, &ntf); err != nil {
		return err
	}
	if ntf.GetEnablePeng() != canPeng {
		return errors.New("碰标志错误")
	}
	if ntf.GetEnableDianpao() != canDianpao {
		return errors.New("点炮标志错误")
	}
	if ntf.GetEnableMinggang() != canGang {
		return errors.New("杠标志错误")
	}
	return nil
}

// WaitChupaiWenxunNtf0 等待出牌问询通知(包含检测弃动作)
func WaitChupaiWenxunNtf0(desk *DeskData, seat int, canPeng bool, canDianpao bool, canGang bool, canQi bool) error {
	player := GetDeskPlayerBySeat(seat, desk)
	expector, _ := player.Expectors[msgid.MsgID_ROOM_CHUPAIWENXUN_NTF]

	ntf := room.RoomChupaiWenxunNtf{}
	if err := expector.Recv(global.DefaultWaitMessageTime, &ntf); err != nil {
		return err
	}
	if ntf.GetEnablePeng() != canPeng {
		return errors.New("碰标志错误")
	}
	if ntf.GetEnableDianpao() != canDianpao {
		return errors.New("点炮标志错误")
	}
	if ntf.GetEnableMinggang() != canGang {
		return errors.New("杠标志错误")
	}
	if ntf.GetEnableQi() != canQi {
		return errors.New("弃标志错误")
	}
	return nil
}
