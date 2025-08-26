package enums

import (
	"errors"
)

/* 回收类型 */
// 无、回收、兑换
type ProductRecycleType int

const (
	ProductRecycleTypeNone     ProductRecycleType = iota + 1 // 无
	ProductRecycleTypeRecycle                                // 回收
	ProductRecycleTypeExchange                               // 兑换
)

var ProductRecycleTypeMap = map[ProductRecycleType]string{
	ProductRecycleTypeNone:     "无",
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
	if err := p.InMap(); err != nil {
		return ProductRecycleTypeMap[ProductRecycleTypeNone]
	}
	return ProductRecycleTypeMap[p]
}
