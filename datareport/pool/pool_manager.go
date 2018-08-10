package pool

import (
	"steve/datareport/pool/abs"
	"steve/datareport/pool/impl"
	"github.com/spf13/viper"
)

type PoolManager struct {
	Pool abs.TaskPool
}

var mgr *PoolManager

func init() {
	//runtime.GOMAXPROCS(4) //先不设置
	maxChannelNum := viper.GetInt("max_channel_num")
	maxTaskQueueSize := viper.GetInt("max_task_queue_size")
	waringNum := viper.GetInt("waring_num")
	pool := impl.NewChannelTaskPool("weqTest", maxChannelNum, maxTaskQueueSize, waringNum)
	mgr = &PoolManager{
		Pool:pool,
	}
	pool.Start()
}

func GetTaskPool() *PoolManager {
	return mgr
}
