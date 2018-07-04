package fantype

// 计算函数 ID 定义
const (
	// 平胡计算函数 ID
	pinghuFuncID int = 0
	// 清一色计算函数 ID
	qingyiseFuncID int = 1
	// 七对
	qiduiFuncID int = 2
	// 清七对
	qingqiduiFuncID int = 3

	// 龙七对
	// 清龙七对

	// 碰碰胡
	pengpenghuFuncID int = 6
)

// checkFunc 检测函数
type checkFunc func(tc *typeCalculator) bool

var checkFuncs map[int]checkFunc

func init() {
	checkFuncs = map[int]checkFunc{
		pinghuFuncID:     checkPinghu,
		qingyiseFuncID:   checkQingyise,
		qiduiFuncID:      checkQidui,
		qingqiduiFuncID:  checkQingqidui,
		pengpenghuFuncID: checkPengpenghu,
	}
}
