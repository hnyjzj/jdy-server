package model

import "jdy/enums"

type Staff struct {
	SoftDelete

	Phone *string `json:"phone" gorm:"uniqueIndex;size:11;comment:手机号"`

	Nickname string       `json:"nickname" gorm:"column:nickname;index;comment:姓名"`
	Avatar   string       `json:"avatar" gorm:"size:255;comment:头像"`
	Email    string       `json:"email" gorm:"index;comment:邮箱"`
	Gender   enums.Gender `json:"gender" gorm:"column:gender;type:tinyint(1);comment:性别;"` // 性别

	IsDisabled bool `json:"is_disabled" gorm:"comment:是否禁用"`

	Account  *Account  `json:"account" gorm:"foreignKey:StaffId;references:Id;"`
	Accounts []Account `json:"accounts" gorm:"foreignKey:StaffId;references:Id;"`

	Stores []Store `json:"stores" gorm:"many2many:stores_staffs;"`
}

func init() {
	// 注册模型
	RegisterModels(
		&Staff{},
	)
	// 重置表
	RegisterRefreshModels(
	// &Staff{},
	)
}
