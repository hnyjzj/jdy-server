package store

import (
	"fmt"
	"jdy/errors"
	"jdy/logic/store"
	"jdy/types"

	"github.com/gin-gonic/gin"
)

// 创建门店
func (con StoreController) Create(ctx *gin.Context) {
	var (
		req   types.StoreCreateReq
		logic = store.StoreLogic{}
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		fmt.Printf("err.Error(): %v\n", err.Error())
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	// 创建门店
	if err := logic.Create(ctx, &req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "创建门店成功", nil)
}
