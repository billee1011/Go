package conn

import (
	"net"
	"strconv"
	"github.com/Sirupsen/logrus"
)

type ReportConn struct {
	id int
	config     Config
	connection *net.TCPConn
	failCount  int //连续失败次数
}

func NewReportConn(id int,config Config) *ReportConn {
	return &ReportConn{
		id:id,
		config: config,
	}
}

func (handle *ReportConn) Sender(reportContent string) bool {
	_, error := handle.connection.Write([]byte(reportContent))
	if error != nil {
		handle.failCount++
		logrus.Info("send fail id=",handle.id)
		return false
	}
	if handle.failCount > 0 {
		handle.failCount = 0
	}
	logrus.Info("send suc value=",reportContent)
	return true
}

func (handle *ReportConn) Close() {
	handle.connection.Close()
}

func (handle *ReportConn) Connect() bool {
	address := handle.config.Address + ":" + strconv.Itoa(handle.config.Port)
	tcpAdd, err := net.ResolveTCPAddr("tcp4", address)
	if err != nil {
		logrus.Info("get tcpAddr error", err.Error())
		return false
	}
	conn, err := net.DialTCP("tcp", nil, tcpAdd)
	if err != nil {
		logrus.Info("connection remote address error ", err.Error())
		return false
	}
	conn.SetKeepAlive(false)
	conn.SetNoDelay(true)
	handle.connection = conn
	return true
}
