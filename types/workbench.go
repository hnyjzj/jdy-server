package types

type WorkbenchListRes struct {
	Id string `json:"id"`

	Title string `json:"title" `
	Icon  string `json:"icon" `
	Path  string `json:"path" `

	ParentId string             `json:"parent_id" `
	Children []WorkbenchListRes `json:"children" `
}

type WorkbenchSearchReq struct {
	Keyword string `json:"keyword" binding:"required"`
}

type WorkbenchAddReq struct {
	Title    string `json:"title" binding:"required"`
	Path     string `json:"path" `
	Icon     string `json:"icon" `
	ParentId string `json:"parent_id"`
}

type WorkbenchDelReq struct {
	Id string `json:"id" binding:"required"`
}

type WorkbenchUpdateReq struct {
	WorkbenchDelReq
	WorkbenchAddReq
}
