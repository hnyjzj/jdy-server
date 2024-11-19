package common

import (
	"jdy/controller"
	"jdy/logic/common"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CaptchaController struct {
	controller.BaseController
}

// 获取图片验证码
func (con CaptchaController) Image(c *gin.Context) {
	var (
		logic common.CaptchaLogic
	)

	res, err := logic.ImageCaptcha()

	if err != nil {
		con.Error(c, http.StatusInternalServerError, "获取失败")
		return
	}

	con.Success(c, "ok", res)
}
