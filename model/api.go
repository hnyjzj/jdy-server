package model

type Api struct {
	SoftDelete

	Title  string `json:"title" gorm:"column:title;size:255;comment:标题"`                            // 标题
	Path   string `json:"path" gorm:"uniqueIndex:unique_api;column:path;size:255;comment:接口地址"`     // 接口地址
	Method string `json:"method" gorm:"uniqueIndex:unique_api;column:method;size:255;comment:请求方式"` // 请求方式
}

func init() {
	// 注册模型
	RegisterModels(
		&Api{},
	)
	// 重置表
	RegisterRefreshModels(
	// &Api{},
	)
}
