package enums

import "errors"

/* 打印类型 */
// 销售单、退货单、订金单、维修单、其他
type PrintType int

const (
	PrintTypeSales   PrintType = iota + 1 // 销售单
	PrintTypeReturn                       // 退货单
	PrintTypeDeposit                      // 订金单
	PrintTypeRepair                       // 维修单
	PrintTypeOthers                       // 其他
)

var PrintTypeMap = map[PrintType]string{
	PrintTypeSales:   "销售单",
	PrintTypeReturn:  "退货单",
	PrintTypeDeposit: "订金单",
	PrintTypeRepair:  "维修单",
	PrintTypeOthers:  "其他",
}

func (p PrintType) ToMap() any {
	return PrintTypeMap
}

func (p PrintType) InMap() error {
	if _, ok := PrintTypeMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}

func (p PrintType) String() string {
	return PrintTypeMap[p]
}
