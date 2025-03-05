package scripts

import (
	"log"
	"os"
	"time"

	"github.com/robfig/cron/v3"
)

func (c *CronScript) Funcs() {
	// 每天 9 点执行
	_, _ = c.AddFunc("0 0 9 * * *", func() {
		go SendGoldPriceSetMessage()
	})
}

type CronScript struct {
	*cron.Cron
}

func Cron() {
	// 创建定时服务
	c := CronScript{cron.New(
		cron.WithLogger(cron.VerbosePrintfLogger(log.New(os.Stdout, "cron: ", log.LstdFlags))),
		cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger), cron.Recover(cron.DefaultLogger)),
		cron.WithLocation(time.FixedZone("CST", 8*3600)),
		cron.WithSeconds(),
	)}
	defer c.Stop()

	// 添加定时任务
	c.Funcs()

	// 启动定时任务
	c.Start()

	// 阻塞主线程
	select {}
}
