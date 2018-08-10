package goldclient

import (
	"context"
	"errors"
	"fmt"
	"steve/server_pb/gold"
	"steve/structs"
	"time"

	"google.golang.org/grpc"
)

/*
	功能：金币服的Client API封装,实现调用
	作者： SkyWang
	日期： 2018-7-26
*/

// 交易序列号递增ID
var gSeqIdx int64 = 0

/*
方法：添加玩家金币
参数：uid=玩家ID, goldType=货币类型, changeValue=货币变化值,funcId=功能ID,channel=渠道ID, gameId=游戏ID(0表示大厅), level=场次ID
返回：变化后的金币值,错误信息
*/
func AddGold(uid uint64, goldType int16, changeValue int64, funcId int32, channel int64, gameId int32, level int32) (int64, error) {

	// 得到服务连接
	con, err := getGoldServer(uid, goldType)
	if err != nil || con == nil {
		return 0, errors.New("no connection")
	}

	createTm := time.Now().Unix()
	// 生成唯一的交易流水号:uid-funcId-goldType-createTm-idx
	seq := fmt.Sprintf("%d-%d-%d-%d-%d", uid, funcId, goldType, time.Now().UnixNano(), getIdx())

	// 初始化Request
	item := new(gold.AddItem)
	item.Uid = uid
	item.GoldType = int32(goldType)
	item.ChangeValue = changeValue
	item.Channel = channel
	item.FuncId = funcId
	item.GameId = gameId
	item.Level = level
	item.Seq = seq
	item.Time = createTm

	// 新建Client
	client := gold.NewGoldClient(con)
	// 调用RPC方法
	rsp, err := client.AddGold(context.Background(), &gold.AddGoldReq{
		Item: item,
	})

	// 检测返回值
	if err != nil {
		return 0, err
	}

	if rsp.ErrCode != gold.ResultStat_SUCCEED {
		return 0, errors.New("add gold failed")
	}
	return rsp.CurValue, nil
}

/*
方法：获取玩家金币
参数：uid=玩家ID, goldType=货币类型
返回：当前货币值，错误信息
*/
func GetGold(uid uint64, goldType int16) (int64, error) {

	// 得到服务连接
	con, err := getGoldServer(uid, goldType)
	if err != nil || con == nil {
		return 0, errors.New("no connection")
	}

	// 初始化Request
	item := new(gold.GetItem)
	item.Uid = uid
	item.GoldType = int32(goldType)

	// 新建Client
	client := gold.NewGoldClient(con)
	// 调用RPC方法
	rsp, err := client.GetGold(context.Background(), &gold.GetGoldReq{
		Item: item,
	})

	// 检测返回值
	if err != nil {
		return 0, err
	}

	if rsp.GetErrCode() != gold.ResultStat_SUCCEED {
		return 0, fmt.Errorf("获取失败：%v", rsp.GetErrCode())
	}

	return rsp.GetItem().GetValue(), nil
}

// 获取递增ID
func getIdx() int64 {
	gSeqIdx++
	return gSeqIdx
}

// 根据金币服的路由策略生成服务连接获取方式
func getGoldServer(uid uint64, goldType int16) (*grpc.ClientConn, error) {
	_ = goldType
	e := structs.GetGlobalExposer()
	// 得到服务连接
	// 金币服采用对uid进行一致性hash路由策略.
	con, err := e.RPCClient.GetConnectByServerHashId("gold", uid)
	if err != nil || con == nil {
		return nil, errors.New("no connection")
	}

	return con, nil
}
