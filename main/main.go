package main

import (
	"encoding/json"
	"fmt"
	"goStudy/main/obj"
	"goStudy/queue"
	"io/ioutil"
	"time"
)

const (
	ERROR           = 3
	PRIVILEGEDFLAG  = "0000003fffffffff"
	PRIVILEGEDFLAG2 = "000001ffffffffff"
	SNODEVERSION    = "2.8.2-1.el7.243"
	UPDATE_VERSION  = "2.8.2-1.el7.243"
)

func main() {
	//quantity := minQuantity(10000)
	//fmt.Println("quantity", quantity)
	newQue := queue.NewQueue(10000)
	go inputInQue(newQue)
	go getForQue(newQue)
}

func inputInQue(q *queue.EsQueue) {
	// 写数据
	for i := 0; i < 20000; i++ {
		ok, quantity := q.Put(i)
		fmt.Println("input que status :", ok, "quantity:", quantity, "content:", i)
		time.Sleep(5 * time.Millisecond)
	}
}

func getForQue(q *queue.EsQueue) {
	// 读数据
	for i := 0; i < 20000; i++ {
		val, ok, quantity := q.Get()
		fmt.Println("get que status :", ok, "quantity:", quantity, "content:", val)
		time.Sleep(10 * time.Millisecond)
	}
}

func minQuantity(v uint32) uint32 {
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v++
	return v
}

type (
	Targets struct {
		Machine *obj.Machine // 被探测的设备
		Task    *NetTask
	}

	TargetsForCurl struct {
		Machine *obj.Machine // 被探测的设备
		Task    *NetTaskForCurl
	}

	TargetsForTcp struct {
		Machine *obj.Machine // 被探测的设备
		Task    *NetTaskForTcp
	}

	NetTask struct {
		Control *NetControl // 相关探测参数
		Type    string      // 等同于任务名
	}

	NetTaskForCurl struct {
		Control *curlControl // 相关探测参数
		Type    string       // 等同于任务名
	}

	NetTaskForTcp struct {
		Control *TcpControl // 相关探测参数
		Type    string      // 等同于任务名
	}

	NetControl struct {
		TimeToAggr        int     `json:"time_to_aggr"`        // 几次后聚合，目前上层下来都0，默认会是6次
		TaskName          string  `json:"task_name"`           // 任务名
		SurvivalTime      int32   `json:"survival_time"`       // 存活时间
		IpType            string  `json:"ip_type"`             // 探测IP类型：ipv4-ipv4
		DetectServiceType string  `json:"detect_service_type"` // 探测服务类型  没用到
		DetectProtocol    string  `json:"detect_protocol"`     // 探测协议：ping
		Params            Params2 `json:"params"`              // 参数
	}

	curlControl struct {
		TimeToAggr        int       `json:"time_to_aggr"`        // 几次后聚合，目前上层下来都0，默认会是6次
		TaskName          string    `json:"task_name"`           // 任务名
		SurvivalTime      int32     `json:"survival_time"`       // 存活时间
		IpType            string    `json:"ip_type"`             // 探测IP类型：ipv4-ipv4
		DetectServiceType string    `json:"detect_service_type"` // 探测服务类型  没用到
		DetectProtocol    string    `json:"detect_protocol"`     // 探测协议：ping
		CurlParam         CurlParam `json:"params"`              // 参数
	}

	TcpControl struct {
		TimeToAggr        int       `json:"time_to_aggr"`        // 几次后聚合，目前上层下来都0，默认会是6次
		TaskName          string    `json:"task_name"`           // 任务名
		SurvivalTime      int32     `json:"survival_time"`       // 存活时间
		IpType            string    `json:"ip_type"`             // 探测IP类型：ipv4-ipv4
		DetectServiceType string    `json:"detect_service_type"` // 探测服务类型  没用到
		DetectProtocol    string    `json:"detect_protocol"`     // 探测协议：ping
		TcpParams         TcpParams `json:"params"`              // 参数
	}

	Params2 struct {
		DetectFrequency int  `json:"detect_frequency"`
		PackageSize     int  `json:"package_size"`
		PackageNum      int  `json:"package_num"`
		IsNeedMtr       bool `json:"is_need_mtr"`
		Threshold       int  `json:"threshold"`
		MaxHop          int  `json:"max_hop"`
	}

	CurlParam struct {
		DetectFrequency int    `json:"detect_frequency"`
		FileType        int    `json:"file_type"`
		FileName        string `json:"file_name"`
		HTTPPort        int    `json:"http_port"`
		HTTPType        string `json:"http_type"`
	}

	TcpParams struct {
		TargetPort      int    `json:"target_port"`
		DetectFrequency int    `json:"detect_frequency"`
		RequestType     string `json:"request_type"`
		RequestContent  string `json:"request_content"`
		ResponseType    string `json:"response_type"`
		ResponseContent string `json:"response_content"`
	}
)

func ReadDataFromFile() []string {
	// 解决格式装换时的string问题
	file, err := ioutil.ReadFile("/Users/meifengfeng/workSpace/study/goStudy/goStudy/main/pingTestData1.json")
	if err != nil {
		fmt.Println("无法读取文件1:", err)
		return nil
	}

	var targets []Targets

	err = json.Unmarshal(file, &targets)
	if err != nil {
		fmt.Println("解析JSON时出错1:", err)
		return nil
	}

	var temp []string
	for _, target := range targets {
		paramsByte, err := json.Marshal(target.Task.Control.Params)
		paramsStr := string(paramsByte)
		if err != nil {
			return nil
		}
		tempInner := obj.Targets{}
		Machine := obj.Machine{}
		NT := obj.NetTask{}
		Control := obj.NetControl{}

		Machine = *target.Machine

		Control.Params = paramsStr
		Control.SurvivalTime = target.Task.Control.SurvivalTime
		Control.TaskName = target.Task.Control.TaskName
		Control.DetectProtocol = target.Task.Control.DetectProtocol
		Control.IpType = target.Task.Control.IpType
		Control.DetectServiceType = target.Task.Control.DetectServiceType
		Control.TimeToAggr = target.Task.Control.TimeToAggr

		NT.Control = &Control
		NT.Type = target.Task.Type

		tempInner.Machine = &Machine
		tempInner.Task = &NT

		strv2, err := json.Marshal(tempInner)
		if err != nil {
			return nil
		}
		temp = append(temp, string(strv2))
	}

	fmt.Println("temp", temp)
	resMarshal, err := json.Marshal(temp)
	if err != nil {
		fmt.Println("无法读取文件3:", err)
		return nil
	}

	err = ioutil.WriteFile("/Users/meifengfeng/workSpace/study/goStudy/goStudy/main/pingTestData.json", resMarshal, 0644)
	if err != nil {
		fmt.Println("无法写入文件:", err)
		return nil
	}
	// 读取文件内容

	fileContent, err := ioutil.ReadFile("/Users/meifengfeng/workSpace/study/goStudy/goStudy/main/pingTestData.json")
	if err != nil {
		fmt.Println("无法读取文件2:", err)
		return nil
	}

	// 定义一个切片用于存储解析后的数据
	var jsonStrings []string

	// 解析 JSON 数据
	err = json.Unmarshal(fileContent, &jsonStrings)
	if err != nil {
		fmt.Println("解析JSON时出错2:", err)
		return nil
	}
	return jsonStrings
}

func ReadDataFromFileWithCurl() []string {
	// 解决格式装换时的string问题
	file, err := ioutil.ReadFile("/Users/meifengfeng/workSpace/study/goStudy/goStudy/main/curlTestData1.json")
	if err != nil {
		fmt.Println("无法读取文件1:", err)
		return nil
	}

	var targets []TargetsForCurl

	err = json.Unmarshal(file, &targets)
	if err != nil {
		fmt.Println("解析JSON时出错1:", err)
		return nil
	}

	var temp []string
	for _, target := range targets {
		paramsByte, err := json.Marshal(target.Task.Control.CurlParam)
		paramsStr := string(paramsByte)
		if err != nil {
			return nil
		}
		tempInner := obj.Targets{}
		Machine := obj.Machine{}
		NT := obj.NetTask{}
		Control := obj.NetControl{}

		Machine = *target.Machine

		Control.Params = paramsStr
		Control.SurvivalTime = target.Task.Control.SurvivalTime
		Control.TaskName = target.Task.Control.TaskName
		Control.DetectProtocol = target.Task.Control.DetectProtocol
		Control.IpType = target.Task.Control.IpType
		Control.DetectServiceType = target.Task.Control.DetectServiceType
		Control.TimeToAggr = target.Task.Control.TimeToAggr

		NT.Control = &Control
		NT.Type = target.Task.Type

		tempInner.Machine = &Machine
		tempInner.Task = &NT

		strv2, err := json.Marshal(tempInner)
		if err != nil {
			return nil
		}
		temp = append(temp, string(strv2))
	}

	fmt.Println("temp", temp)
	resMarshal, err := json.Marshal(temp)
	if err != nil {
		fmt.Println("无法读取文件3:", err)
		return nil
	}

	err = ioutil.WriteFile("/Users/meifengfeng/workSpace/study/goStudy/goStudy/main/curlTestData.json", resMarshal, 0644)
	if err != nil {
		fmt.Println("无法写入文件:", err)
		return nil
	}
	// 读取文件内容

	fileContent, err := ioutil.ReadFile("/Users/meifengfeng/workSpace/study/goStudy/goStudy/main/curlTestData.json")
	if err != nil {
		fmt.Println("无法读取文件2:", err)
		return nil
	}

	// 定义一个切片用于存储解析后的数据
	var jsonStrings []string

	// 解析 JSON 数据
	err = json.Unmarshal(fileContent, &jsonStrings)
	if err != nil {
		fmt.Println("解析JSON时出错2:", err)
		return nil
	}
	return jsonStrings
}

func ReadDataFromFileWithTcp(strFile string) []string {
	// 解决格式装换时的string问题
	file, err := ioutil.ReadFile("/Users/meifengfeng/workSpace/study/goStudy/goStudy/main/tcpData/tcpTestData1.json")
	if err != nil {
		fmt.Println("无法读取文件1:", err)
		return nil
	}

	var targets []TargetsForTcp

	err = json.Unmarshal(file, &targets)
	if err != nil {
		fmt.Println("解析JSON时出错1:", err)
		return nil
	}

	var temp []string
	for _, target := range targets {
		paramsByte, err := json.Marshal(target.Task.Control.TcpParams)
		paramsStr := string(paramsByte)
		if err != nil {
			return nil
		}
		tempInner := obj.Targets{}
		Machine := obj.Machine{}
		NT := obj.NetTask{}
		Control := obj.NetControl{}

		Machine = *target.Machine

		Control.Params = paramsStr
		Control.SurvivalTime = target.Task.Control.SurvivalTime
		Control.TaskName = target.Task.Control.TaskName
		Control.DetectProtocol = target.Task.Control.DetectProtocol
		Control.IpType = target.Task.Control.IpType
		Control.DetectServiceType = target.Task.Control.DetectServiceType
		Control.TimeToAggr = target.Task.Control.TimeToAggr

		NT.Control = &Control
		NT.Type = target.Task.Type

		tempInner.Machine = &Machine
		tempInner.Task = &NT

		strv2, err := json.Marshal(tempInner)
		if err != nil {
			return nil
		}
		temp = append(temp, string(strv2))
	}

	resMarshal, err := json.Marshal(temp)
	if err != nil {
		fmt.Println("无法读取文件3:", err)
		return nil
	}

	err = ioutil.WriteFile("/Users/meifengfeng/workSpace/study/goStudy/goStudy/main/tcpTestData.json", resMarshal, 0644)
	if err != nil {
		fmt.Println("无法写入文件:", err)
		return nil
	}
	// 读取文件内容

	fileContent, err := ioutil.ReadFile("/Users/meifengfeng/workSpace/study/goStudy/goStudy/main/tcpTestData.json")
	//fileContent, err := ioutil.ReadFile(strFile)

	if err != nil {
		fmt.Println("无法读取文件2:", err)
		return nil
	}

	// 定义一个切片用于存储解析后的数据
	var jsonStrings []string

	// 解析 JSON 数据
	err = json.Unmarshal(fileContent, &jsonStrings)
	if err != nil {
		fmt.Println("解析JSON时出错2:", err)
		return nil
	}
	return jsonStrings
}
