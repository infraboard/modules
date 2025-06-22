package impl

import (
	"context"

	"github.com/infraboard/mcube/v2/ioc/config/datasource"
	"github.com/infraboard/mcube/v2/types"
	"github.com/infraboard/modules/iam/apps/label"
)

// CreateLabel implements label.Service.
func (i *LabelServiceImpl) CreateLabel(ctx context.Context, in *label.CreateLabelRequest) (*label.Label, error) {
	ins, err := label.NewLabel(in)
	if err != nil {
		return nil, err
	}

	if err := datasource.DBFromCtx(ctx).
		Create(ins).
		Error; err != nil {
		return nil, err
	}
	return ins, nil
}

// QueryLabel implements label.Service.
func (i *LabelServiceImpl) QueryLabel(ctx context.Context, in *label.QueryLabelRequest) (*types.Set[*label.Label], error) {
	set := types.New[*label.Label]()

	query := datasource.DBFromCtx(ctx).Model(&label.Label{})
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

// DeleteLabel implements label.Service.
func (i *LabelServiceImpl) DeleteLabel(ctx context.Context, in *label.DeleteLabelRequest) (*label.Label, error) {
	panic("unimplemented")
}

// DescribeLabel implements label.Service.
func (i *LabelServiceImpl) DescribeLabel(ctx context.Context, in *label.DescribeLabelRequest) (*label.Label, error) {
	panic("unimplemented")
}

// UpdateLabel implements label.Service.
func (i *LabelServiceImpl) UpdateLabel(ctx context.Context, in *label.UpdateLabelRequest) (*label.Label, error) {
	panic("unimplemented")
}
