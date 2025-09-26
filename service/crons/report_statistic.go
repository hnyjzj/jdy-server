package crons

import "jdy/logic/statistic"

func init() {
	RegisterCrons(
		Crons{
			// // 每天晚上10 点半
			Spec: "0 30 22 * * *",
			Func: SendReportStatistic,
		},
	)
}

func SendReportStatistic() {
	logic := statistic.StatisticLogic{}
	logic.SendReportStatistic()
}
