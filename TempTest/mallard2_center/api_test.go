package mallard2_center

import (
	"crypto/tls"
	"fmt"
	"goStudy/tool/utils"
	"io"
	"net/http"
	"testing"
	"time"
)

type (
	// Request 同步请求
	Request struct {
		// center 的地址
		CenterAddr string
		Hash       string
		// 请求超时时间（默认是10秒）
		Timeout time.Duration
		// http.Transport 用于缓存长连接，用于大量HTTP请求的连接复用
		Transport *http.Transport
	}
	// Response 同步请求结果
	Response struct {
		// 从返回头里的Content-Crc取出
		Hash string
		// 当HTTP响应状态码是304的时候为true，意味着有缓存
		NotModified bool
		// 具体的返回值，诸如： center.Endpoints， syncer.StorageQueryList 等等
		Value interface{}
		// 非预期响应时 Error 就不为空
		Error error
		// Value 的长度（针对具体的对象）
		Count int
		// 默认是同步函数名，如：endpoints, endpoints-info等等
		Keyname string
	}
	// SyncFunc 同步函数，发出请求获得结果
	SyncFunc func(Request) Response

	Storage2 struct {
		// 键是 metric 名，值是`存储方案id-数据库id(@hot)`的列表，被分到 Hot 和 Common
		// map[metric.name]*MetricToInstance
		Metrics map[string]*MetricToInstance `json:"metrics"`
		// 主要包含 field 和 tag 的定义
		// map[metric.name][]*MetricField
		MetricStructure map[string][]*MetricField `json:"metric_structure"`
		// 类型映射，每个 metric 的字段都会选择一个全局类型（ type_id ），对应到存储中应该将其映射到实际的类型
		// 数据来自 portal.metadata_type_mapping，第一个键是具体的存储类型（如clickhouse），第二个键是id，值是实际类型
		// map[metadata_type_mapping.db_type][metadata_type_mapping.id]metadata_type_mapping.true_type
		MetricTypes map[string]map[int]string `json:"metric_types"`
		// 键是存储方案id，值是 `存储方案id-数据库id(@hot)`
		// map[metadata_storage_plan.id][]`storagePlan.id-db.id(@hot)`
		PlanDBs map[int][]string `json:"plan_dbs"`
		// 键是 `存储方案id-数据库id(@hot)`，如果存储方案名称中包含了hot字样的话，就在键中加入@hot字符串
		// 并且带hot的 S2Instance.QuerySort 变为 1e9
		// map[`storage_plan.id-db.id(@hot)`]*S2Instance
		Instances map[string]*S2Instance `json:"instances"`
		// 用于 clickhouse 的结构，键是 metric 名称，值是用于 Clickhouse 的字段名和字段类型
		// map[metric.name][]*FieldInfo
		// TODO 成功上线后应删除
		MetricsStructure map[string][]*FieldInfo `json:"metrics_structure"`

		// 以上几个字段的JSON去进行crc32
		CRC uint32 `json:"crc"`

		// hash 桶变化的记录，记录了与旧的相比增加的和被删除的 hash 桶
		changes *S2Change
		// 第一个键是 metric id， 第二个键值存储方案id，具体值是这个存储方案对应的所有数据库
		// map[metric.id]map[storage_plan.id][]db.id
		planMetrics map[int]map[int][]int
	}

	MetricToInstance struct {
		// 值是 `存储方案id-数据库id@hot`
		// []`storage_plan.id-db_id@hot`
		Hot []string `json:"h"`
		// 值是 `存储方案id-数据库id`
		// []`storage_plan.id-db_id`
		Common []string `json:"c"`
	}

	MetricField struct {
		// 为 field 或 tag
		Category string `json:"category"`
		// 字段名
		Name string `json:"name"`
		// 字段类型id
		TypeID int `json:"type_id"`
	}

	S2Instance struct {
		ID         int    `json:"id"`
		DBName     string `json:"db_name"`
		DBURL      string `json:"dburl"`
		User       string `json:"user"`
		Password   string `json:"password"`
		DBType     int    `json:"db_type"`
		DBTypename string `json:"db_typename"`
		// 数据库的唯一标识 label
		Remark string `json:"remark"`
		// 用于匹配metric的label,判断是否应该存入这个数据库
		MatchLabel string `json:"match_label"`
		Status     bool   `json:"status"`
		PlanID     int    `json:"plan_id"`
		PlanName   string `json:"plan_name"`
		// 接收数据离当前的截止时间，单位秒，默认7200秒即2小时
		Expire int64 `json:"expire"`
		// 数据库排序值，只有当 metadata_storage_plan_db.status 可用时才为 metadata_storage_plan.sort
		QuerySort int `json:"query_sort"`
	}

	FieldInfo struct {
		FieldName string `json:"fn"`
		FieldType string `json:"ft"`
	}

	S2Change struct {
		// 新添加的 hash 桶
		Inserts []S2ChangeItem `json:"inserts"`
		// 旧的被删除的 hash 桶
		Deletes []S2ChangeItem `json:"deletes"`
	}

	// S2ChangeItem 数据库 hash 变化
	S2ChangeItem struct {
		PlanID   int `json:"plan_id"`
		MetricID int `json:"metric_id"`
		DbID     int `json:"db_id"`
	}
)

func handleRequest(sReq Request, req *http.Request) (*http.Response, error) {
	client := &http.Client{
		Transport: sReq.Transport,
		Timeout:   sReq.Timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		return nil, fmt.Errorf("bad-status-%d:%s", resp.StatusCode, string(body))
	}
	return resp, nil
}

func ReqStorage2(req Request) (res Response) {
	r, err := http.NewRequest("GET", req.CenterAddr+"/api/storage2?gzip=1&crc="+req.Hash, nil)
	if err != nil {
		res.Error = err
		return
	}
	resp, err := handleRequest(req, r)
	if err != nil {
		res.Error = err
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == 304 {
		res.NotModified = true
		return
	}
	var value Storage2
	if err := utils.UnGzipJSON(resp.Body, &value); err != nil {
		res.Error = err
		return
	}
	// 运算一下内部数据
	res.Value = value
	res.Hash = resp.Header.Get("Content-Crc")
	res.Count = len(value.Metrics) + len(value.Instances)
	return
}

func ReqStorage(req Request) (res Response) {
	r, err := http.NewRequest("GET", req.CenterAddr+"/api/storage?gzip=1&crc="+req.Hash, nil)
	if err != nil {
		res.Error = err
		return
	}
	resp, err := handleRequest(req, r)
	if err != nil {
		res.Error = err
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == 304 {
		res.NotModified = true
		return
	}
	var value Storage2
	if err := utils.UnGzipJSON(resp.Body, &value); err != nil {
		res.Error = err
		return
	}
	return
}

func Test1(t *testing.T) {
	req := Request{
		CenterAddr: "https://mallard-center1.bs58i.baishancdnx.com",
		Hash:       "1209644030",
		Timeout:    time.Second * time.Duration(60),
		Transport: &http.Transport{
			MaxIdleConns:    10,
			MaxConnsPerHost: 5,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	res := ReqStorage2(req)
	fmt.Println(res)

	//storage := ReqStorage(req)
	//fmt.Println(storage)

}
