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
	ID            int            `yaml:"id"`             // 选项 ID
	WallCards     []int          `yaml:"wall_cards"`     //墙牌
	HuGameOver    bool           `yaml:"hu_gameover"`    // 胡牌后是否触发游戏结束
	Hnz           Hnz            `yaml:"hnz"`            //是否需要换N张
	NeedDingque   bool           `yaml:"need_dingque"`   //是否需要定缺
	NeedAddflower bool           `yaml:"need_addflower"` //是否需要补花
	NeedChi       bool           `yaml:"need_chi"`       //是否可以吃牌
	PlayerStates  []XingpaiState `yaml:"player_states"`  //玩家不参与游戏的状态
	// Zhuang        int            `yaml:"zhuang"`         //选庄的方式
}

// XingpaiState 行牌状态
type XingpaiState uint32

const (
	//Run 正常行牌状态
	Run XingpaiState = 0
	//Hu 胡牌状态
	Hu XingpaiState = 1
	//Giveup 认输状态
	Giveup XingpaiState = 2
)

// Hnz 换n张
type Hnz struct {
	Need bool `yaml:"need"` //是否需要换n张
	Num  int  `yaml:"num"`  //换的张数
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
