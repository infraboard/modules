package task

import (
	"context"
	"encoding/json"

	"github.com/infraboard/mcube/v2/tools/pretty"
)

var (
	sync_runners = map[string]SyncRunner{}
)

func RegistrySyncRunner(name string, runner SyncRunner) {
	sync_runners[name] = runner
}

func GetSyncRunner(name string) SyncRunner {
	return sync_runners[name]
}

func ListSyncRunner() (ns []string) {
	for k := range sync_runners {
		ns = append(ns, k)
	}
	return
}

type SyncRunner interface {
	// 同步执行，通过ctx控制取消
	Run(ctx context.Context, task *Task) error
}

type UnimplementedSyncRunner struct{}

func (r *UnimplementedSyncRunner) Start(context.Context, *Task) {}

const (
	RUN_PARAM_TYPE_JSON = "json"
)

type RUN_PARAM_TYPE string

func NewJsonRunParam(jsonStr string) *RunParam {
	return &RunParam{
		Type:  RUN_PARAM_TYPE_JSON,
		Value: jsonStr,
	}
}

type RunParam struct {
	// 参数类型
	Type RUN_PARAM_TYPE `json:"type" gorm:"column:run_param_type;type:varchar(60)" description:"参数类型"`
	// 参数值
	Value string `json:"value" gorm:"column:run_param_value;type:text" description:"参数值"`
}

func (r *RunParam) String() string {
	return pretty.ToJSON(r)
}

func (r *RunParam) Load(v any) error {
	return json.Unmarshal([]byte(r.Value), v)
}
