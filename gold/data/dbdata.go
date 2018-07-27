package data

/*
	功能： 服务数据保存到Mysql.
	作者： SkyWang
	日期： 2018-7-25

*/

// 从DB加载玩家金币
func LoadGoldFromDB(uid uint64) (map[int16]int64, error) {
	m := make(map[int16]int64)
	_ = uid
	m[0] = 5000
	m[1] = 5000
	m[2] = 5000
	m[3] = 5000
	return m, nil
}

// 将玩家金币同步到DB
func SaveGoldToDB(uid uint64, goldList map[int16]int64) error {
	_ = uid
	_ = goldList
	return nil
}
