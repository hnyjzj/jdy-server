package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	SoftDelete

	Phone    *string `json:"phone" gorm:"uniqueIndex;size:11;comment:手机号"`
	Username *string `json:"username" gorm:"index;comment:用户名"`
	Password string  `json:"-" gorm:"size:255;comment:密码"`

	NickName string `json:"nickname" gorm:"index;comment:姓名"`
	Avatar   string `json:"avatar" gorm:"size:255;comment:头像"`
	Email    string `json:"email" gorm:"index;comment:邮箱"`

	LastLoginAt *time.Time `json:"last_login_at" gorm:"size:255;comment:最后登录时间"`
	LastLoginIp *string    `json:"-" gorm:"size:255;comment:最后登录IP"`
	UpdatePwdAt *time.Time `json:"-" gorm:"comment:修改密码时间"`

	IsDisabled bool `json:"is_disabled" gorm:"comment:是否禁用"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.Password != "" {
		if err := u.HashPassword(); err != nil {
			return err
		}
	}

	u.BaseModel.BeforeCreate(tx)
	return
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
	RegisterModels(
		&User{},
	)
	// 重置表
	RegisterRefreshModels(
	// &User{},
	)
}
