package model

import (
	"jdy/types"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Account struct {
	SoftDelete
	Platform types.PlatformType `json:"platform" gorm:"index;comment:平台"`

	Phone    *string `json:"phone" gorm:"uniqueIndex;size:11;comment:手机号"`
	Username *string `json:"username" gorm:"index;comment:用户名"`
	Password *string `json:"-" gorm:"size:255;comment:密码"`

	Nickname *string `json:"nickname" gorm:"size:255;comment:昵称"`
	Avatar   *string `json:"avatar" gorm:"size:255;comment:头像"`
	Email    *string `json:"email" gorm:"size:255;comment:邮箱"`
	Gender   uint    `json:"gender" gorm:"size:255;comment:性别"`

	LastLoginAt *time.Time `json:"last_login_at" gorm:"comment:最后登录时间"`
	LastLoginIp string     `json:"-" gorm:"size:255;comment:最后登录IP"`

	IsDisabled bool `json:"is_disabled" gorm:"comment:是否禁用"`

	StaffId *string `json:"staff_id" gorm:"size:255;comment:员工ID"`
	Staff   Staff   `json:"-" gorm:"foreignKey:StaffId;references:Id"`
}

func (u *Account) BeforeCreate(tx *gorm.DB) (err error) {
	if u.Password != nil {
		if err := u.HashPassword(); err != nil {
			return err
		}
	}

	u.BaseModel.BeforeCreate(tx)
	return
}

// 加密密码
func (u *Account) HashPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(*u.Password), 10)
	if err != nil {
		return err
	}
	password := string(bytes)
	u.Password = &password
	return nil
}

// 校验密码
func (u *Account) VerifyPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(*u.Password), []byte(password))
}
func init() {
	// 注册模型
	RegisterModels(
		&Account{},
	)
	// 重置表
	RegisterRefreshModels(
	// &Account{},
	)
}