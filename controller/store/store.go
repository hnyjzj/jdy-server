package store

import (
	"jdy/controller"
	"jdy/logic/store"

	"github.com/gin-gonic/gin"
)

type StoreController struct {
	controller.BaseController
}

// 门店列表
func (con StoreController) List(ctx *gin.Context) {
	var (
		logic = store.StoreLogic{}
	)

	// 查询门店列表
	list, err := logic.List(ctx)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", list)
}
