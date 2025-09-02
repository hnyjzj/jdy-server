package boos

import "jdy/enums"

type Where struct {
	Duration  enums.Duration `json:"duration" label:"时间范围" find:"true" required:"true" sort:"2" type:"number" input:"radio" preset:"typeMap" binding:"required"` // 时间范围
	StartTime string         `json:"startTime" label:"开始时间" find:"true" required:"false" sort:"3" type:"string" input:"date"`                                    // 开始时间
	EndTime   string         `json:"endTime" label:"结束时间" find:"true" required:"false" sort:"4" type:"string" input:"date"`                                      // 结束时间
}

type DataReq struct {
	Duration  enums.Duration `json:"duration" binding:"required"`
	StartTime string         `json:"startTime"`
	EndTime   string         `json:"endTime"`
}
