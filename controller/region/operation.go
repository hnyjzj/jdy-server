package region

import (
	"jdy/errors"
	"jdy/logic/region"
	"jdy/types"

	"github.com/gin-gonic/gin"
)

// 创建区域
func (con RegionController) Create(ctx *gin.Context) {
	var (
		req   types.RegionCreateReq
		logic = region.RegionLogic{}
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	// 创建区域
	if err := logic.Create(ctx, &req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "创建区域成功", nil)
}

// 更新区域
func (con RegionController) Update(ctx *gin.Context) {
	var (
		req   types.RegionUpdateReq
		logic = region.RegionLogic{}
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	// 更新区域
	if err := logic.Update(ctx, &req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "更新区域成功", nil)
}

// 删除区域
func (con RegionController) Delete(ctx *gin.Context) {
	var (
		req   types.RegionDeleteReq
		logic = region.RegionLogic{}
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	// 删除区域
	if err := logic.Delete(ctx, &req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "删除区域成功", nil)
}
