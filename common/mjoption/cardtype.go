package mjoption

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

// FanType 番型
type FanType struct {
	ID     int   `yaml:"id"`      // 番型 ID
	FuncID int   `yaml:"func_id"` // 计算函数 ID
	Mutex  []int `yaml:"mutex"`   // 互斥番型列表
	Method int   `yaml:"method"`  // 分数计算方式
	Score  int   `yaml:"score"`   // 番数
}

// CardTypeOption 牌型选项
type CardTypeOption struct {
	ID       int             `yaml:"id"`               // 选项 ID
	Fantypes map[int]FanType `yaml:"enable_fan_types"` // 支持的番型
}

// CardTypeOptionManager 选项管理器
type CardTypeOptionManager struct {
	cardTypeOptionMap map[int]*CardTypeOption
}

// GetCardTypeOption 获取牌型选项
func (som *CardTypeOptionManager) GetCardTypeOption(optID int) *CardTypeOption {
	if opt, ok := som.cardTypeOptionMap[optID]; ok {
		return opt
	}
	return nil
}

func (som *CardTypeOptionManager) loadOption(path string) {
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "CardTypeOptionManager.loadOption",
		"path":      path,
	})
	data, err := ioutil.ReadFile(path)
	if err != nil {
		entry.WithError(err).Panicln("读取文件失败")
	}
	opt := CardTypeOption{}
	if err := yaml.Unmarshal(data, &opt); err != nil {
		entry.WithError(err).Panicln("反序列化失败")
	}
	if _, exist := som.cardTypeOptionMap[opt.ID]; exist {
		entry.WithField("id", opt.ID).Panicln("结算选项 ID 重复")
	}
	som.cardTypeOptionMap[opt.ID] = &opt
}

// loadOptions 加载选项文件
func (som *CardTypeOptionManager) loadOptions(optionDir string) {
	som.cardTypeOptionMap = make(map[int]*CardTypeOption)
	filepath.Walk(optionDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			som.loadOption(path)
		}
		return nil
	})
}

// NewCardTypeOptionManager is CardType option manager creator
func NewCardTypeOptionManager(optDir string) *CardTypeOptionManager {
	som := &CardTypeOptionManager{}
	som.loadOptions(optDir)
	return som
}
