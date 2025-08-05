package types

import (
	"jdy/enums"
)

type ProductDamageReq struct {
	Code   string `json:"code" binding:"required"`   // 条码
	Reason string `json:"reason" binding:"required"` // 损坏原因
}

type ProductConversionReq struct {
	Id     string                `json:"id" binding:"required"`   // 产品ID
	Type   enums.ProductTypeUsed `json:"type" binding:"required"` // 仓库类型
	Remark string                `json:"remark"`                  // 备注
}
