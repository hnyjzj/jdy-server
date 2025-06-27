package model

import (
	"fmt"
	"jdy/enums"
	"jdy/types"
	"strings"

	"gorm.io/gorm"
)

type Staff struct {
	SoftDelete

	Phone *string `json:"phone" gorm:"uniqueIndex;size:11;comment:手机号"` // 手机号

	Nickname string       `json:"nickname" gorm:"column:nickname;index;comment:姓名"`        // 姓名
	Avatar   string       `json:"avatar" gorm:"size:255;comment:头像"`                       // 头像
	Email    string       `json:"email" gorm:"index;comment:邮箱"`                           // 邮箱
	Gender   enums.Gender `json:"gender" gorm:"column:gender;type:tinyint(1);comment:性别;"` // 性别

	IsDisabled bool `json:"is_disabled" gorm:"comment:是否禁用"` // 是否禁用

	Account  *Account  `json:"account,omitempty" gorm:"foreignKey:StaffId;references:Id;"`  // 账号
	Accounts []Account `json:"accounts,omitempty" gorm:"foreignKey:StaffId;references:Id;"` // 账号

	Stores  []Store  `json:"stores" gorm:"many2many:store_staffs;"`   // 店铺
	Regions []Region `json:"regions" gorm:"many2many:region_staffs;"` // 区域
	Roles   []Role   `json:"roles" gorm:"many2many:role_staffs;"`     // 角色
}

func (Staff) Get(Id string) (*Staff, error) {
	var staff Staff
	db := DB.Model(staff).Where("id = ?", Id)

	db = db.Preload("Stores")
	db = db.Preload("Roles", func(tx *gorm.DB) *gorm.DB {
		return tx.Preload("Apis").Preload("Routers")
	})

	if err := db.First(&staff).Error; err != nil {
		return nil, err
	}

	return &staff, nil
}

func (S *Staff) HasPermissionApi(path string) bool {
	has := false
	for _, role := range S.Roles {
		for _, api := range role.Apis {
			if api.Path == path {
				has = true
				break
			}
		}
	}

	return has
}

func (Staff) WhereCondition(db *gorm.DB, query *types.StaffWhere) *gorm.DB {
	if query.Phone != "" {
		db = db.Where("phone = ?", query.Phone)
	}
	if query.Nickname != "" {
		db = db.Where("nickname LIKE ?", fmt.Sprintf("%%%s%%", strings.Replace(query.Nickname, "%", "\\%", -1)))
	}
	if query.Gender != 0 {
		db = db.Where("gender = ?", query.Gender)
	}
	if query.IsDisabled {
		db = db.Where("is_disabled = ?", query.IsDisabled)
	}

	return db
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
