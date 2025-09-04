package sync

import (
	"jdy/model"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ApiLogic struct {
	SyncLogic
}

type ApiListReq struct {
	Routes gin.RoutesInfo
}

func (a *ApiLogic) List(req *ApiListReq) error {
	if len(req.Routes) == 0 {
		return nil
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 清空api表
		if err := tx.Where("1 = 1").Delete(&model.Api{}).Error; err != nil {
			return err
		}

		// 遍历路由
		for _, route := range req.Routes {
			// 解析路由路径
			paths := strings.Split(route.Path, "/")

			// 跳过根路由
			if paths[1] != "api" {
				continue
			}

			// 解析路由名称
			names := paths[1:]

			var (
				parent = &model.Api{}
				err    error
			)

			for i := range names {
				parentPaths := "/" + strings.Join(names[:i], "/")
				parentPaths = strings.TrimSpace(parentPaths)

				pr := route
				pr.Path = parentPaths

				parent, err = a.setParent(tx, pr, parent)
				if err != nil {
					return err
				}
			}

			if err := a.setApi(tx, route, parent); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (a *ApiLogic) setParent(tx *gorm.DB, route gin.RouteInfo, p *model.Api) (*model.Api, error) {
	var (
		parent model.Api

		data = model.Api{
			Path: route.Path,
		}
	)

	db := tx.Unscoped()
	if p != nil && p.Id != "" {
		data.ParentId = &p.Id
	} else {
		db = db.Where("parent_id is null")
	}
	db = db.Where(data).Attrs(data)

	if err := db.FirstOrCreate(&parent).Error; err != nil {
		return nil, err
	}

	if parent.DeletedAt.Valid {
		if err := tx.Unscoped().Model(&model.Api{}).Where("id = ?", parent.Id).Update("deleted_at", nil).Error; err != nil {
			return nil, err
		}
	}

	return &parent, nil
}

func (a *ApiLogic) setApi(tx *gorm.DB, route gin.RouteInfo, p *model.Api) error {
	var (
		api model.Api

		data = model.Api{
			Path:     route.Path,
			Method:   route.Method,
			ParentId: &p.Id,
		}
	)

	if err := tx.Model(&model.Api{}).Unscoped().Where(data).Attrs(data).FirstOrCreate(&api).Error; err != nil {
		return err
	}

	if api.DeletedAt.Valid {
		if err := tx.Unscoped().Model(&model.Api{}).Where("id = ?", api.Id).Update("deleted_at", nil).Error; err != nil {
			return err
		}
	}

	return nil
}
