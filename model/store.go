package model

type Store struct {
	SoftDelete

	ParentId *string `json:"parent_id" gorm:"size:255;comment:父级门店id"`

	Name     string `json:"name" gorm:"size:255;comment:名称"`
	Address  string `json:"address" gorm:"size:500;comment:地址"`
	Contact  string `json:"contact" gorm:"size:255;comment:联系人"`
	Logo     string `json:"logo" gorm:"size:255;comment:logo"`
	Order    int    `json:"order" gorm:"size:10;comment:排序"`
	Province string `json:"province" gorm:"size:255;comment:省份"`
	City     string `json:"city" gorm:"size:255;comment:城市"`
	District string `json:"district" gorm:"size:255;comment:区域"`

	SourceId int `json:"store_id" gorm:"size:255;comment:门店id"`

	Children []*Store `json:"children,omitempty" gorm:"-"`
}

func init() {
	// 注册模型
	RegisterModels(
		&Store{},
	)
	// 重置表
	RegisterRefreshModels(
	// &Store{},
	)
}
