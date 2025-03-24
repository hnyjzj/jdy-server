package enums

import "errors"

/* 平台类型 */
// 账号、企业微信
type PlatformType string

const (
	PlatformTypeAccount PlatformType = "account"
	PlatformTypeWxWork  PlatformType = "wxwork"
)

var PlatformTypeMap = map[PlatformType]string{
	PlatformTypeAccount: "账号",
	PlatformTypeWxWork:  "企业微信",
}

func (p PlatformType) ToMap() any {
	return PlatformTypeMap
}

func (p PlatformType) InMap() error {
	if _, ok := PlatformTypeMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}

func (p PlatformType) String() string {
	return string(p)
}
