package goods
/*
	功能： 物品定义: 用于邮件附件物品，商城物品，活动奖励物品，奖励列表是通过定义一系列物品作为奖励内容.
	作者: SkyWang
	日期: 2018-8-8
 */

 // 物品类型
 const (
	 GoodType_Prop = 0							// 道具物品
	GoodType_Gold = 1								// 货币物品

 )

// 物品定义
type Goods struct {
	GoodsType    int16  `json:"goodsType"`  		// 物品类型: 1=货币, 2=道具
	GoodsId		int32  `json:"goodsId"`  			// 物品ID
	GoodsNum  	int32  `json:"goodsNum"` 			// 物品数量
}
