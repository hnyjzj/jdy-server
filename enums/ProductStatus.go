package enums

import "errors"

/* 状态 */
// 全部、正常、报损、旧料、调拨、已售、退货
type ProductStatus int

const (
	ProductStatusAll      ProductStatus = iota // 全部
	ProductStatusNormal                        // 正常
	ProductStatusDamage                        // 报损
	ProductStatusAllocate                      // 调拨
	ProductStatusSold                          // 已售
	ProductStatusReturn                        // 退货
)

var ProductStatusMap = map[ProductStatus]string{
	ProductStatusAll:      "全部",
	ProductStatusNormal:   "正常",
	ProductStatusDamage:   "报损",
	ProductStatusAllocate: "调拨",
	ProductStatusSold:     "已售",
	ProductStatusReturn:   "退货",
}

// 判断状态是否可以转换
func (p ProductStatus) CanTransitionTo(newStatus ProductStatus) error {
	transitions := map[ProductStatus][]ProductStatus{
		ProductStatusNormal: {
			ProductStatusDamage,
			ProductStatusAllocate,
			ProductStatusSold,
		},
		ProductStatusDamage: {
			ProductStatusNormal,
		},
		ProductStatusAllocate: {
			ProductStatusNormal,
			ProductStatusDamage,
		},
		ProductStatusSold: {
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
