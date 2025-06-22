package label

import (
	"github.com/infraboard/mcube/v2/ioc/config/validator"
	"github.com/infraboard/mcube/v2/tools/pretty"
	"github.com/infraboard/modules/iam/apps"
)

func NewLabel(spc *CreateLabelRequest) (*Label, error) {
	if err := spc.Validate(); err != nil {
		return nil, err
	}
	return &Label{
		ResourceMeta: *apps.NewResourceMeta(),
		Spec:         spc,
	}, nil
}

type Label struct {
	// 基础数据
	apps.ResourceMeta
	// 空间定义
	Spec *CreateLabelRequest `json:"spec" bson:",inline" gorm:"embedded"`
}

func (l *Label) TableName() string {
	return "labels"
}

func (l *Label) String() string {
	return pretty.ToJSON(l)
}

func NewCreateLabelRequest() *CreateLabelRequest {
	return &CreateLabelRequest{
		CreateCreateLabelSpec: CreateCreateLabelSpec{
			Resources:   []string{},
			EnumOptions: []*EnumOption{},
			Extras:      map[string]string{},
		},
	}
}

type CreateLabelRequest struct {
	// 创建人
	CreateBy string `json:"create_by" bson:"create_by" gorm:"column:create_by;type:varchar(255)"`
	// 标签的键
	Domain string `json:"domain" bson:"domain" gorm:"column:domain;type:varchar(100)"`
	// 标签的键
	Namespace string `json:"namespace" bson:"namespace" gorm:"column:namespace;type:varchar(100)"`
	// 用户参数
	CreateCreateLabelSpec
}

func (r *CreateLabelRequest) Validate() error {
	return validator.Validate(r)
}

func (r *CreateLabelRequest) AddEnumOption(enums ...*EnumOption) *CreateLabelRequest {
	r.EnumOptions = append(r.EnumOptions, enums...)
	return r
}

type CreateCreateLabelSpec struct {
	// 适用于那些资源
	Resources []string `json:"resources" bson:"resources" gorm:"column:resources;type:json;serializer:json;" description:"适用于那些资源" optional:"true"`

	// 标签的键, 标签的Key不允许修改, 带前缀的 tech.dev.frontend01 tech.dev.backend01
	Key string `json:"key" bson:"key" gorm:"column:key;type:varchar(255)" validate:"required"`
	// 标签的键的描述
	KeyDesc string `json:"key_desc" bson:"key_desc" gorm:"column:key_desc;type:varchar(255)" validate:"required"`
	// 标签的颜色
	Color string `json:"color" bson:"color" gorm:"column:color;type:varchar(100)"`

	// 值类型
	ValueType VALUE_TYPE `json:"value_type" gorm:"column:value_type;type:varchar(20)" bson:"value_type"`
	// 标签默认值
	DefaultValue string `json:"default_value" gorm:"column:default_value;type:text" bson:"default_value"`
	// 值描述
	ValueDesc string `json:"value_desc" gorm:"column:value_desc;type:text" bson:"value_desc"`
	// 是否是多选
	Multiple bool `json:"multiple" bson:"multiple" gorm:"column:multiple;tinyint(1)"`
	// 枚举值的选项
	EnumOptions []*EnumOption `json:"enum_options,omitempty" bson:"enum_options" gorm:"column:enum_options;type:json;serializer:json;"`
	// 基于Http枚举的配置
	HttpEnumConfig HttpEnumConfig `json:"http_enum_config,omitempty" gorm:"embedded" bson:"http_enum_config"`
	// 值的样例
	Example string `json:"example" bson:"example" gorm:"column:example;type:text"`

	// 扩展属性
	Extras map[string]string `json:"extras" bson:"extras" gorm:"column:extras;type:json;serializer:json;"`
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
	Url string `json:"url" bson:"url" gorm:"column:http_enum_config_url;type:text"`
	// Enum Label映射的字段名
	KeyFiled string `json:"enum_label_name" bson:"enum_label_name" gorm:"column:http_enum_config_key_filed;type:varchar(100)"`
	// Enum Value映射的字段名
	ValueFiled string `json:"enum_label_value" bson:"enum_label_value" gorm:"column:http_enum_config_value_filed;type:varchar(100)"`
}
