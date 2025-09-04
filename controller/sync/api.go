package sync

import (
	"jdy/model"
	G "jdy/service/gin"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ApiController struct {
	SyncController
}

func (con ApiController) List(ctx *gin.Context) {
	routes := G.Gin.Routes()

	// 清空api表
	if err := model.DB.Where("1 = 1").Delete(&model.Api{}).Error; err != nil {
		log.Println(err)
		return
	}

	for _, route := range routes {
		names := strings.Split(route.Path, "/")
		if names[1] != "api" {
			continue
		}

		var (
			parent *model.Api
			err    error
		)

		paths := names[1:]
		for i := range paths {
			parentPaths := "/" + strings.Join(paths[:i], "/")
			parentPaths = strings.TrimSpace(parentPaths)
			parent, err = con.setParent(parentPaths, parent)
			if err != nil {
				log.Println(err)
			}
		}

		var (
			api model.Api
		)
		db := model.DB.Model(&api)
		db = db.Unscoped()
		db = db.Where(model.Api{
			Path:   route.Path,
			Method: route.Method,
		})

		if err := db.Find(&api).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				log.Println(err)
				return
			}
		}

		if api.Id == "" {
			data := model.Api{
				Path:     route.Path,
				Method:   route.Method,
				ParentId: &parent.Id,
			}
			if err := model.DB.Create(&data).Error; err != nil {
				log.Println(err)
				return
			}
		} else if api.DeletedAt.Valid {
			if err := model.DB.Unscoped().Model(&api).Update("deleted_at", nil).Error; err != nil {
				log.Println(err)
				return
			}
		}
	}
	con.Success(ctx, "ok", nil)
}

func (con ApiController) setParent(path string, p *model.Api) (*model.Api, error) {
	parent := model.Api{
		Path: path,
	}
	if p != nil {
		parent.ParentId = &p.Id
	}

	if err := model.DB.Unscoped().Where(parent).FirstOrCreate(&parent).Error; err != nil {
		return nil, err
	}
	if parent.DeletedAt.Valid {
		if err := model.DB.Unscoped().Model(&parent).Update("deleted_at", nil).Error; err != nil {
			return nil, err
		}
	}

	return &parent, nil
}
