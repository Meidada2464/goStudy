/**
 * Package cmdPing
 * @Author fengfeng.mei <fengfeng.mei@baishan.com>
 * @Date 2024/12/3 16:05
 */

package main

import (
	"bytes"
	"fmt"
	"math"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	DefPktStr     = "8"  // 默认发包数
	DefPktSizeStr = "32" // 默认包大小
	DefWaitStr    = "1"  // 默认等待时间
	DefTimeout    = 30   // 默认超时时间
	DefTimeoutStr = "30" // 默认超时时间
)

type (
	IpFamily string
)

var (
	IpFamilyV4 IpFamily = "ipv4"
	IpFamilyV6 IpFamily = "ipv6"
)

func main() {
	args := os.Args

	if len(args) > 3 {
		fmt.Println("invalid args")
		return
	}

	ip := net.ParseIP(args[1])
	if ip == nil {
		fmt.Println("invalid ip")
		return
	}

	count := args[2]

	countInt, err := strconv.Atoi(count)
	if err != nil {
		fmt.Println("invalid count")
		return
	}

	for i := 0; i < countInt; i++ {
		fmt.Println("======================================================")
		fmt.Println("第", i+1, "次", "开始ping")

		lost, minRtt, avgRtt, maxRtt, mdevRtt, duration, err := cmdPingOnce(ip)
		if err != nil {
			fmt.Println("cmdPingOnce error :", err)
			fmt.Println("cmd error : lost:", lost, "minRtt:", minRtt, "avgRtt:", avgRtt, "maxRtt:", maxRtt, "mdevRtt:", mdevRtt, "duration:", duration)
		}
		fmt.Println("cmd success : lost:", lost, "minRtt:", minRtt, "avgRtt:", avgRtt, "maxRtt:", maxRtt, "mdevRtt:", mdevRtt, "duration:", duration)

		if avgRtt > 1000 {
			fmt.Println("waring waring waring waring waring waring")
		}
		time.Sleep(time.Second * 5)
	}
}

func cmdPingOnce(ip net.IP) (lost, minRtt, avgRtt, maxRtt, mdevRtt float64, duration string, err error) {
	args := []string{
		"-c", DefPktStr,
		"-s", DefPktSizeStr,
		"-w", DefTimeoutStr,
		"-i", DefWaitStr,
		"-q",
	}

	// 处理ip
	if GetIpFamily(ip, IpFamilyV4, IpFamilyV6, "") == IpFamilyV6 {
		args = append(args, "-6")
	}

	args = append(args, ip.String())
	retryCount := 3

retry:
	out, errOut, err := RunCommand("ping", args, DefTimeout*time.Second+5*time.Second) // ping超时退出后也会返回结果，防止没反回前被kill
	if err != nil {
		if len(errOut) > 0 {
			return 100, 0, 0, 0, 0, duration, fmt.Errorf("%s, stderr: %s", err, errOut)
		}
		lost, minRtt, avgRtt, maxRtt, mdevRtt, duration, err = parsePingOutput(string(out))
		if err != nil {
			return 100, 0, 0, 0, 0, duration, err
		}

		return lost, minRtt, avgRtt, maxRtt, mdevRtt, duration, err
	}
	lost, minRtt, avgRtt, maxRtt, mdevRtt, duration, _ = parsePingOutput(string(out))
	if lost < 0 && retryCount > 0 {
		// 丢包率小于0，ping echo id重复，小概率事件，简单修复做重试。后续在原生ping中优化重复id场景
		retryCount--
		fmt.Println("ping-retry", "ip", ip.String(), "lost", lost, "retry", 3-retryCount)
		goto retry
	}
	return lost, minRtt, avgRtt, maxRtt, mdevRtt, duration, nil
}

func parsePingOutput(output string) (float64, float64, float64, float64, float64, string, error) {
	var (
		packetLoss       float64
		minRTT           float64
		avgRTT           float64
		maxRTT           float64
		mdevRTT          float64
		duration         string
		invaliLossOutput = true
		invaliRTTOutput  = true
	)

	// 根据输出中的关键字进行分割和查找
	for i, line := range strings.Split(output, "\n") {
		if strings.Contains(line, "packet loss") {
			plIdx := strings.LastIndex(line, "% packet loss")
			fields := strings.Fields(line[:plIdx])
			if len(fields) > 5 {
				lost, pErr := strconv.ParseFloat(fields[len(fields)-1], 64)
				if pErr == nil {
					packetLoss = lost
					invaliLossOutput = false
				}
			}
		}

		fmt.Println("i", i, "line : ", line)

		if strings.Contains(line, "min/avg/max/mdev") {
			fields := strings.Split(line, " ")
			if len(fields) > 3 {
				rttValues := strings.Split(fields[3], "/")
				if len(rttValues) == 4 {
					minRtt, errMin := strconv.ParseFloat(rttValues[0], 64)
					avgRtt, errAvg := strconv.ParseFloat(rttValues[1], 64)
					maxRtt, errMax := strconv.ParseFloat(rttValues[2], 64)
					mdevRtt, errMdev := strconv.ParseFloat(rttValues[3], 64)

					fmt.Println("len rtt is 4 , min: %.3f ms, avg: %.3f ms, max: %.3f ms, mdev: %.3f ms\n", minRtt, avgRtt, maxRtt, mdevRtt)

					if errMin == nil && errAvg == nil && errMax == nil && errMdev == nil {
						minRTT = math.Round((minRtt/1000)*10000) / 10000
						avgRTT = math.Round((avgRtt/1000)*10000) / 10000
						maxRTT = math.Round((maxRtt/1000)*10000) / 10000
						mdevRTT = math.Round((mdevRtt/1000)*10000) / 10000
						//minRTT = minRtt / 1e3
						//avgRTT = avgRtt / 1e3
						//maxRTT = maxRtt / 1e3
						//mdevRTT = mdevRtt / 1e3
						invaliRTTOutput = false
						fmt.Println("min: %.3f ms, avg: %.3f ms, max: %.3f ms, mdev: %.3f ms\n", minRTT, avgRTT, maxRTT, mdevRTT)
					} else {
						fmt.Println("解析 RTT 值时出错")
					}
				}
			}
		} else {
			fmt.Println("iping-task cmd执行超时 - 152")
		}

		if strings.Contains(line, "time") {
			pos := strings.Index(line, "time")
			duration = line[pos+4:]
		}
	}

	if invaliRTTOutput || invaliLossOutput {
		return 100, 0, 0, 0, 0, duration, fmt.Errorf("invalid output: %s", output)
	}

	return packetLoss, minRTT, avgRTT, maxRTT, mdevRTT, duration, nil
}

func GetIpFamily[V any](ip net.IP, ipv4, ipv6, unknown V) V {
	if ip.To4() != nil {
		return ipv4
	}
	if ip.To4() == nil && ip.To16() != nil {
		return ipv6
	}
	return unknown
}

func RunCommand(name string, arg []string, timeout time.Duration) ([]byte, []byte, error) {
	var (
		stdout = bytes.NewBuffer(nil)
		stderr = bytes.NewBuffer(nil)
	)
	cmd := exec.Command(name, arg...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	// 强制 kill
	time.AfterFunc(timeout, func() {
		if cmd.Process != nil {
			_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		}
	})

	err := cmd.Run()
	return stdout.Bytes(), stderr.Bytes(), err
}

//func ExactDivision(num, divisor float64) float64 {
//	bigNum := big.NewFloat(num)
//	bigDivisor := big.NewFloat(divisor)
//	result := new(big.Float).Quo(bigNum, bigDivisor)
//	result.SetPrec(64).SetMode(big.ToNearestAway)
//	return result
//}
