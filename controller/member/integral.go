package member

import (
	"jdy/errors"
	"jdy/logic/member"
	"jdy/types"
	"jdy/utils"

	"github.com/gin-gonic/gin"
)

type MemberIntegralController struct {
	MemberController
}

func (con MemberIntegralController) Where(ctx *gin.Context) {
	where := utils.StructToWhere(types.MemberIntegralWhere{})

	con.Success(ctx, "ok", where)
}

func (con MemberIntegralController) List(ctx *gin.Context) {
	var (
		req types.MemberIntegralListReq

		logic = member.MemberIntegralLogic{
			Ctx:   ctx,
			Staff: con.GetStaff(ctx),
		}
	)

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

func (con MemberIntegralController) Change(ctx *gin.Context) {
	var (
		req types.MemberIntegralChangeReq

		logic = member.MemberIntegralLogic{
			Ctx:   ctx,
			Staff: con.GetStaff(ctx),
		}
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	if err := logic.Change(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}
