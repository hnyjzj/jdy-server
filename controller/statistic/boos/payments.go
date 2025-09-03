package boos

import (
	"jdy/logic/statistic/boos/payments"
	"jdy/utils"

	"github.com/gin-gonic/gin"
)

func (con BoosController) PaymentsWhere(ctx *gin.Context) {
	where := utils.StructToWhere(payments.Where{})

	con.Success(ctx, "ok", where)
}

// 订单收支统计
func (con BoosController) PaymentsData(ctx *gin.Context) {
	var (
		req   payments.DataReq
		logic = payments.Logic{}
	)

	// 获取请求参数
	if err := ctx.ShouldBindJSON(&req); err != nil {
		con.Exception(ctx, "参数错误")
		return
	}

	// 获取当前登录用户
	if staff, err := con.GetStaff(ctx); err != nil {
		con.Exception(ctx, "无法获取")
		return
	} else {
		logic.Staff = staff
		logic.Ctx = ctx
	}

	res, err := logic.GetDatas(&req)
	if err != nil {
		con.Exception(ctx, "获取失败")
		return
	}

	con.Success(ctx, "ok", res)
}
