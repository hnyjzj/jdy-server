package model

type Router struct {
	SoftDelete

	Title    string  `json:"title" gorm:"size:255;comment:标题"`
	Icon     string  `json:"icon" gorm:"size:255;comment:图标"`
	Path     string  `json:"path" gorm:"size:255;comment:路径"`
	ParentId *string `json:"parent_id" gorm:"index;size:255;comment:父级ID"`
	Sort     int     `json:"sort" gorm:"type:int(11);default:0;comment:排序"`

	Children []*Router `json:"children,omitempty" gorm:"-"`
}

// 获取树形结构
func (Router) GetTree(Pid *string, inIds []string) ([]*Router, error) {
	var list []*Router
	db := DB
	if Pid != nil {
		db = db.Where(&Router{ParentId: Pid})
	} else {
		db = db.Where("parent_id IS NULL")
	}
	if len(inIds) > 0 {
		db = db.Where("id IN (?)", inIds)
	}
	db = db.Order("sort desc")
	if err := db.Find(&list).Error; err != nil {
		return nil, err
	}
	for _, v := range list {
		children, err := v.GetTree(&v.Id, inIds)
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

func (Router) GetTreeReverse(id string) ([]Router, error) {
	var (
		router Router
		list   []Router
	)
	if err := DB.First(&router, "id = ?", id).Error; err != nil {
		return nil, err
	}

	list = append(list, router)

	if router.ParentId != nil {
		res, err := router.GetTreeReverse(*router.ParentId)
		if err != nil {
			return nil, err
		}
		list = append(list, res...)
	}

	return list, nil
}
func init() {
	// 注册模型
	RegisterModels(
		&Router{},
	)
	// 重置表
	RegisterRefreshModels(
	// &Router{},
	)
}
