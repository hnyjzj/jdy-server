package enums

import "errors"

// 收支来源
// 销售-收款、销售-退货、定金单-收款、定金单-退款、其他收支-其他、其他收支-定金、其他收支-手续费、其他收支-汇回公司
type FinanceSource int

const (
	FinanceSourceSaleReceive    FinanceSource = iota + 1 // 销售-收款
	FinanceSourceSaleRefund                              // 销售-退货
	FinanceSourceDepositReceive                          // 定金单-收款
	FinanceSourceDepositRefund                           // 定金单-退款
	FinanceSourceOtherReceive                            // 其他收支-其他
	FinanceSourceOtherDeposit                            // 其他收支-定金
	FinanceSourceOtherFee                                // 其他收支-手续费
	FinanceSourceOtherReturn                             // 其他收支-汇回公司
)

var FinanceSourceMap = map[FinanceSource]string{
	FinanceSourceSaleReceive:    "销售-收款",
	FinanceSourceSaleRefund:     "销售-退货",
	FinanceSourceDepositReceive: "定金单-收款",
	FinanceSourceDepositRefund:  "定金单-退款",
	FinanceSourceOtherReceive:   "其他收支-其他",
	FinanceSourceOtherDeposit:   "其他收支-定金",
	FinanceSourceOtherFee:       "其他收支-手续费",
	FinanceSourceOtherReturn:    "其他收支-汇回公司",
}

func (p FinanceSource) ToMap() any {
	return FinanceSourceMap
}

func (p FinanceSource) InMap() error {
	if _, ok := FinanceSourceMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}

func (p FinanceSource) String() string {
	return FinanceSourceMap[p]
}
