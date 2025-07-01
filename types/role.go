package types

import "jdy/enums"

type RoleWhere struct {
	Id        string         `json:"id" label:"角色ID" find:"true" update:"true" sort:"4" type:"string" input:"text" required:"false"`                                         // 角色ID
	Name      string         `json:"name" label:"角色名称" find:"true" create:"true" update:"true" sort:"4" type:"string" input:"text" required:"false"`                         // 角色名称
	Desc      string         `json:"desc" label:"角色描述" find:"true" create:"true" update:"true" sort:"4" type:"string" input:"text" required:"false"`                         // 角色描述
	Identity  enums.Identity `json:"identity" label:"角色身份" find:"true" create:"true" update:"true" sort:"4" type:"string" input:"select" required:"false"  preset:"typeMap"` // 角色身份
	IsDefault bool           `json:"is_default" label:"是否是默认" find:"true" create:"true" update:"true" sort:"4" type:"string" input:"text" required:"false"`                  // 是否是默认角色
}

type RoleCreateReq struct {
	Name      string         `json:"name" binding:"required"`     // 角色名称
	Desc      string         `json:"desc"`                        // 角色描述
	Identity  enums.Identity `json:"identity" binding:"required"` // 角色身份
	IsDefault bool           `json:"is_default"`                  // 是否是默认角色
}

type RoleInfoReq struct {
	Id string `json:"id" binding:"required"`
}

type RoleUpdateReq struct {
	Id        string         `json:"id" binding:"required"`
	Name      string         `json:"name"`       // 角色名称
	Desc      string         `json:"desc"`       // 角色描述
	Identity  enums.Identity `json:"identity"`   // 角色身份
	IsDefault bool           `json:"is_default"` // 是否是默认角色

	Apis    []string `json:"apis"`    // 角色API
	Routers []string `json:"routers"` // 角色路由
}

type RoleDeleteReq struct {
	Id string `json:"id" binding:"required"`
}
