package mysql

import (
	"context"
	"time"

	"github.com/infraboard/mcube/v2/exception"
	"github.com/infraboard/mcube/v2/ioc/config/datasource"
	"github.com/infraboard/mcube/v2/types"
	"github.com/infraboard/modules/iam/apps/token"
	"github.com/infraboard/modules/iam/apps/user"
)

// 登录接口(颁发Token)
func (i *TokenServiceImpl) IssueToken(ctx context.Context, req *token.IssueTokenRequest) (*token.Token, error) {
	// 避免同一个用户多次登录
	// 颁发成功后  把之前的Token标记为失效,作业

	issuer := token.GetIssue(req.Issuer)
	if issuer == nil {
		return nil, exception.NewBadRequest("provider %s no support", req.Issuer)
	}
	tk, err := issuer.IssueToken(ctx, req.Parameter)
	if err != nil {
		return nil, err
	}

	tk.SetIssuer(req.Issuer)
	if err := datasource.DBFromCtx(ctx).
		Create(tk).
		Error; err != nil {
		return nil, err
	}

	// 补充用户信息, 只补充了用户的角色
	uDesc := user.NewDescribeUserRequestById(tk.UserIdString())
	_, err = i.user.DescribeUser(ctx, uDesc)
	if err != nil {
		return nil, err
	}

	return tk, nil
}

// 校验Token 是给内部中间层使用 身份校验层
func (i *TokenServiceImpl) ValiateToken(ctx context.Context, req *token.ValiateTokenRequest) (*token.Token, error) {
	// 1. 查询Token (是不是我们这个系统颁发的)
	tk := token.NewToken()
	err := datasource.DBFromCtx(ctx).
		Where("access_token = ?", req.AccessToken).
		First(tk).
		Error
	if err != nil {
		return nil, err
	}

	// 2.1 判断Ak是否过期
	if err := tk.IsAccessTokenExpired(); err != nil {
		// 判断刷新Token是否过期
		if err := tk.IsRreshTokenExpired(); err != nil {
			return nil, err
		}

		// 如果开启了自动刷新
		if i.AutoRefresh {
			tk.SetRefreshAt(time.Now())
			tk.SetExpiredAtByDuration(i.refreshDuration, 4)
			if err := datasource.DBFromCtx(ctx).Save(tk); err != nil {
				i.log.Error().Msgf("auto refresh token error, %s", err.Error)
			}
		}

		return nil, err
	}

	return tk, nil
}

func (i *TokenServiceImpl) DescribeToken(ctx context.Context, in *token.DescribeTokenRequest) (*token.Token, error) {
	query := datasource.DBFromCtx(ctx)
	switch in.DescribeBy {
	case token.DESCRIBE_BY_ACCESS_TOKEN:
		query = query.Where("access_token = ?", in.DescribeValue)
	default:
		return nil, exception.NewBadRequest("unspport describe type %s", in.DescribeValue)
	}

	tk := token.NewToken()
	if err := query.First(tk).Error; err != nil {
		return nil, err
	}
	return tk, nil
}

// 退出接口(销毁Token)
func (i *TokenServiceImpl) RevolkToken(ctx context.Context, in *token.RevolkTokenRequest) (*token.Token, error) {
	tk, err := i.DescribeToken(ctx, token.NewDescribeTokenRequest(in.AccessToken))
	if err != nil {
		return nil, err
	}
	if err := tk.CheckRefreshToken(in.RefreshToken); err != nil {
		return nil, err
	}

	tk.Lock(token.LOCK_TYPE_REVOLK, "user revolk token")
	err = datasource.DBFromCtx(ctx).Model(&token.Token{}).
		Where("access_token = ?", in.AccessToken).
		Where("refresh_token = ?", in.RefreshToken).
		Updates(tk.Status.ToMap()).
		Error
	if err != nil {
		return nil, err
	}
	return tk, err
}

// 查询已经颁发出去的Token
func (i *TokenServiceImpl) QueryToken(ctx context.Context, in *token.QueryTokenRequest) (*types.Set[*token.Token], error) {
	set := types.New[*token.Token]()
	query := datasource.DBFromCtx(ctx).Model(&token.Token{})

	// 查询总量
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
