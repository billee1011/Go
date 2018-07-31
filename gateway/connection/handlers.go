package connection

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"steve/client_pb/gate"
	"steve/client_pb/msgid"
	"steve/structs/exchanger"
	"steve/structs/proto/gate_rpc"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	"github.com/spf13/viper"
)

// HandleHeartBeat 处理心跳
func HandleHeartBeat(clientID uint64, header *steve_proto_gaterpc.Header, req gate.GateHeartBeatReq) (ret []exchanger.ResponseMsg) {
	connection := GetConnectionMgr().GetConnection(clientID)
	if connection == nil {
		return
	}
	connection.HeartBeat()

	response := gate.GateHeartBeatRsp{
		TimeStamp: proto.Uint64(req.GetTimeStamp()),
	}
	logrus.WithFields(logrus.Fields{
		"client_id": clientID,
		"player_id": connection.GetPlayerID(),
		"response":  response.String(),
		"request":   req.String(),
	}).Debugln("心跳")
	return []exchanger.ResponseMsg{
		exchanger.ResponseMsg{
			MsgID: uint32(msgid.MsgID_GATE_HEART_BEAT_RSP),
			Body:  &response,
		},
	}
}

// HandleTransmitHTTPReq transmit http request to platform
func HandleTransmitHTTPReq(clientID uint64, header *steve_proto_gaterpc.Header, req gate.GateTransmitHTTPReq) (ret []exchanger.ResponseMsg) {
	return []exchanger.ResponseMsg{
		exchanger.ResponseMsg{
			MsgID: uint32(msgid.MsgID_GATE_TRANSMIT_HTTP_RSP),
			Body:  commitHTTPReq(&req),
		},
	}
}

func commitHTTPReq(request *gate.GateTransmitHTTPReq) *gate.GateTransmitHTTPRsp {
	response := &gate.GateTransmitHTTPRsp{
		StatusCode: proto.Int32(400),
	}
	httpAddr := viper.GetString("platform_http_address")
	if httpAddr == "" {
		logrus.Errorln("平台服地址没有配置")
		return response
	}
	url := fmt.Sprintf("%s%s", httpAddr, request.GetUrl())
	method := request.GetMethod()
	contentType := request.GetContentType()

	entry := logrus.WithFields(logrus.Fields{
		"url":    url,
		"method": method,
	})
	var httpResponse *http.Response

	if method == "GET" {
		var err error
		httpResponse, err = http.Get(url)
		if err != nil {
			entry.WithError(err).Errorln("请求失败")
			return response
		}
	} else if method == "POST" {
		var err error
		httpResponse, err = http.Post(url, contentType, bytes.NewReader(request.GetData()))
		if err != nil {
			entry.WithError(err).Errorln("请求失败")
			return response
		}
	} else {
		entry.Warningln("不支持的 method")
		return response
	}
	var err error
	response.ResponseBody, err = ioutil.ReadAll(httpResponse.Body)
	httpResponse.Body.Close()
	if err == nil || err == io.EOF {
		response.StatusCode = proto.Int32(int32(httpResponse.StatusCode))
	}
	entry.WithField("status_code", response.GetStatusCode()).Debugln("请求完成")
	return response
}
