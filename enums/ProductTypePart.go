package enums

import (
	"errors"
)

/* 配件类型 */
// 配件、物料、赠品、商品
type ProductTypePart int

const (
	ProductTypePartAccessory ProductTypePart = iota + 1 // 配件
	ProductTypePartMaterial                             // 物料
	ProductTypePartGift                                 // 赠品
	ProductTypePartCommodity                            // 商品
)

var ProductTypePartMap = map[ProductTypePart]string{
	ProductTypePartAccessory: "配件",
	ProductTypePartMaterial:  "物料",
	ProductTypePartGift:      "赠品",
	ProductTypePartCommodity: "商品",
}

func (p ProductTypePart) ToMap() any {
	return ProductTypePartMap
}

func (p ProductTypePart) InMap() error {
	if _, ok := ProductTypePartMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}
