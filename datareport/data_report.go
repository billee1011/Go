package datareport

import (
	"steve/structs"
	"steve/structs/net"
	"steve/server_pb/data_report"
	"steve/structs/service"
	reportservice "steve/datareport/service"
	)

type dataReport struct {
	e   *structs.Exposer
	dog net.WatchDog
}

// GetService 获取服务接口，被 serviceloader 调用
func GetService() service.Service {
	return new(dataReport)
}

// NewService 创建服务
func NewService() service.Service {
	return new(dataReport)
}

func (d *dataReport) Init(e *structs.Exposer, param ...string) error {
	d.e = e
	rpcServer := e.RPCServer
	err := rpcServer.RegisterService(datareport.RegisterReportServiceServer, reportservice.GetReportService())
	if err != nil {
		return err
	}
	return nil
}

func (d *dataReport) Start() error {
	return nil
}