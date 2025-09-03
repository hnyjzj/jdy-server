package enums

import (
	"errors"
)

/* 支付方式 */
// 微信、支付宝、聚合支付、现金、刷卡、商场、储值卡、其他
type OrderPaymentMethod int

const (
	OrderPaymentMethodWeChat     OrderPaymentMethod = iota + 1 // 微信
	OrderPaymentMethodAlipay                                   // 支付宝
	OrderPaymentMethodJuhe                                     // 聚合支付
	OrderPaymentMethodCash                                     // 现金
	OrderPaymentMethodSwipe                                    // 刷卡
	OrderPaymentMethodMall                                     // 商场
	OrderPaymentMethodStoredCard                               // 储值卡
	OrderPaymentMethodOther                                    // 其他

)

var OrderPaymentMethodMap = map[OrderPaymentMethod]string{
	OrderPaymentMethodWeChat:     "微信",
	OrderPaymentMethodAlipay:     "支付宝",
	OrderPaymentMethodJuhe:       "聚合支付",
	OrderPaymentMethodCash:       "现金",
	OrderPaymentMethodSwipe:      "刷卡",
	OrderPaymentMethodMall:       "商场",
	OrderPaymentMethodStoredCard: "储值卡",
	OrderPaymentMethodOther:      "其他",
}

func (p OrderPaymentMethod) ToMap() any {
	return OrderPaymentMethodMap
}

func (p OrderPaymentMethod) InMap() error {
	if _, ok := OrderPaymentMethodMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}

func (p OrderPaymentMethod) String() string {
	return OrderPaymentMethodMap[p]
}
