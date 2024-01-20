package Goroutine

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
)

func TestReadFile(t *testing.T) {
	//	创建一个闭包函数并赋值,查看运行时的GC状态
	memConsumed := func() uint64 {
		runtime.GC()
		var s runtime.MemStats
		runtime.ReadMemStats(&s)
		return s.Sys
	}

	var c <-chan interface{}
	var wg sync.WaitGroup
	noop := func() { wg.Done(); <-c }

	const numGoroutines = 1e4
	wg.Add(numGoroutines)
	before := memConsumed()

	for i := numGoroutines; i > 0; i-- {
		go noop()
	}
	wg.Wait()

	after := memConsumed()
	fmt.Printf("%.3fkb", float64(after-before)/numGoroutines/1e3)
}
