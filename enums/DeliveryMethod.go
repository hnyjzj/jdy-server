package enums

import (
	"errors"
)

/* 收货方式 */
// 自提、邮寄
type DeliveryMethod int

const (
	DeliveryMethodSelfPickup DeliveryMethod = iota + 1 // 自提
	DeliveryMethodMail                                 // 邮寄
)

var DeliveryMethodMap = map[DeliveryMethod]string{
	DeliveryMethodSelfPickup: "自提",
	DeliveryMethodMail:       "邮寄",
}

func (p DeliveryMethod) ToMap() any {
	return DeliveryMethodMap
}

func (p DeliveryMethod) InMap() error {
	if _, ok := DeliveryMethodMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}
