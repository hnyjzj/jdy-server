package model

import (
	"errors"
	"fmt"
	"jdy/enums"
	"jdy/types"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Staff struct {
	SoftDelete

	Phone    string `json:"phone" gorm:"size:11;comment:手机号"`  // 手机号
	Username string `json:"username" gorm:"index;comment:用户名"` // 用户名
	Password string `json:"-" gorm:"size:255;comment:密码"`      // 密码

	Nickname string       `json:"nickname" gorm:"column:nickname;index;comment:昵称"`        // 昵称
	Avatar   string       `json:"avatar" gorm:"size:255;comment:头像"`                       // 头像
	Email    string       `json:"email" gorm:"index;comment:邮箱"`                           // 邮箱
	Gender   enums.Gender `json:"gender" gorm:"column:gender;type:tinyint(1);comment:性别;"` // 性别

	IsDisabled  bool       `json:"is_disabled" gorm:"comment:是否禁用"`     // 是否禁用
	LastLoginAt *time.Time `json:"last_login_at" gorm:"comment:最后登录时间"` // 最后登录时间
	LastLoginIp string     `json:"-" gorm:"size:255;comment:最后登录IP"`    // 最后登录IP

	Identity enums.Identity `json:"identity" gorm:"type:tinyint(1);not null;comment:身份"`    // 身份
	RoleId   string         `json:"role_id" gorm:"type:varchar(255);not NULL;comment:角色ID"` // 角色ID
	Role     *Role          `json:"role" gorm:"foreignKey:RoleId;references:Id;comment:角色"` // 角色

	Stores          []Store  `json:"stores" gorm:"many2many:store_staffs;"`               // 店铺
	StoreSuperiors  []Store  `json:"store_superiors" gorm:"many2many:store_superiors;"`   // 负责的店铺
	Regions         []Region `json:"regions" gorm:"many2many:region_staffs;"`             // 区域
	RegionSuperiors []Region `json:"region_superiors" gorm:"many2many:region_superiors;"` // 负责的区域
}

// 加密密码
func (Staff) HashPassword(password *string) (string, error) {
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
func (u *Staff) VerifyPassword(password string) error {
	if u.Password == "" || password == "" {
		return errors.New("password is nil")
	}
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

// 更新登录信息
func (u *Staff) UpdateLoginInfo(ip string) {
	now := time.Now()
	u.LastLoginIp = ip
	u.LastLoginAt = &now
}

func (Staff) Get(Id, Username *string) (*Staff, error) {
	var staff Staff
	db := DB.Model(staff)

	if Id != nil {
		db = db.Where("id = ?", *Id)
	}
	if Username != nil {
		db = db.Where("username = ?", *Username)
	}

	db = db.Preload("Stores")
	db = db.Preload("Role", func(tx *gorm.DB) *gorm.DB {
		return tx.Preload("Apis").Preload("Routers")
	})

	if err := db.First(&staff).Error; err != nil {
		return nil, err
	}

	return &staff, nil
}

func (S *Staff) HasPermissionApi(path string) bool {
	has := false
	for _, api := range S.Role.Apis {
		if api.Path == path {
			has = true
			break
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

func (Staff) Preloads(db *gorm.DB) *gorm.DB {
	db = db.Preload("Role", func(tx *gorm.DB) *gorm.DB {
		return tx.Preload("Apis").Preload("Routers")
	})
	db = db.Preload("Stores")
	db = db.Preload("StoreSuperiors")
	db = db.Preload("Regions", func(tx *gorm.DB) *gorm.DB {
		return tx.Preload("Stores")
	})
	db = db.Preload("RegionSuperiors", func(tx *gorm.DB) *gorm.DB {
		return tx.Preload("Stores")
	})

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
