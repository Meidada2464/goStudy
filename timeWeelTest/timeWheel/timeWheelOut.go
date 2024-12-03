/**
 * Package timeWeelTest
 * @Author fengfeng.mei <fengfeng.mei@baishan.com>
 * @Date 2024/11/28 16:37
 */

package timeWheel

import (
	"fmt"
	"github.com/go-errors/errors"
	"time"
)

func (tw *TimeWheel) AddTask(interval time.Duration, times int, key interface{}, cbTask Task) error {
	if interval <= 0 || key == nil || times < -1 || times == 0 {
		return errors.New("invalid param")
	}

	_, ok := tw.taskRecord.Load(key)
	if ok {
		return errors.New("key already exist")
	}

	delaySeconds := int(interval.Seconds())
	intervalSeconds := int(tw.interval.Seconds())
	step := delaySeconds / intervalSeconds
	if delaySeconds%intervalSeconds > 0 {
		step++
	}

	tw.addTaskChannel <- &twTask{
		interval: interval,
		times:    times,
		key:      key,
		cbTask:   cbTask,
		start:    tw.currentPos,
		step:     step,
	}

	fmt.Println("add task success")
	return nil
}

func (tw *TimeWheel) RemoveTask(key interface{}) error {
	if key == nil {
		return errors.New("invalid param")
	}

	value, ok := tw.taskRecord.Load(key)
	if !ok {
		return errors.New("key not exist")
	}
	task := value.(*twTask)
	task.cbTask.Release()
	tw.taskRecord.Delete(key)
	return nil
}

func (tw *TimeWheel) updateTask(key interface{}, interval time.Duration, cbTask Task) error {
	if key == nil {
		return errors.New("invalid param")
	}

	value, ok := tw.taskRecord.Load(key)
	if !ok {
		return errors.New("task does not exist")
	}
	task := value.(*twTask)
	task.cbTask = cbTask
	task.interval = interval
	return nil
}

func (tw *TimeWheel) Len() int {
	var c int
	tw.taskRecord.Range(func(key, value any) bool {
		c++
		return true
	})
	return c
}
