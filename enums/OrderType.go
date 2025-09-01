package enums

import (
	"errors"
)

/* 订单类型 */
// 销售单、定金单、维修单、其他收支单、退货单
type OrderType int

const (
	OrderTypeSales   OrderType = iota + 1 // 销售单
	OrderTypeDeposit                      // 定金单
	OrderTypeRepair                       // 维修单
	OrderTypeOthers                       // 其他收支单
	OrderTypeReturn                       // 退货单
)

var OrderTypeMap = map[OrderType]string{
	OrderTypeSales:   "销售单",
	OrderTypeDeposit: "定金单",
	OrderTypeRepair:  "维修单",
	OrderTypeOthers:  "其他收支单",
	OrderTypeReturn:  "退货单",
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
