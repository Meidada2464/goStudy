package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type (
	ResBox struct {
		Code int          `json:"code"`
		Msg  string       `json:"msg"`
		Data []BoxMachine `json:"data"`
	}

	BoxMachine struct {
		Sn         string `json:"sn"`
		Hostname   string `json:"hostname"`
		Country    string `json:"country_en"`
		Isp        string `json:"isp_en"`
		Province   string `json:"province_en"`
		City       string `json:"city_en"`
		SupplierID int32  `json:"supplier_id"`
	}
)

func main() {
	// 请求路径与处理函数
	http.HandleFunc("/api/snode/detector", handler)

	// 启动http服务器
	log.Println("http server start")
	err := http.ListenAndServe(":9999", nil)
	if err != nil {
		log.Println("http server start error:", err)
		return
	}
}

func handler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("request time:", time.Now().Unix())

	// 读取本地 JSON 文件
	file, err := os.Open("/Users/meifengfeng/workSpace/study/goStudy/goStudy/MockTree/moke_data.json")
	if err != nil {
		http.Error(w, "Error opening file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// 解析 JSON 文件
	var resBox ResBox
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&resBox)
	if err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	// 设置响应类型为JSON
	w.Header().Set("Content-Type", "application/json")

	// 将响应数据编码为 JSON 并发送
	json.NewEncoder(w).Encode(resBox)
}
