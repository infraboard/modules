package config

import (
	"context"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/validator"
	"github.com/infraboard/mcube/v2/types"
)

const (
	AppName = "system_config"
)

func GetService() Service {
	return ioc.Controller().Get(AppName).(Service)
}

type Service interface {
	// 添加配置
	AddConfig(context.Context, *AddConfigRequest) (*types.Set[*ConfigItem], error)
	// 查询配置项
	QueryConfig(context.Context, *QueryConfigRequest) (*types.Set[*ConfigItem], error)
	// 查询配置详情
	DescribeConfig(context.Context, *DescribeConfigRequest) (*ConfigItem, error)
	// 更新配置
	UpdateConfig(context.Context, *UpdateConfigRequest) (*ConfigItem, error)
}

func NewAddConfigRequest() *AddConfigRequest {
	return &AddConfigRequest{
		Items: []*KVItem{},
	}
}

type AddConfigRequest struct {
	Items []*KVItem `json:"items"`
}

func (r *AddConfigRequest) Validate() error {
	return validator.Validate(r)
}

func (r *AddConfigRequest) AddKVItem(items ...*KVItem) *AddConfigRequest {
	r.Items = append(r.Items, items...)
	return r
}

func NewQueryConfigRequest() *QueryConfigRequest {
	return &QueryConfigRequest{}
}

type QueryConfigRequest struct {
	Group string `json:"group"`
}

func NewDescribeConfigRequestById(id string) *DescribeConfigRequest {
	return &DescribeConfigRequest{
		DescribeBy:    DESCRIBE_BY_ID,
		DescribeValue: id,
	}
}

func NewDescribeConfigRequestByKey(key string) *DescribeConfigRequest {
	return &DescribeConfigRequest{
		DescribeBy:    DESCRIBE_BY_KEY,
		DescribeValue: key,
	}
}

type DescribeConfigRequest struct {
	DescribeBy    DESCRIBE_BY
	DescribeValue string
}

func NewUpdateConfigRequest(id string) *UpdateConfigRequest {
	return &UpdateConfigRequest{
		Id: id,
	}
}

type UpdateConfigRequest struct {
	Id string `json:"id"`
	KVItem
}
