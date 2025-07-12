package cronjob

import (
	"context"

	"github.com/infraboard/mcube/v2/http/request"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/types"
)

const (
	APP_NAME = "cronjobs"
)

func GetService() Service {
	return ioc.Controller().Get(APP_NAME).(Service)
}

type Service interface {
	// 添加CronJob
	AddCronJob(context.Context, *CronJobSpec) (*CronJob, error)
	// 查询列表
	QueryCronJob(context.Context, *QueryCronJobRequest) (*types.Set[*CronJob], error)
	// 查询详情
	DescribeCronJob(context.Context, *DescribeCronJobRequest) (*CronJob, error)
	// 更新Cronjob
	UpdateCronJob(context.Context, *UpdateCronJobRequest) (*CronJob, error)
	// 删除Cronjob
	DeleteCronJob(context.Context, *DeleteCronJobRequest) (*CronJob, error)
}

func NewQueryCronJobRequest() *QueryCronJobRequest {
	return &QueryCronJobRequest{
		PageRequest: request.NewDefaultPageRequest(),
	}
}

type QueryCronJobRequest struct {
	// 分页参数
	*request.PageRequest
	// 任务名称
	Name string `json:"name"`
}

type UpdateCronJobRequest struct {
}

type DescribeCronJobRequest struct {
	// 任务Id
	Id string `json:"id" description:"Id"`
}

type DeleteCronJobRequest struct {
	DescribeCronJobRequest
}
