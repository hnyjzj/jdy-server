package enums

import "errors"

/* 调拨类型 */
// 门店调拨、调拨出库、调至区域
type ProductAccessorieAllocateMethod int

const (
	ProductAccessorieAllocateMethodStore  ProductAccessorieAllocateMethod = iota + 1 // 门店调拨
	ProductAccessorieAllocateMethodOut                                               // 调拨出库
	ProductAccessorieAllocateMethodRegion                                            // 调至区域
)

var ProductAccessorieAllocateMethodMap = map[ProductAccessorieAllocateMethod]string{
	ProductAccessorieAllocateMethodStore:  "门店调拨",
	ProductAccessorieAllocateMethodOut:    "调拨出库",
	ProductAccessorieAllocateMethodRegion: "调至区域",
}

func (p ProductAccessorieAllocateMethod) ToMap() any {
	return ProductAccessorieAllocateMethodMap
}

func (p ProductAccessorieAllocateMethod) InMap() error {
	if _, ok := ProductAccessorieAllocateMethodMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}
