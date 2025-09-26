package old_stock

import (
	"errors"
)

type Types int

const (
	TypesCount        Types = iota + 1 // 件数
	TypesWeightMetal                   // 金重
	TypesRecyclePrice                  // 回收价
)

var TypesMap = map[Types]string{
	TypesCount:        "件数",
	TypesWeightMetal:  "金重",
	TypesRecyclePrice: "回收价",
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
