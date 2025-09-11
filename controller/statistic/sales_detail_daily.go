package statistic

import (
	"jdy/enums"
	"jdy/logic/statistic"
	"jdy/types"
	"time"

	"github.com/gin-gonic/gin"
)

func (con StatisticController) SalesDetailDaily(ctx *gin.Context) {
	var (
		req   types.StatisticSalesDetailDailyReq
		logic = statistic.StatisticLogic{}
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, "参数错误")
		return
	}

	// 校验参数
	if err := req.Validate(); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	// 获取当前登录用户
	if staff, err := con.GetStaff(ctx); err != nil {
		con.ExceptionWithAuth(ctx, err)
		return
	} else {
		logic.Staff = staff
	}

	// 如果是区域经理以下的身份
	if logic.Staff.Identity < enums.IdentityAreaManager {
		now := time.Now()
		// 开始时间不能小于 2 个月前那一天的 0 点
		start_limit := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).AddDate(0, -2, 0)
		if req.StartTime.Before(start_limit) {
			con.Exception(ctx, "开始时间不能小于 2 个月前")
			return
		}
	}

	res, err := logic.SalesDetailDaily(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", res)
}
