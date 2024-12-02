package types

type StoreCreateReq struct {
	ParentId *string `json:"parent_id"`            // 父级门店id
	Sort     int     `json:"sort" binding:"min=0"` // 排序

	Name     string `json:"name" binding:"required"`     // 门店名称
	Province string `json:"province" binding:"required"` // 省份
	City     string `json:"city" binding:"required"`     // 城市
	District string `json:"district" binding:"required"` // 区域
	Address  string `json:"address" binding:"required"`  // 门店地址
	Contact  string `json:"contact" binding:"required"`  // 联系方式
	Logo     string `json:"logo"`                        // 门店logo

	SyncWxwork bool `json:"sync_wxwork"` // 是否同步到企业微信
	WxworkId   int  `json:"wxwork_id"`
}

type StoreUpdateReq struct {
	Id string `json:"id" binding:"required"`

	StoreCreateReq
}

type StoreDeleteReq struct {
	Id string `json:"id" binding:"required"`

	SyncWxwork bool `json:"sync_wxwork"` // 是否同步到企业微信
}

type StoreInfoReq struct {
	Id string `json:"id" binding:"required"`
}

type StoreWhereReq struct {
	Name     *string `json:"name"`     // 门店名称
	Province *string `json:"province"` // 省份
	City     *string `json:"city"`     // 城市
	District *string `json:"district"` // 区域
	Address  string  `json:"address"`  // 门店地址
	Contact  string  `json:"contact"`  // 联系方式

	WxworkId int     `json:"wxwork_id"`
	ParentId *string `json:"parent_id"` // 父级门店id
}

type StoreListReq struct {
	PageReq
	Where StoreWhereReq `json:"where"`
}
