package payments

import "errors"

type Types int

const (
	TypesSurplus Types = iota + 1 // 结余
	TypesIncome                   // 收入
	TypesExpense                  // 支出
)

var TypesMap = map[Types]string{
	TypesSurplus: "结余",
	TypesIncome:  "收入",
	TypesExpense: "支出",
}

func (p Types) ToMap() any {
	return TypesMap
}

func (p Types) InMap() error {
	if _, ok := TypesMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}
