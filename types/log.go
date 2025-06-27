package types

type OnCaptureScreenReq struct {
	Username  string `json:"username" binding:"required"`
	Storename string `json:"storename"`
	Url       string `json:"url" binding:"required"`
}
