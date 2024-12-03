/**
 * Package timeWeelTest
 * @Author fengfeng.mei <fengfeng.mei@baishan.com>
 * @Date 2024/11/28 17:19
 */

package main

import (
	"fmt"
	"sync"
	"time"
	"timeWheel/timeWheel"
	"timeWheel/twUtils"
)

type (
	ModTask struct {
		interval time.Duration
		times    int
		key      interface{}
	}
)

func (m *ModTask) Run() {
	fmt.Println("work is running")
}

func (m *ModTask) Release() {
	fmt.Println("work is release")
}

func main() {
	// 创建一个时间轮，每秒刷新一次，每个小时是一个周期
	wg := sync.WaitGroup{}
	tw := timeWheel.New(time.Second, 60)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		tw.Start()
	}(&wg)

	mt := &ModTask{
		interval: twUtils.GetInterval(60),
		times:    60,
		key:      "123",
	}
	err := tw.AddTask(twUtils.GetInterval(60), 60, mt.key, mt)
	if err != nil {
		return
	}
	wg.Wait()

}
