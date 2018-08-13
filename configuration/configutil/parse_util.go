package configutil

import (
	"steve/entity/config"

	"github.com/W1llyu/ourjson"
)

func main() {
	json := `[{ 
"gameID":1,
"levelID":1,
"name":"新手场",
"fee":1,
"baseScores":1,
"lowScores":100,
"highScores":1000000,
"realOnlinePeople":1,
"showOnlinePeople":1,
"status":1,
"tag":null,
"isAlms":1,
"remark":null
},
{ 
"gameID":2,
"levelID":1,
"name":"新手场",
"fee":1,
"baseScores":1,
"lowScores":0,
"highScores":1000000,
"realOnlinePeople":1,
"showOnlinePeople":1,
"status":1,
"tag":null,
"isAlms":1,
"remark":null
},
{ 
"gameID":3,
"levelID":1,
"name":"新手场",
"fee":1,
"baseScores":1,
"lowScores":0,
"highScores":1000000,
"realOnlinePeople":1,
"showOnlinePeople":1,
"status":1,
"tag":null,
"isAlms":1,
"remark":null
},
{ 
"gameID":4,
"levelID":1,
"name":"新手场",
"fee":1,
"baseScores":2,
"lowScores":0,
"highScores":1000000,
"realOnlinePeople":1,
"showOnlinePeople":1,
"status":1,
"tag":null,
"isAlms":1,
"remark":null
}]`
	ParseToGameLevelConfigMap(json)
}

func ParseToGameLevelConfigMap(json string) []config.GameLevelConfig {
	arr, _ := ourjson.ParseArray(json)
	jsonObjs := arr.Values()
	result := make([]config.GameLevelConfig, len(jsonObjs))
	for index, obj := range jsonObjs {
		configObj := parseToGameLevelConfig(obj.JsonObject())
		result[index] = configObj
	}
	return result
}

func ParseToGameConfigMap(json string) []config.GameConfig {
	arr, _ := ourjson.ParseArray(json)
	jsonObjs := arr.Values()
	result := make([]config.GameConfig, len(jsonObjs))
	for index, obj := range jsonObjs {
		configObj := parseToGameConfig(obj.JsonObject())
		result[index] = configObj
	}
	return result
}
func parseToGameConfig(json *ourjson.JsonObject) config.GameConfig {
	gameId, _ := json.GetInt("gameID")
	name, _ := json.GetString("name")
	typ, _ := json.GetInt("type")
	minPeople, _ := json.GetInt("minPeople")
	maxPeople, _ := json.GetInt("maxPeople")
	playform, _ := json.GetInt("playform")
	countryID, _ := json.GetInt("countryID")
	provinceID, _ := json.GetInt("provinceID")
	cityID, _ := json.GetInt("cityID")
	channelID, _ := json.GetInt("channelID")
	return config.GameConfig{
		GameID:     gameId,
		Name:       name,
		Type:       typ,
		MinPeople:  minPeople,
		MaxPeople:  maxPeople,
		Playform:   playform,
		CountryID:  countryID,
		ProvinceID: provinceID,
		CityID:     cityID,
		ChannelID:  channelID,
	}
}

func parseToGameLevelConfig(json *ourjson.JsonObject) config.GameLevelConfig {
	gameId, _ := json.GetInt("gameID")
	levelID, _ := json.GetInt("levelID")
	name, _ := json.GetString("name")
	fee, _ := json.GetInt("fee")
	baseScores, _ := json.GetInt("baseScores")
	lowScores, _ := json.GetInt("lowScores")
	highScores, _ := json.GetInt("highScores")
	realOnlinePeople, _ := json.GetInt("realOnlinePeople")
	showOnlinePeople, _ := json.GetInt("showOnlinePeople")
	status, _ := json.GetInt("status")
	tag, _ := json.GetInt("tag")
	isAlms, _ := json.GetInt("isAlms")
	remark, _ := json.GetString("remark")
	return config.GameLevelConfig{
		GameID:           gameId,
		LevelID:          levelID,
		Name:             name,
		BaseScores:       baseScores,
		Fee:              fee,
		LowScores:        lowScores,
		HighScores:       highScores,
		RealOnlinePeople: realOnlinePeople,
		ShowOnlinePeople: showOnlinePeople,
		Status:           status,
		Tag:              tag,
		IsAlms:           isAlms,
		Remark:           remark,
	}
}
