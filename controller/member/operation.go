package member

import (
	"fmt"
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
		fmt.Printf("err.Error(): %v\n", err.Error())
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

func (con MemberController) Update(ctx *gin.Context) {
	con.Success(ctx, "ok", nil)
}