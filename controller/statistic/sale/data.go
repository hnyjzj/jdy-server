package sale

import (
	"jdy/enums"
	"jdy/logic/statistic/sale"

	"github.com/gin-gonic/gin"
)

func (con SaleController) Data(ctx *gin.Context) {
	var (
		req   sale.DataReq
		logic = sale.StatisticSaleLogic{}

		onlyself bool
	)

	// 获取请求参数
	if err := ctx.ShouldBindJSON(&req); err != nil {
		con.Exception(ctx, "参数错误")
		return
	}

	// 获取当前登录用户
	if staff, err := con.GetStaff(ctx); err != nil {
		con.ExceptionWithAuth(ctx, err)
		return
	} else {
		logic.Staff = staff
		logic.Ctx = ctx

		if staff.Identity < enums.IdentityAreaManager && req.StoreId == "" {
			con.Exception(ctx, "参数错误")
			return
		}

		if staff.Identity < enums.IdentityShopkeeper {
			onlyself = true
		}
	}

	res, err := logic.Data(&req, onlyself)
	if err != nil {
		con.Exception(ctx, "获取失败")
		return
	}

	con.Success(ctx, "ok", res)
}
