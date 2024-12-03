package promexport

type (
	Config struct {
		Switch    bool       `json:"switch"`
		Resources []Resource `json:"resources"`
	}

	Resource struct {
		ScrapeAim      string `json:"scrape_aim"`
		ScrapePrefix   string `json:"scrape_prefix"`
		ScrapeUrl      string `json:"scrape_url"`
		ScrapeInterval int64  `json:"scrape_interval"`
	}

	PromMetric struct {
		Name     string    `json:"name"`
		Time     int64     `json:"time"`
		Help     string    `json:"help,omitempty"`
		JobName  string    `json:"job_name"`
		Endpoint string    `json:"endpoint"`
		Type     string    `json:"type"`
		Metrics  []Metrics `json:"metrics"`
	}

	/*
	* prometheus 共5种数据类型：UNTYPED、COUNTER、GAUGE、SUMMARY、HISTOGRAM
	*	这里将 UNTYPED、COUNTER、GAUGE 归为默认类
	*	默认类型可能有：Labels、Value、TimestampMs
	*	SUMMARY可能有：Labels、TimestampMs、Quantiles、Count、Sum
	*	HISTOGRAM可能有：Labels、TimestampMs、Buckets、Count、Sum
	 */
	Metrics struct {
		Labels      map[string]string `json:"labels,omitempty"`
		TimestampMs string            `json:"timestamp_ms,omitempty"`
		Value       string            `json:"value,omitempty"`
		Quantiles   map[string]string `json:"quantiles,omitempty"`
		Buckets     map[string]string `json:"buckets,omitempty"`
		Count       string            `json:"count,omitempty"`
		Sum         string            `json:"sum,omitempty"`
	}
)

func (p *PromMetric) GetName() string {
	if p != nil && p.Name != "" {
		return p.Name
	}
	return ""
}

func (p *PromMetric) GetEndpoint() string {
	if p != nil && p.Endpoint != "" {
		return p.Endpoint
	}
	return ""
}

func (p *PromMetric) GetType() string {
	if p != nil && p.Type != "" {
		return p.Type
	}
	return ""
}

func (p *PromMetric) GetJobName() string {
	if p != nil && p.JobName != "" {
		return p.JobName
	}
	return ""
}
