package enums

import (
	"errors"
)

/* 订单类型 */
// 全部、销售单、定金单、维修单、其他
type OrderType int

const (
	OrderTypeAll     OrderType = iota // 全部
	OrderTypeSales                    // 销售单
	OrderTypeDeposit                  // 定金单
	OrderTypeRepair                   // 维修单
	OrderTypeOthers                   // 其他
)

var OrderTypeMap = map[OrderType]string{
	OrderTypeAll:     "全部",
	OrderTypeSales:   "销售单",
	OrderTypeDeposit: "定金单",
	OrderTypeRepair:  "维修单",
	OrderTypeOthers:  "其他",
}

func (p OrderType) ToMap() any {
	return OrderTypeMap
}

func (p OrderType) InMap() error {
	if _, ok := OrderTypeMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}
