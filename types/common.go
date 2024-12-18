package types

type PageReq struct {
	Page  int `json:"page" form:"page" uri:"page" binding:"required"`
	Limit int `json:"limit" form:"limit" uri:"limit" binding:"required"`
}

type PageRes[T any] struct {
	Total int64 `json:"total"`
	List  []T   `json:"list"`
}
