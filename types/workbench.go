package types

type WorkbenchListRes struct {
	Id string `json:"id"`

	Title string `json:"title" `
	Icon  string `json:"icon" `
	Path  string `json:"path" `

	ParentId string             `json:"parent_id" `
	Children []WorkbenchListRes `json:"children" `
}

type WorkbenchAddReq struct {
	Title    string `json:"title" binding:"required"`
	Icon     string `json:"icon"  binding:"required"`
	Path     string `json:"path"  binding:"required"`
	ParentId string `json:"parent_id"`
}

type WorkbenchDelReq struct {
	Id string `json:"id" binding:"required"`
}

type WorkbenchUpdateReq struct {
	WorkbenchDelReq
	WorkbenchAddReq
}
