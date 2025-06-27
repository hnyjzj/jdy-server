package model

type Region struct {
	SoftDelete

	IdWx  string `json:"id_wx" gorm:"size:255;comment:微信ID"` // 微信ID
	Name  string `json:"name" gorm:"size:255;comment:名称"`    // 名称
	Order int    `json:"order" gorm:"comment:排序"`            // 排序

	Stores    []Store `json:"stores" gorm:"many2many:region_stores;"`       // 门店
	Staffs    []Staff `json:"staffs" gorm:"many2many:region_staffs;"`       // 员工
	Superiors []Staff `json:"superiors" gorm:"many2many:region_superiors;"` // 负责人
}

func init() {
	// 注册模型
	RegisterModels(
		&Region{},
	)
	// 重置表
	RegisterRefreshModels(
	// &Region{},
	)
}
