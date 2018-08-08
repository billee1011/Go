package impl

import (
	"fmt"
	"steve/datareport/bean"
	"steve/datareport/conn"
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
	conn.GetConManager().GetConnection().Sender(sendValue)
	return err
}
