package product

import (
	"jdy/controller"
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
