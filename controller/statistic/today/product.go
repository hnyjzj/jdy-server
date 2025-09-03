package today

import (
	"jdy/logic/statistic/today"

	"github.com/gin-gonic/gin"
)

// 今日货品
func (con ToDayController) Product(ctx *gin.Context) {
	var (
		req   today.ProductReq
		logic = today.ToDayLogic{}
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, "参数错误")
		return
	}

	// 获取当前登录用户
	if staff, err := con.GetStaff(ctx); err != nil {
		con.Exception(ctx, "无法获取")
		return
	} else {
		logic.Staff = staff
	}

	res, err := logic.Product(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", res)
}
