package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"steve/gutils"
	hs "steve/room/hszswitch/hszswitch"
	"strconv"

	"github.com/hashicorp/consul/api"

	"google.golang.org/grpc"

	"github.com/Sirupsen/logrus"

	"github.com/spf13/viper"
)

type server struct {
}

func (s *server) GetHSZSwitch(context context.Context, dn *hs.DoNothing) (*hs.HszSwitch, error) {
	return &hs.HszSwitch{
		Hsz: Open,
	}, nil
}

const (
	// HszSwitchAddr 代表监听客户端的IP地址，默认值为 127.0.0.1
	HszSwitchAddr = "hsz_switch_addr"
	// HszSwitch 换三张开关关键字
	HszSwitch = "hszswitch"
	// HszSwitchSerIP 换三张grpc服务地ip
	HszSwitchSerIP = "hsz_switch_ser_ip"
	// HszSwitchSerPort 换三张grpc服务地端口
	HszSwitchSerPort = "hsz_switch_ser_port"
	// ConfigName 配置文件名字
	ConfigName = "config"
)

//Open 开关
var Open = true

func handle(resp http.ResponseWriter, req *http.Request) {
	value := req.FormValue(HszSwitch)
	if len(value) == 0 {
		respMSG(resp, fmt.Sprintf("开关关键字switch有误"), 404)
		return
	}
	open, err := strconv.ParseBool(value)
	if err != nil {
		respMSG(resp, fmt.Sprintf("switch对应的值有误:%v", err), 404)
		return
	}
	Open = open
	respMSG(resp, fmt.Sprintf("配置换三张开关成功,当前为:%v", Open), 200)
}

func respMSG(resp http.ResponseWriter, message string, code int) {
	resp.WriteHeader(code)
	resp.Write([]byte(message))
	switch code {
	case 200:
		logrus.Infoln(message)
	default:
		logrus.Debugln(message)
	}
}

func init() {
	initDefaultConfig()
}

func initDefaultConfig() {
	viper.SetDefault(HszSwitchAddr, "127.0.0.1:8081")
	viper.SetDefault(HszSwitchSerIP, "127.0.0.1")
	viper.SetDefault(HszSwitchSerPort, 8082)
	viper.SetConfigName(ConfigName)
	viper.AddConfigPath("./")
}

func hszGrpc() {
	s := grpc.NewServer()
	hs.RegisterSwitchHandlerServer(s, &server{})
	viper.ReadInConfig()
	listenIP := viper.GetString(HszSwitchSerIP)
	listenPort := viper.GetInt(HszSwitchSerPort)
	listenAddr := fmt.Sprintf("%v:%v", listenIP, listenPort)
	logrus.WithFields(
		logrus.Fields{
			"listenAddr": listenAddr,
		}).Infoln("启动换三张开关grpc服务")
	registerToConsul(listenIP, listenPort)
	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		logrus.Fatalf("failed to listen hszGrpcServer:%v", err)
	}
	if err := s.Serve(lis); err != nil {
		logrus.Fatalf("failed to start hszGrpcServer:%v", err)
	}
}

func registerToConsul(ip string, port int) error {
	consul, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return err
	}
	agent := consul.Agent()
	if err = agent.ServiceRegister(&api.AgentServiceRegistration{
		ID:      generateServiceID(gutils.XuezhanOptionService, ip, port),
		Name:    gutils.XuezhanOptionService,
		Port:    port,
		Address: ip,
	}); err != nil {
		return err
	}
	return nil
}

func generateServiceID(serviceName string, ip string, port int) string {
	return fmt.Sprintf("%v_%v:%v", serviceName, ip, port)
}

func main() {
	go hszGrpc()
	http.HandleFunc("/", handle)
	viper.ReadInConfig()
	listenAddr := viper.GetString(HszSwitchAddr)
	logrus.WithFields(
		logrus.Fields{
			"listenAddr": listenAddr,
		}).Infoln("启动换三张开关服务器")
	if err := http.ListenAndServe(listenAddr, nil); err != nil {
		logrus.Debugln(fmt.Sprintf("启动服务器失败:%v", err))
	}
}
