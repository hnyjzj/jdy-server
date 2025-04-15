package enums

import (
	"errors"
)

/* 小数点控制 */
// 不保留小数、保留一位小数、保留两位小数、保留三位小数
type DecimalPoint int

const (
	DecimalPointNone DecimalPoint = iota + 1 // 不保留小数
	DecimalPoint1                            // 保留一位小数
	DecimalPoint2                            // 保留两位小数
	DecimalPoint3                            // 保留三位小数
)

var DecimalPointMap = map[DecimalPoint]string{
	DecimalPointNone: "不保留小数",
	DecimalPoint1:    "保留一位小数",
	DecimalPoint2:    "保留两位小数",
	DecimalPoint3:    "保留三位小数",
}

func (p DecimalPoint) ToMap() any {
	return DecimalPointMap
}

func (p DecimalPoint) InMap() error {
	if _, ok := DecimalPointMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}
