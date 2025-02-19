package member

import (
	"jdy/errors"
	"jdy/logic/member"
	"jdy/types"

	"github.com/gin-gonic/gin"
)

func (con MemberController) Create(ctx *gin.Context) {
	var (
		req   types.MemberCreateReq
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

	// 创建会员
	if err := logic.Create(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "创建成功", nil)
}

// 更新会员
func (con MemberController) Update(ctx *gin.Context) {
	var (
		req   types.MemberUpdateReq
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

	// 创建会员
	if err := logic.Update(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "更新成功", nil)
}
