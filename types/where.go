package types

type WhereForm struct {
	Label string `json:"label"` // 名称
	Type  string `json:"type"`  // text, number, select, date, datetime, time, radio, checkbox, switch, slider, cascader, tree, upload, editor

	Required bool `json:"required"` // 是否必填
	Preset   any  `json:"preset"`   // 预设：[value1, value2]|{value1: label1, value2: label2}|{value1: {label1: label2}}
}
