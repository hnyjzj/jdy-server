package scripts

import (
	"github.com/robfig/cron/v3"
)

func Cron() {
	// 创建定时任务
	c := cron.New()

	// 启动定时任务
	c.Start()
}
