package boos

import (
	"jdy/logic/statistic/boos/old_sales"
	"jdy/utils"

	"github.com/gin-gonic/gin"
)

func (con BoosController) OldSalesWhere(ctx *gin.Context) {
	where := utils.StructToWhere(old_sales.Where{})

	con.Success(ctx, "ok", where)
}

func (con BoosController) OldSalesTitles(ctx *gin.Context) {

	var (
		logic = old_sales.Logic{}
	)

	res := logic.GetTitles()

	con.Success(ctx, "ok", res)
}

// 旧料库存统计
func (con BoosController) OldSalesData(ctx *gin.Context) {
	var (
		req   old_sales.DataReq
		logic = old_sales.Logic{}
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
	}

	res, err := logic.GetDatas(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", res)
}
