package finished_stock

import (
	"fmt"
	"jdy/enums"
)

func (l *Logic) GetTitles() []TitleRes {
	var titles []TitleRes
	titles = append(titles, TitleRes{
		Title:     "门店",
		Key:       "name",
		Width:     "100px",
		Fixed:     "left",
		ClassName: "age",
		Align:     "center",
	})
	titles = append(titles, TitleRes{
		Title: "总",
		Key:   "total",
		Width: "100px",
		Fixed: "left",
		Align: "center",
	})

	for k, v := range enums.ProductClassFinishedMap {
		titles = append(titles, TitleRes{
			Title: v,
			Key:   fmt.Sprint(k),
			Width: "100px",
			Align: "center",
		})
	}

	return titles
}
