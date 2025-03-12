package product

import (
	"jdy/controller"
	"jdy/errors"
	"jdy/logic/product"
	"jdy/types"
	"jdy/utils"

	"github.com/gin-gonic/gin"
)

type ProductEnterController struct {
	controller.BaseController
}

// 入库单筛选条件
func (con ProductEnterController) Where(ctx *gin.Context) {
	where := utils.StructToWhere(types.ProductEnterWhere{})

	con.Success(ctx, "ok", where)
}

// 入库单
func (con ProductEnterController) Create(ctx *gin.Context) {
	var (
		req types.ProductEnterCreateReq

		logic = product.ProductEnterLogic{
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
	res, err := logic.Create(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", res)
}

// 入库单列表
func (con ProductEnterController) List(ctx *gin.Context) {
	var (
		req types.ProductEnterListReq

		logic = product.ProductEnterLogic{
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
	res, err := logic.EnterList(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", res)
}

// 入库单详情
func (con ProductEnterController) Info(ctx *gin.Context) {
	var (
		req types.ProductEnterInfoReq

		logic = product.ProductEnterLogic{
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
	res, err := logic.EnterInfo(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", res)
}

// 入库单添加产品
func (con ProductEnterController) AddProduct(ctx *gin.Context) {
	var (
		req types.ProductEnterAddProductReq

		logic = product.ProductEnterLogic{
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
	res, err := logic.AddProduct(&req)
	if err != nil {
		con.ExceptionWithResult(ctx, err.Error(), res)
		return
	}

	con.Success(ctx, "ok", res)
}

// 入库单编辑产品
func (con ProductEnterController) EditProduct(ctx *gin.Context) {
	var (
		req types.ProductEnterEditProductReq

		logic = product.ProductEnterLogic{
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
	if err := logic.EditProduct(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}

// 入库单删除产品
func (con ProductEnterController) DelProduct(ctx *gin.Context) {
	var (
		req types.ProductEnterDelProductReq

		logic = product.ProductEnterLogic{
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
	if err := logic.DelProduct(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}

// 入库单完成
func (con ProductEnterController) Finish(ctx *gin.Context) {
	var (
		req types.ProductEnterFinishReq

		logic = product.ProductEnterLogic{
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
	if err := logic.Finish(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}

// 入库单取消
func (con ProductEnterController) Cancel(ctx *gin.Context) {
	var (
		req types.ProductEnterCancelReq

		logic = product.ProductEnterLogic{
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
	if err := logic.Cancel(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}
