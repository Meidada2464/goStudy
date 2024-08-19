package main

import (
	"errors"
	"fmt"
	"math"
	"net"
	"os"
	"strconv"
	"syscall"
	"time"
)

const (
	// ECONNREFUSED 0x6f,
	// reference the errors in zerrors_linux_amd64.go
	ECONNREFUSED        = "connection refused"
	ConnectEConnRefused = "connect: connection refused"
)

func main() {
	// 获取到用户自定义的参数
	args := os.Args
	if len(args) > 5 {
		panic("args too long")
	}

	// 本地ip
	localIp := args[1]

	// 目标ip
	dstIp := args[2]

	// 目标端口
	dstPort, err := strconv.Atoi(args[3])
	if err != nil {
		fmt.Println("dstPort error")
		return
	}

	// 发包数
	packetNum, err := strconv.Atoi(args[4])
	if err != nil {
		fmt.Println("packetNum error")
		return
	}

	avgPacketLossRate, avgRtt, err := tcpPing(localIp, dstIp, dstPort, packetNum)
	fmt.Println("ping:", avgPacketLossRate, "avgRtt:", avgRtt, "error:", err)
}

func tcpPing(localIp, ip string, port, packetNum int) (avgPacketLossRate float64, avgRtt float64, err error) {
	addr, err := net.ResolveTCPAddr("tcp6", fmt.Sprintf("[%s]:%d", ip, port))
	if err != nil {
		return
	}

	var (
		packetsSend int
		packetsRecv int
		totalRtt    float64
		rtt         time.Duration
		dialer      *net.Dialer

		onceTimeout = time.Second // 单次探测的超时时间
	)

	if localIp != "" {
		var localAddr *net.TCPAddr
		localAddr, err = net.ResolveTCPAddr("tcp6", fmt.Sprintf("[%s]:0", localIp))
		if err != nil {
			return
		}
		dialer = &net.Dialer{
			LocalAddr: localAddr,
			Timeout:   onceTimeout,
		}
	}

	ticker := time.NewTicker(time.Second) // 用于控制探测间隔
	defer ticker.Stop()

	for i := 0; i < packetNum; i++ {
		packetsSend++

		// blocking dial
		if dialer != nil {
			rtt, err = tcpPingOnceByDialer(dialer, addr.String())
		} else {
			rtt, err = tcpPingOnceByDial(addr.String(), onceTimeout)
		}
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

	err = nil
	return
}

// tcpPingOnceByDialer 单次尝试建连
func tcpPingOnceByDialer(dialer *net.Dialer, addr string) (rtt time.Duration, err error) {
	start := time.Now()
	defer func() {
		rtt = time.Since(start)
	}()

	// TODO: net.Dial will retry one time. If use syscall.Connect,
	// it can't control the timeout(it actually can control, but more complicated..).
	conn, err := dialer.Dial("tcp6", addr)
	if err != nil {
		var opErr *net.OpError
		if errors.As(err, &opErr) &&
			opErr.Err.Error() == ConnectEConnRefused {
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

// tcpPingOnceByDial 单次尝试建连
func tcpPingOnceByDial(addr string, timeout time.Duration) (rtt time.Duration, err error) {
	start := time.Now()
	defer func() {
		rtt = time.Since(start)
	}()

	// TODO: net.Dial will retry one time. If use syscall.Connect,
	// it can't control the timeout(it actually can control, but more complicated..).
	conn, err := net.DialTimeout("tcp6", addr, timeout)
	if err != nil {
		var opErr *net.OpError
		if errors.As(err, &opErr) &&
			opErr.Err.Error() == ConnectEConnRefused {
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

// tcpPingByConnectOnce 单次调用Connect
func tcpPingByConnectOnce(fd int, addr syscall.SockaddrInet4) (rtt time.Duration, err error) {
	start := time.Now()
	defer func() {
		rtt = time.Since(start)
	}()

	err = syscall.Connect(fd, &addr)
	if err != nil {
		if err.Error() == ECONNREFUSED {
			err = nil
			return
		}
		return
	}

	// 这种情况虽然没有问题，但是应该预警
	//err = fmt.Errorf("%s's port is open", addr)
	return
}

// 还有一种是直接RawSocket的方式，但是问题在于会收到IP层的所有数据包，需要自行比对五元组进行筛选数据包，
// 实现起来会更加复杂，并且有可能性能损耗也会多一点点。但是好处也是显而易见的，定制化程度最高，
// 并且可以自行对数据包内容做一定程度的调控，如果配合服务端的话，甚至可以在payload里定义一个基于4层传输层之上的“应用协议”，
// 该协议的好处在于，可以更加细粒度的感知数据包收发情况，甚至后续进一步扩展到7层做速度的衡量也是适用的，
// 但是这种方式更加接近于一种In-House Fine-Grained的方式，需要有公司层面的支持以及业务的需求（如调度、质量）再考虑去迭代，
// 否则得不偿失。

type float interface {
	~float32 | ~float64
}

// ToFixed rounds down to a float to a given number of decimal places
//
// digits must be greater than or equal to 0
func ToFixed[T float](n T, digits int) T {
	if digits < 0 {
		digits = 0
	}
	pow := math.Pow10(digits)
	i := int64(float64(n) * pow)
	return T(float64(i) / pow)
}
