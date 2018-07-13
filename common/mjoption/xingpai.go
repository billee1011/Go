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
	ID                   int            `yaml:"id"`                      // 选项 ID
	WallCards            []int          `yaml:"wall_cards"`              // 墙牌
	HuGameOver           bool           `yaml:"hu_gameover"`             // 胡牌后是否触发游戏结束
	Hnz                  Hnz            `yaml:"hnz"`                     // 是否开启换N张
	EnableDingque        bool           `yaml:"enable_dingque"`          // 是否开启定缺
	EnableAddflower      bool           `yaml:"enable_addflower"`        // 是否开启补花
	EnableChi            bool           `yaml:"enable_chi"`              // 是否可以吃牌
	EnableKaijuAddflower bool           `yaml:"enable_kaiju_addflower"`  //是否开局补花
	PlayerNoNormalStates int32          `yaml:"player_no_normal_states"` // 玩家不参与游戏的不正常状态
	PlayerNum            int            `yaml:"player_num"`              //玩家人数
	FapaiType            FapaiType      `yaml:"fapai_type"`              //发牌方式
	TianhuCardType       TianhuCardType `yaml:"tianhu_card_type"`        //天胡胡哪张牌
	EnableGangSettle     bool           `yaml:"enable_gang_settle"`      //是否需要杠结算
	// Zhuang        int            `yaml:"zhuang"`         //选庄的方式
}

// FapaiType 发牌类型
type FapaiType int

// TianhuCardType 天胡胡哪张牌
type TianhuCardType int

const (
	// NomarlFapai 正常发牌，庄家14张，闲家13张
	NomarlFapai FapaiType = 1
	// ErrenFapai 二人麻将类型的发牌，所有人都只发13张
	ErrenFapai FapaiType = 2
	// MostTingsCard 听牌最多的那张牌
	MostTingsCard TianhuCardType = 1
	// RightCard 最右边的牌
	RightCard TianhuCardType = 2
	// MoCard 摸上来的牌
	MoCard TianhuCardType = 3
)

// Hnz 换n张
type Hnz struct {
	Enable bool `yaml:"enable"` //是否开启换n张
	Num    int  `yaml:"num"`    //换的张数
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
