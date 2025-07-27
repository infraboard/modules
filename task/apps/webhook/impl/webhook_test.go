package impl_test

import (
	"testing"
	"time"

	"github.com/infraboard/modules/task/apps/webhook"
)

func TestRunWebHook(t *testing.T) {
	spec := webhook.NewWebHookSpec()
	ins, err := svc.Run(t.Context(), spec)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ins)

	// 等待消费者处理
	time.Sleep(60 * time.Second)
}

func TestQueryWebHook(t *testing.T) {
	ins, err := svc.QueryWebHook(t.Context(), webhook.NewQueryWebHookRequest())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ins)
}
