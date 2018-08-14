package service

import (
	"steve/server_pb/data_report"
	"context"
	_ "steve/datareport/conn"
	_ "steve/datareport/pool"
	"steve/datareport/bean"
	"steve/datareport/pool/impl"
	"steve/datareport/pool"
	"time"
	lock2 "steve/common/lock"
	"steve/entity/cache"
	"steve/common/data/redis"
	"strings"
	"steve/datareport/fixed"
)

type ReportService struct {
	/*queue *queue.LogQueue
	connMgr *conn.ReportConnManager
	pool *pool.PoolManager*/
}

var service *ReportService

func init() {
	RunTimeReport(func() {
		lockKey := cache.FmtRedisLockKeyReport()
		lock:= lock2.LockRedis(lockKey)
		if lock {
			cli := redis.GetRedisClient()
			keys := cli.Keys(cache.FmtGameReportKeyGame())
			for _,key := range keys.Val() {
				val:= strings.Split(strings.Split(key,":")[1],"-")
				redisVal := cli.Get(key).Val()
				gameId := val[0]
				level := val[1]
				count := redisVal
				logVal := gameId+"|"+level+"|"+count
				log := bean.CreateLogBean(int32(fixed.LOG_TYPE_GAME_PERSON_NUM),0,0,0,0,logVal)
				task := impl.NewLogReportTask(log)
				pool.GetTaskPool().Pool.Execute(task)
				lock2.UnLockRedis(lockKey)
			}
		}
	},5 * time.Minute)
	service = &ReportService{
	}
}

func GetReportService() *ReportService {
	return service
}

func (service *ReportService) Report(ctx context.Context, request *datareport.ReportRequest) (*datareport.ReportResponse, error) {
	log := pb2Bean(request)
	task := impl.NewLogReportTask(log)
	pool.GetTaskPool().Pool.Execute(task)
	return &datareport.ReportResponse{ErrCode: 0}, nil
}

func pb2Bean(request *datareport.ReportRequest) *bean.LogBean {
	return bean.CreateLogBean(
		request.GetLogType(),
		request.GetProvince(),
		request.GetCity(),
		request.GetChannel(),
		request.GetPlayerId(),
		request.GetValue(),
	)
}