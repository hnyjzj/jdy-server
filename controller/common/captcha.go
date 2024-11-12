package common

import (
	"jdy/controller"
	commonlogic "jdy/logic/common"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CaptchaController struct {
	controller.BaseController
}

// 获取图片验证码
func (con CaptchaController) GetImage(c *gin.Context) {
	var (
		logic commonlogic.CaptchaLogic
	)

	res, err := logic.ImageCaptcha()

	if err != nil {
		con.ErrorJson(c, http.StatusInternalServerError, "获取失败")
		return
	}

	con.SuccessJson(c, "ok", res)
}
