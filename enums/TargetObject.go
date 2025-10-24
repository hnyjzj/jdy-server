package enums

import "errors"

/* 销售目标统计对象 */
// 分组、个人
type TargetObject int

const (
	TargetObjectGroup    TargetObject = iota + 1 // 分组
	TargetObjectPersonal                         // 个人
)

var TargetObjectMap = map[TargetObject]string{
	TargetObjectGroup:    "分组",
	TargetObjectPersonal: "个人",
}

func (p TargetObject) ToMap() any {
	return TargetObjectMap
}

func (p TargetObject) InMap() error {
	if _, ok := TargetObjectMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}

func (p TargetObject) String() string {
	return TargetObjectMap[p]
}
