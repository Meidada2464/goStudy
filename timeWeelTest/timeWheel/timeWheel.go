/**
 * Package timeWeelTest
 * @Author fengfeng.mei <fengfeng.mei@baishan.com>
 * @Date 2024/11/28 14:50
 */

package timeWheel

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

type (
	Task interface {
		Run()
		Release()
	}

	TimeWheel struct {
		interval        time.Duration // 每次跳动的时间
		ticker          *time.Ticker
		slotNum         int // 插槽数
		slots           []*list.List
		taskNumsPerSlot []int
		currentPos      int
		stopChannel     chan bool
		taskRecord      *sync.Map
		lock            sync.RWMutex

		addTaskChannel chan *twTask
	}

	twTask struct {
		interval time.Duration
		times    int
		circle   int
		start    int
		step     int
		key      any
		cbTask   Task
	}
)

var (
	ErrTaskNotFind = fmt.Errorf("task not find")
)

// New create time wheel
func New(interval time.Duration, slotNum int) *TimeWheel {
	// 判断是否符合条件
	if interval <= 0 || slotNum <= 0 {
		fmt.Println("invalid param")
		return nil
	}

	tw := &TimeWheel{
		interval:        interval,
		slots:           make([]*list.List, slotNum), // 初始化插槽
		taskNumsPerSlot: make([]int, slotNum),        // 最外层记录一共有多少个插槽
		currentPos:      0,                           // 默认从第一个插槽开始
		slotNum:         slotNum,

		addTaskChannel: make(chan *twTask),
		stopChannel:    make(chan bool),
		taskRecord:     &sync.Map{},
	}

	// 初始化每个插槽中的链表
	for i := 0; i < tw.slotNum; i++ {
		tw.slots[i] = list.New()
		tw.taskNumsPerSlot[i] = 0
	}

	return tw
}

func (tw *TimeWheel) Start() {
	// 每次跳动的时间间隔
	tw.ticker = time.NewTicker(tw.interval)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		for {
			select {
			case <-tw.ticker.C: // 到了执行的时间，开始固定地执行
				fmt.Println("ticker.C inner", tw.currentPos)
				fmt.Println("Slot len", len(tw.slots))
				tw.tickHandler()
			case task := <-tw.addTaskChannel:
				fmt.Println("addTaskChannel", task)
				tw.addTask(task)
			case <-tw.stopChannel:
				tw.ticker.Stop() // 收到停止信号后停止ticker,防止协程泄露
				return
			}
		}
	}(&wg)
	wg.Wait()
}

func (tw *TimeWheel) Stop() {
	tw.stopChannel <- true
}

func (tw *TimeWheel) tickHandler() {
	l := tw.slots[tw.currentPos] // 找出当前坐标下的链表，该链表内是一组待执行的任务
	fmt.Println("tickHandler currentPos ", tw.currentPos, "l len", l.Len())

	tw.scanAddRunTask(l, tw.currentPos)
	if tw.currentPos == tw.slotNum-1 { // 每执行完一个时间轮，应该停止这个时间轮的运行
		tw.Stop()
	} else {
		tw.currentPos++
	}
}

func (tw *TimeWheel) scanAddRunTask(l *list.List, nowPos int) {
	if l == nil || l.Len() == 0 {
		fmt.Println("scanAddRunTask now Pos ", nowPos, "l len", l.Len())
		return
	}

	for item := l.Front(); item != nil; { // 队列内的内容从对头push出来，依次执行
		task := item.Value.(*twTask) // 断言队列中内容的合法性

		if task.times == 0 {
			next := item.Next() // 这里为什么要获取下一个节点？
			l.Remove(item)
			tw.taskRecord.Delete(task.key)
			item = next
			tw.taskNumsPerSlot[nowPos] = l.Len()
			continue
		}

		if task.circle > 0 {
			task.circle-- // 减少一次圈数
			item = item.Next()
			continue
		}

		// 开协程去真正执行
		go task.cbTask.Run()
		next := item.Next()

		l.Remove(item)
		tw.taskNumsPerSlot[nowPos] = l.Len()
		item = next

		if task.times == 1 {
			task.times = 0
			tw.taskRecord.Delete(task.key)
		} else {
			if task.times > 0 {
				task.times--
			}
		}
	}
}

// AddTask add task , 考虑的是如何将任务分发到不同的插槽，比如如何实现均衡分配等问题
func (tw *TimeWheel) addTask(task *twTask) {
	tw.lock.Lock()
	defer tw.lock.Unlock()

	if task.times == 0 {
		fmt.Println("task times should not be 0")
		return
	}

	// 再随机选择插槽的位置
	task.start += task.step
	task.start = task.start % tw.slotNum
	pos, circle := tw.getPositionAndCircle(tw.interval, task.start, task.step)
	task.circle = circle

	tw.slots[pos].PushBack(task)
	tw.taskNumsPerSlot[pos] = tw.slots[pos].Len()
	tw.taskRecord.Store(task.key, task)
}

func (tw *TimeWheel) getPositionAndCircle(d time.Duration, start, step int) (pos, circle int) {
	delaySeconds := int(d.Seconds())              // 延迟执行的时间
	intervalSeconds := int(tw.interval.Seconds()) // 时间轮定时执行的时间
	circle = delaySeconds / intervalSeconds / tw.slotNum
	pos = tw.getMinPos(start, step)
	return
}

func (tw *TimeWheel) getMinPos(start, step int) int {
	min := start
	step = step % tw.slotNum
	end := start + step
	if end > tw.slotNum { // 跨圈
		for i := start; i < tw.slotNum; i++ {
			if tw.taskNumsPerSlot[min] > tw.taskNumsPerSlot[i] {
				min = i
			}
		}
		for i := 0; i < end-tw.slotNum; i++ {
			if tw.taskNumsPerSlot[min] > tw.taskNumsPerSlot[i] {
				min = i
			}
		}
	} else { // 本圈
		for i := start; i < end; i++ {
			if tw.taskNumsPerSlot[min] > tw.taskNumsPerSlot[i] {
				min = i
			}
		}
	}
	// 每次都会执行
	if step == 0 {
		for i := 0; i < tw.slotNum; i++ {
			if tw.taskNumsPerSlot[min] > tw.taskNumsPerSlot[i] {
				min = i
			}
		}
	}
	return min
}
