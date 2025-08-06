package types

import (
	"jdy/enums"
	"time"
)

type ProductHistoryWhere struct {
	Code    string                `json:"code" label:"条码" input:"text" type:"string" find:"true" sort:"1" required:"false"`                      // 产品
	Type    enums.ProductTypeUsed `json:"type" label:"产品类型" input:"select" type:"number" find:"true" sort:"2" required:"false" preset:"typeMap"` // 产品类型
	StoreId string                `json:"store_id" label:"门店" input:"text" type:"string" find:"false" sort:"3" required:"false"`                 // 门店
	Action  enums.ProductAction   `json:"action" label:"操作" input:"select" type:"number" find:"true" sort:"4" required:"false" preset:"typeMap"` // 操作

	StartTime *time.Time `json:"start_time" label:"开始时间" input:"date" type:"date" find:"true" sort:"5" required:"false"` // 开始时间
	EndTime   *time.Time `json:"end_time" label:"结束时间" input:"date" type:"date" find:"true" sort:"6" required:"false"`   // 结束时间
}

type ProductAccessorieHistoryWhere struct {
	Name    string              `json:"name" label:"配件名称" input:"text" type:"string" find:"true" sort:"1" required:"false"`                    // 产品
	StoreId string              `json:"store_id" label:"门店" input:"text" type:"string" find:"false" sort:"2" required:"false"`                 // 门店
	Action  enums.ProductAction `json:"action" label:"操作" input:"select" type:"number" find:"true" sort:"3" required:"false" preset:"typeMap"` // 操作

	StartTime *time.Time `json:"start_time" label:"开始时间" input:"date" type:"date" find:"true" sort:"4" required:"false"` // 开始时间
	EndTime   *time.Time `json:"end_time" label:"结束时间" input:"date" type:"date" find:"true" sort:"5" required:"false"`   // 结束时间
}

type ProductHistoryListReq struct {
	PageReq
	Where ProductHistoryWhere `json:"where" binding:"required"`
}
type ProductAccessorieHistoryListReq struct {
	PageReq
	Where ProductAccessorieHistoryWhere `json:"where" binding:"required"`
}

type ProductHistoryInfoReq struct {
	Id string `json:"id" binding:"required"`
}
