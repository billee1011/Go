package loader

import (
	"errors"
	"steve/structs/redisfactory"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/hashicorp/consul/api"
)

var errNewConsulAgent = errors.New("创建 consul agent 失败")
var errRegisterFailed = errors.New("向 consul 注册服务失败")
var errAllocServerID = errors.New("分配服务 ID 失败")
var errNewRedisClient = errors.New("创建 redis 客户端失败")

// registerParams 服务注册参数
type registerParams struct {
	serverName   string
	addr         string
	port         int
	redisFactory redisfactory.RedisFactory
}

// registerServer 注册服务
func registerServer(rp *registerParams) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name":   "registerServer",
		"server_name": rp.serverName,
		"addr":        rp.addr,
		"port":        rp.port,
	})

	serverID, err := allocServerID(logEntry, rp.redisFactory)
	if err != nil {
		logEntry.Panicln(err)
	}
	if err := registerToConsul(logEntry, rp.serverName, rp.addr, rp.port, serverID); err != nil {
		logEntry.Panicln(err)
	}
}

// registerToConsul 向 consul 注册服务
func registerToConsul(logEntry *logrus.Entry, serverName string, addr string, port int, serverID string) error {
	logEntry = logEntry.WithFields(logrus.Fields{
		"func_name":   "registerToConsul",
		"server_name": serverName,
		"addr":        addr,
		"port":        port,
		"server_id":   serverID,
	})
	agent := createConsulAgent(logEntry)
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

// allocServerID 分配服务 ID
func allocServerID(logEntry *logrus.Entry, redisFactory redisfactory.RedisFactory) (string, error) {
	redisCli, err := redisFactory.NewClient()
	if err != nil {
		logEntry.WithError(err).Errorln(errNewRedisClient)
		return "", errNewRedisClient
	}
	result := redisCli.Incr("serviceloader:service:maxid")
	if result.Err() != nil {
		logEntry.WithError(result.Err()).Errorln(errAllocServerID)
		return "", errAllocServerID
	}
	return strconv.FormatInt(result.Val(), 10), nil
}

// getConsulAgent 获取 consul 代理
func createConsulAgent(logEntry *logrus.Entry) *api.Agent {
	config := api.DefaultConfig()
	consul, err := api.NewClient(config)
	if err != nil {
		logEntry.WithError(err).Errorln("创建 consul api 客户端失败")
		return nil
	}
	return consul.Agent()
}
