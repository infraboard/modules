package webhook

import (
	"context"
	"encoding/json"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/infraboard/mcube/v2/ioc/config/bus"
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
	// WebHook 的唯一标识
	Id string `json:"id"`
	// 创建时间
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;type:timestamp;default:current_timestamp;not null;index;" description:"创建时间"`
	// WebHook定义
	WebHookSpec
	// 状态与统计
	WebHookStatus
}

func (h *WebHook) TableName() string {
	return "webhooks"
}

func (e *WebHook) LoadFromEvent(event *bus.Event) error {
	return json.Unmarshal(event.Data, e)
}

func (h *WebHook) SetDefault() {
	if h.ContentType == "" {
		h.ContentType = "application/json"
	}
	if h.Timeout == 0 {
		h.Timeout = 30
	}
	if h.Id == "" {
		h.Id = uuid.NewString()
	}
}

func (h *WebHook) Run(ctx context.Context) {
	h.SetDefault()
	h.SetUpdateAt(time.Now())
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

func NewWebHookSpec() *WebHookSpec {
	return &WebHookSpec{
		Headers:    map[string]string{},
		Conditions: []string{},
	}
}

type WebHookSpec struct {
	// WebHook 名称
	Name string `json:"name"`
	// 描述
	Description string `json:"description"`

	// 关联Task
	RefTaskId string `json:"ref_task_id"`

	// 目标配置
	// 接收 HTTP 请求的 URL
	TargetURL string `json:"target_url"`
	// HTTP 方法，如 POST、GET 等
	Method string `json:"method"`
	// 自定义请求头
	Headers map[string]string `json:"headers" gorm:"column:headers;type:json;serializer:json;"`

	// 安全验证
	// 用于签名验证的密钥（HMAC）
	Secret string `json:"secret"`
	// 签名头的名称，如 "X-Hub-Signature"
	SignatureHeader string `json:"signature_header"`

	// 哪些条件下触发
	Conditions []string `json:"conditions" gorm:"column:conditions;type:json;serializer:json;"`

	// 请求内容配置
	// 请求体类型，如 "application/json"
	ContentType string `json:"content_type"`
	// 自定义请求体模板（可选）
	Payload string `json:"payload"`

	// 重试与超时
	// 失败重试次数
	RetryCount int `json:"retry_count"`
	// 请求超时时间（毫秒）
	Timeout int `json:"timeout"`
}

func (s *WebHookSpec) AddCondtion(conditions ...string) {
	if s.Conditions == nil {
		s.Conditions = []string{}
	}
	s.Conditions = append(s.Conditions, conditions...)
}

func (s *WebHookSpec) IsCondtionOk(condition string) bool {
	return slices.Contains(s.Conditions, condition)
}

func NewWebHookStatus() *WebHookStatus {
	return &WebHookStatus{}
}

type WebHookStatus struct {
	// 更新时间
	UpdateAt *time.Time `json:"update_at" gorm:"column:update_at;type:timestamp;default:current_timestamp;not null;index;" description:"更新时间"`
	// 状态与统计
	Status STATUS `json:"status" gorm:"column:status;type:varchar(32);not null;" description:"状态"`
	// 失败信息
	Message string `json:"message" gorm:"column:message;type:text;" description:"失败信息"`
}

func (w WebHookStatus) TableName() string {
	return "webhooks"
}

func (w *WebHookStatus) SetUpdateAt(v time.Time) *WebHookStatus {
	w.UpdateAt = &v
	return w
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
