package enums

import "errors"

/* 调拨状态 */
// 全部、盘点中、调拨中、待接收、已完成、已取消、已驳回
type ProductAllocateStatus int

const (
	ProductAllocateStatusAll       ProductAllocateStatus = iota // 全部
	ProductAllocateStatusInventory                              // 盘点中 (正在添加产品)
	ProductAllocateStatusAllocate                               // 调拨中 (确认调拨)
	ProductAllocateStatusCompleted                              // 已完成 (已接收)
	ProductAllocateStatusCanceled                               // 已取消 (取消调拨)
	ProductAllocateStatusRejected                               // 已驳回 (驳回调拨)
)

var ProductAllocateStatusMap = map[ProductAllocateStatus]string{
	ProductAllocateStatusAll:       "全部",
	ProductAllocateStatusInventory: "盘点中",
	ProductAllocateStatusAllocate:  "调拨中",
	ProductAllocateStatusCompleted: "已完成",
	ProductAllocateStatusCanceled:  "已取消",
	ProductAllocateStatusRejected:  "已驳回",
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
