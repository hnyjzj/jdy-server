package enums

import (
	"errors"
)

/* 类型 */
// 全部、成品、旧料、配件
type ProductType int

const (
	ProductTypeAll         ProductType = iota // 全部
	ProductTypeFinished                       // 成品
	ProductTypeOld                            // 旧料
	ProductTypeAccessories                    // 配件
)

var ProductTypeMap = map[ProductType]string{
	ProductTypeAll:         "全部",
	ProductTypeFinished:    "成品",
	ProductTypeOld:         "旧料",
	ProductTypeAccessories: "配件",
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

// 判断状态是否可以转换
func (p ProductType) CanTransitionTo(n ProductType) error {
	transitions := map[ProductType][]ProductType{
		ProductTypeFinished: {ProductTypeOld},
		ProductTypeOld:      {ProductTypeFinished},
	}

	if allowed, ok := transitions[p]; ok {
		for _, o := range allowed {
			if o == n {
				return nil
			}
		}
	}

	return errors.New("非法的状态转换")
}
