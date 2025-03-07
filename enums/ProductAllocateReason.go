package enums

import "errors"

/* 调拨原因 */
// 旧料返厂、特殊渠道售出、其他
type ProductAllocateReason int

const (
	ProductAllocateReasonOldFactory     ProductAllocateReason = iota + 1 // 旧料返厂
	ProductAllocateReasonSpecialChannel                                  // 特殊渠道售出
	ProductAllocateReasonOther                                           // 其他
)

var ProductAllocateReasonMap = map[ProductAllocateReason]string{
	ProductAllocateReasonOldFactory:     "旧料返厂",
	ProductAllocateReasonSpecialChannel: "特殊渠道售出",
	ProductAllocateReasonOther:          "其他",
}

func (p ProductAllocateReason) ToMap() any {
	return ProductAllocateReasonMap
}

func (p ProductAllocateReason) InMap() error {
	if _, ok := ProductAllocateReasonMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}
