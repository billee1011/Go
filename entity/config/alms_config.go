package config

type AlmsConfig struct {
	almsCountDown    int `json:"almsCountDown"`
	depositCountDown int `json:"depositCountDown"`
	getNorm          int `json:"getNorm"`
	getTimes         int `json:"getTimes"`
	getNumber        int `json:"getNumber"`
	version          int `json:"version"`
}
