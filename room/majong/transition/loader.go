package transition

import (
	"errors"

	yaml "gopkg.in/yaml.v2"
)

type stateTran struct {
	Events    []string `yaml:"events"`
	NextState string   `yaml:"next_state"`
}

type transition struct {
	GameID int `yaml:"game_id"`
	States []struct {
		CurState string      `yaml:"state"`
		Trans    []stateTran `yaml:"transition"`
	} `yaml:"states"`
}

var errOpenFile = errors.New("打开文件失败")
var errReadFile = errors.New("读取文件失败")

func loadTransition(fileName string) ([]transition, error) {
	tt := []transition{}
	// data, err := ioutil.ReadFile(fileName)
	// if err != nil {
	// 	logrus.WithError(err).Error(errOpenFile)
	// 	return nil, errOpenFile
	// }
	if err := yaml.Unmarshal([]byte(transitionCfg), &tt); err != nil {
		return nil, err
	}
	return tt, nil
}
