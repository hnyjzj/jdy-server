package controller

import (
	"jdy/errors"
	"jdy/model"
	"net/http"

	"github.com/acmestack/gorm-plus/gplus"
	"github.com/gin-gonic/gin"
)

type BaseController struct{}

// 获取 token 中的用户信息
func (BaseController) GetUser(ctx *gin.Context) *model.User {
	userInfo := ctx.MustGet("user").(model.User)
	user, db := gplus.SelectById[model.User](userInfo.Id)

	if db.Error != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    errors.ErrUserNotFound.Code,
			"message": errors.ErrUserNotFound.Message,
		})
		ctx.Abort()
		return nil
	}

	if user.IsDisabled {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    errors.ErrUserDisabled.Code,
			"message": errors.ErrUserDisabled.Message,
		})
		ctx.Abort()
		return nil
	}

	return user
}
