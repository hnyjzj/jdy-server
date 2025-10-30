package enums

import "errors"

/* 销售目标统计方式 */
// 金额、件数
type TargetMethod int

const (
	TargetMethodAmount   TargetMethod = iota + 1 // 金额
	TargetMethodQuantity                         // 件数
)

var TargetMethodMap = map[TargetMethod]string{
	TargetMethodAmount:   "金额",
	TargetMethodQuantity: "件数",
}

func (p TargetMethod) ToMap() any {
	return TargetMethodMap
}

func (p TargetMethod) InMap() error {
	if _, ok := TargetMethodMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}

func (p TargetMethod) String() string {
	return TargetMethodMap[p]
}
