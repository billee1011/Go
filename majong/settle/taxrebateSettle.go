package settle

import majongpb "steve/server_pb/majong"

// TaxRebateSettle 退税结算
type TaxRebateSettle struct {
}

//SettleTaxRebate 退税结算 呼叫转移的杠，不用退税
func (s *TaxRebateSettle) SettleTaxRebate(context *majong.MajongContext) []*majongpb.SettleInfo {
}
