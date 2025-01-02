package enums

import "errors"

/* 会员状态 */
// 全部、正常、禁用
type MemberStatus int

const (
	MemberStatusAll     MemberStatus = iota // 全部
	MemberStatusNormal                      // 正常
	MemberStatusDisable                     // 禁用
)

var MemberStatusMap = map[MemberStatus]string{
	MemberStatusAll:     "全部",
	MemberStatusNormal:  "正常",
	MemberStatusDisable: "禁用",
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
