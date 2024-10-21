package trait

import (
	"context"
	"time"
)

type WorkerPool interface {
	Push(task func()) error
	BatchPush(tasks ...func()) (int, error)
	PushWithTimeOut(timeout time.Duration, task func()) error
	BatchPushWithTimeOut(timeout time.Duration, tasks ...func()) (int, error)
	PushWithContext(ctx context.Context, task func()) error
	BatchPushWithContext(ctx context.Context, tasks ...func()) (int, error)
	Stop()
}
