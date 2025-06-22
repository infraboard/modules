package label

// 值类型
type VALUE_TYPE string

const (
	// 文本
	VALUE_TYPE_TEXT VALUE_TYPE = "text"
	// 布尔值, 只能是ture或者false
	VALUE_TYPE_BOOLEAN VALUE_TYPE = "bool"
	// 枚举
	VALUE_TYPE_ENUM VALUE_TYPE = "enum"
	// 基于url的远程选项拉去, 仅存储URL地址, 前端自己处理
	VALUE_TYPE_HTTP_ENUM VALUE_TYPE = "http_enum"
)
