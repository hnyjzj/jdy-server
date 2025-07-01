package region

import (
	"jdy/errors"
	"jdy/logic/region"
	"jdy/types"

	"github.com/gin-gonic/gin"
)

type RegionSuperiorController struct {
	RegionController
}

// 区域负责人列表
func (con RegionSuperiorController) List(ctx *gin.Context) {
	var (
		req   types.RegionSuperiorListReq
		logic = region.RegionSuperiorLogic{
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
	list, err := logic.List(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", list)
}

func (con RegionSuperiorController) Add(ctx *gin.Context) {
	var (
		req   types.RegionSuperiorAddReq
		logic = region.RegionSuperiorLogic{
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

	// 添加区域
	if err := logic.Add(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}

func (con RegionSuperiorController) Del(ctx *gin.Context) {
	var (
		req   types.RegionSuperiorDelReq
		logic = region.RegionSuperiorLogic{
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

	// 删除区域
	if err := logic.Del(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}
