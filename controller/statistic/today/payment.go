package today

import (
	"jdy/enums"
	"jdy/logic/statistic/today"

	"github.com/gin-gonic/gin"
)

// 今日收支
func (con ToDayController) Payment(ctx *gin.Context) {
	var (
		req   today.PaymentReq
		logic = today.ToDayLogic{}
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
		con.Exception(ctx, "无法获取")
		return
	} else {
		logic.Staff = staff

		if staff.Identity < enums.IdentityAreaManager && req.StoreId == "" {
			con.Exception(ctx, "参数错误")
			return
		}
	}

	res, err := logic.Payment(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", res)
}
