package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 成功响应
func (con BaseController) SuccessJson(c *gin.Context, message string, data interface{}) {
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
func (con BaseController) ErrorJson(c *gin.Context, code int, message string) {
	response := gin.H{
		"code":    code,
		"message": message,
	}

	c.JSON(http.StatusOK, response)
}

// 异常响应
func (con BaseController) ExceptionJson(c *gin.Context, message string) {
	response := gin.H{
		"code":    http.StatusInternalServerError,
		"message": message,
	}

	c.JSON(http.StatusOK, response)
	c.Abort()
}
