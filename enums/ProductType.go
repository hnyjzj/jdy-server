package enums

import "errors"

/* 类型 */
// 全部、成品、旧料
type ProductType int

const (
	ProductTypeAll      ProductType = iota // 全部
	ProductTypeFinished                    // 成品
	ProductTypeOld                         // 旧料
)

var ProductTypeMap = map[ProductType]string{
	ProductTypeAll:      "全部",
	ProductTypeFinished: "成品",
	ProductTypeOld:      "旧料",
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
