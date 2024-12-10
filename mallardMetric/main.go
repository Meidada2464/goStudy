/**
 * Package main
 * @Author fengfeng.mei <fengfeng.mei@baishan.com>
 * @Date 2024/12/4 11:03
 */

package main

import (
	"encoding/json"
	"fmt"
	"mallardMetric/Metric"
)

func main() {
	var ms []*Metric.Metric
	err := json.Unmarshal([]byte(`[{"name": "ares_iping_log","time": 1733292405,"fields": {"lost": 0,"max_rtt": 0.0016870000000000001,"mdev_rtt": 0.000201,"min_rtt": 0.0007970000000000001,"rtt": 0.0014830000000000002},"tags": {"cachegroup": "KH-PhnomPenh-PhnomPenh-10-cache-1","dst_asn": "38235","dst_country": "Cambodia","dst_ip": "116.212.132.30","dst_province": "Cambodia","sertypes": "cache|cache_cache|live|live_edge|parent|static_parent|xdp|xdp_director","src_country": "Cambodia","src_node": "KH-PhnomPenh-PhnomPenh-jd-10","src_province": "PhnomPenh"},"endpoint": "KH-PhnomPenh-PhnomPenh-10-203-144-91-230","agent_time": 1733292406088525}]`), &ms)
	if err != nil {
		fmt.Println("unmarshal error:", err)
		return
	}

	metrics := Metric.ConvertMetrics(ms)
	fmt.Println("metrics :", metrics)
}
