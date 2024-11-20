package types

type UserReq struct {
	Username string `json:"username" binding:"required,min=2,max=50,regex=^[a-zA-Z0-9_]+$"` // 用户名
	Phone    string `json:"phone" binding:"required,min=11,max=11,regex=^1\\d{10}$"`        // 手机号
	Password string `json:"password" binding:"required"`                                    // 密码

	NickName string `json:"nickname" binding:"required,min=2,max=50,regex=^[\u4e00-\u9fa5]+$"` // 姓名
	Avatar   string `json:"avatar"`                                                            // 头像
	Email    string `json:"email"`                                                             // 邮箱
}

type UserRes struct {
	Id       string `json:"id"`       // 用户ID
	Username string `json:"username"` // 用户名
	Phone    string `json:"phone"`    // 手机号

	NickName string `json:"nickname"` // 姓名
	Avatar   string `json:"avatar"`   // 头像
	Email    string `json:"email"`    // 邮箱
}
