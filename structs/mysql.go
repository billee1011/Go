package structs

import (
	"github.com/go-xorm/xorm"
)

// MysqlEngineMgr Mysql 管理
type MysqlEngineMgr interface {
	// GetEngine 根据名字获取 *xorm.Engine
	//
	// 在 config 中配置 mysql 数据库列表
	// 如：
	// mysql_list:
	//   mysql1:
	//     user: root
	//     passwd: 123456
	//     addr: 127.0.0.1:3306
	//     db: database1
	//     params:
	//       charset: utf8
	//   mysql2:
	//     user: root
	//     passwd: 123456
	//     addr: 127.0.0.1:3306
	//     db: database1
	//     params:
	//       charset: utf8
	GetEngine(name string) (*xorm.Engine, error)
}
