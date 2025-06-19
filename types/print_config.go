package types

import "jdy/enums"

type PrintConfig struct {
	Size  PrintSize     `json:"size"`  // 打印纸张大小
	Base  PrintPosition `json:"base"`  // 基础
	Phone PrintPosition `json:"phone"` // 手机号
	List  PrintPosition `json:"list"`  // 列表
	Total PrintPosition `json:"total"` // 总价
	More  PrintPosition `json:"more"`  // 更多
}

func (PrintConfig) Default() PrintConfig {
	def := PrintConfig{
		Size: PrintSize{
			Width:    200,
			Height:   140,
			FontSize: 3,
		},
		Base: PrintPosition{
			Width:  30,
			Height: 15,
			Top:    0,
			Bottom: 0,
			Left:   0,
			Right:  6,
		},
		Phone: PrintPosition{
			Width:  0,
			Height: 0,
			Top:    42,
			Bottom: 0,
			Left:   0,
			Right:  8,
		},
		List: PrintPosition{
			Width:  0,
			Height: 50,
			Top:    56,
			Bottom: 0,
			Left:   8,
			Right:  7,
		},
		Total: PrintPosition{
			Width:  0,
			Height: 0,
			Top:    108,
			Bottom: 0,
			Left:   0,
			Right:  10,
		},
		More: PrintPosition{
			Width:  0,
			Height: 0,
			Top:    0,
			Bottom: 10,
			Left:   20,
			Right:  10,
		},
	}

	return def
}

type PrintSize struct {
	Width    int `json:"width"`    // 宽度
	Height   int `json:"height"`   // 高度
	FontSize int `json:"fontSize"` // 字体大小
}

type PrintPosition struct {
	Width  int `json:"width"`  // 宽度
	Height int `json:"height"` // 高度
	Top    int `json:"top"`    // 顶部
	Bottom int `json:"bottom"` // 底部
	Left   int `json:"left"`   // 左侧
	Right  int `json:"right"`  // 右侧
}

type PrintReq struct {
	StoreId string          `json:"store_id" binding:"required"`
	Name    string          `json:"name" binding:"required"`
	Type    enums.PrintType `json:"type" binding:"required"`
	Config  PrintConfig     `json:"config" binding:"required"`
}

type PrintWhere struct {
	Id      string          `json:"id"`
	StoreId string          `json:"store_id" binding:"required"`
	Name    string          `json:"name"`
	Type    enums.PrintType `json:"type"`
}

type PrintListReq struct {
	PageReq
	Where PrintWhere `json:"where"`
}

type PrintInfoReq struct {
	Id      string          `json:"id"`
	Type    enums.PrintType `json:"type"`
	StoreId string          `json:"store_id"`
}

type PrintUpdateReq struct {
	Id string `json:"id" binding:"required"`

	Data PrintReq `json:"data" binding:"required"`
}

type PrintDeleteReq struct {
	Id string `json:"id" binding:"required"`
}

type PrintCopyReq struct {
	Id      string `json:"id" binding:"required"`
	StoreId string `json:"store_id" binding:"required"`
	Name    string `json:"name"`
}
