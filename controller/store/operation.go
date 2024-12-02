package store

import (
	"jdy/errors"
	"jdy/logic/store"
	"jdy/types"

	"github.com/gin-gonic/gin"
)

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

// 更新门店
func (con StoreController) Update(ctx *gin.Context) {
	var (
		req   types.StoreUpdateReq
		logic = store.StoreLogic{}
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	// 更新门店
	if err := logic.Update(ctx, &req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "更新门店成功", nil)
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

	// 删除门店
	if err := logic.Delete(ctx, &req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "删除门店成功", nil)
}
