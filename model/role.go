package model

type Role struct {
	SoftDelete

	Name string `json:"name" gorm:"column:name;size:255;not null;comment:角色名称"` // 角色名称
	Desc string `json:"desc" gorm:"column:desc;size:255;comment:角色描述"`          // 角色描述

	OperatorId string `json:"operator_id" gorm:"size:255;not null;comment:操作员ID;"`              // 操作员ID
	Operator   Staff  `json:"operator" gorm:"foreignKey:OperatorId;references:Id;comment:操作员;"` // 操作员
	IP         string `json:"ip" gorm:"size:255;not null;comment:IP地址;"`                        // IP地址

	Apis    []Api    `json:"apis" gorm:"many2many:role_apis;comment:角色API;"`      // 角色API
	Routers []Router `json:"routers" gorm:"many2many:role_routers;comment:角色路由;"` // 角色路由
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
