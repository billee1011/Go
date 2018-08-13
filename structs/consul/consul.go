package consul
/*
	功能: consul请求接口，提供获取当前服务是否是主节点，当前主节点服务ID等consul相关接口定义.
	作者: SkyWang
	日期: 2018-8-11
 */

type Requester interface {

	// 是否是主节点
	IsMasterNode() bool
	// 获取当前主节点的服务ID
	GetMasterNodeId() (string, error)


}
