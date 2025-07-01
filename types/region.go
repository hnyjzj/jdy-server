package types

type RegionCreateReq struct {
	Sort int `json:"sort" binding:"min=0"` // 排序

	Name string `json:"name" binding:"required"` // 门店名称
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
	Name *string `json:"name" label:"门店名称" find:"true" sort:"1" type:"string" input:"text"`
}

type RegionStaffListReq struct {
	RegionId string `json:"id" binding:"required"` // 门店id
}

type RegionStaffAddReq struct {
	RegionId string   `json:"id" binding:"required"`       // 门店id
	StaffId  []string `json:"staff_id" binding:"required"` // 用户id
}

type RegionStaffDelReq struct {
	RegionId string   `json:"id" binding:"required"`       // 门店id
	StaffId  []string `json:"staff_id" binding:"required"` // 用户id
}
