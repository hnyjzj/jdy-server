package types

type WhereForm struct {
	Name     string `json:"name"`     // 字段名
	Label    string `json:"label"`    // 名称
	Sort     int    `json:"sort"`     // 排序
	Type     string `json:"type"`     // 值类型： string, number, boolean, object, string[], number[], boolean[]
	Input    string `json:"input"`    // 输入框类型：text, number, password, textarea, email, url, tel, search, range, color, date, datetime, time, week, month, quarter, year, file, image, video, audio, editor
	Required bool   `json:"required"` // 是否必填
	Show     bool   `json:"show"`     // 是否显示
	Preset   any    `json:"preset"`   // 预设：[value1, value2]|{value1: label1, value2: label2}|{value1: {label1: label2}}
}
