package enums

import (
	"errors"
)

/* 回收类型 */
// 回收、兑换
type ProductRecycleType int

const (
	ProductRecycleTypeRecycle  ProductRecycleType = iota + 1 // 回收
	ProductRecycleTypeExchange                               // 兑换
)

var ProductRecycleTypeMap = map[ProductRecycleType]string{
	ProductRecycleTypeRecycle:  "回收",
	ProductRecycleTypeExchange: "兑换",
}

func (p ProductRecycleType) ToMap() any {
	return ProductRecycleTypeMap
}

func (p ProductRecycleType) InMap() error {
	if _, ok := ProductRecycleTypeMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}

func (p ProductRecycleType) String() string {
	return ProductRecycleTypeMap[p]
}
