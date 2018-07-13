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
	// --单局结算
	EnableGang      bool `yaml:"enable_gang"`                   // 是否开启杠结算
	EnableChaHuaZhu bool `yaml:"enable_tenable_chahuazhuuisui"` // 是否开启查花猪
	EnableChaDaJiao bool `yaml:"enable_chadajiao"`              // 是否开启查大叫
	EnableTuisui    bool `yaml:"enable_tuisui"`                 // 是否开启退税

	NeedBillDetails bool `yaml:"need_bill_deatils"` // 是否需要单局结算详情

	// --杠结算
	GangValue GangValue `yaml:"gang_value"` // 杠倍数
	// --胡结算
	HuValue HuTypeValue `yaml:"hu_type_value"` // 胡牌倍数

	// --其他
	HuPlayerSettle HuPlayerSettle `yaml:"huPlayer_settle"` //已胡牌玩家(未退出)可参与的结算

	HuQuitPlayerSettle HuQuitPlayerSettle `yaml:"huquit_player_settle"` // 已胡牌玩家(退出)可参与的结算

	GiveUpPlayerSettle GiveUpPlayerSettle `yaml:"giveup_player_settle"` // 认输玩家可参与的结算

	GangInstantSettle bool `yaml:"gang_instant_settle"` // 杠是否可以立即结算
	HuInstantSettle   bool `yaml:"hu_instant_settle"`   // 胡是否可以立即结算

}

// HuTypeValue 胡牌倍数
type HuTypeValue struct {
	HuGanghoupao        uint32 `yaml:"hu_ganghoupao"`        //杠后炮倍数
	HuQiangGangHu       uint32 `yaml:"hu_qiangganghu"`       //枪杠胡倍数
	HuDianPao           uint32 `yaml:"hu_dianpao"`           //点炮倍数
	HuGangKai           uint32 `yaml:"hu_gangkai"`           //杠开倍数
	HuHaidDiLao         uint32 `yaml:"hu_haidilao"`          //海底捞倍数
	HuGangShangHaiDiLao uint32 `yaml:"hu_gangshanghaidilao"` //杠上海底捞倍数
	HuZiMo              uint32 `yaml:"hu_zimo"`              //自摸倍数
	HuTianHu            uint32 `yaml:"hu_tianhu"`            //天胡倍数
	HuDiHu              uint32 `yaml:"hu_dihu"`              //地胡倍数
}

// GangValue 杠倍数
type GangValue struct {
	AnGangValue   uint32 `yaml:"angang_value"`   //暗杠倍数
	MingGangValue uint32 `yaml:"minggang_value"` //明杠倍数
	BuGangValue   uint32 `yaml:"bugang_value"`   //补杠倍数
}

// HuPlayerSettle 已胡牌玩家(未退出)可参与的结算
type HuPlayerSettle struct {
	HuPlayerGangSettle  bool `yaml:"huPlayer_gang_settle"`   //杠结算
	HuPlayeHuSettle     bool `yaml:"huPlayer_hu_settle"`     //胡结算
	HuPlayerRoundSettle bool `yaml:"huPlayer_round_settele"` //单局结算

}

// HuQuitPlayerSettle 已胡牌玩家(退出)可参与的结算
type HuQuitPlayerSettle struct {
	HuQuitPlayerGangSettle  bool `yaml:"huquit_player_gang_settle"`                         //杠结算
	HuQuitPlayeHuSettle     bool `yaml:"huquit_playe_hu_settle"`                            //胡结算
	HuQuitPlayerRoundSettle bool `yaml:"huPlayer_round_huquit_player_round_settelesettele"` //单局结算

}

// GiveUpPlayerSettle 认输玩家可参与的结算
type GiveUpPlayerSettle struct {
	GiveUpPlayerGangSettle  bool `yaml:"giveup_player_gang_settle"`   //杠结算
	GiveUpPlayerHuSettle    bool `yaml:"giveup_player_hu_settle"`     //胡结算s
	GiveUpPlayerRoundSettle bool `yaml:"giveup_player_round_settele"` //单局结算

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
