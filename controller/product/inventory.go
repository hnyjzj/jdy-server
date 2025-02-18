package product

import (
	"jdy/controller"
	"jdy/errors"
	"jdy/logic/product"
	"jdy/types"
	"jdy/utils"

	"github.com/gin-gonic/gin"
)

type ProductInventoryController struct {
	controller.BaseController
}

func (con ProductInventoryController) Where(ctx *gin.Context) {
	where := utils.StructToWhere(types.ProductInventoryWhere{})

	con.Success(ctx, "ok", where)
}

func (con ProductInventoryController) Create(ctx *gin.Context) {
	var (
		req types.ProductInventoryCreateReq

		logic = product.ProductInventoryLogic{
			ProductLogic: product.ProductLogic{
				Ctx:   ctx,
				Staff: con.GetStaff(ctx),
			},
		}
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}
	if err := req.Validate(); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	// 调用逻辑层
	res, err := logic.Create(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", res)
}

func (con ProductInventoryController) List(ctx *gin.Context) {
	var (
		req types.ProductInventoryListReq

		logic = product.ProductInventoryLogic{
			ProductLogic: product.ProductLogic{
				Ctx:   ctx,
				Staff: con.GetStaff(ctx),
			},
		}
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	// 调用逻辑层
	res, err := logic.List(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", res)
}
