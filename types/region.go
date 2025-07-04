package types

type RegionCreateReq struct {
	Name string `json:"name" binding:"required"` // 区域名称
}

type RegionUpdateReq struct {
	Id string `json:"id" binding:"required"`

	RegionCreateReq
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
	Where RegionWhere `json:"where"`
}

type RegionWhere struct {
	Name *string `json:"name" label:"区域名称" find:"true" sort:"1" type:"string" input:"text"`
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
