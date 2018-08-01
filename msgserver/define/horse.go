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

 type HorseRace struct{
	Id		int64			// 唯一编号
	Channel int64			// 渠道ID
	Prov	int32			// 省份ID
	City	int32			// 城市ID
	bUse	int8			// 是否启用
	bUseParent	int8		// 是否启用上级配置
	tickTime	int32		// 两条跑马灯播放间隔，秒
	sleepTime	int32		// 一轮结束后等待时间,秒
	content []string		// 跑马灯内容列表
	startTime string		// 开始时间
	endTime   string		// 结束时间
}