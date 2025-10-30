package types

type WhereForm struct {
	Name      string           `json:"name"`      // 字段名
	Label     string           `json:"label"`     // 名称
	Sort      int              `json:"sort"`      // 排序
	Type      string           `json:"type"`      // 值类型： string, number, boolean, object, string[], number[], boolean[]
	Input     string           `json:"input"`     // 输入框类型：text, number, password, textarea, email, url, tel, search, range, color, date, datetime, time, week, month, quarter, year, file, image, video, audio, editor
	Required  bool             `json:"required"`  // 是否必填
	Find      bool             `json:"find"`      // 查询是否显示
	Create    bool             `json:"create"`    // 创建是否显示
	Update    bool             `json:"update"`    // 更新是否显示
	List      bool             `json:"list"`      // 列表是否显示
	Info      bool             `json:"info"`      // 详情是否显示
	Preset    any              `json:"preset"`    // 预设：[value1, value2]|{value1: label1, value2: label2}|{value1: {label1: label2}}
	Condition []WhereCondition `json:"condition"` // 条件：[{key: string, value: string, operator: string}]
}

type WhereCondition struct {
	Key      string `json:"key"`      // 字段名
	Value    any    `json:"value"`    // 值
	Operator string `json:"operator"` // 操作符：=, !=, >, >=, <, <=, in, not in, like, not like, is null, is not null
}
