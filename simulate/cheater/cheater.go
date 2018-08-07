package cheater

import (
	"fmt"
	"net/http"
	"steve/external/configclient"
	configpb "steve/server_pb/config"
	"steve/simulate/config"
	"sync"

	"google.golang.org/grpc"
)

var configGRPCCli *grpc.ClientConn
var configGRPCCliInit sync.Once

// SetPlayerCoin 设置玩家金豆数
func SetPlayerCoin(playerID uint64, coin uint64) error {
	url := fmt.Sprintf("%s/setgold?player_id=%v&gold=%v", config.GetPeipaiURL(), playerID, coin)
	if _, err := http.DefaultClient.Get(url); err != nil {
		return fmt.Errorf("访问设置金币服务失败:%v", err)
	}
	return nil
}

// MockConfigClient mock config client
func MockConfigClient() error {
	configGRPCCliInit.Do(func() {
		addr := config.GetConfigRPCAddr()
		var err error
		configGRPCCli, err = grpc.Dial(addr, grpc.WithInsecure())
		if err != nil {
			fmt.Printf("连接配置服失败:%s", err.Error())
			configGRPCCli = nil
		}
	})
	if configGRPCCli == nil {
		return fmt.Errorf("配置服未连接")
	}
	configclient.ConfigCliGetter = func() (configpb.ConfigClient, error) {
		return configpb.NewConfigClient(configGRPCCli), nil
	}
	return nil
}
