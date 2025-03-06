package service

import "jdy/service/crons"

func Start() {
	// 启动定时任务
	crons.Cron()
}
