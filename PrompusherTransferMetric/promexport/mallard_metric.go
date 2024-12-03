package promexport

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
)

// Metric is one value of monitoring metric
type Metric struct {
	Name     string                 `json:"name,omitempty"`
	Time     int64                  `json:"time,omitempty"`
	Value    float64                `json:"value,omitempty"`
	Fields   map[string]interface{} `json:"fields,omitempty"`
	Tags     map[string]string      `json:"tags,omitempty"`
	Endpoint string                 `json:"endpoint,omitempty"`
	Step     int                    `json:"step,omitempty"`
	From     string                 `json:"-"`
	SortKey  string                 `json:"-"`
}

// Hash generates unique hash of the metric
func (m *Metric) Hash(withEndpoint bool) string {
	hash0 := MD5String(m.Name)
	hash1 := MD5String(m.TagString(true))
	if withEndpoint {
		return hash0[:4] + hash1[:28] + MD5String(m.Endpoint)[:8]
	}
	return hash0[:12] + hash1[:28]
}

// FieldValue gets field value by name
func (m *Metric) FieldValue(name string) (float64, error) {
	if name == "value" {
		return m.Value, nil
	}
	fieldv, ok := m.Fields[name]
	if !ok {
		return 0, fmt.Errorf("field '%s' is missing", name)
	}
	return ToFloat(fieldv)
}

var metricHashSkipTag = map[string]bool{
	"sertypes":       true,
	"cachegroup":     true,
	"storagegroup":   true,
	"hang_status":    true,
	"use_status":     true,
	"fault_status":   true,
	"agent_endpoint": true,
}

// TagString returns tags and endpoint string as keyword for the metric
func (m *Metric) TagString(withEndpoint bool) string {
	str := make([]string, 0, len(m.Tags)+2)
	str = append(str, "name="+m.Name)
	if withEndpoint {
		str = append(str, "endpoint="+m.Endpoint)
	}
	for tag, value := range m.Tags {
		if metricHashSkipTag[tag] {
			continue
		}
		str = append(str, tag+"="+value)
	}
	sort.Sort(sort.StringSlice(str))
	return strings.Join(str, ",")
}

// String prints memory friendly
func (m Metric) String() string {
	s, err := json.Marshal(m)
	if err != nil {
		return fmt.Sprintf("#ERROR#{{{ %s %#v}}}", err.Error(), m)
	}
	return string(s)
}

// FillTags fill addon tag values
func (m *Metric) FillTags(sertypes, cachegroup, storagegroup, gendpoint string) {
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
}

// FullTags return all tags with serv and endpoint data
func (m *Metric) FullTags() map[string]string {
	fullTags := make(map[string]string, len(m.Tags)+1)
	for k, v := range m.Tags {
		fullTags[k] = v
	}
	fullTags["endpoint"] = m.Endpoint
	return fullTags
}

// ExtractTags extracts tag string to map
func ExtractTags(s string) (map[string]string, error) {
	if s == "" {
		return nil, nil
	}
	tags := make(map[string]string, strings.Count(s, ","))
	tagSlice := strings.Split(s, ",")
	for _, tag := range tagSlice {
		pair := strings.SplitN(tag, "=", 2)
		if len(pair) != 2 {
			return nil, fmt.Errorf("bad tag %s", tag)
		}
		k := strings.TrimSpace(pair[0])
		v := strings.TrimSpace(pair[1])
		tags[k] = v

	}
	return tags, nil
}

// ExtractFields extracts field string to map
func ExtractFields(s string) (map[string]interface{}, error) {
	if s == "" {
		return nil, nil
	}
	fields := make(map[string]interface{}, strings.Count(s, ","))
	fieldSlice := strings.Split(s, ",")
	for _, field := range fieldSlice {
		pair := strings.SplitN(field, "=", 2)
		if len(pair) != 2 {
			return nil, fmt.Errorf("bad-field-'%s'-in-'%s'", field, s)
		}
		k := strings.TrimSpace(pair[0])
		pair[1] = strings.TrimSpace(pair[1])
		vv, err := strconv.ParseFloat(pair[1], 64)
		if err != nil {
			err = nil
			vs, err := strconv.Unquote(pair[1])
			if err == nil {
				fields[k] = vs
			} else {
				fields[k] = pair[1]
			}
			continue
		}
		if math.IsNaN(vv) || math.IsInf(vv, 0) {
			return nil, fmt.Errorf("bad-field-value-'%s'-in-'%s'", field, s)
		}
		fields[k] = vv
		continue
	}
	return fields, nil
}

// MetricRaw is old metric data struct
type MetricRaw struct {
	Endpoint  string      `json:"endpoint"`
	Metric    string      `json:"metric"`
	Value     interface{} `json:"value"`
	Fields    string      `json:"fields"`
	Step      int         `json:"step"`
	Type      string      `json:"counterType"`
	Tags      string      `json:"tags"`
	Timestamp int64       `json:"timestamp"`
}

// ToNew converts old Metric to new Metric object
func (m *MetricRaw) ToNew() (*Metric, error) {
	floatValue, err := ToFloat(m.Value)
	if err != nil {
		return nil, err
	}
	mc := &Metric{
		Name:     m.Metric,
		Time:     m.Timestamp,
		Step:     m.Step,
		Value:    floatValue,
		Endpoint: m.Endpoint,
	}
	if mc.Fields, err = ExtractFields(m.Fields); err != nil {
		return nil, err
	}
	if mc.Tags, err = ExtractTags(m.Tags); err != nil {
		return nil, err
	}
	return mc, nil
}
