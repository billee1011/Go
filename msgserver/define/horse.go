package define

/*
	功能: 跑马灯结构定义
	作者: Skywang
	日期: 2018-8-1


  `n_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `n_channel` bigint(20) NOT NULL COMMENT '渠道ID',
  `n_prov` bigint(20) DEFAULT NULL COMMENT '省包ID',
  `n_city` bigint(20) DEFAULT NULL COMMENT '城市ID',
  `n_bUse` tinyint(1) DEFAULT '1' COMMENT '是否启用',
  `n_bUseParent` tinyint(1) DEFAULT '1' COMMENT '是否启用上级配置',
  `n_tick` int(11) DEFAULT NULL COMMENT '两条跑马灯播放间隔，秒',
  `n_sleep` int(11) DEFAULT NULL COMMENT '一轮结束后等待时间,秒',
  `n_startTime1` varchar(100) DEFAULT NULL COMMENT '跑马灯1开始时间',
  `n_endTime1` varchar(100) DEFAULT NULL COMMENT '跑马灯1结束时间',
  `n_type` int(11) DEFAULT NULL COMMENT '时间类型: 1: 指定时间， 2：每天 3.每周',
  `n_horse1` varchar(200) DEFAULT NULL COMMENT '跑马灯1',
  `n_horse2` varchar(200) DEFAULT NULL COMMENT ' 跑马灯2',
  `n_horse3` varchar(200) DEFAULT NULL COMMENT '跑马灯3',
  `n_horse4` varchar(200) DEFAULT NULL COMMENT '跑马灯4',
  `n_horse5` varchar(200) DEFAULT NULL COMMENT '跑马灯5',
 */

 type HorseContent struct {
 	PlayType	int8 				// 时间类型 1=循环播放, 2=指定时间
	WeekDate	map[int8]bool		// 周N列表, 循环播放时选择周列表
	BeginDate	string				// 开始日期 2018-07-30
	EndDate		string				// 结束日期 2018-08-30
	BeginTime   string				// 开始时间 15:00
	EndTime		string				// 结束时间 20:00
	Content     string				// 跑马灯内容
	CheckStatus int8				// 检测状态
 }

 type HorseRace struct{
	Id		int64			// 唯一编号
	Channel int64			// 渠道ID
	Prov	int64			// 省份ID
	City	int64			// 城市ID
	IsUse	int8			// 是否启用: 0=不启用, 1=启用
	IsUseParent	int8		// 是否启用上级配置: 0=不启用, 1=启用
	TickTime	int32		// 两条跑马灯播放间隔，秒
	SleepTime	int32		// 一轮结束后等待时间,秒
	CheckStatus int8		// 检测状态
	Content []*HorseContent	// 跑马灯内容列表
}

type HorseContentJson struct {
	PlayType	int8 				`json:"playType"` // 时间类型 1=循环播放, 2=指定时间
	WeekDate	[]int8				`json:"weekDate"` // 周N列表, 循环播放时选择周列表,周日=0, 周1=1，周2=2，周6=6
	BeginDate	string				`json:"beginDate"` // 开始日期 2018-07-30
	EndDate		string				`json:"endDate"` // 结束日期 2018-08-30
	BeginTime   string				`json:"beginTime"` // 开始时间 15:00
	EndTime		string				`json:"endTime"` // 结束时间 20:00
	Content     string				`json:"content"` // 跑马灯内容
}

type HorseRaceJson struct{
	TickTime	int32		`json:"tickTime"` // 两条跑马灯播放间隔，秒
	SleepTime	int32		`json:"sleepTime"` // 一轮结束后等待时间,秒
	Horse []*HorseContentJson	`json:"horse"` // 跑马灯内容列表
}
