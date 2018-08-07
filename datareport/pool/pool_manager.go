package pool

import (
	"steve/datareport/pool/abs"
	"steve/datareport/pool/impl"
)

type PoolManager struct {
	Pool abs.TaskPool
}

var mgr *PoolManager

func init() {
	//runtime.GOMAXPROCS(4) //先不设置
	pool := impl.NewChannelTaskPool("weqTest", 100, 5000, 2000)
	mgr = &PoolManager{
		Pool:pool,
	}
	pool.Start()
}

func GetTaskPool() *PoolManager {
	return mgr
}
