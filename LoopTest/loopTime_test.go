package LoopTest

import (
	"fmt"
	"goStudy/tool/utils"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

//
//func TestLoop(t *testing.T) {
//	Loop()
//}
//
//func TestGOTemp(t *testing.T) {
//	GOTemp()
//}
//
//func TestGOTemp2(t *testing.T) {
//	TimeT(time.Now())
//}
//
//func TimeT(nowTime time.Time) {
//	var (
//		agentPeriod  = nowTime.Unix()/60%int64(2) == 0
//		serverPeriod = nowTime.Unix() / 60
//	)
//
//	fmt.Println("agentPeriod", agentPeriod)
//	fmt.Println("agentPeriod", serverPeriod)
//
//}

func TestReadFile(t *testing.T) {
	ReadFile()
}

func TestCollectMetric(t *testing.T) {
	utils.TimeLoop(time.Minute, func(now time.Time) {
		fmt.Println("now", now.Unix())
		// 收集指标
		CollectMetric(now)
	})

	WaitKill()
}

// WaitSignals wait signals capture
func WaitSignals(signals ...os.Signal) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, signals...)
	<-c
}

// WaitKill is short name for WaitInterrupt,
// capture kill or interrupt signal
func WaitKill() {
	WaitSignals(os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
}
