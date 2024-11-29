package store

import (
	"fmt"
	"jdy/controller"
	"jdy/errors"
	"jdy/logic/store"
	"jdy/types"

	"github.com/gin-gonic/gin"
)

type StoreController struct {
	controller.BaseController
}

// 创建门店
func (con StoreController) Create(ctx *gin.Context) {
	var (
		req   types.StoreCreateReq
		logic = store.StoreLogic{}
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	// 创建门店
	if err := logic.Create(ctx, &req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "创建门店成功", nil)
}

// 删除门店
func (con StoreController) Delete(ctx *gin.Context) {
	var (
		req   types.StoreDeleteReq
		logic = store.StoreLogic{}
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	fmt.Printf("req: %v\n", req)

	// 删除门店
	if err := logic.Delete(ctx, &req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "删除门店成功", nil)
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
