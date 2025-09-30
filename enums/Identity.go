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

var IdentityMapExternal = map[Identity]string{
	IdentityClerk:       "销售顾问",
	IdentityShopkeeper:  "店长",
	IdentityAreaManager: "区域经理",
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

func (p Identity) String() string {
	return IdentityMap[p]
}
func (p Identity) StringExternal() string {
	if value, ok := IdentityMapExternal[p]; ok {
		return value
	} else {
		return ""
	}
}

// 获取比当前身份小的身份
func (p Identity) GetMinMap() any {
	min := make(map[Identity]string)
	for key, value := range IdentityMap {
		if key < p {
			min[key] = value
		}
	}

	var last Identity
	for key := range IdentityMap {
		last = key
	}
	if p == last {
		min[last] = IdentityMap[last]
	}

	return min
}
