package runtime_1

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestContext(t *testing.T) {
	ctx := context.Background()
	context.WithTimeout(ctx, time.Second)
	cancel, cancelFunc := context.WithCancel(ctx)

	context.WithDeadline(cancel, time.Now())
	cancelFunc()
}

func TestLock(t *testing.T) {
	// 互斥锁
	mutex := sync.Mutex{}
	mutex.Lock()

	mutex.Unlock()
}
