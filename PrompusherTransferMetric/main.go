/**
 * Package PrompusherTransferMetric
 * @Author fengfeng.mei <fengfeng.mei@baishan.com>
 * @Date 2024/11/21 16:12
 */

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"goStudy/PrompusherTransferMetric/promexport"
	"io/fs"
	"path/filepath"
	"strings"
)

func main() {
	err := filepath.Walk("/tmp/plug", func(fpath string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Println("walk error	1:", err)
			return err
		}
		if info.IsDir() {
			fmt.Println("walk error	2:", err)
			return err
		}
		if !isAllowFilename(info.Name()) {
			fmt.Println("walk error	3:", err)
			return err
		}
		logFile := filepath.Join("/tmp/plugLog", strings.Replace(filepath.ToSlash(fpath), "/", "_", -1)) + ".log"
		fmt.Println("exec file:", fpath)
		plugin, err := NewPlugin(fpath, logFile, info)
		dataBytes, err := plugin.Exec()
		fmt.Println("exec Done file:", plugin.File)
		// 执行错误
		if err != nil {
			fmt.Println("exec error,file:", plugin.File, " error:", err)
			return err
		}

		// 没有结果
		mLen := len(dataBytes)
		if mLen == 0 {
			fmt.Println("exec fruitless,file:", plugin.File)
			return err
		}

		// 处理后
		//fmt.Println(string(dataBytes))
		res := fn(plugin.File, dataBytes)
		if res.Error != nil {
			// 如果数据错误，记录一下
			if !res.IsZero {
				badData := dataBytes
				if len(badData) > 256 { // 限制长度
					badData = badData[:255]
				}
				plugin.writeLog(badData)
			}
			fmt.Println("exec done but res is err,file:", plugin.File, " error:", res.Error)
			return err
		}
		return nil
	})
	if err != nil {
		fmt.Println("err:", err)
		return
	}
}

func fn(file string, data []byte) ExecResult {
	metrics, err := tryToNewMetrics(data)
	if len(metrics) == 0 {
		fmt.Println("tryToNewMetrics")
		if metrics, err = fromOldMetrics(data); err != nil {
			return ExecResult{Error: err}
		}
	}
	if len(metrics) == 0 {
		return ExecResult{IsZero: true, Error: err}
	}

	queue := make([]interface{}, 0, len(metrics))
	for _, item := range metrics {
		queue = append(queue, item)
	}

	return ExecResult{Metrics: metrics}
}

func isAllowFilename(filename string) bool {
	if len(filename) == 0 {
		return false
	}
	firstByte := filename[0]
	if firstByte < 49 || firstByte > 57 {
		return false
	}
	ext := filepath.Ext(filename)
	if ext == ".py" || ext == ".sh" {
		return true
	}
	return false
}

func tryToNewMetrics(data []byte) ([]*promexport.Metric, error) {
	var metrics []*promexport.Metric
	if err := json.Unmarshal(data, &metrics); err != nil {
		return nil, err
	}
	for _, m := range metrics {
		if m.Name == "" {
			return nil, errors.New("metric-name-empty")
		}
	}
	return metrics, nil
}

func fromOldMetrics(data []byte) ([]*promexport.Metric, error) {
	var (
		metricOlds []*promexport.MetricRaw
		metrics    []*promexport.Metric
	)
	if err := json.Unmarshal(data, &metricOlds); err != nil {
		return nil, err
	}
	for _, old := range metricOlds {
		m, err := old.ToNew()
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, m)
	}
	return metrics, nil
}
