package region

import (
	"jdy/controller"
	"jdy/errors"
	"jdy/logic/region"
	"jdy/types"
	"jdy/utils"

	"github.com/gin-gonic/gin"
)

type RegionController struct {
	controller.BaseController
}

// 区域筛选条件
func (con RegionController) Where(ctx *gin.Context) {
	where := utils.StructToWhere(types.RegionWhere{})

	con.Success(ctx, "ok", where)
}

// 区域列表
func (con RegionController) List(ctx *gin.Context) {
	var (
		req   types.RegionListReq
		logic = region.RegionLogic{}
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	// 查询区域列表
	list, err := logic.List(ctx, &req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", list)
}

// 我的区域列表
func (con RegionController) My(ctx *gin.Context) {
	var (
		req   types.RegionListMyReq
		logic = region.RegionLogic{
			Ctx: ctx,
		}
	)

	if staff, err := con.GetStaff(ctx); err != nil {
		con.ExceptionWithAuth(ctx, err)
		return
	} else {
		logic.Staff = staff
	}

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	// 查询区域列表
	list, err := logic.My(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", list)
}

// 区域详情
func (con RegionController) Info(ctx *gin.Context) {
	var (
		req   types.RegionInfoReq
		logic = region.RegionLogic{}
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	// 查询区域详情
	info, err := logic.Info(ctx, &req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", info)
}
