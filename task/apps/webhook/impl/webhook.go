package impl

import (
	"context"

	"github.com/infraboard/mcube/v2/exception"
	"github.com/infraboard/mcube/v2/ioc/config/bus"
	"github.com/infraboard/mcube/v2/ioc/config/datasource"
	"github.com/infraboard/mcube/v2/types"
	"github.com/infraboard/modules/task/apps/webhook"
	"gorm.io/gorm"
)

// Run implements webhook.Service.
func (i *WebHookServiceImpl) Run(ctx context.Context, in *webhook.WebHookSpec) (*webhook.WebHook, error) {
	hook := webhook.NewWebHook(*in)
	// 保存Hook执行记录
	err := datasource.DBFromCtx(ctx).Save(hook).Error
	if err != nil {
		return nil, err
	}

	// 发送到队列, 只发送HookId
	err = bus.GetService().Publish(ctx, hook.ToBusEvent(i.WebhookQueue))
	if err != nil {
		return nil, err
	}
	return hook, nil
}

// 查询WebHook具体执行状态
func (i *WebHookServiceImpl) DescribeWebHook(ctx context.Context, in *webhook.DescribeWebHookRequest) (*webhook.WebHook, error) {
	query := datasource.DBFromCtx(ctx)

	ins := &webhook.WebHook{}
	if err := query.Where("id = ?", in.Id).First(ins).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, exception.NewNotFound("task %s not found", in.Id)
		}
		return nil, err
	}

	return ins, nil
}

// QueryWebHook implements webhook.Service.
func (i *WebHookServiceImpl) QueryWebHook(ctx context.Context, in *webhook.QueryWebHookRequest) (*types.Set[*webhook.WebHook], error) {
	set := types.NewSet[*webhook.WebHook]()

	query := datasource.DBFromCtx(ctx).Model(&webhook.WebHook{})
	if in.RefTaskId != "" {
		query = query.Where("ref_task_id = ?", in.RefTaskId)
	}

	err := query.Count(&set.Total).Error
	if err != nil {
		return nil, err
	}
	err = query.
		Order("created_at desc").
		Find(&set.Items).
		Error
	if err != nil {
		return nil, err
	}

	return set, nil
}
