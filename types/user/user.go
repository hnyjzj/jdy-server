package usertype

type UserRes struct {
	Id       string `json:"id"`
	UserName string `json:"username"`
	Phone    string `json:"phone"`

	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Email  string `json:"email"`
}
