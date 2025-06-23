package model

type Api struct {
	SoftDelete

	Title    string  `json:"title" gorm:"column:title;size:255;comment:标题"`                            // 标题
	Path     string  `json:"path" gorm:"uniqueIndex:unique_api;column:path;size:255;comment:接口地址"`     // 接口地址
	Method   string  `json:"method" gorm:"uniqueIndex:unique_api;column:method;size:255;comment:请求方式"` // 请求方式
	ParentId *string `json:"parent_id" gorm:"size:255;comment:父级ID"`                                   // 父级ID

	Children []*Api `json:"children,omitempty" gorm:"-"`
}

// 获取树形结构
func (Api) GetTree(Path, Pid *string) ([]*Api, error) {
	var list []*Api
	db := DB.Model(&Api{})

	if Path != nil {
		db = db.Where(&Api{Path: *Path})
	}
	if Pid != nil {
		db = db.Where(&Api{ParentId: Pid})
	}
	if Path == nil && Pid == nil {
		db = db.Where("parent_id IS NULL")
	}

	if err := db.Find(&list).Error; err != nil {
		return nil, err
	}

	for _, v := range list {
		children, err := v.GetTree(nil, &v.Id)
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
		&Api{},
	)
	// 重置表
	RegisterRefreshModels(
	// &Api{},
	)
}
