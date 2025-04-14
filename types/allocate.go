package types

import (
	"errors"
	"jdy/enums"
	"time"
)

type ProductAllocateCreateReq struct {
	Method      enums.ProductAllocateMethod `json:"method" binding:"required"` // 调拨类型
	Type        enums.ProductType           `json:"type" binding:"required"`   // 仓库类型
	Reason      enums.ProductAllocateReason `json:"reason" binding:"required"` // 调拨原因
	Remark      string                      `json:"remark"`                    // 备注
	FromStoreId string                      `json:"from_store_id"`             // 调出门店
	ToStoreId   string                      `json:"to_store_id"`               // 调入门店

	EnterId string `json:"enter_id"` // 入库单号
}

func (req *ProductAllocateCreateReq) Validate() error {
	if req.Method == enums.ProductAllocateMethodStore && req.ToStoreId == "" {
		return errors.New("调拨门店不能为空")
	}

	return nil
}

type ProductAllocateWhere struct {
	Method      enums.ProductAllocateMethod `json:"method" label:"调拨类型" input:"select" type:"number" find:"true" create:"true" sort:"1" required:"true" preset:"typeMap"` // 调拨类型
	Type        enums.ProductType           `json:"type" label:"仓库类型" input:"select" type:"number" find:"true" create:"true" sort:"2" required:"true" preset:"typeMap"`   // 仓库类型
	Reason      enums.ProductAllocateReason `json:"reason" label:"调拨原因" input:"select" type:"number" find:"true" create:"true" sort:"3" required:"true" preset:"typeMap"` // 调拨原因
	FromStoreId string                      `json:"from_store_id" label:"调出门店" input:"search" type:"string" sort:"4" required:"false"`                                    // 调出门店
	ToStoreId   string                      `json:"to_store_id" label:"调入门店" input:"search" type:"string" find:"true" create:"true"  sort:"4" required:"false"`           // 调入门店
	Status      enums.ProductAllocateStatus `json:"status" label:"调拨状态" input:"select" type:"number" find:"true" create:"true" sort:"5" required:"true" preset:"typeMap"` // 调拨状态

	StartTime *time.Time `json:"start_time" label:"开始时间" input:"date" type:"date" find:"true" sort:"6" required:"false"` // 开始时间
	EndTime   *time.Time `json:"end_time" label:"结束时间" input:"date" type:"date" find:"true" sort:"6" required:"false"`   // 结束时间

	StoreId string `json:"store_id" label:"门店" input:"search" type:"string" find:"true" sort:"7" required:"false"` // 门店
}

func (req *ProductAllocateWhere) Validate() error {
	if req.Method == enums.ProductAllocateMethodStore && req.ToStoreId == "" {
		return errors.New("调拨门店不能为空")
	}

	return nil
}

type ProductAllocateListReq struct {
	PageReq
	Where ProductAllocateWhere `json:"where"`
}

type ProductAllocateInfoReq struct {
	Id string `json:"id" binding:"required"`
}

type ProductAllocateAddReq struct {
	Id        string   `json:"id" binding:"required"` // 调拨单ID
	ProductId []string `json:"product_id" binding:"required"`
}

type ProductAllocateRemoveReq struct {
	Id        string `json:"id" binding:"required"` // 调拨单ID
	ProductId string `json:"product_id" binding:"required"`
}

type ProductAllocateConfirmReq struct {
	Id string `json:"id" binding:"required"` // 调拨单ID
}

type ProductAllocateCancelReq struct {
	Id string `json:"id" binding:"required"` // 调拨单ID
}

type ProductAllocateCompleteReq struct {
	Id string `json:"id" binding:"required"` // 调拨单ID
}
