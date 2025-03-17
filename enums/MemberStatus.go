package enums

import "errors"

/* 会员状态 */
// 待完善、正常、禁用、员工删除、用户删除
type MemberStatus int

const (
	MemberStatusPending        MemberStatus = iota + 1 // 待完善
	MemberStatusNormal                                 // 正常
	MemberStatusDisable                                // 禁用
	MemberStatusEmployeeDelete                         // 员工删除
	MemberStatusUserDelete                             // 用户删除

)

var MemberStatusMap = map[MemberStatus]string{
	MemberStatusPending:        "待完善",
	MemberStatusNormal:         "正常",
	MemberStatusDisable:        "禁用",
	MemberStatusEmployeeDelete: "员工删除",
	MemberStatusUserDelete:     "用户删除",
}

func (p MemberStatus) ToMap() any {
	return MemberStatusMap
}

func (p MemberStatus) InMap() error {
	if _, ok := MemberStatusMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}
