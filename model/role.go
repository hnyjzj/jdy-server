package model

import (
	"jdy/enums"
	"jdy/types"

	"gorm.io/gorm"
)

type Role struct {
	SoftDelete

	Name     string         `json:"name" gorm:"column:name;index:unique;size:255;not null;comment:角色名称"` // 角色名称
	Desc     string         `json:"desc" gorm:"column:desc;size:255;comment:角色描述"`                       // 角色描述
	Identity enums.Identity `json:"identity" gorm:"column:identity;not null;comment:角色身份"`               // 角色身份

	IsDefault bool `json:"is_default" gorm:"column:is_default;default:0;comment:是否默认"` // 是否默认

	OperatorId string `json:"operator_id" gorm:"size:255;comment:操作员ID;"`                                 // 操作员ID
	Operator   Staff  `json:"operator,omitempty" gorm:"foreignKey:OperatorId;references:Id;comment:操作员;"` // 操作员
	IP         string `json:"ip" gorm:"size:255;comment:IP地址;"`                                           // IP地址

	Apis    []Api    `json:"apis" gorm:"many2many:role_apis;comment:角色API;"`      // 角色API
	Routers []Router `json:"routers" gorm:"many2many:role_routers;comment:角色路由;"` // 角色路由
	Staffs  []Staff  `json:"staffs" gorm:"many2many:role_staffs;comment:角色员工;"`   // 角色员工
}

func (Role) Default(Identity enums.Identity) (*Role, error) {
	var role Role
	if err := DB.Model(&Role{}).Where(&Role{
		Identity:  Identity,
		IsDefault: true,
	}).First(&role).Error; err != nil {
		return nil, err
	}

	return &role, nil
}

func (Role) SetDefault(db *gorm.DB, staff_id string, Identity enums.Identity) error {
	var role Role
	if err := db.Model(&Role{}).Where(&Role{
		Identity:  Identity,
		IsDefault: true,
	}).First(&role).Error; err != nil {
		return err
	}

	if err := db.Model(&Staff{}).Where("id = ?", staff_id).Updates(Staff{
		RoleId: role.Id,
	}).Error; err != nil {
		return err
	}

	return nil
}

func (Role) WhereCondition(db *gorm.DB, query *types.RoleWhere) *gorm.DB {
	if query.Id != "" {
		db = db.Where("id = ?", query.Id)
	}
	if query.Identity != 0 {
		db = db.Where("identity = ?", query.Identity)
	}

	return db
}

func (Role) Preloads(db *gorm.DB, query *types.RoleWhere) *gorm.DB {
	db = db.Preload("Apis")
	db = db.Preload("Routers")

	return db
}

func (Role) Init() error {
	for identity, name := range enums.IdentityMap {
		var role Role
		if err := DB.Model(&Role{}).Where(&Role{
			Identity:  identity,
			IsDefault: true,
		}).Attrs(Role{
			Name: name,
			Desc: "默认角色",
		}).FirstOrCreate(&role).Error; err != nil {
			return err
		}
	}

	return nil
}

func init() {
	// 注册模型
	RegisterModels(
		&Role{},
	)
	// 重置表
	RegisterRefreshModels(
	// &Role{},
	)
}
