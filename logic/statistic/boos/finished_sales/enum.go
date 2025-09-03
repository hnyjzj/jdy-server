package finished_sales

import (
	"errors"
)

type Types int

const (
	TypesSales       Types = iota + 1 // 销售
	TypesCount                        // 件数
	TypesWeightMetal                  // 金重
)

var TypesMap = map[Types]string{
	TypesSales:       "销售",
	TypesCount:       "件数",
	TypesWeightMetal: "金重",
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
