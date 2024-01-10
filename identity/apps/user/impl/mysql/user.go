package imp

import (
	"context"

	"github.com/infraboard/mcube/v2/exception"
	"github.com/infraboard/modules/identity/apps/user"
	"gorm.io/gorm"
)

// 创建用户
func (i *UserServiceImpl) CreateUser(
	ctx context.Context,
	req *user.CreateUserRequest) (
	*user.User, error) {
	// 1. 校验用户参数
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// 2. 生成一个User对象(ORM对象)
	ins := user.NewUser(req)

	if err := i.db.
		WithContext(ctx).
		Create(ins).
		Error; err != nil {
		return nil, err
	}

	// 4. 返回结果
	return ins, nil
}

// 删除用户
func (i *UserServiceImpl) DeleteUser(
	ctx context.Context,
	req *user.DeleteUserRequest) error {
	_, err := i.DescribeUserRequest(ctx,
		user.NewDescribeUserRequestById(req.IdString()))
	if err != nil {
		return err
	}

	return i.db.
		WithContext(ctx).
		Where("id = ?", req.Id).
		Delete(&user.User{}).
		Error
}

// 怎么查询一个用户
func (i *UserServiceImpl) DescribeUserRequest(
	ctx context.Context,
	req *user.DescribeUserRequest) (
	*user.User, error) {

	query := i.db.WithContext(ctx)

	// 1. 构造我们的查询条件
	switch req.DescribeBy {
	case user.DESCRIBE_BY_ID:
		query = query.Where("id = ?", req.DescribeValue)
	case user.DESCRIBE_BY_USERNAME:
		query = query.Where("username = ?", req.DescribeValue)
	}

	// SELECT * FROM `users` WHERE username = 'admin' ORDER BY `users`.`id` LIMIT 1
	ins := user.NewUser(user.NewCreateUserRequest())
	if err := query.First(ins).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, exception.NewNotFound("user %s not found", req.DescribeValue)
		}
		return nil, err
	}

	// 数据库里面存储的就是Hash
	ins.SetIsHashed()

	return ins, nil
}
