package LoopTest

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

//func Loop() {
//	var stopFlag *atomic.Bool
//
//	fmt.Println("Loop begin")
//	utils.TimeLoop(time.Minute, func(now time.Time) {
//		if stopFlag.Load() {
//			fmt.Println("stop")
//			return
//		}
//		// 收集指标
//		collectMetric(now)
//	})
//	time.Sleep(time.Second)
//}

func CollectMetric(nowTime time.Time) {

	fmt.Println("nowTime:", nowTime.Unix())

	var (
		agentPeriod  = nowTime.Unix()/60%int64(3) == 0
		serverPeriod = nowTime.Unix()/60%int64(1) == 0
	)

	//fmt.Println("agentPeriod:", nowTime.Unix()%int64(60))
	//fmt.Println("serverPeriod:", nowTime.Unix()%int64(120))

	if !agentPeriod && !serverPeriod {
		return
	}

	time.Sleep(time.Second)

	if agentPeriod {
		var now = time.Now().Unix()

		fmt.Println("nowTime-now:", now)
		fmt.Println("nowTime-agentPeriod:", nowTime.Unix())

	}

	if serverPeriod {
		var now = time.Now().Unix()
		fmt.Println("nowTime-now", now)
		fmt.Println("nowTime-serverPeriod:", nowTime.Unix())
	}

	time.Sleep(time.Second)

}

func GOTemp() {
	go func(time2 int) {
		fmt.Println("ingo：", time2)
	}(12)

	go fmt.Println("ingo2：")
	time.Sleep(time.Second)
}

func GOTemp2() {
	type T struct {
		I []int
	}
	go func() {

	}()
}

func ReadFile() {

	fileConter := map[string]interface{}{}

	file, err := os.ReadFile("hello.json")
	if err != nil {
		return
	}

	err2 := json.Unmarshal(file, &fileConter)
	if err2 != nil {
		return
	}

	fmt.Println("fileConter", fileConter)

}
