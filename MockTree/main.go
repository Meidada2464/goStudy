package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	treeData = "/Users/meifengfeng/workSpace/study/goStudy/goStudy/MockTree/moke_data.json"
	fogData  = "/Users/meifengfeng/workSpace/study/goStudy/goStudy/MockTree/moke_fog_data.json"
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

	FogData struct {
		Code      int    `json:"code"`
		RequestID string `json:"request_id"`
		Data      []struct {
			ID              int    `json:"id"`
			NodeID          int    `json:"node_id"`
			Hostname        string `json:"hostname"`
			Sn              string `json:"sn"`
			ManagementIP    string `json:"management_ip"`
			UseStatus       int    `json:"use_status"`
			FaultStatus     int    `json:"fault_status"`
			ProvinceID      int    `json:"province_id"`
			CityID          int    `json:"city_id"`
			SupplierID      int    `json:"supplier_id"`
			FaultType       int    `json:"fault_type"`
			FaultRemark     string `json:"fault_remark"`
			HTTPPort        string `json:"http_port"`
			HTTPSPort       string `json:"https_port"`
			PutOnTime       string `json:"put_on_time"`
			PutOffTime      string `json:"put_off_time"`
			PutOffRemark    string `json:"put_off_remark"`
			IsIntranet      int    `json:"is_intranet"`
			RecruitSn       string `json:"recruit_sn"`
			BindIP          string `json:"bind_ip"`
			NetType         int    `json:"net_type"`
			ChargeWay       int    `json:"charge_way"`
			BusinessType    string `json:"business_type"`
			DeployBusiness  string `json:"deploy_business"`
			Source          int    `json:"source"`
			FogSourceType   int    `json:"fog_source_type"`
			IPV6            string `json:"ip_v6"`
			Ipv4Type        int    `json:"ipv4_type"`
			Ipv6Type        int    `json:"ipv6_type"`
			FogOnlineStatus int    `json:"fog_online_status"`
			Supplier        struct {
				ID    int    `json:"id"`
				Name  string `json:"name"`
				SynID int    `json:"syn_id"`
			} `json:"supplier"`
			City struct {
				ID    int    `json:"id"`
				Name  string `json:"name"`
				Ename string `json:"ename"`
				Code  string `json:"code"`
			} `json:"city"`
			Province struct {
				ID        int    `json:"id"`
				Name      string `json:"name"`
				Ename     string `json:"ename"`
				Code      string `json:"code"`
				RegionID  int    `json:"region_id"`
				CountryID int    `json:"country_id"`
				Region    struct {
					ID    int    `json:"id"`
					Name  string `json:"name"`
					Ename string `json:"ename"`
				} `json:"region"`
				Country struct {
					ID    int    `json:"id"`
					Name  string `json:"name"`
					Ename string `json:"ename"`
					Code  string `json:"code"`
					Area  string `json:"area"`
				} `json:"country"`
			} `json:"province"`
			Node struct {
				ID    int    `json:"id"`
				Name  string `json:"name"`
				Ename string `json:"ename"`
			} `json:"node"`
			Svcs []struct {
				Name           string `json:"name"`
				SvrID          int    `json:"svr_id"`
				TypeModule     int    `json:"type_module"`
				FaultStatus    int    `json:"fault_status"`
				Remark         any    `json:"remark"`
				FaultTime      int    `json:"fault_time"`
				TypeModuleName string `json:"type_module_name"`
			} `json:"svcs"`
			ServiceRestrictionType int `json:"service_restriction_type"`
			Accounts               []struct {
				ID              int    `json:"id"`
				AccountName     string `json:"account_name"`
				CardName        string `json:"card_name"`
				IP              string `json:"ip"`
				IPStatus        int    `json:"ip_status"`
				AccountStatus   int    `json:"account_status"`
				FaultType       int    `json:"fault_type"`
				FaultRemark     string `json:"fault_remark"`
				RedialTime      int    `json:"redial_time"`
				SupplierLimit   int    `json:"supplier_limit"`
				MaxLimit        int    `json:"max_limit"`
				Plan            int    `json:"plan"`
				SvrID           int    `json:"svr_id"`
				InnerIP         string `json:"inner_ip"`
				InnerIpv6       string `json:"inner_ipv6"`
				HTTPPort        int    `json:"http_port"`
				HTTPSPort       int    `json:"https_port"`
				UpdateTime      string `json:"update_time"`
				V4HTTPPort      int    `json:"v4_http_port"`
				V4HTTPSPort     int    `json:"v4_https_port"`
				V6HTTPPort      int    `json:"v6_http_port"`
				V6HTTPSPort     int    `json:"v6_https_port"`
				TCPNatType      int    `json:"tcp_nat_type"`
				UDPNatType      int    `json:"udp_nat_type"`
				UpnpAvailable   bool   `json:"upnp_available"`
				PunchTunnelSucc bool   `json:"punch_tunnel_succ"`
				IPType          int    `json:"ip_type"`
				Ipv6            string `json:"ipv6"`
				Ipv6Type        int    `json:"ipv6_type"`
				Ipv6Status      int    `json:"ipv6_status"`
			} `json:"accounts"`
			JumpSvrIP string `json:"jump_svr_ip"`
			Host      string `json:"host"`
			Isp       struct {
				DisplayName string `json:"display_name"`
				Code        string `json:"code"`
				ID          int    `json:"id"`
			} `json:"isp"`
			PutonTime  int   `json:"puton_time"`
			PutoffTime int64 `json:"putoff_time"`
			AccountMax int   `json:"account_max"`
		} `json:"data"`
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
	file, err := os.Open(fogData)
	if err != nil {
		http.Error(w, "Error opening file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// 解析 JSON 文件
	//var resBox ResBox
	var f FogData
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&f)
	if err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	// 设置响应类型为JSON
	w.Header().Set("Content-Type", "application/json")

	// 将响应数据编码为 JSON 并发送
	json.NewEncoder(w).Encode(f)
}
