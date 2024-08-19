package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {

	var data []map[string]interface{}

	// 要发送的数据
	data = append(data, map[string]interface{}{
		"endpoint": "",
		"name":     "mcdnsgcc_cdn_ns",
		"tags": map[string]string{
			"ip":     "210.77.178.32",
			"domain": "www.sgid.sgcc.com.cn",
			"statue": "NOERROR",
		},
		"fields": map[string]interface{}{
			"status": 0,
		},
		"value": 0,
		"step":  60,
		"time":  1688094180,
	})

	// 将数据编码为JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	// 创建HTTP客户端
	client := &http.Client{}

	// 创建POST请求
	req, err := http.NewRequest("POST", "https://mallard-transfer-open.bs58i.baishancdnx.com/open/metric?v2=1", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// 设置请求头
	req.Header.Set("Token-Value", "a3594dfe97ccab10d951ec11a0878211")
	req.Header.Set("Token-User", "detect-agent-box")
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应
	fmt.Println("Response Status:", resp.Status)
}
