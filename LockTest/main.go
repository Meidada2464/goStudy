/**
 * Package LockTest
 * @Author fengfeng.mei <fengfeng.mei@baishan.com>
 * @Date 2024/11/19 12:48
 */

package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	instance *Ping
	once     sync.Once
)

type Ping struct {
	lock       sync.RWMutex
	updateChan chan int
}

func GetInstance() *Ping {
	once.Do(func() {
		instance = &Ping{
			updateChan: make(chan int),
		}
	})
	return instance
}

func main() {
	wg := &sync.WaitGroup{}
	GetInstance()
	wg.Add(2)
	go instance.fnSend(wg)
	go instance.fnReceive(wg)
	wg.Wait()
}

func (p *Ping) fnSend(wg *sync.WaitGroup) {
	startTime := time.Now()
	for i := 0; i < 10; i++ {
		p.updateChan <- i
		fmt.Println("send i:", i)
		time.Sleep(time.Second)
	}
	endTime := time.Now()
	fmt.Println("cost:", endTime.Sub(startTime))
	wg.Done()
}

func (p *Ping) fnReceive(wg *sync.WaitGroup) {
	for {
		i := <-p.updateChan
		fmt.Println("receive i:", i)
		if i == 9 {
			wg.Done()
		}
		time.Sleep(time.Second * 2)
	}
}
