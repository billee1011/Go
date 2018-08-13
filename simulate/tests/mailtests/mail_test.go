package mailtests

import (
	"testing"
	"steve/simulate/utils"
	"github.com/stretchr/testify/assert"
	"steve/client_pb/msgid"
	"steve/simulate/global"
	"steve/client_pb/mailserver"
)


func Test_GetUnReadMailSum(t *testing.T) {

	reqCmd := msgid.MsgID_MAILSVR_GET_UNREAD_SUM_REQ
	rspCmd := msgid.MsgID_MAILSVR_GET_UNREAD_SUM_RSP
	req := &mailserver.MailSvrGetUnReadSumReq{}
	rsp := &mailserver.MailSvrGetUnReadSumRsp{}

	player, err := utils.LoginNewPlayer()
	assert.Nil(t, err)
	assert.NotNil(t, player)



	player.AddExpectors(rspCmd)
	player.GetClient().SendPackage(utils.CreateMsgHead(reqCmd), req)
	expector := player.GetExpector(rspCmd)


	assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, rsp))
	assert.Zero(t, rsp.GetErrCode())

	t.Logf("Test_GetUnReadMailSum win:", rsp)

}


var mailId uint64 = 0

func Test_GetMailList(t *testing.T) {

	reqCmd := msgid.MsgID_MAILSVR_GET_MAIL_LIST_REQ
	rspCmd := msgid.MsgID_MAILSVR_GET_MAIL_LIST_RSP
	req := &mailserver.MailSvrGetMailListReq{}
	rsp := &mailserver.MailSvrGetMailListRsp{}

	player, err := utils.LoginNewPlayer()
	assert.Nil(t, err)
	assert.NotNil(t, player)


	player.AddExpectors(rspCmd)
	player.GetClient().SendPackage(utils.CreateMsgHead(reqCmd), req)
	expector := player.GetExpector(rspCmd)

	assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, rsp))
	assert.Zero(t, rsp.GetErrCode())

	t.Logf("Test_GetMailList win:", rsp)

	if len(rsp.MailList) > 0 {
		mailId = rsp.MailList[0].GetMailId()
		t.Logf("Test_GetMailList mailId:", mailId)
		//getMailDetail(t,player, id)
	}

}

func Test_GetMailDetail(t *testing.T) {

	reqCmd := msgid.MsgID_MAILSVR_GET_MAIL_DETAIL_REQ
	rspCmd := msgid.MsgID_MAILSVR_GET_MAIL_DETAIL_RSP
	req := &mailserver.MailSvrGetMailDetailReq{}
	rsp := &mailserver.MailSvrGetMailDetailRsp{}
	req.MailId = &mailId

	player, err := utils.LoginNewPlayer()
	assert.Nil(t, err)
	assert.NotNil(t, player)

	player.AddExpectors(rspCmd)
	player.GetClient().SendPackage(utils.CreateMsgHead(reqCmd), req)
	expector := player.GetExpector(rspCmd)

	assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, rsp))
	assert.Zero(t, rsp.GetErrCode())

	t.Logf("getMailDetail win:", rsp)

}

func Test_SetMailReadTag(t *testing.T) {

	reqCmd := msgid.MsgID_MAILSVR_SET_READ_TAG_REQ
	rspCmd := msgid.MsgID_MAILSVR_SET_READ_TAG_RSP
	req := &mailserver.MailSvrSetReadTagReq{}
	rsp := &mailserver.MailSvrSetReadTagRsp{}
	req.MailId = &mailId

	player, err := utils.LoginNewPlayer()
	assert.Nil(t, err)
	assert.NotNil(t, player)


	player.AddExpectors(rspCmd)
	player.GetClient().SendPackage(utils.CreateMsgHead(reqCmd), req)
	expector := player.GetExpector(rspCmd)

	assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, rsp))
	assert.Zero(t, rsp.GetErrCode())

	t.Logf("Test_SetMailReadTag win:", rsp)


	reqCmd = msgid.MsgID_MAILSVR_AWARD_ATTACH_REQ
	rspCmd = msgid.MsgID_MAILSVR_AWARD_ATTACH_RSP
	req2 := &mailserver.MailSvrAwardAttachReq{}
	rsp2 := &mailserver.MailSvrAwardAttachRsp{}
	req2.MailId = &mailId
	player.AddExpectors(rspCmd)
	player.GetClient().SendPackage(utils.CreateMsgHead(reqCmd), req2)
	expector = player.GetExpector(rspCmd)

	assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, rsp2))
	assert.Zero(t, rsp2.GetErrCode())

	return

	reqCmd = msgid.MsgID_MAILSVR_DEL_MAIL_REQ
	rspCmd = msgid.MsgID_MAILSVR_DEL_MAIL_RSP
	req3 := &mailserver.MailSvrDelMailReq{}
	rsp3 := &mailserver.MailSvrDelMailRsp{}
	req3.MailId = &mailId
	player.AddExpectors(rspCmd)
	player.GetClient().SendPackage(utils.CreateMsgHead(reqCmd), req3)
	expector = player.GetExpector(rspCmd)

	assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, rsp3))
	assert.Zero(t, rsp3.GetErrCode())

}


