package mysql

import (
	"context"

	"github.com/infraboard/mcube/v2/exception"
	"github.com/infraboard/modules/iam/apps/token"
	"github.com/infraboard/modules/iam/apps/user"
)

// 登录接口(颁发Token)
func (i *TokenServiceImpl) Login(
	ctx context.Context, req *token.IssueTokenRequest) (
	*token.Token, error) {
	provider := token.GetProvider(req.Provider)
	if provider == nil {
		return nil, exception.NewBadRequest("provider %s no support", req.Provider)
	}
	tk, err := provider.IssueToken(ctx, req.Parameter)
	if err != nil {
		return nil, err
	}

	if err := i.db.
		WithContext(ctx).
		Create(tk).
		Error; err != nil {
		return nil, err
	}

	// 补充用户信息, 只补充了用户的角色
	uDesc := user.NewDescribeUserRequestById(tk.UserId)
	_, err = i.user.DescribeUser(ctx, uDesc)
	if err != nil {
		return nil, err
	}

	// 避免同一个用户多次登录
	// 4. 颁发成功后  把之前的Token标记为失效,作业
	return tk, nil
}

// 校验Token 是给内部中间层使用 身份校验层
func (i *TokenServiceImpl) ValiateToken(
	ctx context.Context,
	req *token.ValiateToken) (*token.Token, error) {
	// 1. 查询Token (是不是我们这个系统颁发的)
	tk := token.NewToken()
	err := i.db.
		WithContext(ctx).
		Where("access_token = ?", req.AccessToken).
		First(tk).
		Error
	if err != nil {
		return nil, err
	}

	// 2. 判断Token的合法性:
	// 2.1 判断Ak是否过期
	if err := tk.IsExpired(); err != nil {
		return nil, err
	}
	return tk, nil
}

// 退出接口(销毁Token)
func (i *TokenServiceImpl) Logout(
	ctx context.Context,
	req *token.RevolkTokenRequest) (*token.Token, error) {
	// 1. 查询Token (是不是我们这个系统颁发的)
	tk := token.NewToken()
	err := i.db.WithContext(ctx).
		Where("access_token = ?", req.AccessToken).
		First(tk).
		Error
	if err != nil {
		return nil, err
	}

	err = i.db.WithContext(ctx).
		Where("access_token = ?", req.AccessToken).
		Where("refresh_token = ?", req.RefreshToken).
		Delete(&token.Token{}).
		Error
	return tk, err
}
