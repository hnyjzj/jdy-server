package model

import (
	"github.com/acmestack/gorm-plus/gplus"
	"gorm.io/gorm"
)

type Router struct {
	SoftDelete

	Title string `json:"title" gorm:"size:255;comment:标题"`
	Icon  string `json:"icon" gorm:"size:255;comment:图标"`
	Path  string `json:"path" gorm:"size:255;comment:路径"`

	ParentId *string   `json:"parent_id" gorm:"comment:父级ID"`
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
	query, u := gplus.NewQuery[Router]()
	if Pid != nil {
		query.Eq(&u.ParentId, Pid)
	} else {
		query.IsNull(&u.ParentId)
	}

	list, db := gplus.SelectList(query)
	if db.Error != nil && db.Error != gorm.ErrRecordNotFound {
		return nil, db.Error
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
