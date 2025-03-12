package snippets

import (
	"context"
	"fmt"
	"time"
)

func DoSomething(ctx context.Context) {
	select {
	case <-time.After(1 * time.Second):
		fmt.Println("任务完成")
	case <-ctx.Done():
		fmt.Println("任务被取消:", ctx.Err())
	}
}
