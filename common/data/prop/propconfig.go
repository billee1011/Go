package prop

import (
	"encoding/json"
	"steve/entity/constant"
	"steve/external/configclient"
)

// GetPropsConfig 获取道具配置信息
func GetPropsConfig() (propConfig []constant.PropAttr, err error) {
	// 现在直接从数据库获取，后面改为先从redis获取；订阅更新消息，更新时删掉redis数据 TODO
	val, err := configclient.GetConfig(constant.PropKey, constant.PropSubKey)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(val), propConfig)
	if err != nil {
		return nil, err
	}
	return
}

// GetSomePropsConfig 获取某些道具配置信息
func GetSomePropsConfig(propIDs []int32) (propConfig []constant.PropAttr, err error) {
	// 现在直接从数据库获取，后面改为先从redis获取；订阅更新消息，更新时删掉redis数据 TODO
	val, err := configclient.GetConfig(constant.PropKey, constant.PropSubKey)
	if err != nil {
		return nil, err
	}

	var allConfig []constant.PropAttr
	err = json.Unmarshal([]byte(val), allConfig)
	if err != nil {
		return nil, err
	}
	propConfig = make([]constant.PropAttr, len(propIDs))
	for index, id := range propIDs {
		for _, config := range allConfig {
			if id == config.PropID {
				propConfig[index] = config
			}
		}
	}

	return
}
