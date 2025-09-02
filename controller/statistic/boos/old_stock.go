package boos

import (
	"jdy/logic/statistic/boos/old_stock"
	"jdy/utils"

	"github.com/gin-gonic/gin"
)

func (con BoosController) OldStockWhere(ctx *gin.Context) {
	where := utils.StructToWhere(old_stock.Where{})

	con.Success(ctx, "ok", where)
}

// 旧料库存统计
func (con BoosController) OldStockData(ctx *gin.Context) {
	var (
		req   old_stock.DataReq
		logic = old_stock.Logic{}
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
