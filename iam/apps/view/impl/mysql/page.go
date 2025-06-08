package mysql

import (
	"context"

	"github.com/infraboard/mcube/v2/exception"
	"github.com/infraboard/mcube/v2/ioc/config/datasource"
	"github.com/infraboard/mcube/v2/types"
	"github.com/infraboard/modules/iam/apps/view"
	"gorm.io/gorm"
)

// 添加页面
func (i *ViewServiceImpl) CreatePage(ctx context.Context, in *view.CreatePageRequest) (*view.Page, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	ins := view.NewPage()
	ins.CreatePageRequest = *in

	if err := datasource.DBFromCtx(ctx).
		Create(ins).
		Error; err != nil {
		return nil, err
	}
	return ins, nil
}

// 页面列表
func (i *ViewServiceImpl) QueryPage(ctx context.Context, in *view.QueryPageRequest) (*types.Set[*view.Page], error) {
	set := types.New[*view.Page]()

	query := datasource.DBFromCtx(ctx).Model(&view.Page{})
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

// 页面详情
func (i *ViewServiceImpl) DescribePage(ctx context.Context, in *view.DescribePageRequest) (*view.Page, error) {
	query := datasource.DBFromCtx(ctx)

	ins := &view.Page{}
	if err := query.Where("id = ?", in.Id).First(ins).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, exception.NewNotFound("page %d not found", in.Id)
		}
		return nil, err
	}

	return ins, nil
}

// 更新页面
func (i *ViewServiceImpl) UpdatePage(ctx context.Context, in *view.UpdatePageRequest) (*view.Page, error) {
	descReq := view.NewDeletePageRequest()
	descReq.SetId(in.Id)
	ins, err := i.DeletePage(ctx, descReq)
	if err != nil {
		return nil, err
	}

	ins.CreatePageRequest = in.CreatePageRequest
	return ins, datasource.DBFromCtx(ctx).Where("id = ?", in.Id).Updates(ins).Error
}

// 移除页面
func (i *ViewServiceImpl) DeletePage(ctx context.Context, in *view.DeletePageRequest) (*view.Page, error) {
	descReq := view.NewDescribePageRequest()
	descReq.SetId(in.Id)
	ins, err := i.DescribePage(ctx, descReq)
	if err != nil {
		return nil, err
	}

	return ins, datasource.DBFromCtx(ctx).
		Where("id = ?", in.Id).
		Delete(&view.Page{}).
		Error
}
