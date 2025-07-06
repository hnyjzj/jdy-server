package statistic

import (
	"jdy/controller"
	"jdy/logic/statistic"
	"jdy/types"

	"github.com/gin-gonic/gin"
)

type StatisticController struct {
	controller.BaseController
}

// 门店销售统计
func (con StatisticController) StoreSalesTotal(ctx *gin.Context) {
	var (
		req   types.StatisticStoreSalesTotalReq
		logic = statistic.StatisticLogic{}
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, "参数错误")
		return
	}

	res, err := logic.StoreSalesTotal(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", res)
}

// 今日销售统计
func (con StatisticController) TodaySales(ctx *gin.Context) {
	var (
		req   types.StatisticTodaySalesReq
		logic = statistic.StatisticLogic{}
	)

	// 获取当前登录用户
	if staff, err := con.GetStaff(ctx); err != nil {
		con.ExceptionWithAuth(ctx, err)
		return
	} else {
		logic.Staff = staff
	}

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, "参数错误")
		return
	}

	res, err := logic.TodaySales(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", res)
}

// 今日货品
func (con StatisticController) TodayProduct(ctx *gin.Context) {
	var (
		req   types.StatisticTodayProductReq
		logic = statistic.StatisticLogic{}
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, "参数错误")
		return
	}

	res, err := logic.TodayProduct(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", res)
}
