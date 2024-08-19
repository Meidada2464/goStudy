package main

import (
	"fmt"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}
	// 并发安全的读和写
	channelOwner := func() <-chan int {
		ints := make(chan int, 10)
		// 使用闭包的方式，将处理数据的工作放到单独的协程里
		go func() {
			defer close(ints)
			for i := 0; i < 10; i++ {
				ints <- i
			}
		}()
		return ints
	}

	customer := func(ints <-chan int) {
		defer wg.Done()
		for i := range ints {
			fmt.Println("the value is ", i)
		}
	}

	for i := 0; i < 3; i++ {
		ints := channelOwner()
		wg.Add(1)
		go customer(ints)
		wg.Wait()
	}
}
