package mysql

import (
	"fmt"
	"steve/structs"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cast"

	// 使用 mysql 数据库引擎
	"github.com/go-sql-driver/mysql"

	"github.com/spf13/viper"

	"github.com/go-xorm/xorm"
)

type mysqlEngineMgr struct {
	engines sync.Map // name -> engine
	configs map[string]mysql.Config
}

// CreateMysqlEngineMgr create structs.MysqlEngineMgr
func CreateMysqlEngineMgr() structs.MysqlEngineMgr {
	v := &mysqlEngineMgr{
		configs: make(map[string]mysql.Config, 8),
	}
	v.init()
	return v
}

// init 初始化
func (mg *mysqlEngineMgr) init() {
	allconfigs := viper.GetStringMap("mysql_list")
	for name, _configs := range allconfigs {
		configs := cast.ToStringMap(_configs)

		mg.configs[name] = mysql.Config{
			User:   cast.ToString(configs["user"]),
			Passwd: cast.ToString(configs["passwd"]),
			Net:    "tcp",
			Addr:   cast.ToString(configs["addr"]),
			DBName: cast.ToString(configs["db"]),
			AllowNativePasswords:true,
			Params: cast.ToStringMapString(configs["params"]),
		}
	}
	logrus.WithField("configs", mg.configs).Infoln("mysql 配置列表加载完成")
}

// GetEngine 根据名字获取 *xorm.Engine
func (mg *mysqlEngineMgr) GetEngine(name string) (*xorm.Engine, error) {
	_engine, ok := mg.engines.Load(name)
	if ok {
		return _engine.(*xorm.Engine), nil
	}
	conf, ok := mg.configs[name]
	if !ok {
		return nil, fmt.Errorf("数据库配置不存在")
	}
	engine, err := mg.createEngine(&conf)
	if err != nil {
		return nil, err
	}
	actual, loaded := mg.engines.LoadOrStore(name, engine)
	if loaded {
		engine.Close()
	}
	return actual.(*xorm.Engine), nil
}

func (mg *mysqlEngineMgr) createEngine(conf *mysql.Config) (*xorm.Engine, error) {
	dns := conf.FormatDSN()
	engine, err := xorm.NewEngine("mysql", dns)
	if err != nil {
		return nil, fmt.Errorf("创建失败: %v", err)
	}
	if err := engine.Ping(); err != nil {
		logrus.WithField("dns", dns).Errorln("连接失败")
		return nil, fmt.Errorf("连接失败： %v", err)
	}
	return engine, nil
}
