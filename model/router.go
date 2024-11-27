package model

type Router struct {
	SoftDelete

	Title    string  `json:"title" gorm:"size:255;comment:标题"`
	Icon     string  `json:"icon" gorm:"size:255;comment:图标"`
	Path     string  `json:"path" gorm:"size:255;comment:路径"`
	ParentId *string `json:"parent_id" gorm:"size:255;comment:父级ID"`

	Children []*Router `json:"children,omitempty" gorm:"-"`
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

// 获取树形结构
func (Router) GetTree(Pid *string) ([]*Router, error) {
	var list []*Router
	db := DB
	if Pid != nil {
		db = db.Where(&Router{ParentId: Pid})
	} else {
		db = db.Where("parent_id IS NULL")
	}
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
