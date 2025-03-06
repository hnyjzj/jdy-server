package main

import (
	"jdy/config"
	"jdy/model"
	"jdy/router"
	"jdy/service"
	"jdy/service/gin"
	"jdy/service/redis"
	"log"
)

func main() {
	// 设置日志格式
	log.SetFlags(log.LstdFlags | log.Llongfile)
	// 初始化gin
	g := gin.Init()
	// 路由初始化
	router.Init(g)
	// 启动服务
	go service.Start()
	// 启动 gin http 服务
	gin.Run(g)
}

func init() {
	// 配置初始化
	config.Init()
	// 数据库初始化
	model.Init()
	// 初始化redis
	redis.Init()
}
