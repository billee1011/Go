package msg

import (
	"steve/structs/proto/gate_rpc"
	"steve/structs/exchanger"
	"github.com/golang/protobuf/proto"
	"steve/client_pb/common"
	"steve/client_pb/msgid"
	"steve/mailserver/logic"
	"github.com/Sirupsen/logrus"
	"steve/client_pb/mailserver"
)

/*
 功能：
		1. 完成从GateWay(网关）过来的所有Client的请求消息的处理。
 		2. 通过core.coreConfig配置需要处理的消息列表。
		3. 需要在GateWay配置消息ID开始~ 结束区间 关联到当前服务名,GateWay才会把消息转发到此服务
*/

// 获取未读消息总数请求
func ProcessGetUnReadSumReq(playerID uint64, header *steve_proto_gaterpc.Header, req mailserver.MailSvrGetUnReadSumReq) (ret []exchanger.ResponseMsg) {

	logrus.Debugln("ProcessGetUnReadSumReq req", req)

	response := &mailserver.MailSvrGetUnReadSumRsp{
		ErrCode: proto.Int32(0),
		ErrDesc: proto.String("成功"),
	}
	//
	ret = []exchanger.ResponseMsg{{
		MsgID: uint32(msgid.MsgID_MAILSVR_GET_UNREAD_SUM_RSP),
		Body:  response,
	}}

	sum, err := logic.GetGetUnReadSum(playerID)
	if err != nil {
		response.ErrCode = proto.Int32(int32(common.ErrCode_EC_FAIL))
		response.ErrDesc = proto.String("失败")
		logrus.Errorln("ProcessGetUnReadSumReq err:", err)
		return nil
	}
	response.Sum = &sum
	logrus.Debugln("ProcessGetUnReadSumReq resp", response)
	return ret
}



// 获取邮件消息列表请求
func ProcessGetMailListReq(playerID uint64, header *steve_proto_gaterpc.Header, req mailserver.MailSvrGetMailListReq) (ret []exchanger.ResponseMsg) {

	logrus.Debugln("ProcessGetMailListReq req", req)

	response := &mailserver.MailSvrGetMailListRsp{
		ErrCode: proto.Int32(0),
		ErrDesc: proto.String("成功"),
	}
	//
	ret = []exchanger.ResponseMsg{{
		MsgID: uint32(msgid.MsgID_MAILSVR_GET_MAIL_LIST_RSP),
		Body:  response,
	}}

	mailList, err := logic.GetMailList(playerID)
	if err != nil {
		response.ErrCode = proto.Int32(int32(common.ErrCode_EC_FAIL))
		response.ErrDesc = proto.String("失败")
		logrus.Errorln("ProcessGetMailListReq err:", err)
		return nil
	}

	response.MailList = mailList

	logrus.Debugln("ProcessGetMailListReq resp", response)
	return ret
}


// 获取指定邮件详情请求
func ProcessGetMailDetailReq(playerID uint64, header *steve_proto_gaterpc.Header, req mailserver.MailSvrGetMailDetailReq) (ret []exchanger.ResponseMsg) {

	logrus.Debugln("ProcessGetMailDetailReq req", req)

	response := &mailserver.MailSvrGetMailDetailRsp{
		ErrCode: proto.Int32(0),
		ErrDesc: proto.String("成功"),
	}
	//
	ret = []exchanger.ResponseMsg{{
		MsgID: uint32(msgid.MsgID_MAILSVR_GET_MAIL_DETAIL_RSP),
		Body:  response,
	}}

	detail, err := logic.GetMailDetail(playerID, req.GetMailId())
	if err != nil {
		response.ErrCode = proto.Int32(int32(common.ErrCode_EC_FAIL))
		response.ErrDesc = proto.String("失败")
		logrus.Errorln("ProcessGetMailDetailReq err:", err)
		return nil
	}

	response.Detail = detail
	logrus.Debugln("ProcessGetMailDetailReq resp", response)
	return ret
}

// 标记邮件为已读请求
func ProcessSetReadTagReq(playerID uint64, header *steve_proto_gaterpc.Header, req mailserver.MailSvrSetReadTagReq) (ret []exchanger.ResponseMsg) {

	logrus.Debugln("ProcessSetReadTagReq req", req)

	response := &mailserver.MailSvrSetReadTagRsp{
		ErrCode: proto.Int32(0),
		ErrDesc: proto.String("成功"),
	}
	//
	ret = []exchanger.ResponseMsg{{
		MsgID: uint32(msgid.MsgID_MAILSVR_SET_READ_TAG_RSP),
		Body:  response,
	}}

	err := logic.SetReadTag(playerID, req.GetMailId())
	if err != nil {
		response.ErrCode = proto.Int32(int32(common.ErrCode_EC_FAIL))
		response.ErrDesc = proto.String("失败")
		logrus.Errorln("ProcessSetReadTagReq err:", err)
		return nil
	}

	logrus.Debugln("ProcessSetReadTagReq resp", response)
	return ret
}


// 删除邮件请求
func ProcessDelMailReq(playerID uint64, header *steve_proto_gaterpc.Header, req mailserver.MailSvrDelMailReq) (ret []exchanger.ResponseMsg) {

	logrus.Debugln("ProcessDelMailReq req", req)

	response := &mailserver.MailSvrDelMailRsp{
		ErrCode: proto.Int32(0),
		ErrDesc: proto.String("成功"),
	}
	//
	ret = []exchanger.ResponseMsg{{
		MsgID: uint32(msgid.MsgID_MAILSVR_DEL_MAIL_RSP),
		Body:  response,
	}}

	err := logic.DelMail(playerID, req.GetMailId())
	if err != nil {
		response.ErrCode = proto.Int32(int32(common.ErrCode_EC_FAIL))
		response.ErrDesc = proto.String("失败")
		logrus.Errorln("ProcessDelMailReq err:", err)
		return nil
	}

	logrus.Debugln("ProcessDelMailReq resp", response)
	return ret
}


// 领取附件奖励请求
func ProcessAwardAttachReq(playerID uint64, header *steve_proto_gaterpc.Header, req mailserver.MailSvrAwardAttachReq) (ret []exchanger.ResponseMsg) {

	logrus.Debugln("ProcessAwardAttachReq req", req)

	response := &mailserver.MailSvrAwardAttachRsp{
		ErrCode: proto.Int32(0),
		ErrDesc: proto.String("成功"),
	}
	//
	ret = []exchanger.ResponseMsg{{
		MsgID: uint32(msgid.MsgID_MAILSVR_AWARD_ATTACH_RSP),
		Body:  response,
	}}

	info, err := logic.AwardAttach(playerID, req.GetMailId())
	if err != nil {
		response.ErrCode = proto.Int32(int32(common.ErrCode_EC_FAIL))
		response.ErrDesc = proto.String("失败")
		logrus.Errorln("ProcessAwardAttachReq err:", err)
		return nil
	}
	response.AwardInfo = &info

	logrus.Debugln("ProcessAwardAttachReq resp", response)
	return ret
}





