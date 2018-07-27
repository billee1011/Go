package data

import (
	"fmt"
	"steve/common/data/redis"
	"strconv"
	"strings"
)

/*
	功能： 服务数据保存到redis.
	作者： SkyWang
	日期： 2018-7-25

*/

// 从redis加载玩家金币
func LoadGoldFromRedis(uid uint64) (map[int16]int64, error) {

	r := redis.GetRedisClient()
	key := fmtPlayerKey(uid)
	cmd := r.HGetAll(key)
	if cmd.Err() != nil {
		//logic.ErrNoUser.WithError(cmd.Err()).Errorln(errRedisOperation)
		return nil, fmt.Errorf("get redis err:%v", cmd.Err())
	}
	m := cmd.Val()
	if len(m) == 0 {
		return nil,  fmt.Errorf("redis no user: uid=%d", uid)
	}
	list := make(map[int16]int64, len(m))
	for k, v := range m {
		sp := strings.Split(k, "_")
		if len(sp) == 2 {
			k = sp[1]
		}
		t, err := strconv.ParseInt(k, 10, 16)
		if err != nil {
			return nil, fmt.Errorf("get redis ret err1:%v", m)
		}
		g, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("get redis ret err2:%v", m)
		}
		list[int16(t)] = g
	}

	return list, nil
}

// 保存玩家金币到Redis
func SaveGoldToRedis(uid uint64, goldList map[int16]int64) error {
	r := redis.GetRedisClient()
	key := fmtPlayerKey(uid)

	list := make(map[string]interface{}, len(goldList))
	for k, v := range goldList {
		strKey := fmt.Sprintf("%d", k)
		list[strKey] = v
	}
	cmd := r.HMSet(key, list)
	if cmd.Err() != nil {
		//logic.ErrNoUser.WithError(cmd.Err()).Errorln(errRedisOperation)
		return fmt.Errorf("set redis err:%v", cmd.Err())
	}
	return nil
}

// 格式化Redis Key
func fmtPlayerKey(uid uint64) string {
	return fmt.Sprintf("gold_%v", uid)
}
