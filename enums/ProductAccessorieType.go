package enums

import (
	"errors"
)

/* 配件类型 */
// 配件、物料、赠品、商品
type ProductAccessorieType int

const (
	ProductAccessorieTypeAccessory ProductAccessorieType = iota + 1 // 配件
	ProductAccessorieTypeMaterial                                   // 物料
	ProductAccessorieTypeGift                                       // 赠品
	ProductAccessorieTypeCommodity                                  // 商品
)

var ProductAccessorieTypeMap = map[ProductAccessorieType]string{
	ProductAccessorieTypeAccessory: "配件",
	ProductAccessorieTypeMaterial:  "物料",
	ProductAccessorieTypeGift:      "赠品",
	ProductAccessorieTypeCommodity: "商品",
}

func (p ProductAccessorieType) ToMap() any {
	return ProductAccessorieTypeMap
}

func (p ProductAccessorieType) InMap() error {
	if _, ok := ProductAccessorieTypeMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}
