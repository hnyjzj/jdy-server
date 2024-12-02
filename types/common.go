package types

type PageReq struct {
	Page  int `json:"page" form:"page" uri:"page"`
	Limit int `json:"limit" form:"limit" uri:"limit"`
}

type PageRes struct {
	Total int64 `json:"total"`
}
