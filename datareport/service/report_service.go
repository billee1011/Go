package service

import (
	"steve/server_pb/data_report"
	"context"
	_ "steve/datareport/conn"
	_ "steve/datareport/pool"
	"steve/datareport/bean"
	"steve/datareport/pool/impl"
	"steve/datareport/pool"
)

type ReportService struct {
	/*queue *queue.LogQueue
	connMgr *conn.ReportConnManager
	pool *pool.PoolManager*/
}

var service *ReportService

func init() {
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