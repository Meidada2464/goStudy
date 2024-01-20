package realtime_task

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

func Test(t *testing.T) {
	//m := GetInstance()
	//err := m.Start()
	//fmt.Println("start success")
	//if err != nil {
	//	fmt.Println("start error")
	//	return
	//}
	//
	//go func() {
	//	for {
	//		singleUuid := uuid.New().String()
	//		m.putRealTask(singleUuid)
	//		time.Sleep(time.Second * 10)
	//	}
	//}()
	//
	//group := sync.WaitGroup{}
	//group.Add(1)
	//group.Wait()

	loss, rtt := ParsePingOutput("2 packets transmitted, 0 received, +2 errors, 100% packet loss, time 999ms")
	fmt.Println("loss:", loss, "rtt", rtt)
}

type (
	Manager struct {
		stopChan    chan int    // 停止信号
		dealSign    chan int    // 任务触发信号
		TaskChannel chan string // 传递创建的即时任务的uuid
	}
)

var (
	instance *Manager
	once     sync.Once
)

// GetInstance 获取设备信息管理实例
func GetInstance() *Manager {
	once.Do(func() {
		instance = &Manager{
			stopChan:    make(chan int, 1),
			dealSign:    make(chan int, 1e2),
			TaskChannel: make(chan string, 1e3),
		}
	})
	return instance
}

func (m *Manager) Start() error {
	go m.run()
	return nil
}

func (m *Manager) Stop() {
	m.stopChan <- 1
}

func (m *Manager) run() {
	fmt.Println("manager run begin")
	// 处理即时任务
	for {
		select {
		case <-m.dealSign:
			// 获取及时探测的任务并调用grpc下发
			m.dealRealTimeTask()
		default:

		}
		time.Sleep(time.Second)
	}
}

func (m *Manager) dealRealTimeTask() {
	task := m.popRealTask()
	// 1、数据转换
	fmt.Println("task:", task)
	// 2、grpc调用
}

func (m *Manager) putRealTask(uuid string) {
	fmt.Println("putRealTask-success")
	m.TaskChannel <- uuid
	m.dealSign <- 1
}

func (m *Manager) popRealTask() string {
	fmt.Println("popRealTask-success")
	return <-m.TaskChannel
}

// parsePingOutput 解析输出结果
func ParsePingOutput(output string) (packetLoss, avgRTT float64) {
	// 根据输出中的关键字进行分割和查找
	for _, line := range strings.Split(output, "\n") {
		if strings.Contains(line, "packet loss") {
			// 如果找到包含 packet loss 的行，则进行解析
			// Darwin:
			// 5 packets transmitted, 5 packets received, 0.0% packet loss
			// 2 packets transmitted, 0 packets received, 100.0% packet loss
			// Linux:
			// 2 packets transmitted, 2 received, 0% packet loss, time 1001ms
			// 2 packets transmitted, 0 received, +2 errors, 100% packet loss, time 999ms
			plIdx := strings.LastIndex(line, "% packet loss")
			fields := strings.Fields(line[:plIdx])
			if len(fields) > 5 {
				lost, err := strconv.ParseFloat(fields[len(fields)-1], 64)
				if err == nil {
					packetLoss = lost
				}
			}
		}
		if strings.Contains(line, "min/avg/max/") {
			// 可能没有，比如不同的的时候
			// Darwin:
			// round-trip min/avg/max/stddev = 51.559/52.567/54.262/1.006 ms
			// Linux:
			// rtt min/avg/max/mdev = 12.020/12.029/12.039/0.110 ms
			fields := strings.Split(line, "/")
			if len(fields) == 7 {
				rtt, err := strconv.ParseFloat(fields[4], 64)
				if err == nil {
					avgRTT = rtt
				}
			}
		}
	}

	return packetLoss, avgRTT
}
