package impl_test

import (
	"testing"
	"time"

	"github.com/infraboard/modules/task/apps/event"
)

func TestAddEvent(t *testing.T) {
	spec := event.NewEventSpec()
	spec.Message = "test"
	ins, err := svc.AddEvent(t.Context(), spec)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ins)

	// 等待消费者处理
	time.Sleep(60 * time.Second)
}

func TestQueryEvent(t *testing.T) {
	ins, err := svc.QueryEvent(t.Context(), event.NewQueryEventRequest())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ins)
}
