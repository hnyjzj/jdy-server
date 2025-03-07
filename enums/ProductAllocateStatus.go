package enums

import "errors"

/* 调拨状态 */
// 盘点中、调拨中、待接收、已完成、已取消
type ProductAllocateStatus int

const (
	ProductAllocateStatusInventory ProductAllocateStatus = iota + 1 // 盘点中 (正在添加产品)
	ProductAllocateStatusAllocate                                   // 调拨中 (确认调拨)
	ProductAllocateStatusCompleted                                  // 已完成 (已接收)
	ProductAllocateStatusCanceled                                   // 已取消 (取消调拨)
)

var ProductAllocateStatusMap = map[ProductAllocateStatus]string{
	ProductAllocateStatusInventory: "盘点中",
	ProductAllocateStatusAllocate:  "调拨中",
	ProductAllocateStatusCompleted: "已完成",
	ProductAllocateStatusCanceled:  "已取消",
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
