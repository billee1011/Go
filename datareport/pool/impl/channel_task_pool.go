package impl

import (
	"time"
	"sync/atomic"
	"steve/datareport/pool/abs"
	"github.com/Sirupsen/logrus"
)

/**
	基于信道
 */

type channelTaskPool struct {
	Name               string
	MaxChannelNum      int
	MaxTaskQueueSize   int
	WaringNum          int
	taskChannel        chan abs.Task
	waitingStopChannel chan bool
	workNum            int32
	isClose            bool
}

func NewChannelTaskPool(name string, maxChannelNum int, MaxTaskQueueSize int, WaringNum int) abs.TaskPool {
	var instance abs.TaskPool = &channelTaskPool{
		Name:             name,
		MaxChannelNum:    maxChannelNum,
		MaxTaskQueueSize: MaxTaskQueueSize,
		WaringNum:        WaringNum,
	}
	return instance
}

func (pool *channelTaskPool) GetName() string {
	return pool.Name
}

func (pool *channelTaskPool) Start() {
	pool.taskChannel = make(chan abs.Task, pool.MaxTaskQueueSize)
	pool.waitingStopChannel = make(chan bool)
	for i := 0; i < pool.MaxChannelNum; i++ {
		go pool.run()
	}
	time.Sleep(2000 * time.Millisecond)
	logrus.Info("task pool start success total num ", pool.workNum)
}

func (pool *channelTaskPool) Stop() {
	pool.isClose = true
	pool.waitingStopChannel <- true
}

func (pool *channelTaskPool) StopNow() {
	pool.isClose = true
	close(pool.taskChannel)
}

func (pool *channelTaskPool) Execute(task abs.Task) {
	if pool.isClose {
		println("close pool send task")
		return
	}
	pool.taskChannel <- task
	queueLen := len(pool.taskChannel)
	if queueLen >= pool.WaringNum {
		logrus.Error("waring !!!!  task size = ", queueLen)
	}
}

func (pool *channelTaskPool) OnTaskError(task abs.Task, err error) {
	logrus.Error("task failed Task:", task, "  error:", err)
}

func (pool *channelTaskPool) run() {
	isClose := false
	pool.workNum = atomic.AddInt32(&pool.workNum, 1)
	myId := pool.workNum
	//println("start desk worker currNum=", pool.workNum)
	for !isClose {
		//println(myId, " wait for task")
		select {
		case task, isOpen := <-pool.taskChannel:
			if isOpen {
				logrus.Info(myId, " got desk task")
				err := task.DoTask()
				if err != nil {
					pool.OnTaskError(task, err)
				}
			} else {
				logrus.Info("task channel is close")
				isClose = true
			}

		case closes := <-pool.waitingStopChannel:
			if closes {
				//启动检测任务完成
				go pool.checkOver()
			}
		}
	}
	logrus.Info(pool.GetName(), myId, " is stop")
	pool.workNum = atomic.AddInt32(&pool.workNum, -1)
	if pool.workNum == 0 {
		logrus.Info(pool.GetName(), " is shutdown")
	}
}

func (pool *channelTaskPool) checkOver() {
	logrus.Info("start wait task over")
	for true {
		if len(pool.taskChannel) <= 0 {
			logrus.Info("task all over")
			break
		}
		time.Sleep(1 * time.Second)
	}
	close(pool.taskChannel)
	pool.isClose = true
}
