package enums

import "errors"

/* 配件零售方式 */
// 计件、计重
type ProductAccessorieRetailType int

const (
	ProductAccessorieRetailTypePiece  ProductAccessorieRetailType = iota + 1 // 计件
	ProductAccessorieRetailTypeWeight                                        // 计重
)

var ProductAccessorieRetailTypeMap = map[ProductAccessorieRetailType]string{
	ProductAccessorieRetailTypePiece:  "计件",
	ProductAccessorieRetailTypeWeight: "计重",
}

func (p ProductAccessorieRetailType) ToMap() any {
	return ProductAccessorieRetailTypeMap
}

func (p ProductAccessorieRetailType) InMap() error {
	if _, ok := ProductAccessorieRetailTypeMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}
