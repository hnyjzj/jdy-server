package enums

import "errors"

/* 调拨状态 */
// 草稿、在途、已完成、已取消
type ProductAllocateStatus int

const (
	ProductAllocateStatusDraft     ProductAllocateStatus = iota + 1 // 草稿
	ProductAllocateStatusOnTheWay                                   // 在途
	ProductAllocateStatusCompleted                                  // 已完成
	ProductAllocateStatusCancelled                                  // 已取消
)

var ProductAllocateStatusMap = map[ProductAllocateStatus]string{
	ProductAllocateStatusDraft:     "草稿",
	ProductAllocateStatusOnTheWay:  "在途",
	ProductAllocateStatusCompleted: "已完成",
	ProductAllocateStatusCancelled: "已取消",
}

func (p ProductAllocateStatus) ToMap() any {
	return ProductAllocateStatusMap
}

func (p ProductAllocateStatus) InMap() error {
	if _, ok := ProductAllocateStatusMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}
