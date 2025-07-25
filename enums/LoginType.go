package enums

import (
	"errors"
)

/* 登录方式 */
// 手机、授权、扫码
type LoginType string

const (
	LoginTypePhone LoginType = "phone" // 手机
	LoginTypeAuth  LoginType = "auth"  // 授权
	LoginTypeScan  LoginType = "scan"  // 扫码
)

var LoginTypeMap = map[LoginType]string{
	LoginTypePhone: "手机",
	LoginTypeAuth:  "授权",
	LoginTypeScan:  "扫码",
}

func (p LoginType) ToMap() any {
	return LoginTypeMap
}

func (p LoginType) InMap() error {
	if _, ok := LoginTypeMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}

func (p LoginType) String() string {
	return string(p)
}
