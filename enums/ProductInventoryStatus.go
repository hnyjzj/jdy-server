package enums

import (
	"errors"
	"slices"
)

/* 盘点状态 */
// 草稿、盘点中、待验证、盘点完成、盘点异常、盘点取消、异常修复
type ProductInventoryStatus int

const (
	ProductInventoryStatusDraft          ProductInventoryStatus = iota + 1 // 草稿
	ProductInventoryStatusInventorying                                     // 盘点中
	ProductInventoryStatusToBeVerified                                     // 待验证
	ProductInventoryStatusCompleted                                        // 盘点完成
	ProductInventoryStatusAbnormal                                         // 盘点异常
	ProductInventoryStatusCancelled                                        // 盘点取消
	ProductInventoryStatusAbnormalRepair                                   // 异常修复
)

var ProductInventoryStatusMap = map[ProductInventoryStatus]string{
	ProductInventoryStatusDraft:          "草稿",
	ProductInventoryStatusInventorying:   "盘点中",
	ProductInventoryStatusToBeVerified:   "待验证",
	ProductInventoryStatusCompleted:      "盘点完成",
	ProductInventoryStatusAbnormal:       "盘点异常",
	ProductInventoryStatusCancelled:      "盘点取消",
	ProductInventoryStatusAbnormalRepair: "异常修复",
}

func (p ProductInventoryStatus) ToMap() any {
	return ProductInventoryStatusMap
}

func (p ProductInventoryStatus) InMap() error {
	if _, ok := ProductInventoryStatusMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}

func (p ProductInventoryStatus) String() string {
	return ProductInventoryStatusMap[p]
}

func (p ProductInventoryStatus) IsOver() bool {
	return slices.Contains([]ProductInventoryStatus{
		ProductInventoryStatusToBeVerified,
		ProductInventoryStatusCompleted,
		ProductInventoryStatusAbnormal,
		ProductInventoryStatusCancelled,
		ProductInventoryStatusAbnormalRepair,
	}, p)
}

// 判断状态是否可以转换
func (p ProductInventoryStatus) CanTransitionTo(n ProductInventoryStatus) error {
	transitions := map[ProductInventoryStatus][]ProductInventoryStatus{
		ProductInventoryStatusDraft: { // 草稿->
			ProductInventoryStatusInventorying, // 盘点中
			ProductInventoryStatusCancelled,    // 盘点取消
		},
		ProductInventoryStatusInventorying: { // 盘点中->
			ProductInventoryStatusToBeVerified, // 待验证
			ProductInventoryStatusAbnormal,     // 盘点异常
			ProductInventoryStatusCancelled,    // 盘点取消
		},
		ProductInventoryStatusToBeVerified: { // 待验证->
			ProductInventoryStatusCompleted, // 盘点完成
			ProductInventoryStatusAbnormal,  // 盘点异常
			ProductInventoryStatusCancelled, // 盘点取消
		},
		ProductInventoryStatusAbnormal: { // 盘点异常->
			ProductInventoryStatusInventorying,   // 盘点中
			ProductInventoryStatusCancelled,      // 盘点取消
			ProductInventoryStatusAbnormalRepair, // 异常修复
		},
		ProductInventoryStatusCancelled: { // 盘点取消->
			ProductInventoryStatusInventorying, // 盘点中
		},
	}

	if allowed, ok := transitions[p]; ok {
		if slices.Contains(allowed, n) {
			return nil
		}
	}

	return errors.New("非法的状态转换")
}

// 权限判断
func (p ProductInventoryStatus) CanEdit(status ProductInventoryStatus, StaffId string, InventoryPersonId []string, InspectorId string) bool {
	type Condition struct {
		P, S ProductInventoryStatus
	}

	condition := Condition{p, status}

	switch condition {
	case Condition{ProductInventoryStatusDraft, ProductInventoryStatusInventorying}: // 开始盘点: 草稿->盘点中 : 盘点人
		return slices.Contains(InventoryPersonId, StaffId)
	case Condition{ProductInventoryStatusInventorying, ProductInventoryStatusInventorying}: // 继续盘点: 盘点中->盘点中 : 盘点人
		return slices.Contains(InventoryPersonId, StaffId)
	case Condition{ProductInventoryStatusDraft, ProductInventoryStatusCancelled}: // 取消盘点: 草稿->盘点取消 : 盘点人
		return slices.Contains(InventoryPersonId, StaffId)
	case Condition{ProductInventoryStatusInventorying, ProductInventoryStatusCancelled}: // 取消盘点: 盘点中->盘点取消 : 盘点人
		return slices.Contains(InventoryPersonId, StaffId)
	case Condition{ProductInventoryStatusInventorying, ProductInventoryStatusToBeVerified}: // 开始盘点/结束盘点: 盘点中->待验证 : 盘点人
		return slices.Contains(InventoryPersonId, StaffId)
	case Condition{ProductInventoryStatusToBeVerified, ProductInventoryStatusCompleted}: // 盘点完成: 待验证->盘点完成 : 监盘人
		return StaffId == InspectorId
	case Condition{ProductInventoryStatusToBeVerified, ProductInventoryStatusAbnormal}: // 盘点异常: 待验证->盘点异常 : 监盘人
		return StaffId == InspectorId
	case Condition{ProductInventoryStatusAbnormal, ProductInventoryStatusAbnormalRepair}: // 异常修复: 盘点异常->异常修复 : 监盘人
		return StaffId == InspectorId
	}

	return false
}
