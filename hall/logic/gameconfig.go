package logic

import (
	"context"
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
func InitGameConfig(ctx context.Context) {
	var err error
	for {
		select {
		case <-ctx.Done():
			logrus.Debugf("hall服启动加载的游戏玩法，GameConf:(%v)\n，LevelConf：（%v）", GameConf, LevelConf)
			return
		default:
			GameConf, err = configclient.GetGameConfigMap()
			if err != nil {
				continue
			}
			LevelConf, err = configclient.GetGameLevelConfigMap()
			if err != nil {
				continue
			}
			logrus.Debugf("hall服启动加载的游戏玩法，GameConf:(%v)\n，LevelConf：（%v）", GameConf, LevelConf)
			return
		}
	}
}

// InitChargeConfig 初始化charge配置
func InitChargeConfig(ctx context.Context) {
	var err error
	for {
		select {
		case <-ctx.Done():
			logrus.Debugf("hall服启动始化charge配置失败")
			return
		default:
			err = charge.LoadItemList()
			if err != nil {
				continue
			}
			err = charge.LoadMaxCharge()
			if err != nil {
				continue
			}
			logrus.Debugf("hall服启动始化charge配置成功")
			return
		}
	}
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
