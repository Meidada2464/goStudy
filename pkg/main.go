package main

import (
	"collecter/pkg/bsmallard"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	doris_query          = "https://test:test@mallard-query.bs58i.baishancdnx.com/api/doris/query"
	lost_params_template = `{"metric": "arespinglog", "sql": "SELECT avg(lost) as avg_lost, sProv, d2Node FROM arespinglog a WHERE taskName = 'client_ping_idc_edge' and timestamp >= '%s' and timestamp <= '%s' GROUP BY sProv, d2Node"}`
)

type Data struct {
	SrcProvince string `json:"src_province"`
	DstNode     string `json:"dst_node"`
	AvgLost     string `json:"avg_lost"`
}

type Result struct {
	Data []Data `json:"data"`
}

func main() {
	ticker := time.NewTicker(5 * time.Minute)
	for {
		fmt.Println("Start collecting data...")
		if err := collectData(); err != nil {
			fmt.Println("Error collecting data:", err)
			continue
		}

		fmt.Println("Finish collecting data...")
		<-ticker.C
	}
}

func collectData() error {
	endTime := time.Now()
	startTimeFmt := endTime.Add(-5 * time.Minute).Format("2006-01-02 15:04:05")
	endTimeFmt := endTime.Format("2006-01-02 15:04:05")
	params := fmt.Sprintf(lost_params_template, startTimeFmt, endTimeFmt)

	result, err := getLostData(doris_query, params)
	if err != nil {
		fmt.Println("Error getting SN data:", err)
		return err
	}

	data, err := classifyData(result)

	fmt.Println("data-collocted:", data)

	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	if data == nil {
		return errors.New("no data to push")
	}

	return pushDataToBSMallard(data)
}

func getLostData(url, params string) (*Result, error) {
	resp, err := http.Post(url, "application/json", strings.NewReader(params))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result Result
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func classifyData(result *Result) ([]*bsmallard.MetricRaw, error) {
	time := time.Now().Unix()
	var metric []*bsmallard.MetricRaw
	for _, data := range result.Data {
		s, err := strconv.ParseFloat(data.AvgLost, 32)
		if err != nil {
			return nil, err
		}
		metric = append(metric, &bsmallard.MetricRaw{
			Name: "idc_ping_idc_edge",
			Tags: map[string]string{
				"sProv":  data.SrcProvince,
				"d2Node": data.DstNode,
			},
			Value:    s,
			Time:     time,
			Endpoint: "",
			Step:     5 * 60,
		})
	}

	return metric, nil
}

func pushDataToBSMallard(metric []*bsmallard.MetricRaw) error {
	mClient, err := bsmallard.MustNewClient(&bsmallard.Option{})
	if err != nil {
		return err
	}
	// 推送上报数据
	if err := mClient.Report(metric); err != nil {
		return err
	}

	return nil
}
