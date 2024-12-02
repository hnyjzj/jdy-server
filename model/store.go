package model

import (
	"jdy/types"

	"gorm.io/gorm"
)

type Store struct {
	SoftDelete

	ParentId *string `json:"parent_id" gorm:"size:255;comment:父级门店id"`

	Name     string `json:"name" gorm:"size:255;comment:名称"`
	Address  string `json:"address" gorm:"size:500;comment:地址"`
	Contact  string `json:"contact" gorm:"size:255;comment:联系人"`
	Logo     string `json:"logo" gorm:"size:255;comment:logo"`
	Sort     int    `json:"sort" gorm:"size:10;comment:排序"`
	Province string `json:"province" gorm:"size:255;comment:省份"`
	City     string `json:"city" gorm:"size:255;comment:城市"`
	District string `json:"district" gorm:"size:255;comment:区域"`

	WxworkId int `json:"wxwork_id" gorm:"size:10;comment:企业微信id"`

	Children []*Store `json:"children,omitempty" gorm:"-"`

	Staffs []Staff `json:"staffs" gorm:"many2many:stores_staffs;"`
}

func (Store) WhereCondition(db *gorm.DB, query *types.StoreWhereReq) *gorm.DB {
	if query.ParentId != nil {
		db = db.Where("parent_id = ?", query.ParentId)
	}
	if query.Name != nil {
		db = db.Where("name LIKE ?", "%"+*query.Name+"%")
	}
	if query.Address != "" {
		db = db.Where("address LIKE ?", "%"+query.Address+"%")
	}
	if query.Contact != "" {
		db = db.Where("contact LIKE ?", "%"+query.Contact+"%")
	}
	if query.Province != nil {
		db = db.Where("province LIKE ?", "%"+*query.Province+"%")
	}
	if query.City != nil {
		db = db.Where("city LIKE ?", "%"+*query.City+"%")
	}
	if query.District != nil {
		db = db.Where("district LIKE ?", "%"+*query.District+"%")
	}
	if query.WxworkId != 0 {
		db = db.Where("wxwork_id = ?", query.WxworkId)
	}

	return db
}

// 获取树形结构
func (Store) GetTree(Pid *string) ([]*Store, error) {
	var list []*Store
	db := DB
	if Pid != nil {
		db = db.Where(&Store{ParentId: Pid})
	} else {
		db = db.Where("parent_id IS NULL")
	}
	db = db.Order("sort DESC")
	if err := db.Find(&list).Error; err != nil {
		return nil, err
	}
	for _, v := range list {
		children, err := v.GetTree(&v.Id)
		if err != nil {
			return nil, err
		}
		if len(children) == 0 {
			continue
		}
		v.Children = children
	}

	return list, nil
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
