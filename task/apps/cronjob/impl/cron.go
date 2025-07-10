package impl

import (
	"context"

	"github.com/infraboard/mcube/v2/types"
	"github.com/infraboard/modules/task/apps/cronjob"
)

// AddCronJob implements cronjob.Service.
func (i *CronServiceImpl) AddCronJob(ctx context.Context, in *cronjob.CronJobSpec) (*cronjob.CronJob, error) {
	panic("unimplemented")
}

// DeleteCronJob implements cronjob.Service.
func (i *CronServiceImpl) DeleteCronJob(ctx context.Context, in *cronjob.DeleteCronJobRequest) (*cronjob.CronJob, error) {
	panic("unimplemented")
}

// DescribeCronJob implements cronjob.Service.
func (i *CronServiceImpl) DescribeCronJob(ctx context.Context, in *cronjob.DescribeCronJobRequest) (*cronjob.CronJob, error) {
	panic("unimplemented")
}

// QueryCronJob implements cronjob.Service.
func (i *CronServiceImpl) QueryCronJob(ctx context.Context, in *cronjob.QueryCronJobRequest) (*types.Set[*cronjob.CronJob], error) {
	panic("unimplemented")
}

// UpdateCronJob implements cronjob.Service.
func (i *CronServiceImpl) UpdateCronJob(ctx context.Context, in *cronjob.UpdateCronJobRequest) (*cronjob.CronJob, error) {
	panic("unimplemented")
}
