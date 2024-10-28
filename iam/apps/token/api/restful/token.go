package restful

import (
	"net/http"
	"net/url"

	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/v2/http/restful/response"
	"github.com/infraboard/mcube/v2/ioc/config/application"
	"github.com/infraboard/modules/iam/apps/token"
)

func (h *TokenRestulApiHandler) Login(r *restful.Request, w *restful.Response) {
	// 1. 获取用户的请求参数， 参数在Body里面
	req := token.NewIssueTokenRequest()

	err := r.ReadEntity(req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// 2. 执行逻辑
	tk, err := h.svc.IssueToken(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// access_token 通过SetCookie 直接写到浏览器客户端(Web)
	http.SetCookie(w, &http.Cookie{
		Name:     token.ACCESS_TOKEN_COOKIE_NAME,
		Value:    url.QueryEscape(tk.AccessToken),
		MaxAge:   0,
		Path:     "/",
		Domain:   application.Get().Domain(),
		SameSite: http.SameSiteDefaultMode,
		Secure:   false,
		HttpOnly: true,
	})
	// 在Header头中也添加Token
	w.Header().Set(token.ACCESS_TOKEN_RESPONSE_HEADER_NAME, tk.AccessToken)

	// 3. Body中返回Token对象
	response.Success(w, tk)
}

func (h *TokenRestulApiHandler) ChangeNamespce(r *restful.Request, w *restful.Response) {
	// 1. 获取用户的请求参数， 参数在Body里面
	req := token.NewChangeNamespceRequest()
	err := r.ReadEntity(req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	tk := token.GetTokenFromCtx(r.Request.Context())
	req.UserId = tk.UserId

	// 2. 执行逻辑
	tk, err = h.svc.ChangeNamespce(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// 3. Body中返回Token对象
	response.Success(w, tk)
}

// Logout HandleFunc
func (h *TokenRestulApiHandler) Logout(r *restful.Request, w *restful.Response) {
	req := token.NewRevolkTokenRequest(
		token.GetAccessTokenFromHTTP(r.Request),
		token.GetRefreshTokenFromHTTP(r.Request),
	)

	tk, err := h.svc.RevolkToken(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// access_token 通过SetCookie 直接写到浏览器客户端(Web)
	http.SetCookie(w, &http.Cookie{
		Name:     token.ACCESS_TOKEN_COOKIE_NAME,
		Value:    "",
		MaxAge:   0,
		Path:     "/",
		Domain:   application.Get().Domain(),
		SameSite: http.SameSiteDefaultMode,
		Secure:   false,
		HttpOnly: true,
	})

	// 3. 返回响应
	response.Success(w, tk)
}
