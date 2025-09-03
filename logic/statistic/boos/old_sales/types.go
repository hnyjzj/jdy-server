package old_sales

import "jdy/logic/statistic/boos"

type Where struct {
	boos.Where

	Type Types `json:"type" label:"类型" find:"true" required:"true" sort:"1" type:"number" input:"radio" preset:"typeMap" binding:"required"` // 类型
}

type DataReq struct {
	boos.DataReq

	Type Types `json:"type" binding:"required"`
}
