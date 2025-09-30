package types

type StoreCreateReq struct {
	RegionId string `json:"region_id" binding:"required"` // 区域ID

	Name  string `json:"name" binding:"required"`  // 名称
	Alias string `json:"alias" binding:"required"` // 别名
	Phone string `json:"phone"`                    // 电话
	Order int    `json:"order" binding:"min=0"`    // 排序
}

type StoreUpdateReq struct {
	Id string `json:"id" binding:"required"`

	RegionId string `json:"region_id"`             // 区域ID
	Name     string `json:"name"`                  // 名称
	Alias    string `json:"alias"`                 // 别名
	Phone    string `json:"phone"`                 // 电话
	Order    int    `json:"order" binding:"min=0"` // 排序
}

type StoreDeleteReq struct {
	Id string `json:"id" binding:"required"`
}

type StoreInfoReq struct {
	Id string `json:"id" binding:"required"`
}

type StoreListReq struct {
	PageReq
	Where StoreWhere `json:"where"`
}

type StoreAliasReq struct {
	IsHeadquarters bool `json:"is_headquarters"` // 是否是总部
}

type StoreListMyReq struct {
	HasAll bool       `json:"has_all"` // 是否包含全部
	Where  StoreWhere `json:"where"`
}

type StoreWhere struct {
	Name     string `json:"name" label:"名称" find:"true" sort:"1" type:"string" input:"text"`        // 名称
	Alias    string `json:"alias" label:"别名" find:"true" sort:"2" type:"string" input:"text"`       // 别名
	Phone    string `json:"phone" label:"电话" find:"true" sort:"3" type:"string" input:"text"`       // 电话
	RegionId string `json:"region_id" label:"区域" find:"true" sort:"4" type:"string" input:"search"` // 区域
}

type StoreStaffListReq struct {
	StoreId string `json:"id" binding:"required"` // 门店ID
}

type StoreStaffAddReq struct {
	StoreId string   `json:"id" binding:"required"`       // 门店ID
	StaffId []string `json:"staff_id" binding:"required"` // 员工ID
}

type StoreStaffDelReq struct {
	StoreId string   `json:"id" binding:"required"`       // 门店ID
	StaffId []string `json:"staff_id" binding:"required"` // 员工ID
}

type StoreStaffIsInReq struct {
	StoreId string `json:"id" binding:"required"` // 门店ID
	StaffId string `json:"staff_id"`              // 员工ID
}

type StoreSuperiorListReq struct {
	StoreId string `json:"id" binding:"required"` // 门店ID
}

type StoreSuperiorAddReq struct {
	StoreId    string   `json:"id" binding:"required"`          // 门店ID
	SuperiorId []string `json:"superior_id" binding:"required"` // 负责人ID
}

type StoreSuperiorDelReq struct {
	StoreId    string   `json:"id" binding:"required"`          // 门店ID
	SuperiorId []string `json:"superior_id" binding:"required"` // 负责人ID
}

type StoreSuperiorIsInReq struct {
	StoreId string `json:"id" binding:"required"` // 门店ID
	StaffId string `json:"staff_id"`              // 员工ID
}

type StoreAdminListReq struct {
	StoreId string `json:"id" binding:"required"` // 门店ID
}

type StoreAdminAddReq struct {
	StoreId string   `json:"id" binding:"required"`       // 门店ID
	AdminId []string `json:"admin_id" binding:"required"` // 负责人ID
}

type StoreAdminDelReq struct {
	StoreId string   `json:"id" binding:"required"`       // 门店ID
	AdminId []string `json:"admin_id" binding:"required"` // 负责人ID
}

type StoreAdminIsInReq struct {
	StoreId string `json:"id" binding:"required"` // 门店ID
	StaffId string `json:"staff_id"`              // 员工ID
}
