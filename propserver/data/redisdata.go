package data

import (
	"fmt"
	"steve/common/data/redis"
	"strconv"
	"time"
	"github.com/Sirupsen/logrus"
)

/*
	功能： 服务数据保存到redis.必须对Redis设置过期时间
	作者： SkyWang
	日期： 2018-7-25

*/

// redis 过期时间
var redisTimeOut time.Duration = time.Minute * 60 * 24 * 30

// 格式化Redis Key
func fmtPlayerKey(uid uint64) string {
	return fmt.Sprintf("props_%v", uid)
}


// 从redis加载玩家道具
func LoadPropsFromRedis(uid uint64) (map[uint64]int64, error) {

	r := redis.GetRedisClient()
	key := fmtPlayerKey(uid)
	cmd := r.HGetAll(key)
	if cmd.Err() != nil {
		//logic.ErrNoUser.WithError(cmd.Err()).Errorln(errRedisOperation)
		return nil, fmt.Errorf("LoadPropsFromRedis redis err:%v", cmd.Err())
	}
	m := cmd.Val()
	if len(m) == 0 {
		return nil,  fmt.Errorf("LoadPropsFromRedis no user: uid=%d", uid)
	}
	list := make(map[uint64]int64, len(m))
	for k, v := range m {
		t, err := strconv.ParseInt(k, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("LoadPropsFromRedis ret err1:%v", m)
		}
		if t == 0 {
			logrus.Errorf("prop id = 0 err: uid=%d", uid)
			continue
		}
		g, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("LoadPropsFromRedis ret err2:%v", m)
		}
		list[uint64(t)] = g
	}

	logrus.Debugf("LoadPropsFromRedis win: uid=%d, propslist=%v ", uid, list)
	return list, nil
}

// 保存玩家道具到Redis
func SavePropsToRedis(uid uint64, propsList map[uint64]int64) error {
	// 没有道具，不需要写Redis
	if len(propsList) == 0 {
		return nil
	}

	r := redis.GetRedisClient()
	key := fmtPlayerKey(uid)

	list := make(map[string]interface{}, len(propsList))
	for k, v := range propsList {
		strKey := fmt.Sprintf("%d", k)
		list[strKey] = v
	}
	cmd := r.HMSet(key, list)
	if cmd.Err() != nil {
		//logic.ErrNoUser.WithError(cmd.Err()).Errorln(errRedisOperation)
		logrus.Errorf("SavePropsToRedis err:key=%s,err=%s", key, cmd.Err())
		return fmt.Errorf("set redis err:%v", cmd.Err())
	}
	r.Expire(key, redisTimeOut)
	return nil
}

