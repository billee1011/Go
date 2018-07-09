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
	longqiduiFuncID int = 4
	// 清龙七对
	qinglongqiduiFuncID int = 5
	// 碰碰胡
	pengpenghuFuncID int = 6
	// 清碰
	qingpengFuncID int = 7
	// 金钩钓
	jingoudiaoFuncID int = 8
	// 清金钩钓
	qingjingoudiaoFuncID int = 9
	// 十八罗汉
	shibaluohanFuncID int = 10
	// 清十八罗汉
	qingshibaluohanFuncID int = 11

	// 大四喜
	dasixiFuncID int = 12
	// 大三元
	dasanyuanFuncID int = 13
	// 九莲宝灯
	jiulianbaodengFuncID int = 14
	// 大于五
	dayuwuFuncID int = 15
	// 小于五
	xiaoyuwuFuncID int = 16
	// 大七星
	daqixingFuncID int = 17
	// 连七对
	lianqiduiFuncID int = 18
	// 四杠
	sigangFuncID int = 19
	// 小四喜
	xiaosixiFuncID int = 20
	// 小三元
	xiaosanyuanFuncID int = 21
	// 双龙会
	shuanglonghuiFuncID int = 22
	// 字一色
	ziyiseFuncID int = 23
	// 四暗刻
	siankeFuncID int = 24
	// 四同顺
	sitongshunFuncID int = 25
	// 三元七对
	sanyuanqiduiFuncID int = 26
	// 四喜七对
	sixiqiduiFuncID int = 27
	// 四连刻
	siliankeFuncID int = 28
	// 四步高
	sibugaoFuncID int = 29
	// 混幺九
	hunyaojiuFuncID int = 30
	// 三杠
	sangangFuncID int = 31
	// 天听
	tiantingFuncID int = 32
	// 四字刻
	sizikeFuncID int = 33
	// 大三风
	dasanfengFuncID int = 34
	// 三同顺
	santongshunFuncID int = 35
	// 三连刻
	sanliankeFuncID int = 36
	// 清龙
	qinglongFuncID int = 37
	// 三步高
	sanbugaoFuncID int = 38
	// 全花
	quanhuaFuncID int = 39
	// 三暗刻
	sanankeFuncID int = 40
	// 妙手回春
	miaoshouhuichunFuncID int = 41
	// 海底捞月
	haidilaoyueFuncID int = 42
	// 抢扛胡
	qiangganghuFuncID int = 43
	// 杠上开花
	gangshangkaihuaFuncID int = 44
	// 小三风
	xiaosanfengFuncID int = 45
	// 双箭刻
	shuangjiankeFuncID int = 46
	// 双暗杠
	shuangangangFuncID int = 47
	// 混一色
	hunyiseFuncID int = 48
	// 全求人
	quanqiurenFuncID int = 49
	// 全带幺
	quandaiyaoFuncID int = 50
	// 双明杠
	shuangminggangFuncID int = 51
	// 和绝张
	hujuezhangFuncID int = 52
	// 报听
	baotingFuncID int = 53
	// 报听一发
	baotingyifaFuncID int = 54
	// 春夏秋冬
	chunxiaqiudongFuncID int = 55
	// 梅兰竹菊
	meilanzhujuFuncID int = 56
	// 无花牌
	wuhuapaiFuncID int = 57
	// 门风刻
	menfengkeFuncID int = 58
	// 圈风刻
	quanfengkeFuncID int = 59
	// 箭刻
	jiankeFuncID int = 60
	// 二人麻将平胡
	errenpinghuFuncID int = 61
	// 四归一
	siguiyiFuncID int = 62
	// 断幺
	duanyaoFuncID int = 63
	// 双暗刻
	shuangankeFuncID int = 64
	// 暗杠
	angangFuncID int = 65
	// 门前清
	menqianqingFuncID int = 66
	// 一般高
	yibangaoFuncID int = 67
	// 连六
	lianliuFuncID int = 68
	// 老少副
	laoshaofuFuncID int = 69
	// 明杠
	minggangFuncID int = 70
	// 边张
	bianzhangFuncID int = 71
	// 坎张
	kanzhangFuncID int = 72
	// 不求人
	buqiurenFuncID int = 73
	// 单钓将
	dangdiaojiangFuncID int = 74
	// 自摸
	zimoFuncID int = 75
	// 杠后炮
	ganghoupaoFuncID int = 76
	// 杠上海底捞
	gangshanghaidilaoFuncID int = 77
	// 天胡
	tianhuFuncID int = 78
	// 地胡
	dihuFuncID int = 79
	// 人胡
	renhuFuncID int = 80
	// 点炮
	dianpaoFuncID int = 81
)

// checkFunc 检测函数
type checkFunc func(tc *typeCalculator) bool

var checkFuncs map[int]checkFunc

func init() {
	checkFuncs = map[int]checkFunc{
		// pinghuFuncID:          checkPinghu,
		// qingyiseFuncID:        checkQingyise,
		// qiduiFuncID:           checkQidui,
		// qingqiduiFuncID:       checkQingqidui,
		// longqiduiFuncID:       checkLongQiDui,
		// qinglongqiduiFuncID:   checkQingLongQiDui,
		// pengpenghuFuncID:      checkPengpenghu,
		// qingpengFuncID:        checkQingPeng,
		// jingoudiaoFuncID:      checkJinGouDiao,
		// qingjingoudiaoFuncID:  checkQingJinGouDiao,
		// shibaluohanFuncID:     checkShiBaLuoHan,
		// qingshibaluohanFuncID: checkQingShiBaLuoHan,
		// dasixiFuncID:          checkDaSiXi,
		// dasanyuanFuncID:       checkDaSanYuan,
		// jiulianbaodengFuncID:  checkJiuLianBaoDeng,
		// dayuwuFuncID:          checkDaYuWu,
		// xiaoyuwuFuncID:        checkXiaoYuWu,
		// xiaosixiFuncID:        checkXiaoSiXi,
		// xiaosanyuanFuncID:     checkXiaoSanYuan,
		// siankeFuncID:          checkSiAnKe,
		// sitongshunFuncID:      checkSiTongShun,
		// sanyuanqiduiFuncID:    checkSanYuanQiDui,
		// sibugaoFuncID:         checkSiBuGao,
		// hunyaojiuFuncID:       checkHunYaoJiu,
		// sizikeFuncID:          checkSiZiKe,
		// dasanfengFuncID:       checkDaSanFeng,
		// qinglongFuncID:        checkQingLong,
		// sanbugaoFuncID:        checkSanBuGao,
		// miaoshouhuichunFuncID: checkMiaoShouHuiChun,
		// haidilaoyueFuncID:     checkHaiDiLaoYue,
		// xiaosanfengFuncID:     checkXiaoSanFeng,
		// shuangjiankeFuncID:    checkShuangJianKe,
		// quandaiyaoFuncID:      checkQuanDaiYao,
		// shuangminggangFuncID:  checkShuangMingGang,
		// buqiurenFuncID:        checkBuQiuRen,
		// hujuezhangFuncID:      checkHuJueZhang,
		// menfengkeFuncID:       checkMenFengKe,
		// quanfengkeFuncID:      checkQuanFengKe,
		// jiankeFuncID:          checkJianKe,
		// errenpinghuFuncID:     checkErrRenPingHe,
		// siguiyiFuncID:         checkSiGuiYi,
		// yibangaoFuncID:        checkYiBangGao,
		// lianliuFuncID:         checkLianLiu,
		// laoshaofuFuncID:       checkLaoShaoFu,
		daqixingFuncID:        checkDaQiXing,
		lianqiduiFuncID:       checkLianQiDui,
		sigangFuncID:          checkSiGang,
		tianhuFuncID:          checkTianHu,
		dihuFuncID:            checkDiHu,
		shuangjiankeFuncID:    checkShuanLongHui,
		ziyiseFuncID:          checkZiYiSe,
		renhuFuncID:           checkRenHu,
		sixiqiduiFuncID:       checkSiXiQiDui,
		siliankeFuncID:        checkSiLianKe,
		sangangFuncID:         checkSanGang,
		tiantingFuncID:        checkTianTing,
		santongshunFuncID:     checkSanTongShun,
		qiduiFuncID:           checkQidui,
		sanliankeFuncID:       checkSanLianKe,
		quanhuaFuncID:         checkQuanHua,
		sanankeFuncID:         checkSanAnKe,
		qiangganghuFuncID:     checkQiangGangHu,
		gangshangkaihuaFuncID: checkGangShangKaiHua,
		shuangangangFuncID:    checkShuanAnGang,
		hunyiseFuncID:         checkHunYiSe,
		quanqiurenFuncID:      checkQuanQiuRen,
		baotingFuncID:         checkBaoTing,
		baotingyifaFuncID:     checkBaoTingYiFa,
		chunxiaqiudongFuncID:  checkChunXiaQiuDong,
		meilanzhujuFuncID:     checkMeiLanZhuJiu,
		wuhuapaiFuncID:        checkWuHuaPai,
		siguiyiFuncID:         checkSiGuiYi,
		duanyaoFuncID:         checkDuanYao,
		shuangankeFuncID:      checkShuanAnKe,
		angangFuncID:          checkAnGang,
		menqianqingFuncID:     checkMengQianQing,
		minggangFuncID:        checkMingGang,
		bianzhangFuncID:       checkBianZhang,
		kanzhangFuncID:        checkKanZhang,
		dangdiaojiangFuncID:   checkDanDiaoJiang,
		zimoFuncID:            checkZiMo,
		dianpaoFuncID:         checkDianPao,
	}
}
