package conn

import (
	"steve/datareport/fixed"
	"os"
	"math/rand"
	"github.com/Sirupsen/logrus"
)

type ReportConnManager struct {
	connMap map[int]*ReportConn
}

var conManager *ReportConnManager

func GetConManager() *ReportConnManager{
	return conManager
}

//随机获取一个链接
func (manager *ReportConnManager) GetConnection() *ReportConn{
	return manager.getConnection(0)
}

func (manager *ReportConnManager) getConnection(num int) *ReportConn{
	if num >= fixed.MAX_CONN_NUM{
		return nil
	}
	conn := manager.connMap[rand.Intn(fixed.MAX_CONN_NUM)]
	if conn.failCount >= fixed.MAX_CONN_NUM {
		//重新获取链接
		return manager.getConnection(num+1)
	}
	return conn
}

func init() {
	config := GetReportClientConfig()
	logrus.Info("connected to big data ",config)
	conManager = &ReportConnManager{
		connMap:make(map[int]*ReportConn,fixed.MAX_CONN_NUM),
	}
	i := 0
	tryNum := 0
	for ; i < fixed.MAX_CONN_NUM; {
		conn := NewReportConn(i,config)
		suc := conn.Connect()
		if suc {
			conManager.connMap [i] = conn
			i++
			logrus.Info("create connection id=",conn.id)
			continue
		}
		if tryNum >= 100 {
			logrus.Info("retry num > 100 exit")
			os.Exit(1)
		}
		tryNum++
		logrus.Info("connection fail retry i=",i,",retrynum=",tryNum)
	}

}
