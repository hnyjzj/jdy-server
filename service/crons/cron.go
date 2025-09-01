package crons

import (
	"log"
	"os"
	"time"

	"github.com/robfig/cron/v3"
)

func (c *CronScript) Funcs() {
	for _, v := range CronsList {
		cp := v
		_, err := c.AddFunc(cp.Spec, func() {
			go cp.Func()
		})
		if err != nil {
			log.Printf("定时任务添加失败: %v", err)
		}
	}
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

type Crons struct {
	Spec string
	Func func()
}

var CronsList []Crons = []Crons{}

// 注册定时任务
func RegisterCrons(cron ...Crons) {
	CronsList = append(CronsList, cron...)
}
