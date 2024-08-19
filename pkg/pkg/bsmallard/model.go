package bsmallard

import "net/http"

// mallard option
type Option struct {
	Switch              bool            `json:"switch"`
	ReportUrl           string          `json:"report_url"`
	ReportFilePath      string          `json:"report_file_path"`
	ReportHttpTransport *http.Transport `json:"report_http_transport"`
}

// MetricRaw is metric value of
type MetricRaw struct {
	Name     string                 `json:"name"`
	Time     int64                  `json:"time"`
	Value    float64                `json:"value"`
	Fields   map[string]interface{} `json:"fields,omitempty"`
	Tags     map[string]string      `json:"tags,omitempty"`
	Step     int64                  `json:"step"`
	Endpoint string                 `json:"endpoint"`
}

// Output a map of the number of metric points in time
// used to view metric's statistical logs
func ProfMetricDetailInfo(metrics []*MetricRaw) map[string]map[int64]int64 {
	record := make(map[string]map[int64]int64)
	for _, m := range metrics {
		if _, ok := record[m.Name]; !ok {
			record[m.Name] = make(map[int64]int64)
		}
		record[m.Name][m.Time] = record[m.Name][m.Time] + 1
	}
	return record
}
