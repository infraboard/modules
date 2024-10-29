package mysql

import (
	"context"

	"github.com/infraboard/mcube/v2/exception"
	"github.com/infraboard/mcube/v2/ioc/config/datasource"
	"github.com/infraboard/mcube/v2/types"
	"github.com/infraboard/modules/iam/apps/policy"
	"github.com/infraboard/modules/iam/apps/view"
	"gorm.io/gorm"
)

// 创建策略
func (i *PolicyServiceImpl) CreatePolicy(ctx context.Context, in *policy.CreatePolicyRequest) (*policy.Policy, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	ins := policy.NewPolicy()
	ins.CreatePolicyRequest = *in

	if err := datasource.DBFromCtx(ctx).
		Create(ins).
		Error; err != nil {
		return nil, err
	}
	return ins, nil
}

// 查询策略列表
func (i *PolicyServiceImpl) QueryPolicy(ctx context.Context, in *policy.QueryPolicyRequest) (*types.Set[*policy.Policy], error) {
	set := types.New[*policy.Policy]()

	query := datasource.DBFromCtx(ctx).Model(&policy.Policy{})
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

// 查询详情
func (i *PolicyServiceImpl) DescribePolicy(ctx context.Context, in *policy.DescribePolicyRequest) (*policy.Policy, error) {
	query := datasource.DBFromCtx(ctx)

	ins := &policy.Policy{}
	if err := query.First(ins).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, exception.NewNotFound("policy %d not found", in.Id)
		}
		return nil, err
	}

	return ins, nil
}

// 更新策略
func (i *PolicyServiceImpl) UpdatePolicy(ctx context.Context, in *policy.UpdatePolicyRequest) (*policy.Policy, error) {
	descReq := policy.NewDescribePolicyRequest()
	descReq.SetId(in.Id)
	ins, err := i.DescribePolicy(ctx, descReq)
	if err != nil {
		return nil, err
	}

	ins.CreatePolicyRequest = in.CreatePolicyRequest
	return ins, datasource.DBFromCtx(ctx).Where("id = ?", in.Id).Updates(ins).Error
}

// 删除策略
func (i *PolicyServiceImpl) DeletePolicy(ctx context.Context, in *policy.DeletePolicyRequest) (*policy.Policy, error) {
	descReq := policy.NewDescribePolicyRequest()
	descReq.SetId(in.Id)
	ins, err := i.DescribePolicy(ctx, descReq)
	if err != nil {
		return nil, err
	}

	return ins, datasource.DBFromCtx(ctx).
		Where("id = ?", in.Id).
		Delete(&view.Menu{}).
		Error
}
