package datareportclient

import (
	"steve/structs"
	"errors"
	"steve/server_pb/data_report"
	"context"
	"steve/datareport/fixed"
)

func DataReport(logType  fixed.LogType, province int, city int, channel int, playerId uint64, value string) (int, error) {
	e := structs.GetGlobalExposer()
	conn, err := e.RPCClient.GetConnectByServerName("datareport")
	if err != nil || conn == nil {
		return 1, errors.New("no connection")
	}
	rpcClient := datareport.NewReportServiceClient(conn)
	response, resErr := rpcClient.Report(context.Background(), &datareport.ReportRequest{
		LogType:  int32(logType),
		Province: int32(province),
		City:     int32(city),
		Channel:  int32(channel),
		PlayerId: playerId,
		Value:    value,
	})
	if resErr != nil {
		return 2, resErr
	}
	return int(response.ErrCode), nil
}
