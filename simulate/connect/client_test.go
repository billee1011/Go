package connect

import (
	"fmt"
	"steve/majong_pb/msgId"
	"steve/majong_pb/pb"
	"time"

	"testing"

	"github.com/golang/protobuf/proto"
)

func roomLogin(c Client, t *testing.T) error {
	req := &pb.RoomLoginReq{
		UserName: proto.String("大哥大"),
	}
	head := Head{MsgID: uint32(msgId.MsgID_RoomLogin)}
	resp, err := c.Request(SendHead{Head: head}, req, 1*time.Second)
	if err != nil {
		t.Error(err)
		return err
	}
	if resp.Body.(*pb.RoomLoginRsp).GetLoginRes() != pb.RoomLoginRsp_Success {
		return fmt.Errorf("登录失败")
	}
	return nil
}

func roomJoin(c Client, t *testing.T) error {
	req := &pb.RoomJoinReq{
		SceneId: pb.SceneID_XLHSZ.Enum(),
	}
	head := Head{MsgID: uint32(msgId.MsgID_RoomJoin)}
	resp, err := c.Request(SendHead{Head: head}, req, 1*time.Second)
	if err != nil {
		t.Error(err)
		return err
	}
	if resp.Body.(*pb.RoomJoinRsp).GetResult() != pb.RoomJoinResult_Success {
		return fmt.Errorf("进入游戏失败")
	}
	return nil
}

func roomSitDown(c Client, t *testing.T) error {

	return nil
}

func TestGame(t *testing.T) {
	c := NewTestClient("192.168.8.215:6666", "1.0")
	if err := roomLogin(c, t); err != nil {
		panic(err)
	}
	t.Log("登录成功")
	if err := roomJoin(c, t); err != nil {
		panic(err)
	}
	t.Log("进入游戏成功")
	if err := roomSitDown(c, t); err != nil {
		panic(err)
	}
	t.Log("入座成功")
}

func TestLogin(t *testing.T) {

}
