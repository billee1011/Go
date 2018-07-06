package loader

import (
	"errors"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/hashicorp/consul/api"
)

var errNewConsulAgent = errors.New("创建 consul agent 失败")
var errRegisterFailed = errors.New("向 consul 注册服务失败")
var errAllocServerID = errors.New("分配服务 ID 失败")
var errNewRedisClient = errors.New("创建 redis 客户端失败")

// registerParams 服务注册参数
type registerParams struct {
	serverName string
	addr       string
	port       int
	consulAddr string // consul 地址
}

// registerServer 注册服务
func registerServer(rp *registerParams) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name":   "registerServer",
		"server_name": rp.serverName,
		"addr":        rp.addr,
		"port":        rp.port,
		"consul_addr": rp.consulAddr,
	})
	if rp.serverName == "" {
		logEntry.Infoln("服务名为空，不注册服务")
		return
	}
	serverID := allocServerIDNew(rp)
	logEntry = logEntry.WithField("server_id", serverID)
	if err := registerToConsul(logEntry, rp.serverName, rp.addr, rp.port, serverID, rp.consulAddr); err != nil {
		logEntry.Panicln(err)
	}
	logEntry.Infoln("注册服务到 consul 完成")
}

// allocServerID 分配服务 ID
func allocServerIDNew(rp *registerParams) string {
	return fmt.Sprintf("%s-%s-%d", rp.serverName, rp.addr, rp.port)
}

// registerToConsul 向 consul 注册服务
func registerToConsul(logEntry *logrus.Entry, serverName string, addr string, port int, serverID string, consulAddr string) error {
	logEntry = logEntry.WithFields(logrus.Fields{
		"func_name":   "registerToConsul",
		"server_name": serverName,
		"addr":        addr,
		"port":        port,
		"server_id":   serverID,
		"consul_addr": consulAddr,
	})
	agent := createConsulAgent(logEntry, consulAddr)
	if agent == nil {
		return errNewConsulAgent
	}

	registration := &api.AgentServiceRegistration{
		ID:      serverID,
		Name:    serverName,
		Tags:    []string{},
		Port:    port,
		Address: addr,
	}
	if err := agent.ServiceRegister(registration); err != nil {
		logEntry.Errorln(err)
		return errRegisterFailed
	}
	return nil
}

// getConsulAgent 获取 consul 代理
func createConsulAgent(logEntry *logrus.Entry, consulAddr string) *api.Agent {
	config := api.DefaultConfig()
	config.Address = consulAddr
	consul, err := api.NewClient(config)
	if err != nil {
		logEntry.WithError(err).Errorln("创建 consul api 客户端失败")
		return nil
	}
	return consul.Agent()
}
