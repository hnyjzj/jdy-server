package types

type WorkbenchListReq struct {
	Title    string `json:"title" binding:"required"`
	Icon     string `json:"icon"  binding:"required"`
	Path     string `json:"path"  binding:"required"`
	ParentId string `json:"parent_id"`
}

type WorkbenchListRes struct {
	Id string `json:"id"`

	Title string `json:"title" `
	Icon  string `json:"icon" `
	Path  string `json:"path" `

	ParentId string             `json:"parent_id" `
	Children []WorkbenchListRes `json:"children" `
}
