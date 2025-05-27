package product

import (
	"jdy/controller"
	"jdy/errors"
	"jdy/logic/product"
	"jdy/types"
	"jdy/utils"
	"log"

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
			Ctx: ctx,
		}
	)

	staff, err := con.GetStaff(ctx)
	if err != nil {
		con.ExceptionWithAuth(ctx, err.Error())
		return
	}
	logic.Staff = staff

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		log.Printf("err.Error(): %v\n", err.Error())
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}
	if err := req.Validate(); err != nil {
		log.Printf("err.Error(): %v\n", err.Error())
		con.Exception(ctx, err.Error())
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
			Ctx: ctx,
		}
	)

	staff, err := con.GetStaff(ctx)
	if err != nil {
		con.ExceptionWithAuth(ctx, err.Error())
		return
	}
	logic.Staff = staff

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

func (con ProductInventoryController) Info(ctx *gin.Context) {
	var (
		req types.ProductInventoryInfoReq

		logic = product.ProductInventoryLogic{
			Ctx: ctx,
		}
	)

	staff, err := con.GetStaff(ctx)
	if err != nil {
		con.ExceptionWithAuth(ctx, err.Error())
		return
	}
	logic.Staff = staff

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	// 调用逻辑层
	res, err := logic.Info(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", res)
}

func (con ProductInventoryController) Add(ctx *gin.Context) {
	var (
		req types.ProductInventoryAddReq

		logic = product.ProductInventoryLogic{
			Ctx: ctx,
		}
	)

	staff, err := con.GetStaff(ctx)
	if err != nil {
		con.ExceptionWithAuth(ctx, err.Error())
		return
	}
	logic.Staff = staff

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	// 调用逻辑层
	if err := logic.Add(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}

func (con ProductInventoryController) Change(ctx *gin.Context) {
	var (
		req types.ProductInventoryChangeReq

		logic = product.ProductInventoryLogic{
			Ctx: ctx,
		}
	)

	staff, err := con.GetStaff(ctx)
	if err != nil {
		con.ExceptionWithAuth(ctx, err.Error())
		return
	}
	logic.Staff = staff

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	// 调用逻辑层
	if err := logic.Change(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}
