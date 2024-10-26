package mysql

import (
	"context"

	"github.com/infraboard/mcube/v2/types"
	"github.com/infraboard/modules/iam/apps/view"
)

// 页面列表
func (i *MenuServiceImpl) QueryPage(ctx context.Context, in *view.QueryPageRequest) (*types.Set[*view.Page], error) {
	return nil, nil
}

// 页面详情
func (i *MenuServiceImpl) DescribePage(ctx context.Context, in *view.DescribePageRequest) (*view.Page, error) {
	return nil, nil
}

// 添加页面
func (i *MenuServiceImpl) CreatePage(ctx context.Context, in *view.CreatePageRequest) (*view.Page, error) {
	return nil, nil
}

// 移除页面
func (i *MenuServiceImpl) DeletePage(ctx context.Context, in *view.DeletePageRequest) (*view.Page, error) {
	return nil, nil
}

// 更新页面
func (i *MenuServiceImpl) UpdatePage(ctx context.Context, in *view.UpdatePageRequest) (*view.Page, error) {
	return nil, nil
}
