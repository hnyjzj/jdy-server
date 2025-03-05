package service

import (
	"jdy/scripts"
)

func Start() {
	// 启动定时任务
	scripts.Cron()
}
