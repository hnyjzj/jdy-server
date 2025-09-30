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

	Nickname string       `json:"nickname" gorm:"column:nickname;index;comment:昵称"`     // 昵称
	Avatar   string       `json:"avatar" gorm:"size:255;comment:头像"`                    // 头像
	Email    string       `json:"email" gorm:"index;comment:邮箱"`                        // 邮箱
	Gender   enums.Gender `json:"gender" gorm:"column:gender;type:int(11);comment:性别;"` // 性别

	IsDisabled  bool       `json:"is_disabled" gorm:"comment:是否禁用"`     // 是否禁用
	LastLoginAt *time.Time `json:"last_login_at" gorm:"comment:最后登录时间"` // 最后登录时间
	LastLoginIp string     `json:"-" gorm:"size:255;comment:最后登录IP"`    // 最后登录IP

	Identity enums.Identity `json:"identity" gorm:"index;type:int(11);not null;comment:身份"` // 身份

	LeaderId string `json:"leader_id" gorm:"index;type:varchar(255);default:null;comment:上级ID"` // 上级ID
	Leader   *Staff `json:"leader" gorm:"foreignKey:LeaderId;references:Id;comment:上级"`         // 上级

	TagId string     `json:"tag_id" gorm:"index;type:varchar(255);default:null;comment:标签ID"` // 标签ID
	Tag   StaffTag   `json:"tag" gorm:"foreignKey:TagId;references:Id;comment:标签"`            // 标签
	Logs  []StaffLog `json:"logs" gorm:"foreignKey:StaffId;references:Id;comment:日志"`         // 日志

	RoleId string `json:"role_id" gorm:"type:varchar(255);default:null;comment:角色ID"` // 角色ID
	Role   *Role  `json:"role" gorm:"foreignKey:RoleId;references:Id;comment:角色"`     // 角色

	StoreIds        []string `json:"store_ids" gorm:"-"`                                  // 店铺ID
	Stores          []Store  `json:"stores" gorm:"many2many:store_staffs;"`               // 店铺
	StoreSuperiors  []Store  `json:"store_superiors" gorm:"many2many:store_superiors;"`   // 负责的店铺
	StoreAdmins     []Store  `json:"store_admins" gorm:"many2many:store_admins;"`         // 管理的店铺
	RegionIds       []string `json:"region_ids" gorm:"-"`                                 // 区域ID
	Regions         []Region `json:"regions" gorm:"many2many:region_staffs;"`             // 区域
	RegionSuperiors []Region `json:"region_superiors" gorm:"many2many:region_superiors;"` // 负责的区域
	RegionAdmins    []Region `json:"region_admins" gorm:"many2many:region_admins;"`       // 管理的区域
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
	db := DB.Model(&staff)

	if Id != nil {
		db = db.Where("id = ?", *Id)
	}
	if Username != nil {
		db = db.Where("username = ?", *Username)
	}

	db = staff.Preloads(db)

	if err := db.First(&staff).Error; err != nil {
		return nil, err
	}
	if staff.Identity >= enums.IdentityAdmin {
		if err := DB.Model(Store{}).Find(&staff.StoreAdmins).Error; err != nil {
			return nil, err
		}
		if err := DB.Model(Region{}).Find(&staff.RegionAdmins).Error; err != nil {
			return nil, err
		}
	}

	for _, store := range staff.Stores {
		staff.StoreIds = append(staff.StoreIds, store.Id)
	}
	for _, store := range staff.StoreSuperiors {
		staff.StoreIds = append(staff.StoreIds, store.Id)
	}
	for _, store := range staff.StoreAdmins {
		staff.StoreIds = append(staff.StoreIds, store.Id)
	}
	for _, region := range staff.Regions {
		staff.RegionIds = append(staff.RegionIds, region.Id)
		for _, store := range region.Stores {
			staff.StoreIds = append(staff.StoreIds, store.Id)
		}
	}
	for _, region := range staff.RegionSuperiors {
		staff.RegionIds = append(staff.RegionIds, region.Id)
		for _, store := range region.Stores {
			staff.StoreIds = append(staff.StoreIds, store.Id)
		}
	}
	for _, region := range staff.RegionAdmins {
		staff.RegionIds = append(staff.RegionIds, region.Id)
		for _, store := range region.Stores {
			staff.StoreIds = append(staff.StoreIds, store.Id)
		}
	}

	if staff.Role == nil {
		role, err := Role{}.Default(staff.Identity)
		if err != nil {
			return nil, err
		}
		staff.Role = role
	}

	return &staff, nil
}

func (S *Staff) HasPermissionApi(path string) error {
	if S.Role == nil {
		return nil
	}

	for _, api := range S.Role.Apis {
		if api.Path == path {
			return nil
		}
	}

	var api Api
	if err := DB.Model(&api).Where(&Api{
		Path: path,
	}).First(&api).Error; err != nil {
		return fmt.Errorf("无权限访问: %v", path)
	}

	return fmt.Errorf("暂无权限: %v", api.Title)
}

func (Staff) WhereCondition(db *gorm.DB, query *types.StaffWhere) *gorm.DB {
	if query.Nickname != "" {
		db = db.Where("nickname LIKE ?", fmt.Sprintf("%%%s%%", strings.ReplaceAll(query.Nickname, "%", "\\%")))
	}
	if query.Phone != "" {
		db = db.Where("phone = ?", query.Phone)
	}
	if query.Username != "" {
		db = db.Where("username = ?", query.Username)
	}
	if query.Email != "" {
		db = db.Where("email LIKE ?", fmt.Sprintf("%%%s%%", strings.ReplaceAll(query.Email, "%", "\\%")))
	}
	if query.Gender != 0 {
		db = db.Where("gender = ?", query.Gender)
	}
	if query.IsDisabled {
		db = db.Where("is_disabled = ?", query.IsDisabled)
	}
	if query.Identity != 0 {
		db = db.Where("identity = ?", query.Identity)
	}
	if query.StoreId != "" {
		db = db.Where("id IN (SELECT staff_id FROM store_staffs WHERE store_id = ?)", query.StoreId)
	}

	return db
}

func (Staff) Preloads(db *gorm.DB) *gorm.DB {
	db = db.Preload("Role", func(tx *gorm.DB) *gorm.DB {
		return tx.Preload("Apis").Preload("Routers")
	})
	db = db.Preload("Leader")
	db = db.Preload("Stores")
	db = db.Preload("StoreSuperiors")
	db = db.Preload("StoreAdmins")
	db = db.Preload("Regions", func(tx *gorm.DB) *gorm.DB {
		return tx.Preload("Stores")
	})
	db = db.Preload("RegionSuperiors", func(tx *gorm.DB) *gorm.DB {
		return tx.Preload("Stores")
	})
	db = db.Preload("RegionAdmins", func(tx *gorm.DB) *gorm.DB {
		return tx.Preload("Stores")
	})

	return db
}

// 员工标签
type StaffTag struct {
	Model

	TagId int64  `json:"tag_id" gorm:"index;type:int(11);default:null;comment:标签ID;"` // 标签ID
	Name  string `json:"name" gorm:"type:varchar(255);default:null;comment:标签名称;"`    // 标签名称
}

// 员工日志
type StaffLog struct {
	Model

	Type enums.StaffLogType `json:"type" gorm:"type:tinyint(1);comment:类型;"` // 类型

	StaffId string `json:"staff_id" gorm:"index;type:varchar(255);default:null;comment:员工ID;"` // 员工ID
	Staff   Staff  `json:"staff" gorm:"foreignKey:StaffId;references:Id;comment:员工;"`          // 员工

	OldValue Staff `json:"old_value" gorm:"type:text;serializer:json;comment:旧值;"` // 旧值
	NewValue Staff `json:"new_value" gorm:"type:text;serializer:json;comment:新值;"` // 新值
}

func init() {
	// 注册模型
	RegisterModels(
		&Staff{},
		&StaffTag{},
		&StaffLog{},
	)
	// 重置表
	RegisterRefreshModels(
	// &Staff{},
	// &StaffTag{},
	// &StaffLog{},
	)
}
