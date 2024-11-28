package controller

import (
	"jdy/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 成功响应
func (con BaseController) Success(c *gin.Context, message string, data any) {
	response := gin.H{
		"code":    http.StatusOK,
		"message": message,
	}

	if data != nil {
		response["data"] = data
	}

	c.JSON(http.StatusOK, response)
}

// 失败响应
func (con BaseController) Error(c *gin.Context, code int, message string) {
	response := gin.H{
		"code":    code,
		"message": message,
	}

	c.JSON(http.StatusOK, response)
}

// 逻辑失败响应
func (con BaseController) ErrorLogic(c *gin.Context, err *errors.Errors) {
	response := gin.H{
		"code":    err.Code,
		"message": err.Message,
	}

	c.JSON(http.StatusOK, response)
}

// 异常响应
func (con BaseController) Exception(c *gin.Context, message string) {
	response := gin.H{
		"code":    http.StatusInternalServerError,
		"message": message,
	}

	c.JSON(http.StatusOK, response)
	c.Abort()
}
