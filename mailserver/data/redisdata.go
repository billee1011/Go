package data

import (
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"steve/common/data/redis"
	"steve/mailserver/define"
	"time"
)

/*
	功能： 服务数据保存到redis.必须对Redis设置过期时间
	作者： SkyWang
	日期： 2018-7-25
*/
// redis 过期时间
var redisTimeOut time.Duration = time.Minute * 60 * 24 * 7

// 格式化Redis Key
func fmtPlayerKey(uid uint64) string {
	return fmt.Sprintf("mail_%v", uid)
}

// 从redis加载玩家邮件列表
func LoadUserMailListFromRedis(uid uint64) (map[uint64]*define.PlayerMail, error) {

	r := redis.GetRedisClient()
	key := fmtPlayerKey(uid)
	cmd := r.HGetAll(key)
	if cmd.Err() != nil {
		//logic.ErrNoUser.WithError(cmd.Err()).Errorln(errRedisOperation)
		return nil, fmt.Errorf("LoadUserMailListFromRedis err:%v", cmd.Err())
	}
	m := cmd.Val()
	if len(m) == 0 {
		return nil, fmt.Errorf(" LoadUserMailListFromRedis no user: uid=%d", uid)
	}
	list := make(map[uint64]*define.PlayerMail, len(m))
	for k, v := range m {
		if k == "0" {
			continue
		}
		newMail := new(define.PlayerMail)
		err := json.Unmarshal([]byte(v), newMail)
		if err != nil {
			continue
		}
		if newMail.MailId == 0 {
			continue
		}

		list[newMail.MailId] = newMail
	}

	return list, nil
}

// 从redis加载玩家指定邮件
func LoadTheMailFromRedis(uid uint64, mailId uint64) (*define.PlayerMail, error) {

	r := redis.GetRedisClient()
	key := fmtPlayerKey(uid)
	cmd := r.HGet(key, fmt.Sprintf("%d", mailId))
	if cmd.Err() != nil {

		return nil, fmt.Errorf("LoadTheMailFromRedis err:%v", cmd.Err())
	}
	strJson := cmd.Val()
	if len(strJson) == 0 {
		return nil, fmt.Errorf(" LoadTheMailFromRedis no user: uid=%d", uid)
	}

	newMail := new(define.PlayerMail)
	err := json.Unmarshal([]byte(strJson), newMail)
	if err != nil {
		return nil, fmt.Errorf(" LoadTheMailFromRedis json parse err: uid=%d,err=%v", uid, err)
	}

	return newMail, nil
}

// 保存玩家邮件列表到Redis
func SaveUserMailListToRedis(uid uint64, mailList map[uint64]*define.PlayerMail) error {
	r := redis.GetRedisClient()
	key := fmtPlayerKey(uid)

	list := make(map[string]interface{}, len(mailList))
	if len(mailList) == 0 {
		// 如果无邮件，保存一个空记录到redis，表示玩家无邮件
		list["0"] = ""
	} else {
		for k, v := range mailList {
			data, err := json.Marshal(v)
			if err != nil {
				continue
			}

			strKey := fmt.Sprintf("%d", k)
			list[strKey] = data
		}
	}

	cmd := r.HMSet(key, list)
	if cmd.Err() != nil {
		logrus.Errorf("SaveUserMailListToRedis err:key=%s,err=%s", key, cmd.Err())
		return fmt.Errorf("SaveUserMailListToRedis err:%v", cmd.Err())
	}
	r.Expire(key, redisTimeOut)
	return nil
}

// 保存玩家指定邮件到Redis
func SaveTheMailToRedis(uid uint64, mail *define.PlayerMail) error {
	r := redis.GetRedisClient()
	key := fmtPlayerKey(uid)

	data, err := json.Marshal(mail)
	if err != nil {
		return err
	}

	cmd := r.HSet(key, fmt.Sprintf("%d", mail.MailId), data)
	if cmd.Err() != nil {
		logrus.Errorf("SaveTheMailToRedis err:key=%s,err=%s", key, cmd.Err())
		return fmt.Errorf("SaveTheMailToRedis err:%v", cmd.Err())
	}
	r.Expire(key, redisTimeOut)
	return nil
}
