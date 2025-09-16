package types

type StoreCreateReq struct {
	Order int `json:"order" binding:"min=0"` // 排序

	Name  string `json:"name" binding:"required"`  // 名称
	Alias string `json:"alias" binding:"required"` // 别名
}

type StoreUpdateReq struct {
	Id string `json:"id" binding:"required"`

	StoreCreateReq
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

type StoreListMyReq struct {
	Where StoreWhere `json:"where"`
}

type StoreWhere struct {
	Name     string `json:"name" label:"名称" find:"true" sort:"1" type:"string" input:"text"`      // 名称
	Alias    string `json:"alias" label:"别名" find:"true" sort:"1" type:"string" input:"text"`     // 别名
	RegionId string `json:"region_id" label:"区域" find:"true" sort:"3" type:"string" input:"text"` // 区域
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
