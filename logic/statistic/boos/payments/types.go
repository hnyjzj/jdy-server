package payments

import (
	"jdy/enums"
)

type TitleRes struct {
	Title     string `json:"title"`
	Key       string `json:"key"`
	Width     string `json:"width"`
	Fixed     string `json:"fixed"`
	ClassName string `json:"className"`
	Align     string `json:"align"`
}

type Where struct {
	Type      Types          `json:"type" label:"类型" find:"true" required:"true" sort:"1" type:"number" input:"radio" preset:"typeMap" binding:"required"`       // 类型
	Duration  enums.Duration `json:"duration" label:"时间范围" find:"true" required:"true" sort:"2" type:"number" input:"radio" preset:"typeMap" binding:"required"` // 时间范围
	StartTime string         `json:"startTime" label:"开始时间" find:"true" required:"false" sort:"3" type:"string" input:"date"`                                    // 开始时间
	EndTime   string         `json:"endTime" label:"结束时间" find:"true" required:"false" sort:"4" type:"string" input:"date"`                                      // 结束时间
}

type DataReq struct {
	Type      Types          `json:"type" binding:"required"`
	Duration  enums.Duration `json:"duration" binding:"required"`
	StartTime string         `json:"startTime"`
	EndTime   string         `json:"endTime"`
}
