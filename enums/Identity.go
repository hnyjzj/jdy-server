package enums

import "errors"

/* 身份 */
// 店员、店长、区域经理、总部、管理员、超级管理员
type Identity int

const (
	IdentityClerk        Identity = iota + 1 // 店员
	IdentityShopkeeper                       // 店长
	IdentityAreaManager                      // 区域经理
	IdentityHeadquarters                     // 总部
	IdentityAdmin                            // 管理员
	IdentitySuperAdmin                       // 超级管理员
)

var IdentityMap = map[Identity]string{
	IdentityClerk:        "店员",
	IdentityShopkeeper:   "店长",
	IdentityAreaManager:  "区域经理",
	IdentityHeadquarters: "总部",
	IdentityAdmin:        "管理员",
	IdentitySuperAdmin:   "超级管理员",
}

func (p Identity) ToMap() any {
	return IdentityMap
}

func (p Identity) InMap() error {
	if _, ok := IdentityMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}
