package connection

import (
	"net/http"
	"net/http/httptest"
	"steve/client_pb/gate"
	"testing"

	"github.com/Sirupsen/logrus"

	"github.com/golang/protobuf/proto"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// 测试 POST
func Test_commitHTTPReq_POST(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	}))

	viper.SetDefault("platform_http_address", svr.URL)

	response := commitHTTPReq(&gate.GateTransmitHTTPReq{
		Url:         proto.String(""),
		Method:      proto.String("POST"),
		ContentType: proto.String(""),
		Data:        nil,
	})
	assert.Equal(t, 200, int(response.GetStatusCode()))
	assert.Equal(t, []byte("hello"), response.GetResponseBody())
}

// 测试 GET
func Test_commitHTTPReq_GET(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	}))

	viper.SetDefault("platform_http_address", svr.URL)

	response := commitHTTPReq(&gate.GateTransmitHTTPReq{
		Url:         proto.String(""),
		Method:      proto.String("GET"),
		ContentType: proto.String(""),
		Data:        nil,
	})
	assert.Equal(t, 200, int(response.GetStatusCode()))
	assert.Equal(t, []byte("hello"), response.GetResponseBody())
}

func init() {
	logrus.SetLevel(logrus.DebugLevel)
}
