/**
 * Package Metric
 * @Author fengfeng.mei <fengfeng.mei@baishan.com>
 * @Date 2024/12/4 11:08
 */

package Metric

type (
	Metric struct {
		// metric名称，必须具有唯一性
		Name string `json:"name,omitempty"`
		// 发生时间，当使用的时候如果为空则会被填充为当前时间
		Time int64 `json:"time,omitempty"`
		// 关键数据，用于报警或最关注的重要数据
		Value float64 `json:"value,omitempty"`
		// 主要数据（值推荐是数值），不能使用tags中出现的字段为键名
		Fields map[string]interface{} `json:"fields,omitempty"`
		// 附加标签值（用来区分数据类型的），预定义的字段有sertypes，cachegroup，storagegroup，endpoint；除此之外可以自定义字段
		// map[tagName]value
		Tags map[string]string `json:"tags,omitempty"`
		// endpoint是主机名（也可以通过配置去设置）
		Endpoint string `json:"endpoint,omitempty"`
		// 采集周期，单位是秒
		Step int    `json:"step,omitempty"`
		From string `json:"-"`
		// metric 排序值，为 Time + Name + Endpoint 的字符串
		SortKey           string `json:"-"`
		TransferOpenTime  int64  `json:"transfer_open_time,omitempty"`
		TransferOpenDelay int64  `json:"transfer_open_delay,omitempty"`
		TransferTime      int64  `json:"transfer_time,omitempty"`
		TransferDelay     int64  `json:"transfer_delay,omitempty"`
		StoreTime         int64  `json:"store_time,omitempty"`
		StoreDelay        int64  `json:"store_delay,omitempty"`
		AgentTime         int64  `json:"agent_time,omitempty"`
		AgentDelay        int64  `json:"agent_delay,omitempty"`
	}

	MetricStat struct {
		From    string              `json:"from,omitempty"`
		Created int64               `json:"created,omitempty"`
		Updated int64               `json:"updated,omitempty"`
		Count   int64               `json:"count,omitempty"`
		Tags    map[string]struct{} `json:"tags,omitempty"`
		Fields  map[string]struct{} `json:"fields,omitempty"`
	}
)
