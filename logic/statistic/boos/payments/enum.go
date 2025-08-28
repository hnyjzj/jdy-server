package payments

import "errors"

type Types int

const (
	TypesIncome  Types = iota + 1 // 收入
	TypesExpense                  // 支出
	TypesSurplus                  // 结余
)

var TypesMap = map[Types]string{
	TypesIncome:  "收入",
	TypesExpense: "支出",
	TypesSurplus: "结余",
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
