package charge

import (
	"encoding/json"
	"fmt"
	"steve/entity/constant"
	"steve/external/configclient"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
)

// Item for json unmarshal
/*
	名称 | 类型 | 是否必须 | 默认值 | 备注
	---- | ---- | ---- | ----- | -----
	"name" | string | 是 | - | 商品显示名
	"tag" | string | 否 | '' | 商品标签
	"price" | int | 是 | - | 价格，单位：分
	"coin" | int | 是 | - | 金币数
	"present_coin" | int | 是 | - | 赠送金币数
*/
type Item struct {
	ID          uint64 `json:"item_id"`
	Name        string `json:"name"`
	Tag         string `json:"tag"`
	Price       uint64 `json:"price"`
	Coin        uint64 `json:"coin"`
	PresentCoin uint64 `json:"present_coin"`
}

// city -> items
type cityItems map[string][]Item

// platform -> item
type platformItems map[string]cityItems

var (
	itemLists     = platformItems{}
	itemListsLock sync.RWMutex
	// 每日最大充值数
	maxCharge = struct {
		MaxChargeVal uint64 `json:"max_charge"`
	}{MaxChargeVal: 200000}
	maxChargeLock sync.RWMutex
)

// loadItemList load item list from configuration server
func loadItemList(retry int) error {
	itemListJSON, err := configclient.GetConfigUntilSucc(constant.ChargeItemListKey.Key, constant.ChargeItemListKey.SubKey, retry, 5*time.Second)
	if err != nil {
		return fmt.Errorf("获取商品列表失败:%s", err.Error())
	}
	return parseItemList(itemListJSON)
}

// 解析商品列表
func parseItemList(config string) error {
	var _itemLists platformItems
	if err := json.Unmarshal([]byte(config), &_itemLists); err != nil {
		return fmt.Errorf("反序列化失败：%s", err.Error())
	}
	itemListsLock.Lock()
	itemLists = _itemLists
	logrus.Debugf("充值配置解析成功: %#v", _itemLists)
	itemListsLock.Unlock()
	return nil
}

// loadMaxCharge load max charge from configuration server
func loadMaxCharge(retry int) error {
	maxChargeJSON, err := configclient.GetConfigUntilSucc(constant.ChargeDayMaxKey.Key, constant.ChargeDayMaxKey.SubKey, retry, 5*time.Second)
	if err != nil {
		return fmt.Errorf("获取每日最大充值数失败：%s", err.Error())
	}
	return parseMaxCharge(maxChargeJSON)
}

// parseMaxCharge 解析日最大充值配置
func parseMaxCharge(config string) error {
	maxChargeLock.Lock()
	if err := json.Unmarshal([]byte(config), &maxCharge); err != nil {
		maxChargeLock.Unlock()
		return fmt.Errorf("反序列化失败：%s", err.Error())
	}
	logrus.Debugf("日最大充值配置解析成功: %d", maxCharge.MaxChargeVal)
	maxChargeLock.Unlock()
	return nil
}

// getDayMaxCharge 获取每日充值上限
func getDayMaxCharge() uint64 {
	return maxCharge.MaxChargeVal
}

// getItemList 获取商品列表
func getItemList(city int, platform int) ([]Item, error) {
	platformStr := "android"
	if platform == 2 {
		platformStr = "iphone"
	}
	cityStr := "default"
	if city != 0 {
		cityStr = fmt.Sprintf("city%d", city)
	}

	itemListsLock.RLock()
	cityItems, ok := itemLists[platformStr]
	itemListsLock.RUnlock()

	if !ok {
		return nil, fmt.Errorf("该平台的充值没有配置")
	}
	items, ok := cityItems[cityStr]
	if !ok {
		items, ok = cityItems["default"]
		if !ok {
			return nil, fmt.Errorf("该城市的充值没有配置")
		}
	}
	return items, nil
}

func init() {
	loadItemList(20)
	loadMaxCharge(20)

	configclient.SubConfigChange(constant.ChargeItemListKey.Key, constant.ChargeItemListKey.SubKey, func(key, subkey, val string) error {
		if err := parseItemList(val); err != nil {
			logrus.WithField("val", val).WithError(err).Errorln("商品列表配置解析失败")
		}
		return nil
	})

	configclient.SubConfigChange(constant.ChargeDayMaxKey.Key, constant.ChargeDayMaxKey.SubKey, func(key, subkey, val string) error {
		if err := parseMaxCharge(val); err != nil {
			logrus.WithField("val", val).WithError(err).Errorln("日最大充值配置解析失败")
		}
		return nil
	})
}
