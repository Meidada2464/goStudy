/**
 * Package gorountineTest
 * @Author fengfeng.mei <fengfeng.mei@baishan.com>
 * @Date 2024/11/26 17:16
 */

package gorountineTest

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestCurrent1(t *testing.T) {
	mainSync := sync.WaitGroup{}

	fn1 := func() {
		defer mainSync.Done()
		wg := sync.WaitGroup{}
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				fmt.Println("fn1 i:", i)
			}(i)
			time.Sleep(time.Second)
		}
		wg.Wait()
	}

	fn2 := func() {
		defer mainSync.Done()
		wg := sync.WaitGroup{}
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				fmt.Println("fn2 i:", i)
			}(i)
			time.Sleep(time.Second * 2)
		}
		wg.Wait()
	}

	mainSync.Add(1)
	go fn1()
	mainSync.Add(1)
	go fn2()

	mainSync.Wait()
}
