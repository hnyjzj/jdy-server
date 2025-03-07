package enums

import "errors"

/* 产品状态 */
// 正常、已报损、已调拨、已出售、已定出、盘点中
type ProductStatus int

const (
	ProductStatusNormal   ProductStatus = iota + 1 // 正常
	ProductStatusDamage                            // 已报损
	ProductStatusAllocate                          // 已调拨
	ProductStatusSold                              // 已出售
	ProductStatusReturn                            // 已定出
	ProductStatusCheck                             // 盘点中
)

var ProductStatusMap = map[ProductStatus]string{
	ProductStatusNormal:   "正常",
	ProductStatusDamage:   "已报损",
	ProductStatusAllocate: "已调拨",
	ProductStatusSold:     "已出售",
	ProductStatusReturn:   "已定出",
	ProductStatusCheck:    "盘点中",
}

// 判断状态是否可以转换
func (p ProductStatus) CanTransitionTo(newStatus ProductStatus) error {
	transitions := map[ProductStatus][]ProductStatus{
		ProductStatusNormal: {
			ProductStatusDamage,
			ProductStatusAllocate,
			ProductStatusSold,
			ProductStatusReturn,
			ProductStatusCheck,
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
		ProductStatusCheck: {
			ProductStatusNormal,
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
