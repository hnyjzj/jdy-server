package member

import (
	"jdy/controller"
	"jdy/errors"
	"jdy/logic/member"
	"jdy/types"
	"jdy/utils"

	"github.com/gin-gonic/gin"
)

type MemberController struct {
	controller.BaseController
}

// 会员筛选条件
func (con MemberController) Where(ctx *gin.Context) {
	where := utils.StructToWhere(types.MemberWhere{})

	con.Success(ctx, "ok", where)
}

// 会员列表
func (con MemberController) List(ctx *gin.Context) {
	var (
		req types.MemberListReq

		logic = member.MemberLogic{
			Ctx: ctx,
		}
	)

	staff, err := con.GetStaff(ctx)
	if err != nil {
		con.ExceptionWithAuth(ctx, err.Error())
		return
	}
	logic.Staff = staff

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	res, err := logic.List(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", res)
}

// 会员详情
func (con MemberController) Info(ctx *gin.Context) {
	var (
		req types.MemberInfoReq

		logic = member.MemberLogic{
			Ctx: ctx,
		}
	)

	staff, err := con.GetStaff(ctx)
	if err != nil {
		con.ExceptionWithAuth(ctx, err.Error())
		return
	}
	logic.Staff = staff

	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	if err := req.Validate(); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	res, err := logic.Info(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", res)
}
