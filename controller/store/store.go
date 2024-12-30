package store

import (
	"jdy/controller"
	"jdy/errors"
	"jdy/logic/store"
	"jdy/types"
	"jdy/utils"

	"github.com/gin-gonic/gin"
)

type StoreController struct {
	controller.BaseController
}

// 门店筛选条件
func (con StoreController) Where(ctx *gin.Context) {
	where := utils.StructToWhere(types.StoreWhere{})

	con.Success(ctx, "ok", where)
}

// 门店列表
func (con StoreController) List(ctx *gin.Context) {
	var (
		req   types.StoreListReq
		logic = store.StoreLogic{}
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	// 查询门店列表
	list, err := logic.List(ctx, &req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", list)
}

// 门店详情
func (con StoreController) Info(ctx *gin.Context) {
	var (
		req   types.StoreInfoReq
		logic = store.StoreLogic{}
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	// 查询门店详情
	info, err := logic.Info(ctx, &req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", info)
}
