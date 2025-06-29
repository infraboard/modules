package mysql

import (
	"context"
	"time"

	"github.com/infraboard/mcube/v2/desense"
	"github.com/infraboard/mcube/v2/exception"
	"github.com/infraboard/mcube/v2/ioc/config/datasource"
	"github.com/infraboard/mcube/v2/types"
	"github.com/infraboard/modules/iam/apps/policy"
	"github.com/infraboard/modules/iam/apps/token"
)

// 登录接口(颁发Token)
func (i *TokenServiceImpl) IssueToken(ctx context.Context, in *token.IssueTokenRequest) (*token.Token, error) {
	// 颁发Token
	issuer := token.GetIssuer(in.Issuer)
	if issuer == nil {
		return nil, exception.NewBadRequest("provider %s no support", in.Issuer)
	}
	tk, err := issuer.IssueToken(ctx, in.Parameter)
	if err != nil {
		return nil, err
	}
	tk.SetIssuer(in.Issuer).SetSource(in.Source)

	// 判断当前数据库有没有已经存在的Token
	activeTokenQueryReq := token.NewQueryTokenRequest().
		AddUserId(tk.UserId).
		SetSource(in.Source).
		SetActive(true)
	tks, err := i.QueryToken(ctx, activeTokenQueryReq)
	if err != nil {
		return nil, err
	}
	switch in.Source {
	// 每个端只能有1个活跃登录
	case token.SOURCE_WEB, token.SOURCE_IOS, token.SOURCE_ANDROID, token.SOURCE_PC:
		if tks.Len() > 0 {
			i.log.Debug().Msgf("use exist active token: %s", desense.Default().DeSense(tk.AccessToken, "4", "3"))
			return tks.Items[0], nil
		}
	case token.SOURCE_API:
		if tks.Len() > int(i.MaxActiveApiToken) {
			return nil, exception.NewBadRequest("max active api token overflow")
		}
	}

	// 保持Token
	if err := datasource.DBFromCtx(ctx).
		Create(tk).
		Error; err != nil {
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

	if in.Active != nil {
		if *in.Active {
			query = query.
				Where("lock_at IS NULL AND refresh_token_expired_at > ?", time.Now())
		} else {
			query = query.
				Where("lock_at IS NOT NULL OR refresh_token_expired_at <= ?", time.Now())
		}
	}
	if in.Source != nil {
		query = query.Where("source = ?", *in.Source)
	}
	if len(in.UserIds) > 0 {
		query = query.Where("user_id IN ?", in.UserIds)
	}

	// 查询总量
	err := query.Count(&set.Total).Error
	if err != nil {
		return nil, err
	}

	err = query.
		Order("issue_at desc").
		Offset(int(in.ComputeOffset())).
		Limit(int(in.PageSize)).
		Find(&set.Items).
		Error
	if err != nil {
		return nil, err
	}

	return set, nil
}

// 用户切换空间
func (i *TokenServiceImpl) ChangeNamespce(ctx context.Context, in *token.ChangeNamespceRequest) (*token.Token, error) {
	set, err := i.policy.QueryNamespace(ctx, policy.NewQueryNamespaceRequest().SetUserId(in.UserId).SetNamespaceId(in.NamespaceId))
	if err != nil {
		return nil, err
	}

	ns := set.First()
	if ns == nil {
		return nil, exception.NewPermissionDeny("你没有该空间访问权限")
	}

	// 更新Token
	tk, err := i.DescribeToken(ctx, token.NewDescribeTokenRequest(in.AccessToken))
	if err != nil {
		return nil, err
	}
	tk.SetNamespaceId(ns.Id)
	tk.NamespaceName = ns.Name

	// 保存状态
	if err := datasource.DBFromCtx(ctx).
		Updates(tk).
		Error; err != nil {
		return nil, err
	}
	return tk, nil
}
