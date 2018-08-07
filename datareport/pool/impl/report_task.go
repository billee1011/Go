package impl

import (
	"fmt"
	"steve/datareport/bean"
	"steve/datareport/conn"
	"github.com/Sirupsen/logrus"
)

type LogReportTask struct {
	log *bean.LogBean
}

func NewLogReportTask(log *bean.LogBean) LogReportTask{
	return LogReportTask{
		log:log,
	}
}

func (task LogReportTask) DoTask() error {
	var err error = nil

	defer func() {
		if errs := recover(); errs != nil {
			err = fmt.Errorf("task error")
		}
	}()

	sendValue := task.log.ToReportFormat()
	connection := conn.GetConManager().GetConnection()
	if connection == nil{
		logrus.Error("not get bigdata connection log["+sendValue+"]")
		panic("not get bigdata connection")
	}
	result := connection.Sender(sendValue)
	if !result{
		task.DoTask() //提交失败重试
	}
	return err
}
