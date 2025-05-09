package enums

import (
	"errors"
)

// 收支来源
// 销售-收款、销售-退货、定金单-收款、定金单-退款、其他收支-其他、其他收支-定金、其他收支-手续费、其他收支-汇回公司
type FinanceSourceOther int

const (
	FinanceSourceOtherOtherReceive FinanceSourceOther = iota + 1 // 其他收支-其他
	FinanceSourceOtherOtherDeposit                               // 其他收支-定金
	FinanceSourceOtherOtherFee                                   // 其他收支-手续费
	FinanceSourceOtherOtherReturn                                // 其他收支-汇回公司
)

var FinanceSourceOtherMap = map[FinanceSourceOther]string{
	FinanceSourceOtherOtherReceive: "其他收支-其他",
	FinanceSourceOtherOtherDeposit: "其他收支-定金",
	FinanceSourceOtherOtherFee:     "其他收支-手续费",
	FinanceSourceOtherOtherReturn:  "其他收支-汇回公司",
}

func (p FinanceSourceOther) ToMap() any {
	return FinanceSourceOtherMap
}

func (p FinanceSourceOther) InMap() error {
	if _, ok := FinanceSourceOtherMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}
