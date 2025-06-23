package model

import (
	"jdy/types"

	"gorm.io/gorm"
)

type Role struct {
	SoftDelete

	Name string `json:"name" gorm:"column:name;index:unique;size:255;not null;comment:角色名称"` // 角色名称
	Desc string `json:"desc" gorm:"column:desc;size:255;comment:角色描述"`                       // 角色描述

	IsRoot  bool `json:"is_root" gorm:"column:is_root;comment:是否超级管理员"` // 是否超级管理员
	IsAdmin bool `json:"is_admin" gorm:"column:is_admin;comment:是否管理员"` // 是否管理员

	OperatorId string `json:"operator_id" gorm:"size:255;not null;comment:操作员ID;"`                        // 操作员ID
	Operator   Staff  `json:"operator,omitempty" gorm:"foreignKey:OperatorId;references:Id;comment:操作员;"` // 操作员
	IP         string `json:"ip" gorm:"size:255;not null;comment:IP地址;"`                                  // IP地址

	Apis    []Api    `json:"apis" gorm:"many2many:role_apis;comment:角色API;"`      // 角色API
	Routers []Router `json:"routers" gorm:"many2many:role_routers;comment:角色路由;"` // 角色路由
	Stores  []Store  `json:"stores" gorm:"many2many:role_stores;comment:角色店铺;"`   // 角色店铺
	Staffs  []Staff  `json:"staffs" gorm:"many2many:staff_roles;comment:员工;"`     // 员工
}

func (Role) WhereCondition(db *gorm.DB, query *types.RoleWhere) *gorm.DB {
	if query.Id != "" {
		db = db.Where("id = ?", query.Id)
	}

	return db
}

func (Role) Preloads(db *gorm.DB, query *types.RoleWhere) *gorm.DB {
	db = db.Preload("Apis")
	db = db.Preload("Routers")
	db = db.Preload("Stores")
	db = db.Preload("Staffs")
	db = db.Preload("Operator")

	return db
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
