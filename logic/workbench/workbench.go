package workbench

import (
	"jdy/errors"
	"jdy/logic"
	"jdy/model"
	"jdy/types"
)

type WorkbenchLogic struct {
	logic.Base
}

// 获取路由列表
func (l WorkbenchLogic) GetList() ([]*model.Router, *errors.Errors) {
	list, err := model.Router{}.GetTree(nil)
	if err != nil {
		return nil, errors.New("获取工作台列表失败: " + err.Error())
	}

	return list, nil
}

// 添加路由
func (l WorkbenchLogic) AddRoute(req *types.WorkbenchListReq) (*model.Router, *errors.Errors) {
	route := &model.Router{
		Title: req.Title,
		Path:  req.Path,
		Icon:  req.Icon,
	}

	if req.ParentId != "" {
		route.ParentId = &req.ParentId
	}

	if err := model.DB.Save(route).Error; err != nil {
		return nil, errors.New("添加路由失败: " + err.Error())
	}

	return route, nil
}
