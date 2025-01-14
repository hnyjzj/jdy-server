package member

import (
	"jdy/errors"
	"jdy/logic/member"
	"jdy/types"

	"github.com/gin-gonic/gin"
)

func (con MemberController) Integral(ctx *gin.Context) {
	var (
		req types.MemberIntegralListReq

		logic = member.MemberLogic{
			Ctx:   ctx,
			Staff: con.GetStaff(ctx),
		}
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	res, err := logic.Integral(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", res)
}
