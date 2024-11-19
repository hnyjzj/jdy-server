package types

type UserRes struct {
	Id       string `json:"id"`       // 用户ID
	UserName string `json:"username"` // 用户名
	Phone    string `json:"phone"`    // 手机号

	Name   string `json:"name"`   // 姓名
	Avatar string `json:"avatar"` // 头像
	Email  string `json:"email"`  // 邮箱
}
