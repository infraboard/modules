package config

import (
	"encoding/json"

	"github.com/infraboard/mcube/v2/crypto/cbc"
	"github.com/infraboard/mcube/v2/ioc/config/application"
	"github.com/infraboard/mcube/v2/tools/pretty"
	"github.com/infraboard/modules/iam/apps"
)

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

func (c *ConfigItem) String() string {
	return pretty.ToJSON(c)
}

func (c *ConfigItem) Encrypt() error {
	if c.IsEncrypted {
		cihperText, err := cbc.EncryptToString(c.Value, []byte(application.Get().EncryptKey))
		if err != nil {
			return err
		}
		c.Value = cihperText
	}
	return nil
}

func (c *ConfigItem) Decrypt() error {
	if c.IsEncrypted {
		plainText, err := cbc.DecryptFromString(c.Value, []byte(application.Get().EncryptKey))
		if err != nil {
			return err
		}
		c.Value = plainText
	}
	return nil
}

func (c *ConfigItem) Load(v any) error {
	return json.Unmarshal([]byte(c.Value), v)
}

func NewKVItem(key, value string) *KVItem {
	return &KVItem{
		Key:    key,
		Value:  value,
		Extras: map[string]string{},
	}
}

type KVItem struct {
	// 配置所属组
	Group string `json:"group" bson:"group" validate:"required,lte=64" gorm:"column:group;type:varchar(200);index"`
	// 配置Key名称
	Key string `json:"key" bson:"key" validate:"required,lte=64" gorm:"column:key;type:varchar(200);index"`
	// 配置Key描述
	Desc string `json:"desc" bson:"desc" gorm:"column:desc;type:text"`
	// 配置Key的值
	Value string `json:"value" bson:"value" validate:"required" gorm:"column:value;type:text"`
	// 是否加密
	IsEncrypted bool `json:"is_encrypted" bson:"is_encrypted" validate:"required" gorm:"column:is_encrypted;type:tinyint(1)"`
	// 其他扩展信息
	Extras map[string]string `json:"extras" gorm:"column:extras;serializer:json;type:json" description:"其他扩展信息"`
}

func (i *KVItem) SetGroup(group string) *KVItem {
	i.Group = group
	return i
}

func (i *KVItem) SetDesc(desc string) *KVItem {
	i.Desc = desc
	return i
}

func (i *KVItem) SetIsEncrypted(isEncrypted bool) *KVItem {
	i.IsEncrypted = isEncrypted
	return i
}
