package enums

import (
	"errors"
)

/* 部门 */
// 店、区域、总部
type Department int

const (
	DepartmentStore          Department = iota + 1 // 店
	DepartmentRegion                               // 区域
	DepartmentHeaderquarters                       // 总部
)

var DepartmentMap = map[Department]string{
	DepartmentStore:          "店",
	DepartmentRegion:         "区域",
	DepartmentHeaderquarters: "总部",
}

func (p Department) ToMap() any {
	return DepartmentMap
}

func (p Department) InMap() error {
	if _, ok := DepartmentMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}

func (p Department) String() string {
	return DepartmentMap[p]
}
