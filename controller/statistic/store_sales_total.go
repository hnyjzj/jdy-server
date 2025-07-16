package statistic

import (
	"jdy/enums"
	"jdy/logic/statistic"
	"jdy/types"

	"github.com/gin-gonic/gin"
)

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

	// 获取当前登录用户
	if staff, err := con.GetStaff(ctx); err != nil {
		con.Exception(ctx, "无法获取")
		return
	} else {
		logic.Staff = staff
	}

	if logic.Staff.Identity < enums.IdentityAreaManager {
		con.Exception(ctx, "权限不足")
		return
	}

	_, err := logic.StoreSalesTotal(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	// con.Success(ctx, "ok", res)
	con.Exception(ctx, "暂未开放")
}
