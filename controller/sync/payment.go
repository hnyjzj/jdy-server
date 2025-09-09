package sync

import (
	"jdy/logic/sync"

	"github.com/gin-gonic/gin"
)

type PaymentController struct {
	SyncController
}

// 同步通讯录
func (con PaymentController) SyncPayments(ctx *gin.Context) {
	logic := sync.PaymentLogic{}

	logic.Ctx = ctx

	if err := logic.Payments(); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}
