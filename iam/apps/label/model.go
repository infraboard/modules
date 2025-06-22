package label

import "github.com/infraboard/mcube/v2/pb/resource"

type Label struct {
	// 元信息
	Meta *resource.Meta `json:"meta" bson:",inline"`
	// 空间定义
	Spec *CreateLabelRequest `json:"spec" bson:",inline"`
}

type CreateLabelRequest struct {
	// 标签的键, 标签的Key不允许修改
	Key string `json:"key" bson:"key" validate:"required"`
	// 标签的键的描述
	KeyDesc string `json:"key_desc" bson:"key_desc" validate:"required"`
	// 标签的键
	Domain string `json:"domain" bson:"domain"`
	// 标签的键
	Namespace string `json:"namespace" bson:"namespace"`
	// 是不是必须标签, 如果是必须标签 资源创建时必须添加该标签
	Required bool `json:"required" bson:"required"`
	// 什么条件下必选
	RequiredCondition *RequiredCondition `json:"required_condition" bson:"required_condition"`
	// 标签的颜色
	Color string `json:"color" bson:"color"`
	// 值类型
	ValueType VALUE_TYPE `json:"value_type" bson:"value_type"`
	// 标签默认值
	DefaultValue string `json:"default_value" bson:"default_value"`
	// 值描述
	ValueDesc string `json:"value_desc" bson:"value_desc"`
	// 是否是多选
	Multiple bool `json:"multiple" bson:"multiple"`
	// 枚举值的选项
	EnumOptions []*EnumOption `json:"enum_options,omitempty" bson:"enum_options"`
	// 基于Http枚举的配置
	HttpEnumConfig *HttpEnumConfig `json:"http_enum_config,omitempty" bson:"http_enum_config"`
	// 值的样例
	Example string `json:"example" bson:"example"`
	// 创建人
	CreateBy string `json:"create_by" bson:"create_by"`
	// 角色可见性
	Visiable resource.VISIABLE `json:"visiable" bson:"visiable"`
	// 扩展属性
	Extensions map[string]string `json:"extensions" bson:"extensions"`
}

type RequiredCondition struct {
	// 针对特定资源的必选, 默认针对所有资源
	Resources []string `json:"resources" bson:"resources"`
}

type EnumOption struct {
	// 选项的说明
	Label string `json:"label" bson:"label"`
	// 用户输入
	Input string `json:"input" bson:"input" validate:"required"`
	// 选项的值, 根据parent.input + children.input 自动生成
	Value string `json:"value" bson:"value"`
	// 标签的颜色
	Color string `json:"color" bson:"color"`
	// 是否废弃
	Deprecate bool `json:"deprecate" bson:"deprecate"`
	// 废弃说明
	DeprecateDesc string `json:"deprecate_desc" bson:"deprecate_desc"`
	// 枚举的子选项
	Children []*EnumOption `json:"children,omitempty" bson:"children"`
	// 扩展属性
	Extensions map[string]string `json:"extensions" bson:"extensions"`
}

type HttpEnumConfig struct {
	// 基于枚举的URL, 注意只支持Get方法
	Url string `json:"url" bson:"url"`
	// Enum Label映射的字段名
	EnumLabelName string `json:"enum_label_name" bson:"enum_label_name"`
	// Enum Value映射的字段名
	EnumLabelValue string `json:"enum_label_value" bson:"enum_label_value"`
}
