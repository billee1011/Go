package mjoption

import (
	"io/ioutil"

	"github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// GameOptions 游戏选项
type GameOptions struct {
	GameID           int `yaml:"game_id"`         // 游戏 ID
	SettleOptionID   int `yaml:"settle_option"`   // 结算选项 ID
	CardTypeOptionID int `yaml:"cardtype_option"` // 牌型选项 ID
	XingPaiOptionID  int `yaml:"xingpai_option"`  // 行牌选项 ID
}

// GameOptionManager 游戏选项管理器
type GameOptionManager struct {
	gameOptionMap map[int]*GameOptions
}

// GetGameOptions 获取游戏选项
func (gom *GameOptionManager) GetGameOptions(gameID int) *GameOptions {
	if opt, exist := gom.gameOptionMap[gameID]; exist {
		return opt
	}
	return nil
}

func (gom *GameOptionManager) loadOption(path string) {
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "GameOptionManager.loadOption",
		"path":      path,
	})
	data, err := ioutil.ReadFile(path)
	if err != nil {
		entry.WithError(err).Panicln("读取文件失败")
	}
	options := []*GameOptions{}
	if err := yaml.Unmarshal(data, &options); err != nil {
		entry.WithError(err).Panicln("反序列化失败")
	}
	for _, opt := range options {
		gom.gameOptionMap[opt.GameID] = opt
	}
}

// NewGameOptionManager new game option manager.
// path 为配置选项配置文件路径
func NewGameOptionManager(path string) *GameOptionManager {
	gom := &GameOptionManager{
		gameOptionMap: make(map[int]*GameOptions),
	}
	gom.loadOption(path)
	return gom
}
