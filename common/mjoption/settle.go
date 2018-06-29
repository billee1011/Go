package mjoption

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// SettleOption 结算选项
type SettleOption struct {
	ID int `yaml:"id"` // 选项 ID
	// --杠结算
	AnGangValue   uint32 `yaml:"gang_angang_value"`   // 暗杠倍数
	MingGangValue uint32 `yaml:"gang_minggang_value"` // 明杠倍数
	BuGangValue   uint32 `yaml:"gang_bugang_value"`   // 补杠倍数
	// --胡结算
	HuValue map[string]uint32 `yaml:"hu_value"` // 胡牌倍数
	// --单局结算
	EnableChaHuaZhu bool `yaml:"enable_tenable_chahuazhuuisui"` // 是否开启查花猪
	EnableChaDaJiao bool `yaml:"enable_chadajiao"`              // 是否开启查大叫
	EnableTuisui    bool `yaml:"enable_tuisui"`                 // 是否开启退税

	// --其他结算
	HuPlayerCanSettle     map[string]bool `yaml:"huPlayer_can_settle"`     // 已胡牌玩家可参与的结算
	HuQuitPlayerCanSettle map[string]bool `yaml:"huQuitPlayer_can_settle"` // 已胡牌玩家(退出)可参与的结算
	GiveUpPlayerCanSettle map[string]bool `yaml:"giveUpPlayer_can_settle"` // 认输玩家可参与的结算

	GangInstantSettle bool `yaml:"gang_instant_settle"` // 杠是否可以立即结算
	HuInstantSettle   bool `yaml:"hu_instant_settle"`   // 胡是否可以立即结算

}

// SettleOptionManager 结算选项管理器
type SettleOptionManager struct {
	settleOptionMap map[int]*SettleOption
}

// GetSettleOption 获取结算选项
func (som *SettleOptionManager) GetSettleOption(optID int) *SettleOption {
	if opt, ok := som.settleOptionMap[optID]; ok {
		return opt
	}
	return nil
}

func (som *SettleOptionManager) loadOption(path string) {
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "SettleOptionManager.loadOption",
		"path":      path,
	})
	data, err := ioutil.ReadFile(path)
	if err != nil {
		entry.WithError(err).Panicln("读取文件失败")
	}
	opt := SettleOption{}
	if err := yaml.Unmarshal(data, &opt); err != nil {
		entry.WithError(err).Panicln("反序列化失败")
	}
	if _, exist := som.settleOptionMap[opt.ID]; exist {
		entry.WithField("id", opt.ID).Panicln("结算选项 ID 重复")
	}
	som.settleOptionMap[opt.ID] = &opt
}

// loadOptions 加载选项文件
func (som *SettleOptionManager) loadOptions(optionDir string) {
	som.settleOptionMap = make(map[int]*SettleOption)
	filepath.Walk(optionDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			som.loadOption(path)
		}
		return nil
	})
}

// NewSettleOptionManager is settle option manager creator
func NewSettleOptionManager(optDir string) *SettleOptionManager {
	som := &SettleOptionManager{}
	som.loadOptions(optDir)
	return som
}
