package main1

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

// runCommand 执行命令行
func runCommand(name string, arg []string, timeout time.Duration) ([]byte, []byte, error) {
	var (
		stdout = bytes.NewBuffer(nil)
		stderr = bytes.NewBuffer(nil)
	)
	cmd := exec.Command(name, arg...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	// 强制 kill
	timeStart := time.Now().Unix()
	time.AfterFunc(timeout, func() {
		if cmd.Process != nil {
			timeEnd := time.Now().Unix()
			fmt.Println("time kill :", timeEnd-timeStart)
			_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		}
	})

	timeStart2 := time.Now().Unix()
	err := cmd.Run()
	endStart2 := time.Now().Unix()
	fmt.Println("time run :", endStart2-timeStart2)

	return stdout.Bytes(), stderr.Bytes(), err
}

// parsePingOutput 解析输出结果
func parsePingOutput(output string) (packetLoss, avgRTT float64) {
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

// ping -c 8 -s 32 -w 15 -i 1 -q -I 100.127.84.224 61.241.127.10
// cmdPingOnce 实际发起ping
func cmdPingOnce(ip string) (lost, rtt float64, err error) {
	out, errOut, err := runCommand("ping", []string{"-c", "8", "-i", "1", "-w", "9", "-s", "32", "-q", ip}, 7*time.Second)
	if err != nil {

		fmt.Println("error is this :", err.Error())

		if err.Error() == "exit status 1" {
			fmt.Println("exit status 1 -> ", err)
		}

		if err.Error() == "signal: killed" {
			fmt.Println("signal: killed -> ", err)
		}

		if len(errOut) > 0 {
			return 100, 0, fmt.Errorf("%s, err out: %s", err, errOut)
		}

		return 100, 0, err
	}
	lost, rtt = parsePingOutput(string(out))
	return lost, rtt, nil
}

func main() {
	// 获取命令行参数
	args := os.Args

	// 确保至少有两个参数（第一个参数是程序的名称）
	if len(args) < 2 {
		panic("n")
	}

	// 解析参数为整数
	n, err := strconv.Atoi(args[1])
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			lost, rtt, err := cmdPingOnce("156.238.128.10")
			fmt.Println(i, lost, rtt, err)
		}(i)
	}
	wg.Wait()
}
