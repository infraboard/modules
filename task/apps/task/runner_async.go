package task

import (
	"context"
	"fmt"
)

var (
	async_runners = map[string]AsyncRunner{}
)

func RegistryAsyncRunner(name string, runner AsyncRunner) {
	async_runners[name] = runner
}

func GetAsyncRunner(name string) AsyncRunner {
	return async_runners[name]
}

func ListAsyncRunner() (ns []string) {
	for k := range async_runners {
		ns = append(ns, k)
	}
	return
}

type AsyncRunner interface {
	// 触发异步执行
	Start(ctx context.Context, task *Task) error
	// 获取任务状态
	GetStatus(ctx context.Context, taskID string) (*TaskStatus, error)
	// 取消任务
	Cancel(ctx context.Context, taskID string) error
}

type UnimplementedAsyncRunner struct{}

func (r *UnimplementedAsyncRunner) Start(context.Context, *Task) error {
	return fmt.Errorf("Unimplemented")
}
func (r *UnimplementedAsyncRunner) GetStatus(ctx context.Context, taskID string) (*TaskStatus, error) {
	return nil, fmt.Errorf("Unimplemented")
}
func (r *UnimplementedAsyncRunner) Cancel(ctx context.Context, taskID string) error {
	return fmt.Errorf("Unimplemented")
}
