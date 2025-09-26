package old_exchange

import (
	"errors"
)

type Types int

const (
	TypesRecyclePrice Types = iota + 1 // 抵值
	TypesCount                         // 件数
	TypesWeightMetal                   // 金重
)

var TypesMap = map[Types]string{
	TypesRecyclePrice: "抵值",
	TypesCount:        "件数",
	TypesWeightMetal:  "金重",
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
