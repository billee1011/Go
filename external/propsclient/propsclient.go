package propsclient

import (
	"errors"
	"time"
	"fmt"
	"context"
	"google.golang.org/grpc"
	"steve/structs"
	"steve/server_pb/propserver"
)


/*
	功能：道具服的Client API封装,实现调用
	作者： SkyWang
	日期： 2018-8-14
*/

// 交易序列号递增ID
var gSeqIdx int64 = 0

/*
方法：添加玩家道具
参数：uid=玩家ID, propsList=道具列表(map[道具ID]增加数量),funcId=功能ID,channel=渠道ID, gameId=游戏ID(0表示大厅), level=场次ID
返回：变化后的金币值,错误信息
*/
func AddUserProps(uid uint64, propsList map[uint64]int64, funcId int32, channel int64, gameId int32, level int32) ( error) {

	// 得到服务连接
	con, err := getMyServer(uid)
	if err != nil || con == nil {
		return errors.New("no propserver connection")
	}

	createTm := time.Now().Unix()
	// 生成唯一的交易流水号:uid-funcId-10000-createTm-idx
	// 10000表示道具流水号
	seq := fmt.Sprintf("%d-%d-%d-%d-%d", uid, funcId, 10000, time.Now().UnixNano(), getIdx())

	// 初始化Request
	req := new(props.AddPropsReq)
	req.Uid = uid
	for id, value := range  propsList {
		a := new(props.PropsInfo)
		a.PropsId = id
		a.AddNum = value
		req.PropsList = append(req.PropsList, a)
	}

	req.Channel = channel
	req.FuncId = funcId
	req.GameId = gameId
	req.Level = level
	req.Seq = seq
	req.Time = createTm

	// 新建Client
	client := props.NewPropsClient(con)
	// 调用RPC方法
	rsp, err := client.AddUserProps(context.Background(), req)

	// 检测返回值
	if err != nil {
		return  err
	}

	if rsp.ErrCode != props.ResultStat_SUCCEED {
		return errors.New("add props failed")
	}
	return  nil
}

/*
方法：获取玩家道具
参数：uid=玩家ID, propId=道具ID（为0表示获取玩家所有道具）
返回：道具列表，错误信息
*/
func GetUserProps(uid uint64, propId uint64) ([]*props.GetItem, error) {

	// 得到服务连接
	con, err := getMyServer(uid)
	if err != nil || con == nil {
		return nil, errors.New("no propserver connection")
	}

	// 初始化Request
	req := new(props.GetPropsReq)
	req.Uid = uid
	req.PropsId = propId


	// 新建Client
	client := props.NewPropsClient(con)
	// 调用RPC方法
	rsp, err := client.GetUserProps(context.Background(), req)

	// 检测返回值
	if err != nil {
		return nil, err
	}

	if rsp.GetErrCode() != props.ResultStat_SUCCEED {
		return nil, fmt.Errorf("get props failed：%v", rsp.GetErrCode())
	}

	return rsp.GetPropsList(), nil
}

// 获取递增ID
func getIdx() int64 {
	gSeqIdx++
	return gSeqIdx
}

// 根据金币服的路由策略生成服务连接获取方式
func getMyServer(uid uint64) (*grpc.ClientConn, error) {

	e := structs.GetGlobalExposer()
	// 得到服务连接
	// 道具服采用对uid进行一致性hash路由策略.
	con, err := e.RPCClient.GetConnectByServerHashId("propserver", uid)
	if err != nil || con == nil {
		return nil, errors.New("no connection")
	}

	return con, nil
}
