package types

type PageReq struct {
	Page  int `json:"page" form:"page" uri:"page" binding:"required"`
	Limit int `json:"limit" form:"limit" uri:"limit" binding:"required"`
}

type PageReqNon struct {
	Page  int `json:"page" form:"page" uri:"page"`
	Limit int `json:"limit" form:"limit" uri:"limit"`
}

type PageRes[T any] struct {
	Total int64 `json:"total"`
	List  []T   `json:"list"`
}
