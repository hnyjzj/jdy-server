package scripts

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/robfig/cron/v3"
)

func (c *CronScript) Funcs() {
	// 添加定时任务
	_, _ = c.AddFunc("@every 1h", func() {
		fmt.Println("@every 1h 定时任务执行")
	})
}

type CronScript struct {
	*cron.Cron
}

func Cron(ctx context.Context) error {
	// 创建定时任务
	c := CronScript{cron.New(cron.WithLogger(cron.VerbosePrintfLogger(log.New(os.Stdout, "cron: ", log.LstdFlags))))}

	// 添加定时任务
	c.Funcs()

	// 启动定时任务
	c.Start()

	// 等待上下文取消
	<-ctx.Done()
	c.Stop()

	return nil
}
