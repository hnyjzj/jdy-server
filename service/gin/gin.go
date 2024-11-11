package gin

import (
	"fmt"
	"io"
	"jdy/config"

	"github.com/gin-gonic/gin"
)

var Gin *gin.Engine

func Init() *gin.Engine {
	var (
		config = config.Config
	)
	// 设置gin运行模式
	gin.SetMode(config.Server.Mode)
	// 初始化gin
	Gin = gin.New()
	// 设置gin中间件
	Gin.Use(gin.Recovery())
	// 设置gin输出
	if config.Server.Mode == "debug" {
		Gin.Use(gin.Logger())
	}
	// 设置gin输出
	if config.Server.Mode == "test" {
		Gin.Use(gin.LoggerWithWriter(io.Discard))
	}
	// 设置gin输出
	gin.DefaultWriter = io.Discard

	return Gin
}

func Run(g *gin.Engine) {
	var (
		config = config.Config
	)
	// 启动gin
	g.Run(fmt.Sprintf(":%d", config.Server.Port))
}
