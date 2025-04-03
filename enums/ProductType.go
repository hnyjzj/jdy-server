package enums

import (
	"errors"
)

/* 产品类型 */
// 成品、旧料、配件
type ProductType int

const (
	ProductTypeFinished   ProductType = iota + 1 // 成品
	ProductTypeOld                               // 旧料
	ProductTypeAccessorie                        // 配件
)

var ProductTypeMap = map[ProductType]string{
	ProductTypeFinished:   "成品",
	ProductTypeOld:        "旧料",
	ProductTypeAccessorie: "配件",
}

func (p ProductType) ToMap() any {
	return ProductTypeMap
}

func (p ProductType) InMap() error {
	if _, ok := ProductTypeMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}
