package types

type RemarkWhere struct {
	Id      string `json:"id" label:"备注ID" create:"false" find:"false" update:"true" sort:"1" type:"string" input:"text" required:"false"`        // 备注ID
	StoreId string `json:"store_id" label:"店铺ID" create:"true" find:"false" update:"false" sort:"2" type:"string" input:"text" required:"false"`  // 店铺ID
	Content string `json:"content" label:"备注内容" create:"true" find:"true" update:"true" sort:"3" type:"string" input:"textarea" required:"false"` // 备注内容
}

type RemarkCreateReq struct {
	StoreId string `json:"store_id" binding:"required"` // 店铺ID
	Content string `json:"content" binding:"required"`  // 备注内容
}

type RemarkListReq struct {
	PageReq
	Where RemarkWhere `json:"where"`
}

type RemarkUpdateReq struct {
	Id      string `json:"id" binding:"required"`      // 备注ID
	Content string `json:"content" binding:"required"` // 备注内容
}

type RemarkDeleteReq struct {
	Id string `json:"id" binding:"required"` // 备注ID
}
