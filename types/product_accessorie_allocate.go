package types

import (
	"errors"
	"jdy/enums"
	"time"
)

type ProductAccessorieAllocateCreateReq struct {
	Method      enums.ProductAllocateMethod `json:"method" binding:"required"` // 调拨类型
	FromStoreId string                      `json:"from_store_id"`             // 调出门店
	ToStoreId   string                      `json:"to_store_id"`               // 调入门店
	Remark      string                      `json:"remark"`                    // 备注
}

func (req *ProductAccessorieAllocateCreateReq) Validate() error {
	if req.Method == enums.ProductAllocateMethodStore && req.ToStoreId == "" {
		return errors.New("调拨门店不能为空")
	}
	if req.Method == enums.ProductAllocateMethodOut && req.Remark == "" {
		return errors.New("调拨备注不能为空")
	}

	return nil
}

type ProductAccessorieAllocateWhere struct {
	Status      enums.ProductAllocateStatus `json:"status" label:"调拨状态" input:"select" type:"number" find:"true" create:"false" sort:"1" required:"true" preset:"typeMap"`                                                       // 调拨状态
	Method      enums.ProductAllocateMethod `json:"method" label:"调拨类型" input:"select" type:"number" find:"true" create:"true" sort:"2" required:"true" preset:"typeMap"`                                                        // 调拨类型
	FromStoreId string                      `json:"from_store_id" label:"调出门店" input:"search" type:"string" find:"true" create:"false" sort:"3" required:"false"`                                                                // 调出门店
	ToStoreId   string                      `json:"to_store_id" label:"调入门店" input:"search" type:"string" find:"true" create:"true" sort:"4" required:"true"  condition:"[{\"key\":\"method\",\"operator\":\"=\",\"value\":1}]"` // 调入门店
	Remark      string                      `json:"remark" label:"调拨原因" input:"text" type:"string" find:"true" create:"true" sort:"5" required:"false"`                                                                          // 调拨原因

	StartTime *time.Time `json:"start_time" label:"开始时间" input:"date" type:"date" find:"true" sort:"6" required:"false"` // 开始时间
	EndTime   *time.Time `json:"end_time" label:"结束时间" input:"date" type:"date" find:"true" sort:"6" required:"false"`   // 结束时间
}

func (req *ProductAccessorieAllocateWhere) Validate() error {
	if req.Method == enums.ProductAllocateMethodStore && req.ToStoreId == "" {
		return errors.New("调拨门店不能为空")
	}

	return nil
}

type ProductAccessorieAllocateListReq struct {
	PageReq
	Where ProductAccessorieAllocateWhere `json:"where"`
}

type ProductAccessorieAllocateInfoReq struct {
	Id string `json:"id" binding:"required"` // 调拨单ID
}

type ProductAccessorieAllocateAddReq struct {
	Id       string                                `json:"id" binding:"required"`       // 调拨单ID
	Products []ProductAccessorieAllocateAddProduct `json:"products" binding:"required"` // 产品信息
}

type ProductAccessorieAllocateAddProduct struct {
	ProductId string `json:"product_id" binding:"required"` // 产品ID
	Quantity  int64  `json:"quantity" binding:"required"`   // 数量
}

type ProductAccessorieAllocateRemoveReq struct {
	Id        string `json:"id" binding:"required"`         // 调拨单ID
	ProductId string `json:"product_id" binding:"required"` // 产品ID
}

type ProductAccessorieAllocateConfirmReq struct {
	Id string `json:"id" binding:"required"` // 调拨单ID
}

type ProductAccessorieAllocateCancelReq struct {
	Id string `json:"id" binding:"required"` // 调拨单ID
}

type ProductAccessorieAllocateCompleteReq struct {
	Id string `json:"id" binding:"required"` // 调拨单ID
}
