package main

import (
	"bytes"
	"fmt"
	"sync"
)

func main() {
	// 并发不安全的访问
	printData := func(wg *sync.WaitGroup, data []byte) {
		defer wg.Done()
		// 创建一个buffer，将字符流写入到buffer
		var buff bytes.Buffer
		for _, datum := range data {
			fmt.Fprintf(&buff, "%c", datum)
		}
		fmt.Println(buff.String())
	}

	var wg sync.WaitGroup
	data := []byte("hello word")
	wg.Add(2)
	go printData(&wg, data[3:])
	go printData(&wg, data[:3])
	wg.Wait()
}
