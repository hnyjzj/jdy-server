package types

type PlatformOAuthReq struct {
	Agent string `json:"agent"` // 用户浏览器环境

	Uri string `json:"uri" binding:"required"` // 授权后重定向的回调链接地址，请使用urlencode对链接进行处理
}

type PlatformOAuthRes struct {
	RedirectURL string `json:"redirect_url"` // 重定向链接
}

type PlatformJSSdkReq struct {
	Agent string `json:"agent"` // 用户浏览器环境
}

type PlatformGetUserReq struct {
	UserId string `json:"user_id" binding:"required"` // 用户ID
}

type PlatformDepartmentCreateReq struct {
	Name     string `json:"name" binding:"required"`      // 部门名称
	NameEn   string `json:"name_en" binding:"required"`   // 英文名称
	ParentId int    `json:"parent_id" binding:"required"` // 父部门ID
	Order    int    `json:"order" `                       // 部门排序
}
