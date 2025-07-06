package webhook

import (
	"context"
	"time"

	"github.com/google/uuid"
	"resty.dev/v3"
)

func NewWebHook(spec WebHookSpec) *WebHook {
	return &WebHook{
		WebHookSpec:   spec,
		WebHookStatus: *NewWebHookStatus(),
	}
}

// 任务回调
type WebHook struct {
	// WebHook定义
	WebHookSpec
	// 状态与统计
	WebHookStatus
}

func (h *WebHook) TableName() string {
	return "task_webhooks"
}

func (h *WebHook) SetDefault() {
	if h.ContentType == "" {
		h.ContentType = "application/json"
	}
	if h.Timeout == 0 {
		h.Timeout = 30
	}
	if h.ID == "" {
		h.ID = uuid.NewString()
	}
}

func (h *WebHook) Run(ctx context.Context) {
	h.SetDefault()

	resty.New()
	resp, err := resty.New().R().WithContext(ctx).
		SetAuthToken(h.Secret).
		SetTimeout(time.Second*time.Duration(h.Timeout)).
		SetHeaders(h.Headers).
		SetContentType(h.ContentType).
		SetBody(h.Payload).
		SetRetryCount(h.RetryCount).Execute(h.Method, h.TargetURL)
	if err != nil {
		h.Status = STATUS_FAILED
		h.Message = err.Error()
	} else {
		h.Status = STATUS_SUCCESS
		h.Message = resp.String()
	}
}

type WebHookSpec struct {
	// 基本信息
	ID          string `json:"id"`          // WebHook 的唯一标识
	Name        string `json:"name"`        // WebHook 名称
	Description string `json:"description"` // 描述

	RefTaskId string `json:"ref_task_id"` // 关联Task

	// 目标配置
	TargetURL string            `json:"target_url"` // 接收 HTTP 请求的 URL
	Method    string            `json:"method"`     // HTTP 方法，如 POST、GET 等
	Headers   map[string]string `json:"headers"`    // 自定义请求头

	// 触发条件
	Events []string `json:"events"` // 触发的事件类型列表
	Filter string   `json:"filter"` // 事件过滤条件（可选）

	// 安全验证
	Secret          string `json:"secret"`           // 用于签名验证的密钥（HMAC）
	SignatureHeader string `json:"signature_header"` // 签名头的名称，如 "X-Hub-Signature"

	// 请求内容配置
	ContentType string `json:"content_type"` // 请求体类型，如 "application/json"
	Payload     string `json:"payload"`      // 自定义请求体模板（可选）

	// 重试与超时
	RetryCount int `json:"retry_count"` // 失败重试次数
	Timeout    int `json:"timeout"`     // 请求超时时间（毫秒）
}

func NewWebHookStatus() *WebHookStatus {
	return &WebHookStatus{}
}

type WebHookStatus struct {
	// 状态与统计
	Status STATUS `json:"status"`
	// 失败信息
	Message string `json:"message"`
}

func (w *WebHook) Failed(err error) *WebHook {
	w.Status = STATUS_FAILED
	w.Message = err.Error()
	return w
}

func (w *WebHook) Success() *WebHook {
	w.Status = STATUS_SUCCESS
	return w
}
