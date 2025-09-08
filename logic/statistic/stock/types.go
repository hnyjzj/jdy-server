package stock

type Where struct {
	StoreId string `json:"store_id" label:"门店" find:"true" required:"true" sort:"1" type:"string" input:"text"`
	Day     string `json:"day" label:"日期" find:"true" required:"true" sort:"2" type:"string" input:"day"` // 日期
}

type DataReq struct {
	StoreId string `json:"store_id"`               // 门店
	Day     string `json:"day" binding:"required"` // 日期
}
