package types

import (
	"errors"
	"jdy/enums"
	"time"
)

type ProductAccessorieAllocateCreateReq struct {
	Method      enums.ProductAllocateMethod `json:"method" binding:"required"`        // 调拨类型
	FromStoreId string                      `json:"from_store_id" binding:"required"` // 调出门店
	ToStoreId   string                      `json:"to_store_id"`                      // 调入门店
	Remark      string                      `json:"remark"`                           // 备注
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
	Id          string                      `json:"id" label:"调拨单号" input:"text" type:"string" find:"true" create:"false" sort:"1" required:"false"`                                                                             // 调拨ID
	Status      enums.ProductAllocateStatus `json:"status" label:"调拨状态" input:"select" type:"number" find:"true" create:"false" sort:"1" required:"true" preset:"typeMap"`                                                       // 调拨状态
	Method      enums.ProductAllocateMethod `json:"method" label:"调拨类型" input:"select" type:"number" find:"true" create:"true" sort:"2" required:"true" preset:"typeMap"`                                                        // 调拨类型
	FromStoreId string                      `json:"from_store_id" label:"调出门店" input:"search" type:"string" find:"false" create:"false" sort:"3" required:"false"`                                                               // 调出门店
	ToStoreId   string                      `json:"to_store_id" label:"调入门店" input:"search" type:"string" find:"true" create:"true" sort:"4" required:"true"  condition:"[{\"key\":\"method\",\"operator\":\"=\",\"value\":1}]"` // 调入门店
	Remark      string                      `json:"remark" label:"调拨原因" input:"text" type:"string" find:"true" create:"true" sort:"5" required:"false"`                                                                          // 调拨原因

	StartTime *time.Time `json:"start_time" label:"开始时间" input:"date" type:"date" find:"true" sort:"6" required:"false"` // 开始时间
	EndTime   *time.Time `json:"end_time" label:"结束时间" input:"date" type:"date" find:"true" sort:"6" required:"false"`   // 结束时间

	StoreId string `json:"store_id" label:"门店ID" input:"search" type:"string" find:"false" sort:"7" required:"false"` // 门店ID
}

func (req *ProductAccessorieAllocateWhere) Validate() error {
	if req.Method == enums.ProductAllocateMethodStore && req.ToStoreId == "" {
		return errors.New("调拨门店不能为空")
	}

	if req.Method == enums.ProductAllocateMethodOut && req.Remark == "" {
		return errors.New("调拨备注不能为空")
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

// 添加验证方法
func (req *ProductAccessorieAllocateAddReq) Validate() error {
	if len(req.Products) == 0 {
		return errors.New("产品列表不能为空")
	}

	for _, product := range req.Products {
		if product.Quantity <= 0 {
			return errors.New("产品数量必须大于0")
		}
	}

	return nil
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
