package impl

import (
	"context"

	"github.com/infraboard/mcube/v2/exception"
	"github.com/infraboard/mcube/v2/ioc/config/datasource"
	"github.com/infraboard/mcube/v2/ioc/config/mcron"
	"github.com/infraboard/mcube/v2/types"
	"github.com/infraboard/modules/task/apps/cronjob"
	"github.com/infraboard/modules/task/apps/task"
	"gorm.io/gorm"
)

// AddCronJob implements cronjob.Service.
func (i *CronJobServiceImpl) AddCronJob(ctx context.Context, in *cronjob.CronJobSpec) (*cronjob.CronJob, error) {
	ins := cronjob.NewCronJob(*in)

	// 才真正的创建cron 关联cron实例
	if in.GetEnabled() {
		refId, err := mcron.RunAndAddFunc(in.Cron, func() {
			task.GetService().Run(context.Background(), &ins.TaskSpec)
		})
		if err != nil {
			return nil, err
		}
		ins.RefInstanceId = int(refId)
		ins.Node = i.hostname
	}

	if err := datasource.DBFromCtx(ctx).Save(ins).Error; err != nil {
		return nil, err
	}

	return ins, nil
}

// QueryCronJob implements cronjob.Service.
func (i *CronJobServiceImpl) QueryCronJob(ctx context.Context, in *cronjob.QueryCronJobRequest) (*types.Set[*cronjob.CronJob], error) {
	set := types.New[*cronjob.CronJob]()

	query := datasource.DBFromCtx(ctx).Model(&cronjob.CronJob{})
	if in.Name != "" {
		query = query.Where("name = ?", in.Name)
	}

	err := query.Count(&set.Total).Error
	if err != nil {
		return nil, err
	}

	err = query.
		Order("created_at desc").
		Offset(int(in.ComputeOffset())).
		Limit(int(in.PageSize)).
		Find(&set.Items).
		Error
	if err != nil {
		return nil, err
	}
	return set, nil
}

// DescribeCronJob implements cronjob.Service.
func (i *CronJobServiceImpl) DescribeCronJob(ctx context.Context, in *cronjob.DescribeCronJobRequest) (*cronjob.CronJob, error) {
	query := datasource.DBFromCtx(ctx)

	ins := &cronjob.CronJob{}
	if err := query.Where("id = ?", in.Id).First(ins).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, exception.NewNotFound("cronjob %s not found", in.Id)
		}
		return nil, err
	}
	return ins, nil
}

// DeleteCronJob implements cronjob.Service.
func (i *CronJobServiceImpl) DeleteCronJob(ctx context.Context, in *cronjob.DeleteCronJobRequest) (*cronjob.CronJob, error) {
	panic("unimplemented")
}

// UpdateCronJob implements cronjob.Service.
func (i *CronJobServiceImpl) UpdateCronJob(ctx context.Context, in *cronjob.UpdateCronJobRequest) (*cronjob.CronJob, error) {
	panic("unimplemented")
}
