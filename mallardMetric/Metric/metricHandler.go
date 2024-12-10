/**
 * Package Metric
 * @Author fengfeng.mei <fengfeng.mei@baishan.com>
 * @Date 2024/12/4 11:08
 */

package Metric

import (
	"fmt"
	"time"
)

func ConvertMetrics(values []*Metric) []*Metric {
	metrics := make([]*Metric, 0, len(values)*10)

	values = fillMetrics(values)
	fmt.Println("values:", values)
	metrics = append(metrics, values...)

	return metrics
}

func fillMetrics(metrics []*Metric) []*Metric {
	infos := make(map[string]*MetricStat, 10) // map[string]*MetricStat
	real := make([]*Metric, 0, len(metrics))
	now := time.Now().Unix()
	for i := range metrics {
		metric := metrics[i]
		if metric == nil {
			fmt.Println("nil-metric")
			continue
		}
		if metric.Name == "" {
			fmt.Println("no-name-endpoint", "metric", metric)
			continue
		}
		if metric.Endpoint == "" {
			metric.Endpoint = "ff_test"
		}
		if metric.Time == 0 {
			metric.Time = now
		}
		metric.FillTags("Sertypes",
			"Cachegroup",
			"StorageGroup",
			"Endpoint",
			"Sertags")
		real = append(real, metric)

		info := infos[metric.Name]
		if info == nil {
			info = &MetricStat{
				From:    metric.From,
				Created: metric.Time,
				Updated: metric.Time,
				Count:   1,
				Tags:    make(map[string]struct{}),
				Fields:  make(map[string]struct{}),
			}
			infos[metric.Name] = info
		}
		info.Updated = metric.Time
		info.Count++
		for tagname := range metric.Tags {
			info.Tags[tagname] = struct{}{}
		}
		for fieldName := range metric.Fields {
			info.Fields[fieldName] = struct{}{}
		}
	}
	return real
}

func (m *Metric) FillTags(sertypes, cachegroup, storagegroup, gendpoint, sertags string) {
	if len(m.Tags) == 0 {
		m.Tags = make(map[string]string)
	}
	if m.Tags["sertypes"] == "" && sertypes != "" {
		m.Tags["sertypes"] = sertypes
	}
	if m.Tags["cachegroup"] == "" && cachegroup != "" {
		m.Tags["cachegroup"] = cachegroup
	}
	if m.Tags["storagegroup"] == "" && storagegroup != "" {
		m.Tags["storagegroup"] = storagegroup
	}
	if m.Endpoint != gendpoint {
		m.Tags["agent_endpoint"] = gendpoint
	}
	if m.Tags["sertags"] == "" && sertags != "" {
		m.Tags["sertags"] = sertags
	}
}

//
//func Write(metrics []*Metric) error {
//
//	data, size := toByte(metrics)
//	if size == 0 {
//		return nil
//	}
//
//	n, err := sw.fd.Write(data)
//	if err != nil {
//		return err
//	}
//	sw.size += int64(n)
//	sw.box.Incr("logtool_write_size", int64(n))
//
//	// 兜底，确保文件被删了但句柄没有被释放
//	//if rand.Intn(50) == 1 {
//	//	go sw.redeem()
//	//}
//	return nil
//}
//
//func toByte(metrics []*Metric) ([]byte, int64) {
//	buf := bytes.NewBuffer(nil)
//	for _, m := range metrics {
//		b, err := json.Marshal(m)
//		if err != nil {
//			continue
//		}
//		buf.Write(b)
//		buf.WriteString("\n")
//	}
//	return buf.Bytes(), int64(buf.Len())
//}
