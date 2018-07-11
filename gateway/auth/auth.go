package auth

///////////////////////////////////////////////
// 已经废弃
//////////////////////////////////////////////

// // HandleAuthReq 处理认证请求
// // step 1. 验证请求数据是否合法，包括 token， 过期时间
// // step 2. 保存连接 ID 到玩家 ID 的映射到内存
// func HandleAuthReq(clientID uint64, header *steve_proto_gaterpc.Header, req gate.GateAuthReq) (ret []exchanger.ResponseMsg) {
// 	response := &gate.GateAuthRsp{
// 		ErrCode: gate.ErrCode_ERR_EXPIRE_TOKEN.Enum(),
// 	}
// 	ret = []exchanger.ResponseMsg{{
// 		MsgID: uint32(msgid.MsgID_GATE_AUTH_RSP),
// 		Body:  response,
// 	}}
// 	if !checkRequest(clientID, header, &req, response) {
// 		return
// 	}
// 	checkAnother(req.GetPlayerId())
// 	if !saveConnectPlayerMap(clientID, header, &req, response) {
// 		return
// 	}
// 	response.ErrCode = gate.ErrCode_SUCCESS.Enum()
// 	return
// }

// func checkRequest(clientID uint64, header *steve_proto_gaterpc.Header, req *gate.GateAuthReq, response *gate.GateAuthRsp) bool {
// 	entry := logrus.WithFields(logrus.Fields{
// 		"func_name": "checkRequest",
// 		"client_id": clientID,
// 		"req_token": req.GetToken(),
// 		"expire":    req.GetExpire(),
// 	})
// 	expire := time.Unix(req.GetExpire(), 0)
// 	if time.Now().After(expire) {
// 		response.ErrCode = gate.ErrCode_ERR_EXPIRE_TOKEN.Enum()
// 		entry.Debugln("token 过期")
// 		return false
// 	}
// 	gateIP := viper.GetString(config.ListenClientAddrInquire)
// 	gatePort := viper.GetInt(config.ListenClientPort)
// 	key := viper.GetString(config.AuthKey)
// 	correctToken := auth.GenerateAuthToken(req.GetPlayerId(), gateIP, gatePort, req.GetExpire(), key)

// 	entry = entry.WithFields(logrus.Fields{
// 		"gate_ip":       gateIP,
// 		"gate_port":     gatePort,
// 		"key":           key,
// 		"correct_token": correctToken,
// 	})
// 	if correctToken != req.GetToken() {
// 		response.ErrCode = gate.ErrCode_ERR_INVALID_TOKEN.Enum()
// 		entry.Infoln("token 验证失败")
// 		return false
// 	}
// 	return true
// }
