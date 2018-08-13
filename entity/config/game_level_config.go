package config

// GameLevelConfig 游戏场次配置信息
type GameLevelConfig struct {
	ID               int    `json:"id"`
	GameID           int    `json:"gameID"`
	LevelID          int    `json:"levelID"`
	Name             string `json:"name"`
	Fee              int    `json:"fee"`
	BaseScores       int    `json:"baseScores"`
	LowScores        int    `json:"lowScores"`
	HighScores       int    `json:"highScores"`
	RealOnlinePeople int    `json:"realOnlinePeople"`
	ShowOnlinePeople int    `json:"showOnlinePeople"`
	Status           int    `json:"status"`
	Tag              int    `json:"tag"`
	IsAlms           int    `json:"isAlms"`
	Remark           string `json:"remark"`
}
