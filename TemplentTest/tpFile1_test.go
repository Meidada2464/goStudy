package TemplentTest

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"goStudy/zap"
	"net/http"
	"testing"
	"time"
)

func TestName1(t *testing.T) {
	http.HandleFunc("/", SayHello1)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		fmt.Println("HTTP server failed,err:", err)
		return
	}
}

func difference(slice1, slice2 []string) []string {
	// 创建一个map来存储第二个切片中的元素
	set := make(map[string]bool)

	// 遍历第二个切片，将元素添加到map中
	for _, val := range slice2 {
		set[val] = true
	}

	// 创建一个切片来存储差集
	var result []string

	// 遍历第一个切片，如果元素在map中不存在，则添加到差集中
	for _, val := range slice1 {
		if !set[val] {
			result = append(result, val)
		}
	}

	return result
}

func TestName2(t *testing.T) {
	// 示例用法
	slice1 := []string{"apple", "banana", "orange", "grape"}
	slice2 := []string{"banana", "orange", "kiwi", "pear"}

	result := difference(slice1, slice2)

	fmt.Println("差集:", result)
}

type (
	Option struct {
		Name     string        `yaml:"name"`     // 请求的方法名
		Address  string        `yaml:"address"`  // 接口地址
		Token    string        `yaml:"token"`    // 请求token
		Body     string        `yaml:"body"`     // 请求body
		Method   string        `yaml:"method"`   // 请求的类型，GET/POST
		Timeout  time.Duration `yaml:"timeout"`  // 超时时间(s)
		Interval time.Duration `yaml:"interval"` // 同步时间间隔(s)
	}

	// DLTarget 单个同步器
	DLTarget struct {
		opt    Option                                                   // 参数
		client *http.Client                                             // http连接池
		cb     func(data []byte, apiName string, elapsed time.Duration) // 更新后的回调函数
	}

	// DLManager 所有同步器，定时同步
	DLManager struct {
		log      zap.Logger // 日志，前缀sync
		stopSign chan int   // 停止信号
		targets  []DLTarget // 所有请求
	}
)

func TestName3(t *testing.T) {

	var op []Option
	op = append(op, Option{Name: "water", Address: "https://567d243f-8eaf-4c0e-b56a-5d7be5d687bd.mock.pstmn.io", Token: "", Body: "", Method: "POST", Timeout: 30000000000, Interval: 300000000000})

	NewDLManager(op, ParamUpdate)
}

func NewDLManager(options []Option, cb func(data []byte, apiName string, elapsed time.Duration)) *DLManager {
	targets := make([]DLTarget, 0, len(options))
	zap.Log().Info("NewDLManager", "options-NewDLManager", options)
	fmt.Println("options", options)

	for _, opt := range options {
		zap.Log().Info("NewDLManager", "opt-NewDLManager", opt)
		fmt.Println("opt", opt)
		targets = append(targets, DLTarget{
			opt: opt,
			client: &http.Client{
				Transport: &http.Transport{
					MaxIdleConns:          10,
					MaxIdleConnsPerHost:   5,
					ResponseHeaderTimeout: opt.Timeout,
					IdleConnTimeout:       opt.Interval * 3,
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: true,
					},
				},
				Timeout: opt.Timeout,
			},
			cb: cb,
		})
	}
	zap.Log().Info("NewDLManager", "options-targets", targets[0].opt.Address)
	fmt.Println("targets", targets[0].opt.Address)

	return &DLManager{
		log:      zap.Log().Prefix("dedicatedLineSync"),
		targets:  targets,
		stopSign: make(chan int),
	}
}

func ParamUpdate(data []byte, apiName string, elapsed time.Duration) {
	fmt.Println("aaaccc")
}

type (
	DeleteNetTaskInput struct {
		Mode         int                `json:"mode"`
		DelTask      DeleteNetTask      `json:"task"`
		DelChildTask DeleteNetChildTask `json:"child_task"`
	}

	DeleteNetTask struct {
		Id       int    `json:"id" gorm:"primary_key"`
		TaskName string `json:"task_name,omitempty"`
	}
	DeleteNetChildTask struct {
		TaskId        int           `json:"task_id"`
		Tag           int           `json:"tag"`
		Mode          int32         `json:"mode"`
		TaskName      string        `json:"task_name,omitempty"`
		TaskStatus    int32         `json:"task_status,omitempty"`
		ChildTaskName string        `json:"child_task_name,omitempty"`
		Monitors      []MachineInfo `json:"monitors"`
		Targets       []MachineInfo `json:"targets"`
	}

	MachineInfo struct {
		Type        int32  `json:"type,omitempty"`
		Rule        string `json:"rule,omitempty"`
		Limit       string `json:"limit,omitempty"`
		ServiceType string `json:"service_type,omitempty"`
	}
)

func TestName4(t *testing.T) {
	var test DeleteNetTaskInput
	var mo MachineInfo
	var ta MachineInfo

	mol := make([]MachineInfo, 0)
	tal := make([]MachineInfo, 0)

	mol = append(mol, mo)
	tal = append(tal, ta)

	test.DelChildTask.Monitors = mol
	test.DelChildTask.Targets = tal

	test.Mode = 1
	test.DelTask.TaskName = "dns_hijack"
	test.DelTask.Id = 3493
	test.DelChildTask.ChildTaskName = "lt_guangdong@s3plus-hl01.meituan.net"
	test.DelChildTask.TaskId = 3493
	test.DelChildTask.TaskName = "dns_hijack"
	test.DelChildTask.Mode = 1
	test.DelChildTask.TaskStatus = 1
	test.DelChildTask.Tag = 2
	test.DelChildTask.Monitors[0].Limit = ""
	test.DelChildTask.Monitors[0].Rule = "guangdong"
	test.DelChildTask.Monitors[0].Type = 1
	test.DelChildTask.Monitors[0].ServiceType = "cache_cache"

	marshal, err := json.Marshal(test)
	if err != nil {
		return
	}
	s := string(marshal)
	fmt.Println(s)
}
