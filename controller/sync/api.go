package sync

import (
	"jdy/model"
	G "jdy/service/gin"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ApiController struct {
	SyncController
}

func (con ApiController) List(ctx *gin.Context) {
	routes := G.Gin.Routes()
	var apis []model.Api
	for _, route := range routes {
		var api model.Api
		if err := model.DB.Where(model.Api{
			Path:   route.Path,
			Method: route.Method,
		}).First(&api).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				apis = append(apis, model.Api{
					Path:   route.Path,
					Method: route.Method,
				})
			} else {
				con.Exception(ctx, err.Error())
				return
			}
		}
	}

	var success int64 = 0
	if len(apis) > 0 {
		res := model.DB.Create(&apis)
		if err := res.Error; err != nil {
			con.Exception(ctx, err.Error())
			return
		}
		success = res.RowsAffected
	}

	con.Success(ctx, "ok", success)
}
