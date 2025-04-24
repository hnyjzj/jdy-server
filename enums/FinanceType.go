package enums

import "errors"

// 收支类型
// 收入、支出
type FinanceType int

const (
	FinanceTypeIncome  FinanceType = iota + 1 // 收入
	FinanceTypeExpense                        // 支出
)

var FinanceTypeMap = map[FinanceType]string{
	FinanceTypeIncome:  "收入",
	FinanceTypeExpense: "支出",
}

func (p FinanceType) ToMap() any {
	return FinanceTypeMap
}

func (p FinanceType) InMap() error {
	if _, ok := FinanceTypeMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}
