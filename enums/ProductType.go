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

/* 产品类型(可操作) */
// 成品、旧料
type ProductTypeUsed int

const (
	ProductTypeUsedFinished ProductTypeUsed = iota + 1 // 成品
	ProductTypeUsedOld                                 // 旧料
)

var ProductTypeUsedMap = map[ProductTypeUsed]string{
	ProductTypeUsedFinished: "成品",
	ProductTypeUsedOld:      "旧料",
}

func (p ProductTypeUsed) ToMap() any {
	return ProductTypeUsedMap
}

func (p ProductTypeUsed) InMap() error {
	if _, ok := ProductTypeUsedMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}
