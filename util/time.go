/**
 * Package util
 * @Author fuqiang.li <fuqiang.li@baishan.com>
 * @Date 2023/3/17
 */

package util

import (
	"time"
)

// TimeLoop 循环执行某个函数，立刻执行然后等待周期
func TimeLoop(interval time.Duration, fn func(t time.Time)) func() {
	if fn == nil {
		return nil
	}
	ticker := time.NewTicker(interval)
	go func() {
		for {
			fn(time.Now())
			<-ticker.C
		}
	}()
	return func() {
		ticker.Stop()
	}
}

// TimeLoopThen 循环执行某个函数，等待周期到达再执行
func TimeLoopThen(interval time.Duration, fn func(time.Time)) func() {
	if fn == nil {
		return nil
	}
	ticker := time.NewTicker(interval)
	go func() {
		for {
			now := <-ticker.C
			fn(now)
		}
	}()
	return func() {
		ticker.Stop()
	}
}

// DoWhile 立即执行，然后循环执行
func DoWhile(interval time.Duration, fn func(), stopChan chan struct{}) {
	if fn == nil {
		return
	}
	go func() {
		fn()
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				go fn()
			case <-stopChan:
				return
			}
		}
	}()
}

// While 循环执行
func While(interval time.Duration, fn func(), stopChan chan int) {
	if fn == nil {
		return
	}
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				go fn()
			case <-stopChan:
				return
			}
		}
	}()
}
