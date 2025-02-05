package workbench

import (
	"fmt"
	"jdy/errors"
	"jdy/logic"
	"jdy/model"
	"jdy/types"
)

type WorkbenchLogic struct {
	logic.BaseLogic
}

// 获取路由列表
func (l WorkbenchLogic) GetList() ([]*model.Router, *errors.Errors) {
	list, err := model.Router{}.GetTree(nil)
	if err != nil {
		return nil, errors.New("获取工作台列表失败: " + err.Error())
	}

	return list, nil
}

// 搜索路由
func (l WorkbenchLogic) Search(req *types.WorkbenchSearchReq) ([]*model.Router, *errors.Errors) {
	var list []*model.Router
	if err := model.DB.Debug().Where("title like ?", fmt.Sprintf("%%%s%%", req.Keyword)).Where("path <> ''").Find(&list).Error; err != nil {
		return nil, errors.New("搜索失败")
	}

	return list, nil
}

// 添加路由
func (l WorkbenchLogic) AddRoute(req *types.WorkbenchAddReq) (*model.Router, *errors.Errors) {
	route := model.Router{
		Title: req.Title,
		Path:  req.Path,
		Icon:  req.Icon,
	}

	if req.ParentId != "" {
		route.ParentId = &req.ParentId
	}

	if err := model.DB.Save(&route).Error; err != nil {
		return nil, errors.New("添加路由失败: " + err.Error())
	}

	return &route, nil
}

// 删除路由
func (l WorkbenchLogic) DelRoute(id string) *errors.Errors {
	var route model.Router
	if err := model.DB.First(&route, id).Error; err != nil {
		return errors.New("删除失败，不存在或已被删除")
	}

	var count int64
	if err := model.DB.Model(&model.Router{}).Where("parent_id = ?", id).Count(&count).Error; err != nil || count > 0 {
		return errors.New("无法删除，请先删除下级")
	}

	if err := model.DB.Delete(&route).Error; err != nil {
		return errors.New("删除失败: " + err.Error())
	}
	return nil
}

// 更新路由
func (l WorkbenchLogic) UpdateRoute(req *types.WorkbenchUpdateReq) *errors.Errors {
	var route model.Router
	if err := model.DB.First(&route, req.Id).Error; err != nil {
		return errors.New("更新失败，不存在或已被删除")
	}

	if err := model.DB.Model(&route).Updates(&req).Error; err != nil {
		return errors.New("更新失败: " + err.Error())
	}

	return nil
}
