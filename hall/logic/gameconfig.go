package logic

import (
	entityConf "steve/entity/config"
	"steve/external/configclient"
	"steve/hall/charge"

	"github.com/Sirupsen/logrus"
)

// GameConf 游戏配置
var GameConf []entityConf.GameConfig

// LevelConf 场次配置
var LevelConf []entityConf.GameLevelConfig

// InitGameConfig 初始化游戏配置
func InitGameConfig() {
	var err error
	for {
		// 游戏配置
		GameConf, err = configclient.GetGameConfigMap()
		if err == nil {
			break
		}
	}

	for {
		// 场次配置
		LevelConf, err = configclient.GetGameLevelConfigMap()
		if err == nil {
			break
		}
	}
	logrus.Debugf("hall服启动加载的游戏玩法，GameConf:(%v)\n，LevelConf：（%v）", GameConf, LevelConf)

	return
}

// InitChargeConfig 初始化charge配置
func InitChargeConfig() {
	var err error
	for {
		// 游戏配置
		err = charge.LoadItemList()
		if err == nil {
			break
		}
	}

	for {
		// 场次配置
		err = charge.LoadMaxCharge()
		if err == nil {
			break
		}
	}
	logrus.Debugf("hall服启动始化charge配置")

	return
}

// // InitGameConfig 初始化游戏配置
// func InitGameConfig() error {
// 	var err error
// 	// 游戏配置
// 	GameConf, err = configclient.GetGameConfigMap()
// 	if err != nil {
// 		return fmt.Errorf("InitGameConfig 获取游戏配置失败,error:(%v)", err.Error())
// 	}
// 	// 场次配置
// 	LevelConf, err = configclient.GetGameLevelConfigMap()
// 	if err != nil {
// 		return fmt.Errorf("InitGameConfig 获取游戏级别配置失败,error:(%v)", err.Error())
// 	}
// 	return nil
// }
