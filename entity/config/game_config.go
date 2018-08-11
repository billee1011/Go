package config

// GameConfig 游戏配置
type GameConfig struct {
	ID         int    `json:"id"`
	GameID     int    `json:"gameID"`
	Name       string `json:"name"`
	Type       int    `json:"type"`
	MinPeople  int    `json:"minPeople"`
	MaxPeople  int    `json:"maxPeople"`
	Playform   int    `json:"playform"`
	CountryID  int    `json:"countryID"`
	ProvinceID int    `json:"provinceID"`
	CityID     int    `json:"cityID"`
	ChannelID  int    `json:"channelID"`
}
