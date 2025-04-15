package enums

import (
	"errors"
)

/* 小数点控制 */
// 四舍五入、向下取整、向上取整
type Rounding int

const (
	RoundingRound Rounding = iota + 1 // 四舍五入
	RoundingDown                      // 向下取整
	RoundingUp                        // 向上取整
)

var RoundingMap = map[Rounding]string{
	RoundingRound: "四舍五入",
	RoundingDown:  "向下取整",
	RoundingUp:    "向上取整",
}

func (p Rounding) ToMap() any {
	return RoundingMap
}

func (p Rounding) InMap() error {
	if _, ok := RoundingMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}
