package sync

import (
	"jdy/logic/sync"
	G "jdy/service/gin"

	"github.com/gin-gonic/gin"
)

type ApiController struct {
	SyncController
}

func (con ApiController) List(ctx *gin.Context) {
	var (
		logic = sync.ApiLogic{}
		req   = sync.ApiListReq{}
	)

	req.Routes = G.Gin.Routes()
	logic.Ctx = ctx

	if err := logic.List(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}
