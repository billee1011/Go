package loader

import (
	"errors"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/hashicorp/consul/api"
	"net/http"
)

var errNewConsulAgent = errors.New("创建 consul agent 失败")
var errRegisterFailed = errors.New("向 consul 注册服务失败")
//var errAllocServerID = errors.New("分配服务 ID 失败")
//var errNewRedisClient = errors.New("创建 redis 客户端失败")

// 创建时保存的consul地址
var consulAddress string
// 创建时保存的服务Id
var  svrId string

// registerParams 服务注册参数
type RegisterParams struct {
	serverName string
	addr       string
	port       int
	healthPort int // consul服务健康检查Port
	consulAddr string // consul 地址

}

func RegisterServer2(opt *Option) {
	RegisterServer(&RegisterParams{
		serverName: opt.rpcServerName,
		addr:       opt.rpcAddr,
		port:       opt.rpcPort,
		consulAddr: opt.consulAddr,
		healthPort: opt.healthPort,
	})
}

// registerServer 注册服务
func RegisterServer(rp *RegisterParams) {
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
	if err := registerToConsul(logEntry, rp.serverName, rp.addr, rp.port, serverID, rp.consulAddr, rp.healthPort); err != nil {
		logEntry.Panicln(err)
	}
	logEntry.Infoln("注册服务到 consul 完成")
}

// allocServerID 分配服务 ID
func allocServerIDNew(rp *RegisterParams) string {
	return fmt.Sprintf("%s-%s-%d", rp.serverName, rp.addr, rp.port)
}

// consul对服务进行健康检查
func statusHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Println("check status.")
	fmt.Fprint(w, "status ok!")

}
// consul对服务进行健康检查,通过Http提供检查接口
func startHttpHealth(httPort int) error{
	http.HandleFunc("/status", statusHandler)
	fmt.Println("start listen...")
	addr := fmt.Sprintf(":%d", httPort)
	err := http.ListenAndServe(addr, nil)
	return err
}

// registerToConsul 向 consul 注册服务
func registerToConsul(logEntry *logrus.Entry, serverName string, addr string, port int, serverID string, consulAddr string,healthPort int) error {
	logEntry = logEntry.WithFields(logrus.Fields{
		"func_name":   "registerToConsul",
		"server_name": serverName,
		"addr":        addr,
		"port":        port,
		"server_id":   serverID,
		"consul_addr": consulAddr,
		"health_port": healthPort,
	})
	healthAddr := fmt.Sprintf("%s:%d", addr,  healthPort)
	if healthPort > 0 {
		go startHttpHealth(healthPort)
	}

	agent := createConsulAgent(logEntry, consulAddr)
	if agent == nil {
		return errNewConsulAgent
	}

	// consul对服务进行健康检查
	var ck *api.AgentServiceCheck = nil
	if healthPort > 0 {
		ck = &api.AgentServiceCheck{
			HTTP:     "http://" + healthAddr + "/status",
			Interval: "5s",			// 检查间隔
			Timeout:  "2s",			// 响应超时时间
			DeregisterCriticalServiceAfter: "20s",	// 注销节点超时时间
		}
	}

	svrId = serverID

	registration := &api.AgentServiceRegistration{
		ID:      serverID,
		Name:    serverName,
		Tags:    []string{},
		Port:    port,
		Address: addr,
		Check: ck,
	}
	if err := agent.ServiceRegister(registration); err != nil {
		logEntry.Errorln(err)
		return errRegisterFailed
	}

	//DeleteConsulAgent("match-127.0.0.1-37001")
	//DeleteConsulAgent("match-192.168.8.17-37001")
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
	consulAddress = consulAddr
	return consul.Agent()
}


// getConsulAgent 获取 consul 代理
func DeleteMyConsulAgent() error {

	return DeleteConsulAgent(svrId)
}

// getConsulAgent 获取 consul 代理
func DeleteConsulAgent(sid string) error {
	if len(consulAddress) == 0 {
		return nil
	}
	config := api.DefaultConfig()
	config.Address = consulAddress
	consul, err := api.NewClient(config)
	if err != nil {
		return errors.New("consul connect failed: " + consulAddress)
	}
	if consul == nil {
		return errors.New("consul  = nil: " + consulAddress)
	}
	if err := consul.Agent().ServiceDeregister(sid); err != nil {
		return errors.New("deleteConsulAgent failed: " + svrId)
	}
	return nil
}
