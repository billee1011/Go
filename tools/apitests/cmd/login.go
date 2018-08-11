// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"steve/client_pb/login"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"

	"github.com/spf13/cobra"
)

var loginFlags struct {
	productID int    // 产品 ID
	url       string // 登录 url
	loginType int    // 登录方式
	channel   int    // 渠道 id
	username  string // 用户名
	dymcCode  string // 验证码
	passwd    string // 密码
	proID     int    // 省 ID
	cityID    int    // 市 ID

	// 微信参数
	wxAppID   string
	wxOpenID  string
	wxUnionID string
}

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "登录 api 测试",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		execLogin()
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	loginCmd.Flags().IntVar(&loginFlags.productID, "productid", 0, "产品 ID")
	loginCmd.Flags().StringVar(&loginFlags.url, "url", "", "url")
	loginCmd.Flags().IntVar(&loginFlags.loginType, "logintype", 0, "登录方式")
	loginCmd.Flags().IntVar(&loginFlags.channel, "channel", 0, "渠道")
	loginCmd.Flags().StringVar(&loginFlags.username, "username", "", "username")
	loginCmd.Flags().StringVar(&loginFlags.dymcCode, "dymccode", "", "dymccode")
	loginCmd.Flags().StringVar(&loginFlags.passwd, "passwd", "", "passwd")
	loginCmd.Flags().IntVar(&loginFlags.proID, "proid", 0, "省 ID")
	loginCmd.Flags().IntVar(&loginFlags.cityID, "cityid", 0, "市 ID")
	loginCmd.Flags().StringVar(&loginFlags.wxAppID, "wxappid", "", "wxappid")
	loginCmd.Flags().StringVar(&loginFlags.wxOpenID, "wxopenid", "", "wxopenid")
	loginCmd.Flags().StringVar(&loginFlags.wxUnionID, "wxunionid", "", "wxunionid")
}

func execLogin() {
	loginDataMessage := login.LoginData{
		Type:     login.LoginType(loginFlags.loginType).Enum(),
		Channel:  login.ChannelType(loginFlags.channel).Enum(),
		Username: proto.String(loginFlags.username),
		DymcCode: proto.String(loginFlags.dymcCode),
		Password: proto.String(loginFlags.passwd),
		ProId:    proto.Uint64(uint64(loginFlags.proID)),
		CityId:   proto.Uint64(uint64(loginFlags.cityID)),
		BindInfo: &login.BindInfo{
			WcAppid:   proto.String(loginFlags.wxAppID),
			WcOpenid:  proto.String(loginFlags.wxOpenID),
			WcUnionid: proto.String(loginFlags.wxUnionID),
		},
	}
	loginData, err := proto.Marshal(&loginDataMessage)
	if err != nil {
		logrus.WithError(err).Errorln("序列化失败")
		return
	}

	loginMessage := login.AccountSysLoginRequestData{
		ProductId: proto.Uint64(uint64(loginFlags.productID)),
		Data:      loginData,
	}

	data, err := proto.Marshal(&loginMessage)
	if err != nil {
		logrus.WithError(err).Errorln("序列化失败")
		return
	}

	response, err := http.Post(loginFlags.url, "application/octet-stream", bytes.NewReader(data))
	if err != nil {
		logrus.WithError(err).Errorln("请求失败")
		return
	}
	respData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logrus.WithError(err).Errorln("读取回复数据失败")
		return
	}
	fmt.Println("收到回复：", string(respData))
	return
}
