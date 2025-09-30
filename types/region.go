package types

type RegionCreateReq struct {
	Name  string `json:"name" binding:"required"`  // 名称
	Alias string `json:"alias" binding:"required"` // 别名
}

type RegionUpdateReq struct {
	Id string `json:"id" binding:"required"`

	Name  string `json:"name"`  // 名称
	Alias string `json:"alias"` // 别名
}

type RegionDeleteReq struct {
	Id string `json:"id" binding:"required"`
}

type RegionInfoReq struct {
	Id string `json:"id" binding:"required"`
}

type RegionListReq struct {
	PageReq
	Where RegionWhere `json:"where"`
}

type RegionListMyReq struct {
	HasAll bool        `json:"has_all"` // 是否包含全部
	Where  RegionWhere `json:"where"`
}

type RegionWhere struct {
	Name  string `json:"name" label:"区域名称" find:"true" sort:"1" type:"string" input:"text"`  // 名称
	Alias string `json:"alias" label:"区域别名" find:"true" sort:"1" type:"string" input:"text"` // 别名
}

type RegionStoreListReq struct {
	RegionId string `json:"id" binding:"required"` // 区域id
}

type RegionStoreAddReq struct {
	RegionId string   `json:"id" binding:"required"`       // 区域id
	StoreId  []string `json:"store_id" binding:"required"` // 店铺id
}

type RegionStoreDelReq struct {
	RegionId string   `json:"id" binding:"required"`       // 区域id
	StoreId  []string `json:"store_id" binding:"required"` // 店铺id
}

type RegionStaffListReq struct {
	RegionId string `json:"id" binding:"required"` // 区域id
}

type RegionStaffAddReq struct {
	RegionId string   `json:"id" binding:"required"`       // 区域id
	StaffId  []string `json:"staff_id" binding:"required"` // 用户id
}

type RegionStaffDelReq struct {
	RegionId string   `json:"id" binding:"required"`       // 区域id
	StaffId  []string `json:"staff_id" binding:"required"` // 用户id
}

type RegionSuperiorListReq struct {
	RegionId string `json:"id" binding:"required"` // 区域id
}

type RegionSuperiorAddReq struct {
	RegionId   string   `json:"id" binding:"required"`          // 区域id
	SuperiorId []string `json:"superior_id" binding:"required"` // 上级id
}

type RegionSuperiorDelReq struct {
	RegionId   string   `json:"id" binding:"required"`          // 区域id
	SuperiorId []string `json:"superior_id" binding:"required"` // 上级id
}

type RegionAdminListReq struct {
	RegionId string `json:"id" binding:"required"` // 区域id
}

type RegionAdminAddReq struct {
	RegionId string   `json:"id" binding:"required"`       // 区域id
	AdminId  []string `json:"admin_id" binding:"required"` // 管理员 ID
}

type RegionAdminDelReq struct {
	RegionId string   `json:"id" binding:"required"`       // 区域id
	AdminId  []string `json:"admin_id" binding:"required"` // 管理员 ID
}
