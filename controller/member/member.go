package member

import (
	"jdy/controller"
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
	con.Success(ctx, "ok", nil)
}

// 会员详情
func (con MemberController) Info(ctx *gin.Context) {
	con.Success(ctx, "ok", nil)
}
