package workbench

import (
	"jdy/errors"
	"jdy/logic"
	"jdy/model"
	"jdy/types"

	"github.com/acmestack/gorm-plus/gplus"
)

type WorkbenchLogic struct {
	logic.Base
}

// 获取路由列表
func (l WorkbenchLogic) GetList() ([]*model.Router, *errors.Errors) {
	var (
		workbenchList []*model.Router
	)

	list, err := model.Router{}.GetTree(nil)
	if err != nil {
		return nil, errors.ErrNotFound
	}

	workbenchList = list

	return workbenchList, nil
}

// 添加路由
func (l WorkbenchLogic) AddRoute(req *types.WorkbenchListReq) (*model.Router, *errors.Errors) {
	user := model.Router{
		Title:    req.Title,
		Path:     req.Path,
		Icon:     req.Icon,
		ParentId: &req.ParentId,
	}
	result := gplus.Insert(&user)
	if result.Error != nil {
		return nil, errors.New(result.Error.Error())
	}

	return &user, nil
}