package mysql

import (
	"context"

	"github.com/infraboard/mcube/v2/exception"
	"github.com/infraboard/mcube/v2/ioc/config/datasource"
	"github.com/infraboard/mcube/v2/types"
	"github.com/infraboard/modules/iam/apps/view"
	"gorm.io/gorm"
)

// 创建菜单
func (i *ViewServiceImpl) CreateMenu(ctx context.Context, in *view.CreateMenuRequest) (*view.Menu, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	ins := view.NewMenu()
	ins.CreateMenuRequest = *in

	if err := datasource.DBFromCtx(ctx).
		Create(ins).
		Error; err != nil {
		return nil, err
	}
	return ins, nil
}

// 查询列表
func (i *ViewServiceImpl) QueryMenu(ctx context.Context, in *view.QueryMenuRequest) (*types.Set[*view.Menu], error) {
	set := types.New[*view.Menu]()

	query := datasource.DBFromCtx(ctx).Model(&view.Menu{})
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
func (i *ViewServiceImpl) DescribeMenu(ctx context.Context, in *view.DescribeMenuRequest) (*view.Menu, error) {
	query := datasource.DBFromCtx(ctx)

	ins := &view.Menu{}
	if err := query.Where("id = ?", in.Id).First(ins).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, exception.NewNotFound("menu %d not found", in.Id)
		}
		return nil, err
	}

	return ins, nil
}

// 更新菜单
func (i *ViewServiceImpl) UpdateMenu(ctx context.Context, in *view.UpdateMenuRequest) (*view.Menu, error) {
	descReq := view.NewDescribeMenuRequest()
	descReq.SetId(in.Id)
	ins, err := i.DescribeMenu(ctx, descReq)
	if err != nil {
		return nil, err
	}

	ins.CreateMenuRequest = in.CreateMenuRequest
	return ins, datasource.DBFromCtx(ctx).Where("id = ?", in.Id).Updates(ins).Error
}

// 删除菜单
func (i *ViewServiceImpl) DeleteMenu(ctx context.Context, in *view.DeleteMenuRequest) (*view.Menu, error) {
	descReq := view.NewDescribeMenuRequest()
	descReq.SetId(in.Id)
	ins, err := i.DescribeMenu(ctx, descReq)
	if err != nil {
		return nil, err
	}

	return ins, datasource.DBFromCtx(ctx).
		Where("id = ?", in.Id).
		Delete(&view.Menu{}).
		Error
}
