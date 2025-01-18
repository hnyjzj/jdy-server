package enums

import "errors"

/* 状态 */
// 全部、正常、已报损、已调拨、已出售、已定出
type ProductStatus int

const (
	ProductStatusAll      ProductStatus = iota // 全部
	ProductStatusNormal                        // 正常
	ProductStatusDamage                        // 已报损
	ProductStatusAllocate                      // 已调拨
	ProductStatusSold                          // 已出售
	ProductStatusReturn                        // 已定出
)

var ProductStatusMap = map[ProductStatus]string{
	ProductStatusAll:      "全部",
	ProductStatusNormal:   "正常",
	ProductStatusDamage:   "已报损",
	ProductStatusAllocate: "已调拨",
	ProductStatusSold:     "已出售",
	ProductStatusReturn:   "已定出",
}

// 判断状态是否可以转换
func (p ProductStatus) CanTransitionTo(newStatus ProductStatus) error {
	transitions := map[ProductStatus][]ProductStatus{
		ProductStatusNormal: {
			ProductStatusDamage,
			ProductStatusAllocate,
			ProductStatusSold,
			ProductStatusReturn,
		},
		ProductStatusDamage: {
			ProductStatusNormal,
			ProductStatusAllocate,
		},
		ProductStatusAllocate: {
			ProductStatusNormal,
			ProductStatusDamage,
		},
		ProductStatusSold: {
			ProductStatusNormal,
			ProductStatusReturn,
		},
		ProductStatusReturn: {
			ProductStatusNormal,
			ProductStatusDamage,
		},
	}
	if allowed, ok := transitions[p]; ok {
		for _, status := range allowed {
			if status == newStatus {
				return nil
			}
		}
	}
	return errors.New("非法的状态转换")
}

func (p ProductStatus) ToMap() any {
	return ProductStatusMap
}

func (p ProductStatus) InMap() error {
	if _, ok := ProductStatusMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}
