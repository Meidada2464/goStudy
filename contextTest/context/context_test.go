/**
 * Package context
 * @Author fengfeng.mei <fengfeng.mei@baishan.com>
 * @Date 2024/10/29 13:27
 */

package main

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

// TestContext 使用context对goroutine进行控制
func TestContext1(t *testing.T) {
	var (
		wg = &sync.WaitGroup{}
	)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	wg.Add(1)
	go handle1(ctx, 500*time.Millisecond, wg)
	wg.Wait()
}

func handle1(ctx context.Context, duration time.Duration, wg *sync.WaitGroup) {
	defer wg.Done()
	select {
	case <-time.After(duration):
		fmt.Println("the duration is over")
	case <-ctx.Done():
		fmt.Println("the context is done")
	}
}

func TestContext2(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	go handle2(ctx, 6*time.Second)
	select {
	case <-ctx.Done():
		fmt.Println("the main context is done")
	}
}

func handle2(ctx context.Context, duration time.Duration) {
	select {
	case <-time.After(duration):
		fmt.Println("the duration is over")
	case <-ctx.Done():
		fmt.Println("the context is done")
	}
}

func TestContext3(t *testing.T) {
	var (
		wg = &sync.WaitGroup{}
	)
	ctx, cancelFunc := context.WithCancel(context.Background())
	wg.Add(3)
	go handle31(ctx, wg)
	go handle32(ctx, wg)
	go handle33(ctx, wg)
	wg.Wait()
	cancelFunc()
	select {
	case <-ctx.Done():
		fmt.Println("the main context is done")
	}
}

func handle31(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(20 * time.Second)
	fmt.Println("the handle31 is done")
}

func handle32(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(10 * time.Second)
	fmt.Println("the handle32 is done")
}

func handle33(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(5 * time.Second)
	fmt.Println("the handle33 is done")
}
