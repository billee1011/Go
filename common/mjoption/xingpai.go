package mjoption

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

// XingPaiOption 行牌选项
type XingPaiOption struct {
	ID         int  `yaml:"id"`          // 选项 ID
	HuGameOver bool `yaml:"hu_gameover"` // 胡牌后是否触发游戏结束
}

// XingPaiOptionManager 选项管理器
type XingPaiOptionManager struct {
	xingPaiOptionMap map[int]*XingPaiOption
}

// GetXingPaiOption 获取行牌选项
func (som *XingPaiOptionManager) GetXingPaiOption(optID int) *XingPaiOption {
	if opt, ok := som.xingPaiOptionMap[optID]; ok {
		return opt
	}
	return nil
}

func (som *XingPaiOptionManager) loadOption(path string) {
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "XingPaiOptionManager.loadOption",
		"path":      path,
	})
	data, err := ioutil.ReadFile(path)
	if err != nil {
		entry.WithError(err).Panicln("读取文件失败")
	}
	opt := XingPaiOption{}
	if err := yaml.Unmarshal(data, &opt); err != nil {
		entry.WithError(err).Panicln("反序列化失败")
	}
	if _, exist := som.xingPaiOptionMap[opt.ID]; exist {
		entry.WithField("id", opt.ID).Panicln("结算选项 ID 重复")
	}
	som.xingPaiOptionMap[opt.ID] = &opt
}

// loadOptions 加载选项文件
func (som *XingPaiOptionManager) loadOptions(optionDir string) {
	som.xingPaiOptionMap = make(map[int]*XingPaiOption)
	filepath.Walk(optionDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			som.loadOption(path)
		}
		return nil
	})
}

// NewXingPaiOptionManager is XingPai option manager creator
func NewXingPaiOptionManager(optDir string) *XingPaiOptionManager {
	som := &XingPaiOptionManager{}
	som.loadOptions(optDir)
	return som
}
