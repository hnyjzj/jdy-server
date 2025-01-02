package enums

import "errors"

/* 会员来源 */
// 全部、人工录入、企业微信
type MemberSource int

const (
	MemberSourceAll        MemberSource = iota // 全部
	MemberSourceStaff                          // 员工录入
	MemberSourceWechatWork                     // 企业微信
)

var MemberSourceMap = map[MemberSource]string{
	MemberSourceAll:        "全部",
	MemberSourceStaff:      "员工录入",
	MemberSourceWechatWork: "企业微信",
}

func (p MemberSource) ToMap() any {
	return MemberSourceMap
}

func (p MemberSource) InMap() error {
	if _, ok := MemberSourceMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}
