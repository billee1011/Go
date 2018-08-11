package define

/*
 功能： 基础结构和常量定义
 作者： SkyWang
 日期： 2018-7-24

 */


// 邮件状态: 未发送=0＞审核中=1＞已审核=2＞发送中=3＞发送结束=4＞已拒绝=5＞已撤回=6＞已失效=7
const (
	StateNoSend  int8 = 0			//  未发送
	StateChecking  int8 = 1		//  审核中
	StateChecked  int8 = 2		//  已审核
	StateSending  int8 = 3		//  发送中
	StateSended  int8 = 4			//  发送结束
	StateReject  int8 = 5			//  已拒绝
	StateBack   int8 = 6			//  已撤回
	StateDelete  int8 = 7			//  已失效
)

const (
	SendAll = 0					// 发送给所有玩家
	SendList = 1					// 指定玩家列表
)
