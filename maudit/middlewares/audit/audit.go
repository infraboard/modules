package aduit

import (
	"context"

	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/application"
	"github.com/infraboard/mcube/v2/ioc/config/gorestful"
	ioc_kafka "github.com/infraboard/mcube/v2/ioc/config/kafka"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/infraboard/modules/iam/apps/endpoint"
	permission "github.com/infraboard/modules/iam/permission/restful"
	"github.com/infraboard/modules/maudit/apps/event"
	"github.com/rs/zerolog"
	"github.com/segmentio/kafka-go"
)

func init() {
	ioc.Config().Registry(&auditor{
		Topic: "maudit",
	})
}

func Audit(v bool) (string, bool) {
	return event.META_AUDIT_KEY, v
}

type auditor struct {
	ioc.ObjectImpl

	log *zerolog.Logger

	// 当前这个消费者 配置的topic
	Topic string `toml:"topic" json:"topic" yaml:"topic"  env:"TOPIC"`
	//
	wirter *kafka.Writer
}

func (a *auditor) Name() string {
	return "auditor"
}

func (a *auditor) Init() error {
	a.log = log.Sub("mauditor")
	a.log.Debug().Msgf("maduit topic name: %s", a.Topic)
	a.wirter = ioc_kafka.Producer(a.Topic)

	// 添加到中间件, 加到Root Router里面
	ws := gorestful.RootRouter()
	ws.Filter(a.Audit())
	return nil
}

// 补充中间件函数逻辑
func (a *auditor) Audit() restful.FilterFunction {
	return func(r1 *restful.Request, r2 *restful.Response, fc *restful.FilterChain) {
		sr := r1.SelectedRoute()
		md := NewMetaData(sr.Metadata())

		// 开关打开，则开启审计
		if md.GetBool(event.META_AUDIT_KEY) {

			// 获取当前是否需要审计
			e := event.NewEvent()

			// 用户信息
			tk := permission.GetTokenFromCtx(r1.Request.Context())
			if tk != nil {
				e.Who = tk.UserName
				e.Extras["namespace"] = tk.NamespaceName
			}

			// ioc 里面获取当前应用的名称
			e.Service = application.Get().AppName
			e.ResourceType = md.GetString(endpoint.META_RESOURCE_KEY)
			e.Action = md.GetString(endpoint.META_ACTION_KEY)

			// {id} /:id
			e.ResourceId = r1.PathParameter("id")
			e.UserAgent = r1.Request.UserAgent()
			e.Extras["method"] = sr.Method()
			e.Extras["path"] = sr.Path()
			e.Extras["operation"] = sr.Operation()

			// 补充处理后的数据
			e.StatusCode = r2.StatusCode()
			// 发送给topic, 使用这个中间件的使用者，需要配置kafka
			err := a.wirter.WriteMessages(context.Background(), e.ToKafkaMessage())
			if err != nil {
				a.log.Error().Msgf("send message error, %s", err)
			} else {
				a.log.Debug().Msgf("send audit event ok, who: %s, resource: %s, action: %s", e.Who, e.ResourceType, e.Action)
			}
		}

		// 路有给后续逻辑
		fc.ProcessFilter(r1, r2)
	}
}

func NewMetaData(data map[string]any) *MetaData {
	return &MetaData{
		data: data,
	}
}

type MetaData struct {
	data map[string]any
}

func (m *MetaData) GetString(key string) string {
	if v, ok := m.data[key]; ok {
		return v.(string)
	}
	return ""
}

func (m *MetaData) GetBool(key string) bool {
	if v, ok := m.data[key]; ok {
		return v.(bool)
	}
	return false
}
