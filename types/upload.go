package types

// 模块
type UploadModel string

func (s UploadModel) String() string {
	return string(s)
}

const (
	// 头像
	UploadModelAvatar UploadModel = "avatar"
	// 工作台
	UploadModelWorkbench UploadModel = "workbench"
	// 门店
	UploadModelStore UploadModel = "store"
)

// 类型
type UploadType string

func (s UploadType) String() string {
	return string(s)
}

const (
	UploadTypeImage UploadType = "image"
)
