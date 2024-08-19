package goTest

import (
	"fmt"
	"github.com/pkg/profile"
	"goStudy/util"
	"testing"
	"time"
)

func BenchmarkFib(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Fib(10)
	}
}

func TestFib(t *testing.T) {
	for i := 0; i < 10; i++ {
		t.Log(Fib(i))
	}
}

func TestConcat(t *testing.T) {
	defer profile.Start(profile.MemProfile, profile.MemProfileRate(100)).Stop()
	concat(100)
}

func TestDoWhile(t *testing.T) {
	util.DoWhile(30*time.Second, func() {
		// 分类任务下发
		_ = func() error {
			fmt.Println("at sleep : ", time.Now())
			time.Sleep(time.Second * 60)
			fmt.Println("undo sleep : ", time.Now())
			return nil
		}()

		go func() {
			fmt.Println("in the go : ", time.Now())
		}()
	}, func() chan struct{} {
		var (
			c = make(chan struct{}, 1)
		)
		return c
	}())
}
