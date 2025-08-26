package enums

import (
	"errors"
)

/* 回收方式 */
// 按克、按件
type ProductRecycleMethod int

const (
	ProductRecycleMethod_KG    ProductRecycleMethod = iota + 1 // 按克
	ProductRecycleMethod_PIECE                                 // 按件
)

var ProductRecycleMethodMap = map[ProductRecycleMethod]string{
	ProductRecycleMethod_KG:    "按克",
	ProductRecycleMethod_PIECE: "按件",
}

func (p ProductRecycleMethod) ToMap() any {
	return ProductRecycleMethodMap
}

func (p ProductRecycleMethod) InMap() error {
	if _, ok := ProductRecycleMethodMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}

func (p ProductRecycleMethod) String() string {
	return ProductRecycleMethodMap[p]
}
