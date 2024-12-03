/**
 * Package httpClient
 * @Author fengfeng.mei <fengfeng.mei@baishan.com>
 * @Date 2024/12/3 14:20
 */

package main

import (
	"fmt"
	"httpclient/httpS"
	"runtime"
	"time"
)

var (
	intTemp []int
)

func main() {
	// 开启httpsServer
	httpServerStart()

	testServer()
}

func httpServerStart() {
	instance := httpS.GetInstance()
	err := instance.Start()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func testServer() {
loopSign:
	for i := 0; i < 10; i++ {

		go func() {
			fmt.Println("in go i:", i)
			time.Sleep(20 * time.Second)
		}()

		time.Sleep(2 * time.Second)
		goroutine := runtime.NumGoroutine()
		fmt.Println("i:", i, "goroutine num :", goroutine)

		if i == 9 {
			goto loopSign
		}
	}
}
