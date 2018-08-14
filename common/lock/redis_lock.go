package lock

//reids锁
import (
	"steve/common/data/redis"
	"time"
)

func LockRedis(key string) bool {
	cli := redis.GetRedisClient()
	cmd := cli.SetNX(key, 1, 1*time.Minute)
	return cmd.Val()
}

func UnLockRedis(key string) bool {
	cli := redis.GetRedisClient()
	script := "if redis.call('get', KEYS[1]) == ARGV[1] then return redis.call('del', KEYS[1]) else return 0 end"
	cmd := cli.Eval(script, []string{key}, 1)
	return cmd.Val() == 1
}
