package service

import "jdy/service/redis"

func Start() {
	// 初始化redis
	redis.Init()
}
