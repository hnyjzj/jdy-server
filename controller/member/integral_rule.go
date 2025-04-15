package member

import (
	"jdy/errors"
	"jdy/logic/member"
	"jdy/types"

	"github.com/gin-gonic/gin"
)

type MemberIntegralRuleController struct {
	MemberController
}

func (con MemberIntegralRuleController) Finished(ctx *gin.Context) {
	var (
		req types.MemberIntegralRuleReq

		logic = member.MemberIntegralRuleLogic{
			Ctx: ctx,
		}
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	res, err := logic.Finished(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", res)
}

func (con MemberIntegralRuleController) Old(ctx *gin.Context) {
	var (
		req types.MemberIntegralRuleReq

		logic = member.MemberIntegralRuleLogic{
			Ctx: ctx,
		}
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	res, err := logic.Old(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", res)
}

func (con MemberIntegralRuleController) Accessorie(ctx *gin.Context) {
	var (
		req types.MemberIntegralRuleReq

		logic = member.MemberIntegralRuleLogic{
			Ctx: ctx,
		}
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	res, err := logic.Accessorie(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", res)
}
