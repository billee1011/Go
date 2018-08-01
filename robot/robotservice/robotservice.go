package robotservice

import (
	"context"
	"steve/server_pb/robot"
)

//Robotservice 机器人服务
type Robotservice struct{}

var defaultObject = new(Robotservice)
var _ robot.RobotServiceServer = DefaultRobot()

// DefaultRobot 默认对象
func DefaultRobot() *Robotservice {
	return defaultObject
}

//GetRobotPlayerIDByInfo 根据请求信息获取机器人玩家ID
func (r *Robotservice) GetRobotPlayerIDByInfo(ctx context.Context, request *robot.GetRobotPlayerIDReq) (*robot.GetRobotPlayerIDRsp, error) {
	playerID, err := getRobotPlayerIDByInfo(request)
	rsp := &robot.GetRobotPlayerIDRsp{
		RobotPlayerId: playerID,
	}
	return rsp, err
}

//SetRobotPlayerState 设置机器人玩家状态
func (r *Robotservice) SetRobotPlayerState(ctx context.Context, request *robot.SetRobotPlayerStateReq) (*robot.SetRobotPlayerStateRsp, error) {
	rsp := &robot.SetRobotPlayerStateRsp{
		Result: true,
	}
	if err := setRobotPlayerState(request); err != nil {
		rsp.Result = false
		return rsp, err
	}
	return rsp, nil
}
