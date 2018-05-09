package consul

import (
	"fmt"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/go-redis/redis"
	consulapi "github.com/hashicorp/consul/api"
)

func allocServiceID(redisClient *redis.Client) (string, error) {
	result := redisClient.Incr("serviceloader:service:maxid")
	if result.Err() != nil {
		return "", fmt.Errorf("分配服务 ID 失败: %v", result.Err())
	}
	return strconv.FormatInt(result.Val(), 10), nil
}

func registerService(serviceName string, addr string, port int, redisClient *redis.Client) error {
	entry := logrus.WithFields(logrus.Fields{
		"service_name": serviceName,
		"address":      addr,
		"port":         port,
	})

	if serviceName == "" {
		entry.Warn("未配置服务名称，该服务不会被其他服务发现")
		return nil
	}
	serviceID, err := allocServiceID(redisClient)
	if err != nil {
		entry.WithError(err).Errorln("分配服务 ID 失败")
		return err
	}

	agent := gConsulClient.Agent()
	if err := agent.ServiceRegister(&consulapi.AgentServiceRegistration{
		ID:      serviceID,
		Name:    serviceName,
		Tags:    []string{},
		Port:    port,
		Address: addr,
		// TODO: Checks
	}); err != nil {
		return fmt.Errorf("服务注册失败(%s, %s, %d)：%v", serviceName, addr, port, err)
	}
	return nil
}
