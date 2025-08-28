package old_stock

type TitleRes struct {
	Title     string `json:"title"`
	Key       string `json:"key"`
	Width     string `json:"width"`
	Fixed     string `json:"fixed"`
	ClassName string `json:"className"`
	Align     string `json:"align"`
}

type Where struct {
	Type Types `json:"type" label:"类型" find:"true" required:"true" sort:"1" type:"number" input:"radio" preset:"typeMap" binding:"required"` // 类型
}

type DataReq struct {
	Type Types `json:"type" binding:"required"`
}
