package enums

import (
	"errors"
)

/* 盘点产品状态 */
// 应盘、实盘、盘盈、盘亏
type ProductInventoryProductStatus int

const (
	ProductInventoryProductStatusShould ProductInventoryProductStatus = iota // 应盘
	ProductInventoryProductStatusActual                                      // 实盘
	ProductInventoryProductStatusExtra                                       // 盘盈
	ProductInventoryProductStatusLoss                                        // 盘亏
)

var ProductInventoryProductStatusMap = map[ProductInventoryProductStatus]string{
	ProductInventoryProductStatusShould: "应盘",
	ProductInventoryProductStatusActual: "实盘",
	ProductInventoryProductStatusExtra:  "盘盈",
	ProductInventoryProductStatusLoss:   "盘亏",
}

func (p ProductInventoryProductStatus) ToMap() any {
	return ProductInventoryProductStatusMap
}

func (p ProductInventoryProductStatus) InMap() error {
	if _, ok := ProductInventoryProductStatusMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}
