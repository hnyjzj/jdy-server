package model

import (
	"errors"
	"jdy/enums"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	SoftDelete
	Platform enums.PlatformType `json:"platform" gorm:"index;comment:平台"`

	Phone    *string `json:"phone" gorm:"index;size:11;comment:手机号"`
	Username *string `json:"username" gorm:"index;comment:用户名"`
	Password *string `json:"-" gorm:"size:255;comment:密码"`

	Nickname *string      `json:"nickname" gorm:"size:255;comment:昵称"`
	Avatar   *string      `json:"avatar" gorm:"size:255;comment:头像"`
	Email    *string      `json:"email" gorm:"size:255;comment:邮箱"`
	Gender   enums.Gender `json:"gender" gorm:"column:gender;type:tinyint(1);comment:性别;"` // 性别

	LastLoginAt *time.Time `json:"last_login_at" gorm:"comment:最后登录时间"`
	LastLoginIp string     `json:"-" gorm:"size:255;comment:最后登录IP"`

	IsDisabled bool `json:"is_disabled" gorm:"comment:是否禁用"`

	StaffId *string `json:"staff_id" gorm:"size:255;comment:员工ID"`
	Staff   *Staff  `json:"staff" gorm:"foreignKey:StaffId;references:Id;"`
}

// 加密密码
func (Account) HashPassword(password *string) (string, error) {
	if password == nil {
		return "", errors.New("password is nil")
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(*password), 10)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// 校验密码
func (u *Account) VerifyPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(*u.Password), []byte(password))
}

// 更新登录信息
func (u *Account) UpdateLoginInfo(ip string) {
	now := time.Now()
	u.LastLoginIp = ip
	u.LastLoginAt = &now
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
