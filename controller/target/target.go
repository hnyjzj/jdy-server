package target

import (
	"jdy/controller"
	"jdy/errors"
	"jdy/logic/target"
	"jdy/types"
	"jdy/utils"

	"github.com/gin-gonic/gin"
)

type TargetController struct {
	controller.BaseController
}

// 筛选条件
func (con TargetController) Where(ctx *gin.Context) {
	where := utils.StructToWhere(types.TargetWhere{})

	con.Success(ctx, "ok", where)
}

// 销售目标分组筛选
func (con TargetController) WhereGroup(ctx *gin.Context) {
	where := utils.StructToWhere(types.TargetWhereGroup{})

	con.Success(ctx, "ok", where)
}

// 销售目标个人筛选
func (con TargetController) WherePersonal(ctx *gin.Context) {
	where := utils.StructToWhere(types.TargetWherePersonal{})

	con.Success(ctx, "ok", where)
}

// 创建
func (con TargetController) Create(ctx *gin.Context) {
	var (
		req   types.TargetCreateReq
		logic target.Logic
	)

	if staff, err := con.GetStaff(ctx); err != nil {
		con.ExceptionWithAuth(ctx, err)
		return
	} else {
		logic.Staff = staff
		logic.Ctx = ctx
	}

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	// 创建目标
	if err := logic.Create(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}
