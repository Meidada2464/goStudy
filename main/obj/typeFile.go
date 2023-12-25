package obj

type (
	Machine struct {
		Name       string // 主机名
		Prov       string // 省份
		City       string // 城市
		Status     string // 在itil上的设备状态2101|2007|2201
		Node       string // 所属节点
		CacheGroup string // 所属cache组
		ServerType string // 服务类型
		IP         string // ip
		ISP        string // isp
		RealISP    string
	}
	Targets struct {
		Machine *Machine // 被探测的设备
		Task    *NetTask
	}
	NetTask struct {
		Control *NetControl // 相关探测参数
		Type    string      // 等同于任务名
	}
	NetControl struct {
		TimeToAggr        int    `json:"time_to_aggr"`        // 几次后聚合，目前上层下来都0，默认会是6次
		TaskName          string `json:"task_name"`           // 任务名
		SurvivalTime      int32  `json:"survival_time"`       // 存活时间
		IpType            string `json:"ip_type"`             // 探测IP类型：ipv4-ipv4
		DetectServiceType string `json:"detect_service_type"` // 探测服务类型  没用到
		DetectProtocol    string `json:"detect_protocol"`     // 探测协议：ping
		Params            string `json:"params"`              // 参数
	}
)
