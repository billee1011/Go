package queue

import "steve/datareport/bean"

type LogQueue struct{
	logQueue chan *bean.LogBean
}

var queue *LogQueue

func init(){
	queue = &LogQueue{
		logQueue:make(chan *bean.LogBean,5 * 10000),
	}
}

func GetLogQueue() *LogQueue{
	return queue
}

func (queue *LogQueue) Take() *bean.LogBean{
	log,isClose := <- queue.logQueue
	if isClose{
		return nil
	}
	return log
}

func (queue *LogQueue) Put(logBean *bean.LogBean) {
	queue.logQueue <- logBean
}