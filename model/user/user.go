package usermodel

import (
	"jdy/model"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	model.SoftDelete

	Phone    *string `json:"phone" gorm:"uniqueIndex;size:11;comment:手机号"`
	Username *string `json:"username" gorm:"index;comment:用户名"`
	Password string  `json:"-" gorm:"size:255;comment:密码"`

	Name   string `json:"name" gorm:"index;comment:姓名"`
	Avatar string `json:"avatar" gorm:"size:255;comment:头像"`
	Email  string `json:"email" gorm:"index;comment:邮箱"`

	LastLoginAt *time.Time `json:"last_login_at" gorm:"size:255;comment:最后登录时间"`
	LastLoginIp *string    `json:"-" gorm:"size:255;comment:最后登录IP"`
	UpdatePwdAt *time.Time `json:"-" gorm:"comment:修改密码时间"`

	IsDisabled bool `json:"is_disabled" gorm:"comment:是否禁用"`
}

// 加密密码
func (u *User) HashPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	if err != nil {
		return err
	}
	u.Password = string(bytes)
	return nil
}

// 校验密码
func (u *User) VerifyPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

func init() {
	// 注册模型
	model.RegisterModels(
		&User{},
	)
	// 重置表
	model.RegisterRefreshModels(
	// &User{},
	)
}
