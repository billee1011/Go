package handle

import (
	"net/http"
	"steve/gutils"

	"github.com/Sirupsen/logrus"
)

func respMSG(resp http.ResponseWriter, message string, code int) {
	resp.WriteHeader(code)
	resp.Write([]byte(message))
	switch code {
	case 200:
		logrus.Infoln(message)
	default:
		logrus.Debugln(message)
	}
}

//LogPeiPaiInfos 打印配牌信息
func LogPeiPaiInfos() {
	for k, info := range peiPaiInfos {
		logrus.WithFields(logrus.Fields{
			"game":   info.Key,
			"cards":  info.Cards,
			"len":    info.Len,
			"fx":     info.Fx,
			"zhuang": info.Zhuang,
		}).Info(k)
	}
}

//LogOptionInfos 打印选项信息
func LogOptionInfos() {
	for k, info := range optionInfos {
		logrus.WithFields(logrus.Fields{
			"game": info.Key,
			"fx":   info.Hsz,
		}).Info(k)
	}
}

func idIntToStr(gameID int) string {
	switch gameID {
	case gutils.SCXLGameID:
		return SCXL
	case gutils.SCXZGameID:
		return SCXZ
	case gutils.ERMJGameID:
		return ERMJ
	}
	return ""
}
