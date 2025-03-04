package enums

import "errors"

/* 调拨原因 */
// 全部、旧料返厂、特殊渠道售出、其他
type ProductAllocateReason int

const (
	ProductAllocateReasonAll            ProductAllocateReason = iota // 全部
	ProductAllocateReasonOldFactory                                  // 旧料返厂
	ProductAllocateReasonSpecialChannel                              // 特殊渠道售出
	ProductAllocateReasonOther                                       // 其他
)

var ProductAllocateReasonMap = map[ProductAllocateReason]string{
	ProductAllocateReasonAll:            "全部",
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
