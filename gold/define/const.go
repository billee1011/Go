package define
/*
 功能： 基础结构和常量定义
 作者： SkyWang
 日期： 2018-7-24

 */

import "fmt"

// 错误定义
var (
	ErrGoldType = fmt.Errorf("gold type error")
	ErrNoUser = fmt.Errorf("no user")
	ErrLoadDB = fmt.Errorf("load from db failed")
)

