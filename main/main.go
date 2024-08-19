package main

import (
	"errors"
	"fmt"
	"math"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

// 使用tcping的方式开协程来进行模拟高并发ping探测
func main() {
	// 获取参数对象
	args := os.Args
	if len(args) > 5 {
		panic("args too long")
	}

	_, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Println("args error")
		return
	}

	port, err := strconv.Atoi(args[3])
	if err != nil {
		fmt.Println("port error")
		return
	}

	outTime, err := strconv.Atoi(args[4])
	if err != nil {
		fmt.Println("outTime error")
		return
	}

	var wg sync.WaitGroup

	// 执行tcping的次数
	for {
		wg.Add(1)
		go func(count int) {
			defer wg.Done()
			loss, avgRtt, err := tcpPing(args[2], port, count, outTime)
			fmt.Println("execCount:", count, "dst:", args[2], ":", port, "ping:", loss, "avgRtt:", avgRtt)
			if err != nil {
				fmt.Println("execCount:", count, "dst:", args[2], ":", port, "pingError:", loss, "avgRtt:", avgRtt)
				fmt.Println("execCount:", count, "go tcping error", err)
			}
		}(1)
		time.Sleep(time.Second)
	}
	wg.Wait()
}

func tcpPing(ip string, port int, seq, outTIme int) (avgPacketLossRate float64, avgRtt float64, err error) {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		fmt.Println("tcpPing-Error-57", err)
		return
	}

	var (
		packetsSend int
		packetsRecv int
		totalRtt    float64
		rtt         time.Duration

		// 这两个参数暂时未提参数到外部配置，尽量在上线一段时间后，观察后决定是否提参
		ticker      = time.NewTicker(time.Second)          // 用于控制探测间隔
		onceTimeout = time.Second * time.Duration(outTIme) // 单次探测的超时时间
	)
	defer ticker.Stop()

	for i := 0; i < 8; i++ {
		packetsSend++

		// blocking dial
		rtt, err = tcpPingByDialOnce(addr.String(), onceTimeout, seq, i)
		if onceTimeout.Milliseconds()-rtt.Milliseconds() > 100 {
			// 单次超时和单次请求的时间差值小于100ms认为这次已经接近或者已经超时了
			<-ticker.C
		}

		if err != nil { // 此处考虑对失败情况做打印或者一定收集用于后续的监控分析和迭代赶紧
			continue
		}

		packetsRecv++
		totalRtt += float64(rtt.Microseconds()) / 1e3 // 不能直接加ms，否则有可能会丢失一些精度，例如同机房、单机都可能小于1ms
	}

	// calc result
	if packetsSend == 0 {
		avgPacketLossRate = 100
		avgRtt = 0
	} else {
		avgPacketLossRate = ToFixed(float64(packetsSend-packetsRecv)/float64(packetsSend)*100, 2)
		if packetsRecv == 0 {
			avgRtt = 0
		} else {
			avgRtt = ToFixed(totalRtt/float64(packetsRecv), 2)
		}
	}

	return
}

// tcpPingByDialOnce 单次发SYN包
func tcpPingByDialOnce(addr string, timeout time.Duration, seq, execCount int) (rtt time.Duration, err error) {
	start := time.Now()
	defer func() {
		rtt = time.Since(start)
	}()

	// TODO: net.Dial will retry one time. If use syscall.Connect,
	// it can't control the timeout(it actually can control, but more complicated..).
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		var testError *net.OpError
		errors.As(err, &testError)
		//fmt.Println("tcpPing-Error-118", "seq", seq, "execCount", execCount, "testError", testError, "err", err)
		var opErr *net.OpError
		if errors.As(err, &opErr) &&
			opErr.Err.Error() == "connect: connection refused" {
			// 这种情况符合预期，说明收到了RST
			err = nil
			return
		}
		return
	}

	_ = conn.Close()
	// 这种情况虽然没有问题，但是应该预警
	//err = fmt.Errorf("%s's port is open", addr)
	return
}

type float interface {
	~float32 | ~float64
}

func ToFixed[T float](n T, digits int) T {
	if digits < 0 {
		digits = 0
	}
	pow := math.Pow10(digits)
	i := int64(float64(n) * pow)
	return T(float64(i) / pow)
}
