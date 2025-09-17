package types

import (
	"errors"
	"jdy/enums"
	"time"
)

type ProductAccessorieAllocateCreateReq struct {
	Method           enums.ProductAccessorieAllocateMethod `json:"method" binding:"required"`        // 调拨类型
	FromStoreId      string                                `json:"from_store_id" binding:"required"` // 调出门店
	ToStoreId        string                                `json:"to_store_id"`                      // 调入门店
	ToRegionId       string                                `json:"to_region_id"`                     // 调入区域
	ToHeadquartersId string                                `json:"to_headquarters_id"`               // 调入总部
	Remark           string                                `json:"remark"`                           // 备注
}

func (req *ProductAccessorieAllocateCreateReq) Validate() error {
	if req.Method == enums.ProductAccessorieAllocateMethodStore && req.ToStoreId == "" {
		return errors.New("调拨门店不能为空")
	}
	if req.Method == enums.ProductAccessorieAllocateMethodOut && req.ToHeadquartersId == "" {
		return errors.New("调拨总部不能为空")
	}
	if req.Method == enums.ProductAccessorieAllocateMethodRegion && req.ToRegionId == "" {
		return errors.New("调拨门店不能为空")
	}

	return nil
}

type ProductAccessorieAllocateWhere struct {
	Id               string                                `json:"id" label:"调拨单号" input:"text" type:"string" find:"true" create:"false" info:"true" sort:"1" required:"false"`                                                                                      // 调拨ID
	Status           enums.ProductAllocateStatus           `json:"status" label:"调拨状态" input:"select" type:"number" find:"true" create:"false" info:"true" sort:"2" required:"true" preset:"typeMap"`                                                                // 调拨状态
	Method           enums.ProductAccessorieAllocateMethod `json:"method" label:"调拨类型" input:"select" type:"number" find:"true" create:"true" info:"true" sort:"3" required:"true" preset:"typeMap"`                                                                 // 调拨类型
	FromStoreId      string                                `json:"from_store_id" label:"调出门店" input:"search" type:"string" find:"true" create:"false" info:"true" sort:"4" required:"false"`                                                                         // 调出门店
	ToStoreId        string                                `json:"to_store_id" label:"调入门店" input:"search" type:"string" find:"true" create:"true" info:"true" sort:"5" required:"true"  condition:"[{\"key\":\"method\",\"operator\":\"=\",\"value\":1}]"`          // 调入门店
	ToRegionId       string                                `json:"to_region_id" label:"调入区域" input:"search" type:"string" find:"true" create:"true" info:"true" sort:"6" required:"true"  condition:"[{\"key\":\"method\",\"operator\":\"=\",\"value\":3}]"`         // 调入区域
	ToHeadquartersId string                                `json:"to_headquarters_id" label:"调入总部" input:"search" type:"string" find:"false" create:"true" info:"false" sort:"7" required:"true"  condition:"[{\"key\":\"method\",\"operator\":\"=\",\"value\":2}]"` // 调入总部

	ProductCount int64 `json:"product_count" label:"种类数" info:"true" sort:"7" required:"false"`
	ProductTotal int64 `json:"product_total" label:"总件数" info:"true" sort:"8" required:"false"`

	Remark string `json:"remark" label:"调拨原因" input:"text" type:"string" find:"true" create:"true" info:"true" sort:"9" required:"false"` // 调拨原因

	InitiatorId string `json:"initiator_id" label:"发起人" input:"search" type:"string" find:"true" info:"true" sort:"10" required:"false"` // 发起人
	ReceiverId  string `json:"receiver_id" label:"接收人" input:"search" type:"string" find:"true" info:"true" sort:"11" required:"false"`  // 接收人

	StoreId string `json:"store_id" label:"门店ID" input:"search" type:"string" find:"false" info:"false" sort:"12" required:"false"` // 门店ID

	CreatedAt string     `json:"created_at" label:"创建时间" info:"true" sort:"13" type:"date"`                               // 创建时间
	StartTime *time.Time `json:"start_time" label:"开始时间" input:"date" type:"date" find:"true" sort:"14" required:"false"` // 开始时间
	EndTime   *time.Time `json:"end_time" label:"结束时间" input:"date" type:"date" find:"true" sort:"15" required:"false"`   // 结束时间

}

type ProductAccessorieAllocateListReq struct {
	PageReq
	Where ProductAccessorieAllocateWhere `json:"where"`
}

type ProductAccessorieAllocateDetailsReq struct {
	Where ProductAccessorieAllocateWhere `json:"where"`
}

type ProductAccessorieAllocateInfoReq struct {
	PageReq
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
	Name     string `json:"name" binding:"required"`     // 产品名称
	Quantity int64  `json:"quantity" binding:"required"` // 数量
}

type ProductAccessorieAllocateRemoveReq struct {
	Id        string `json:"id" binding:"required"`         // 调拨单ID
	ProductId string `json:"product_id" binding:"required"` // 产品ID
}

type ProductAccessorieAllocateClearReq struct {
	Id string `json:"id" binding:"required"` // 调拨单ID
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
