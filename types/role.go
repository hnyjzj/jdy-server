package types

type RoleWhere struct {
	Id      string   `json:"id"`       // 角色ID
	Name    string   `json:"name"`     // 角色名称
	Desc    string   `json:"desc"`     // 角色描述
	IsAdmin bool     `json:"is_admin"` // 是否是管理员
	IsRoot  bool     `json:"is_root"`  // 是否是超级管理员
	Apis    []string `json:"apis"`     // 角色API
	Routers []string `json:"routers"`  // 角色路由
	Stores  []string `json:"stores"`   // 角色店铺
	Staffs  []string `json:"staffs"`   // 员工

	StoreId string `json:"store_id"` // 店铺ID
	StaffId string `json:"staff_id"` // 员工ID
}

type RoleCreateReq struct {
	Name    string `json:"name" binding:"required"` // 角色名称
	Desc    string `json:"desc"`                    // 角色描述
	IsAdmin bool   `json:"is_admin"`                // 是否是管理员
}

type RoleInfoReq struct {
	Id string `json:"id" binding:"required"`
}

type RoleUpdateReq struct {
	Id      string   `json:"id" binding:"required"`
	Name    string   `json:"name"`     // 角色名称
	Desc    string   `json:"desc"`     // 角色描述
	IsAdmin bool     `json:"is_admin"` // 是否管理员
	Apis    []string `json:"apis"`     // 角色API
	Routers []string `json:"routers"`  // 角色路由
	Stores  []string `json:"stores"`   // 角色店铺
	Staffs  []string `json:"staffs"`   // 员工
}
