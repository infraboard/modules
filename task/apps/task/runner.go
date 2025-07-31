package task

import (
	"context"
	"encoding/json"
	"fmt"
)

var (
	runners = map[string]Runner{}
)

func RegistryRunner(name string, runner Runner) {
	runners[name] = runner
}

func GetRunner(name string) Runner {
	return runners[name]
}

func ListRunner() (ns []string) {
	for k := range runners {
		ns = append(ns, k)
	}
	return
}

// 执行器接口
type Runner interface {
	Run(context.Context, *RunParam) (fmt.Stringer, error)
}

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

func (r *RunParam) Load(v any) error {
	return json.Unmarshal([]byte(r.Value), v)
}
