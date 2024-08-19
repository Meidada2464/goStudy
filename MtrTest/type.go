package main

type (
	Params struct {
		MaxHops    int `json:"max_hops,omitempty"`    // 最大跳数
		TimeoutMs  int `json:"timeout_ms,omitempty"`  // 超时时间
		PacketSize int `json:"packet_size,omitempty"` // 包大小
		SntSize    int `json:"snt_size,omitempty"`    // 包数量
		SntLimit   int `json:"snt_limit,omitempty"`   //重试次数
	}

	DetectMtrHopDetail struct {
		//标签类信息
		TTL     int    //ttl标识
		AsSign  string //供应商标识
		Address string //目标ip

		//数值类信息Ï
		Loss  float64 //丢包情况
		Snt   int
		Last  float64
		Avg   float64
		Best  float64
		Wrst  float64
		StDev float64
	}

	mtrData struct {
		Tag  map[string]string
		Fild map[string]interface{}
	}
)

const (
	//命令执行超时时间
	execCmdTime = 30

	//解析的每行数据数组长度
	lineDataMinCnt = 9

	//特殊解析下标
	specialAsIndex = 1
	//分割下标
	splitIndex = 8
)
