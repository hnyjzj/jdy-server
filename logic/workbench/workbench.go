package workbench

import (
	"fmt"
	"jdy/enums"
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
	var inIds []string
	if l.Staff.Identity < enums.IdentityAdmin {
		for _, v := range l.Staff.Role.Routers {
			inIds = append(inIds, v.Id)
		}
	}

	list, err := model.Router{}.GetTree(nil, inIds)
	if err != nil {
		return nil, errors.New("获取工作台列表失败: " + err.Error())
	}

	if l.Staff.Identity >= enums.IdentityAdmin {
		return list, nil
	}

	return list, nil
}

// 搜索路由
func (l WorkbenchLogic) Search(req *types.WorkbenchSearchReq) ([]*model.Router, *errors.Errors) {
	var inIds []string
	if l.Staff.Identity < enums.IdentityAdmin {
		for _, v := range l.Staff.Role.Routers {
			inIds = append(inIds, v.Id)
		}
	}
	var (
		list []*model.Router
		db   = model.DB.Model(&model.Router{})
	)
	db = db.Where("title like ?", fmt.Sprintf("%%%s%%", req.Keyword))
	db = db.Where("path <> ''")
	if len(inIds) > 0 {
		db = db.Where("id in (?)", inIds)
	}
	if err := db.Find(&list).Error; err != nil {
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
	if err := model.DB.First(&route, "id = ?", id).Error; err != nil {
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
	if err := model.DB.First(&route, "id = ?", req.Id).Error; err != nil {
		return errors.New("更新失败，不存在或已被删除")
	}

	if err := model.DB.Model(&route).Updates(&req).Error; err != nil {
		return errors.New("更新失败: " + err.Error())
	}

	return nil
}
