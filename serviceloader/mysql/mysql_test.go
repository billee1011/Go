package mysql

import (
	"fmt"
	"testing"

	"github.com/go-sql-driver/mysql"

	"github.com/go-xorm/xorm"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func Test_CreateEngine(t *testing.T) {
	viper.SetConfigFile("./testconfig.yml")
	assert.Nil(t, viper.ReadInConfig())
	mgr := CreateMysqlEngineMgr()
	engine, err := mgr.GetEngine("mysql1")
	assert.Nil(t, err)
	assert.NotNil(t, engine)
}

func Test__(t *testing.T) {
	conf := mysql.Config{
		User:   "root",
		Passwd: "123456",
		Net:    "tcp",
		Addr:   "localhost:3306",
		DBName: "database1",
		Params: map[string]string{"charset": "utf8"},
	}
	dsn := conf.FormatDSN()
	fmt.Println(dsn)
	engine, err := xorm.NewEngine("mysql", dsn)
	assert.Nil(t, err)
	assert.NotNil(t, engine)
	assert.Nil(t, engine.Ping())
}
