package task

import (
	"context"
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
	Type  RUN_PARAM_TYPE `json:"type"`
	Value any            `json:"value"`
}
