package sync

import (
	"jdy/model"
	G "jdy/service/gin"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

type ApiController struct {
	SyncController
}

func (con ApiController) List(ctx *gin.Context) {
	routes := G.Gin.Routes()
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

		var api model.Api
		api.Path = route.Path
		api.Method = route.Method
		api.ParentId = &parent.Id
		if err := model.DB.Where(model.Api{Path: route.Path, Method: route.Method}).Attrs(api).FirstOrCreate(&api).Error; err != nil {
			log.Println(err)
			return
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

	if err := model.DB.Where(parent).FirstOrCreate(&parent).Error; err != nil {
		return nil, err
	}

	return &parent, nil
}
