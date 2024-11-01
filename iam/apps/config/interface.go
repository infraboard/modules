package config

import (
	"context"

	"github.com/infraboard/mcube/v2/ioc"
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
	AddConfig(context.Context, *KVItem) (*ConfigItem, error)
	// 查询配置项
	QueryConfig(context.Context, *QueryConfigRequest) (*types.Set[*ConfigItem], error)
	// 查询配置详情
	DescribeConfig(context.Context, *DescribeConfigRequest) (*ConfigItem, error)
	// 更新配置
	UpdateConfig(context.Context, *UpdateConfigRequest) (*ConfigItem, error)
}

func NewQueryConfigRequest() *QueryConfigRequest {
	return &QueryConfigRequest{}
}

type QueryConfigRequest struct {
	Group string `json:"group"`
}

type DescribeConfigRequest struct {
	Group string `json:"group"`
	Key   string `json:"key"`
}

type UpdateConfigRequest struct {
}
