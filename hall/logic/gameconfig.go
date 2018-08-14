package logic

import (
	"fmt"
	entityConf "steve/entity/config"
	"steve/external/configclient"
)

// GameConf 游戏配置
var GameConf []entityConf.GameConfig

// LevelConf 场次配置
var LevelConf []entityConf.GameLevelConfig

// InitGameConfig 初始化游戏配置
func InitGameConfig() error {
	var err error
	// 游戏配置
	GameConf, err = configclient.GetGameConfigMap()
	if err != nil {
		return fmt.Errorf("InitGameConfig 获取游戏配置失败,error:(%v)", err.Error())
	}
	// 场次配置
	LevelConf, err = configclient.GetGameLevelConfigMap()
	if err != nil {
		return fmt.Errorf("InitGameConfig 获取游戏级别配置失败,error:(%v)", err.Error())
	}
	return nil
}
