package enums

import (
	"errors"
)

/* 员工日志类型 */
// 创建、更新、禁用、删除
type StaffLogType int

const (
	StaffLogTypeCreate  StaffLogType = iota + 1 // 创建
	StaffLogTypeUpdate                          // 更新
	StaffLogTypeDisable                         // 禁用
	StaffLogTypeDelete                          // 删除
)

var StaffLogTypeMap = map[StaffLogType]string{
	StaffLogTypeCreate:  "创建",
	StaffLogTypeUpdate:  "更新",
	StaffLogTypeDisable: "禁用",
	StaffLogTypeDelete:  "删除",
}

func (p StaffLogType) ToMap() any {
	return StaffLogTypeMap
}

func (p StaffLogType) InMap() error {
	if _, ok := StaffLogTypeMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}
