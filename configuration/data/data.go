package data

import (
	"fmt"
	"steve/structs"
)

// GetConfig 获取配置信息
func GetConfig(key, subkey string) (string, error) {
	exposer := structs.GetGlobalExposer()
	engine, err := exposer.MysqlEngineMgr.GetEngine("config")
	if err != nil {
		return "", fmt.Errorf("获取 mysql 引擎失败：%v", err)
	}
	val := ""
	session := engine.Table("t_common_config").Where("`key` = ?", key).And("`subkey` = ?", subkey).Cols("value")
	exist, err := session.Get(&val)
	if err != nil {
		sql, _ := session.LastSQL()
		return "", fmt.Errorf("从数据库获取数据失败：%v, sql:%s", err, sql)
	}
	if !exist {
		return "", fmt.Errorf("配置数据不存在")
	}
	return val, nil
}
