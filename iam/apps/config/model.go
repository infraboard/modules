package config

import "github.com/infraboard/modules/iam/apps"

func NewConfigItem() *ConfigItem {
	return &ConfigItem{
		Meta: *apps.NewMeta(),
		KVItem: KVItem{
			Extras: make(map[string]string),
		},
	}
}

type ConfigItem struct {
	apps.Meta
	KVItem
}

func (c *ConfigItem) TableName() string {
	return "system_config"
}

type KVItem struct {
	// 配置所属组
	Group string `json:"group" bson:"group" validate:"required,lte=64" gorm:"column:group;type:varchar(200);index"`
	// 配置Key名称
	Key string `json:"key" bson:"key" validate:"required,lte=64" gorm:"column:key;type:varchar(200);index"`
	// 配置Key描述
	Desc string `json:"desc" bson:"desc" gorm:"column:desc;type:text"`
	// 配置Key的值
	Value string `json:"value" bson:"value" validate:"required" gorm:"column:key;type:text"`
	// 是否加密
	IsEncrypted bool `json:"is_encrypted" bson:"is_encrypted" validate:"required" gorm:"column:is_encrypted;type:tinyint(1)"`
	// 其他扩展信息
	Extras map[string]string `json:"extras" gorm:"column:extras;serializer:json;type:json" description:"其他扩展信息"`
}
