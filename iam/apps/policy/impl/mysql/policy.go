package mysql

import (
	"context"

	"github.com/infraboard/mcube/v2/types"
	"github.com/infraboard/modules/iam/apps/policy"
)

// 创建策略
func (i *PolicyServiceImpl) CreatePolicy(ctx context.Context, in *policy.CreatePolicyRequest) (*policy.Policy, error) {
	return nil, nil
}

// 查询策略列表
func (i *PolicyServiceImpl) QueryPolicy(ctx context.Context, in *policy.QueryPolicyRequest) (*types.Set[*policy.Policy], error) {
	return nil, nil
}

// 查询详情
func (i *PolicyServiceImpl) DescribePolicy(ctx context.Context, in *policy.DescribePolicyRequest) (*policy.Policy, error) {
	return nil, nil
}

// 更新策略
func (i *PolicyServiceImpl) UpdatePolicy(ctx context.Context, in *policy.UpdatePolicyRequest) (*policy.Policy, error) {
	return nil, nil
}

// 删除策略
func (i *PolicyServiceImpl) DeletePolicy(ctx context.Context, in *policy.DeletePolicyRequest) (*types.Set[*policy.Policy], error) {
	return nil, nil
}
