package types

type StoreCreateReq struct {
	Sort int `json:"sort" binding:"min=0"` // 排序

	Name     string `json:"name" binding:"required"`     // 门店名称
	Province string `json:"province" binding:"required"` // 省份
	City     string `json:"city" binding:"required"`     // 城市
	District string `json:"district" binding:"required"` // 区域
	Address  string `json:"address" binding:"required"`  // 门店地址
	Contact  string `json:"contact" binding:"required"`  // 联系方式
	Logo     string `json:"logo"`                        // 门店logo
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
	Name    *string    `json:"name" label:"门店名称" find:"true" sort:"1" type:"string" input:"text"`
	Region  RegionInfo `json:"region" label:"区域" find:"true" sort:"2" type:"object" input:"region"`
	Address string     `json:"address" label:"门店地址" find:"true" sort:"5" type:"string" input:"text"`
	Contact string     `json:"contact" label:"联系方式" find:"true" sort:"6" type:"string" input:"text"`
	Logo    string     `json:"logo" label:"门店logo" find:"false" sort:"7" type:"string" input:"upload"`
}

type RegionInfo struct {
	Province *string `json:"province"`
	City     *string `json:"city"`
	District *string `json:"district"`
}

type StoreStaffListReq struct {
	StoreId string `json:"id" binding:"required"` // 门店id
}

type StoreStaffAddReq struct {
	StoreId string   `json:"id" binding:"required"`       // 门店id
	StaffId []string `json:"staff_id" binding:"required"` // 用户id
}
