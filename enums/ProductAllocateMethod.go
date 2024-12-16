package enums

import "errors"

/* 调拨方式 */
// 门店调拨、调拨出库
type ProductAllocateMethod int

const (
	ProductAllocateMethodStore ProductAllocateMethod = iota + 1 // 门店调拨
	ProductAllocateMethodOut                                    // 调拨出库
)

var ProductAllocateMethodMap = map[ProductAllocateMethod]string{
	ProductAllocateMethodStore: "门店调拨",
	ProductAllocateMethodOut:   "调拨出库",
}

func (p ProductAllocateMethod) ToMap() any {
	return ProductAllocateMethodMap
}

func (p ProductAllocateMethod) InMap() error {
	if _, ok := ProductAllocateMethodMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}
