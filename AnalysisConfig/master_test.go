package AnalysisConfig

import (
	"encoding/json"
	"fmt"
	"github.com/tealeg/xlsx"
	"log"
	"os/exec"
	"strings"
	"testing"
	"time"
)

var outFile = "./testOut.xlsx"
var outFile2 = "./insertData.xlsx"
var successCount = 0
var failCount = 0

type (
	Student struct {
		Name   string
		Age    int
		Phone  string
		Gender string
		Mail   string
	}

	AdjoinProv struct {
		ispType      string `json:"isp_type"`
		ElementKey   string `json:"element_key"`
		ElementValue string `json:"element_value"`
	}

	MATNodes struct {
		MonitorNodes []string `json:"monitor_nodes"`
		TargetNodes  []string `json:"target_nodes"`
	}

	InsertData struct {
		MonitorIsp  string
		MonitorNode string
		TargetIsp   string
		TargetNode  []string
	}

	response struct {
		Code int `json:"code"`
		Data struct {
			IPingTask struct {
				AddIPingTask struct {
					Code int    `json:"code"`
					Msg  string `json:"msg"`
				} `json:"addIPingTask"`
			} `json:"i_ping_task"`
		} `json:"data"`
		Msg string `json:"msg"`
	}
)

// 导出
func Export(ExportData []outData) {
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("detect_iping_task")
	if err != nil {
		log.Println(err)
		return
	}
	//加列字段名
	WriteTatel(sheet)
	//添加数据
	for _, ed := range ExportData {
		row := sheet.AddRow()
		row.AddCell().Value = ed.DetectType
		row.AddCell().Value = ed.MonitorIsp
		row.AddCell().Value = ed.MonitorNode
		row.AddCell().Value = ed.TargetIsp
		row.AddCell().Value = ed.TargetNode
	}
	err = file.Save(outFile)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("\n\nexport success")

}

func ExportInsertData(insertData []InsertData) {
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("insert_data_detect_iping_task")
	if err != nil {
		log.Println(err)
		return
	}
	//加列字段名
	WriteTatel(sheet)
	//添加数据
	for _, ed := range insertData {
		row := sheet.AddRow()
		row.AddCell().Value = "prov_to_prov"
		row.AddCell().Value = ed.MonitorIsp
		row.AddCell().Value = ed.MonitorNode
		row.AddCell().Value = ed.TargetIsp
		row.AddCell().Value = strings.Join(ed.TargetNode, ",")
	}
	err = file.Save(outFile2)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("\n\nexport success2")

}

func WriteTatel(sheet *xlsx.Sheet) {
	row := sheet.AddRow()
	row.AddCell().Value = "type"
	row.AddCell().Value = "monitor_isp"
	row.AddCell().Value = "monitor_node"
	row.AddCell().Value = "target_isp"
	row.AddCell().Value = "target_node"
}

func ExecCmd(monitorIsp, monitorNode, targetIsp string, targetNode []string) bool {
	var resStruct response

	targetNodeStr := "[\"" + strings.Join(targetNode, "\",\"") + "\"]"

	cmdStr := fmt.Sprintf(`curl 'https://detect-backend.bs58i.baishancdnx.com/graphql/query' \
	-H 'sec-ch-ua: "Chromium";v="118", "Google Chrome";v="118", "Not=A?Brand";v="99"' \
	-H 'sec-ch-ua-mobile: ?0' \
	-H 'Authorization: Bearer eyJhbGciOiJSUzI1NiIsImtpZCI6IkxkU0h3SXNtVlhkRkRpSEIiLCJ0eXAiOiJKV1QifQ.eyJjbGFpbXMiOnsiYWNyIjoiMCIsImF0X2hhc2giOiJVV3o1V1AycTFoTy1qcEJOazZtSFRRIiwiYXVkIjoiemV1cyIsImF1dGhfdGltZSI6MTY5ODIyMzg2OCwiYXpwIjoiemV1cyIsImVtYWlsIjoiZmVuZ2ZlbmcubWVpQGJhaXNoYW4uY29tIiwiZW1haWxfdmVyaWZpZWQiOmZhbHNlLCJleHAiOjE2OTg4Mjg3NzIsImZhbWlseV9uYW1lIjoi5qKF5bOw5bOwIiwiZ2VuZGVyIjoiMSIsImdpdmVuX25hbWUiOiJmZW5nZmVuZy5tZWkiLCJpYXQiOjE2OTgyMjM5NzIsImlzcyI6Imh0dHBzOi8vYWNjb3VudC5iczU4aS5iYWlzaGFuY2RueC5jb20vYXV0aC9yZWFsbXMvbWFzdGVyIiwianRpIjoiMTJkMDM1NTAtMGQxMy00MTdkLTg3NWMtYjhiZTJjZTRjODE1IiwibmFtZSI6ImZlbmdmZW5nLm1laSDmooXls7Dls7AiLCJwcmVmZXJyZWRfdXNlcm5hbWUiOiJmZW5nZmVuZy5tZWkiLCJzZXNzaW9uX3N0YXRlIjoiMGE1ODFhY2EtNzhiNC00ZmU1LWI0MDUtMDEyMzAzZmY2ZjY2Iiwic3NvX2lkIjozMzU1Mywic3ViIjoiYTIwNjRmMzgtNzQ5Ny00YjkyLTkzNTUtMWFkZGFhY2UwMWY4IiwidHlwIjoiSUQifSwiZXhwIjoxNjk4ODI4NzcyLCJyYXciOiJleUpoYkdjaU9pSlNVekkxTmlJc0luUjVjQ0lnT2lBaVNsZFVJaXdpYTJsa0lpQTZJQ0pwV2t4cGNqWkdkMDFzWnpKbFNuTldia1puTVVKcWJVdzRYM0ppUW1SMFZFWkJTREF3TVRsT2VEZEpJbjAuZXlKbGVIQWlPakUyT1RnNE1qZzNOeklzSW1saGRDSTZNVFk1T0RJeU16azNNaXdpWVhWMGFGOTBhVzFsSWpveE5qazRNakl6T0RZNExDSnFkR2tpT2lJeE1tUXdNelUxTUMwd1pERXpMVFF4TjJRdE9EYzFZeTFpT0dKbE1tTmxOR000TVRVaUxDSnBjM01pT2lKb2RIUndjem92TDJGalkyOTFiblF1WW5NMU9Ha3VZbUZwYzJoaGJtTmtibmd1WTI5dEwyRjFkR2d2Y21WaGJHMXpMMjFoYzNSbGNpSXNJbUYxWkNJNklucGxkWE1pTENKemRXSWlPaUpoTWpBMk5HWXpPQzAzTkRrM0xUUmlPVEl0T1RNMU5TMHhZV1JrWVdGalpUQXhaamdpTENKMGVYQWlPaUpKUkNJc0ltRjZjQ0k2SW5wbGRYTWlMQ0p6WlhOemFXOXVYM04wWVhSbElqb2lNR0UxT0RGaFkyRXROemhpTkMwMFptVTFMV0kwTURVdE1ERXlNekF6Wm1ZMlpqWTJJaXdpWVhSZmFHRnphQ0k2SWxWWGVqVlhVREp4TVdoUExXcHdRazVyTm0xSVZGRWlMQ0poWTNJaU9pSXdJaXdpWlcxaGFXeGZkbVZ5YVdacFpXUWlPbVpoYkhObExDSm5aVzVrWlhJaU9pSXhJaXdpYm1GdFpTSTZJbVpsYm1kbVpXNW5MbTFsYVNEbW9vWGxzN0RsczdBaUxDSndjbVZtWlhKeVpXUmZkWE5sY201aGJXVWlPaUptWlc1blptVnVaeTV0WldraUxDSm5hWFpsYmw5dVlXMWxJam9pWm1WdVoyWmxibWN1YldWcElpd2labUZ0YVd4NVgyNWhiV1VpT2lMbW9vWGxzN0RsczdBaUxDSmxiV0ZwYkNJNkltWmxibWRtWlc1bkxtMWxhVUJpWVdsemFHRnVMbU52YlNKOS5SUDZlYkczLUVSeTJMUXB6TTBhdS1DdWhYZmFTaGkzS1BxRkRLX3p1dmZHTEpvOGhMQmotQ2Q1VS1TenJadC1LMHY0VlZUMU1xMURJOFNZb2JiYk9zYk1hUERNMmRLcS13d2xKaHhjZ3pyeWhOc0FUTnNSRlNMOGdYOUJibEhMd0g0V2UwTldZbDNXb05MbjRFeGNlWGNXMjZxOXA5V1NhREhfSTFORlhpbkwzcENHQnZMUGRBTUFuNkdNU2w2d3REbnBlUndiZ1BDQkZ1MHM5ZXlRejl4QjNkSHhPazMyTUtyN1hQbHFLT0ZzWmVHckN3d0FhcWRIeUJOdWZvM1pHd0tDZFJBeDF4bFVsTTRQcng5TXlqdzk2WlpQdXBQUFZpTElNWXd6a2IweDIxVk5xYTBocEZubGpHa0RSZ2FWcUVYbUtyM2dxQlZuOHZDNUNlOWJyOWcifQ.gKHS2fxVOuxkeZmauix-W94dMrjaY9Nz_ehKptM6JrP9vweluZXPmBW4IqIMlUWyRNz2zooQwcUtQuK_t_mCIZMKnGVtBMAkWqnLSWdlFdgoglOLs3vsH42m3r13_VipT8xoejKuvL3aM8CO7ru347BsEvhVl7j6LfXWziLNwoB3g-5vu7IGF2O-zlC93lFHOAugbyrbwc7n6SYU9lj1-6P6IkaABbPb6fBi0TY6m-GOcbDml8H5FWpGlq5rV4CW1K6au4DBeeemKHfMtvPCgsSutf6AAGF2GviVlh5Yyq8o-hXo_EW-rGk9L2CjXvzjengFKg8KTIL56U661U5eJA' \
	-H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36' \
	-H 'Content-Type: application/json;charset=UTF-8' \
	-H 'Accept: application/json, text/plain, */*' \
	-H 'Referer: https://detect-web.bs58i.baishancdnx.com/' \
	-H 'sec-ch-ua-platform: "macOS"' \
	--data-raw $'{"query":"mutation(\\n      $task: AddIPingTaskInputType\u0021\\n      ) {\\n        i_ping_task {\\n          addIPingTask(task: $task) {\\n            code\\n            msg\\n          }\\n        }\\n      }\\n","variables":{"task":{"status":false,"ip_type":"","src_ep_sel_strategy":"province","src_ep_sel_num":5,"src_isp":["%s"],"src_values":["%s"],"dst_ep_sel_strategy":"province","dst_ep_sel_num":5,"dst_isp":["%s"],"dst_values":%s}}}' \
	--compressed`, monitorIsp, monitorNode, targetIsp, targetNodeStr)

	cmd := exec.Command("sh", "-c", cmdStr)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("exec error", err)
	}
	res := string(output)
	err = json.Unmarshal([]byte(res), &resStruct)
	if err != nil || resStruct.Code != 0 {
		fmt.Println("iping任务添加失败，具体信息如下：")
		fmt.Println("curl结果：", res)
		fmt.Println(monitorIsp)
		fmt.Println(monitorNode)
		fmt.Println(targetIsp)
		fmt.Println(targetNodeStr)
		failCount += 1
		failStr := fmt.Sprintf("当前失败数：%d", failCount)
		fmt.Println(failStr)
		fmt.Println("========================")
		return false
	} else {
		fmt.Println("iping任务添加成功！")
		successCount += 1
		sucStr := fmt.Sprintf("当前成功数：%d", successCount)
		fmt.Println(sucStr)
		return true
	}
}

//func addIpingTask(insertData []outData) {
//	for _, data := range insertData {
//		ExecCmd(data.MonitorIsp, data.MonitorNode, data.TargetIsp, data.TargetNode)
//		time.Sleep(time.Second)
//	}
//}

func OneToMore(temp1Data []outData) map[string]map[string][]string {
	//{dx-lt:{"anhui":[guizhou,shanghai...]}}
	MDs := make(map[string]map[string][]string, 1)

	for _, td := range temp1Data {
		//	判断类型,例如：dx-dx,dx-lt,lt-dx ...
		detectType := fmt.Sprintf("%s-%s", td.MonitorIsp, td.TargetIsp)

		// 如果isp-isp类型没有
		if _, ok := MDs[detectType]; !ok {
			tempMap1 := make(map[string][]string, 1)
			var tempArr1 []string
			tempArr1 = append(tempArr1, td.TargetNode)
			tempMap1[td.MonitorNode] = tempArr1
			MDs[detectType] = tempMap1
		} else {
			//	isp-isp类型存在，单是monitor类型不存在
			if _, ok2 := MDs[detectType][td.MonitorNode]; !ok2 {
				var tempArr2 []string
				tempArr2 = append(tempArr2, td.TargetNode)
				MDs[detectType][td.MonitorNode] = tempArr2
			} else {
				var tempArr3 []string
				tempArr3 = MDs[detectType][td.MonitorNode]
				tempArr3 = append(tempArr3, td.TargetNode)
				MDs[detectType][td.MonitorNode] = tempArr3
			}
		}
	}

	// 去重
	for key1, values := range MDs {
		for key2, value := range values {
			temp := removeDuplicates(value)
			MDs[key1][key2] = temp
		}
	}
	return MDs
}

func removeDuplicates(slice []string) []string {
	// 创建一个map，用于记录已经出现过的元素
	seen := make(map[interface{}]bool)

	// 创建一个结果切片
	var result []string

	// 遍历原始切片
	for _, item := range slice {
		// 如果元素已经在map中出现过，跳过该元素
		if seen[item] {
			continue
		}
		// 将元素添加到结果切片中
		result = append(result, item)
		// 将元素添加到map中，表示已经出现过
		seen[item] = true
	}

	return result
}
func printjson(data interface{}) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("JSON encoding error:", err)
		return
	}
	fmt.Println(string(jsonData))
}

func covToInsertData(dates map[string]map[string][]string) []InsertData {
	var res []InsertData
	for key1, ds := range dates {
		//	获取monitor_isp 与 target_isp
		ispArr := strings.Split(key1, "-")
		mIsp := ispArr[0]
		tIsp := ispArr[1]
		for mNode, taNodes := range ds {
			res = append(res, InsertData{
				MonitorIsp:  mIsp,
				MonitorNode: mNode,
				TargetIsp:   tIsp,
				TargetNode:  taNodes,
			})
		}
	}
	return res
}

func addIpingTask(iData []InsertData) {
	for _, d := range iData {
		cmdRes := ExecCmd(d.MonitorIsp, d.MonitorNode, d.TargetIsp, d.TargetNode)
		if cmdRes == false {
			fmt.Println("iping任务添加失败，结束添加！")
			return
		}
		time.Sleep(time.Second)
	}
}

func removeDiffIsp(dates []InsertData) []InsertData {
	var tempData []InsertData
	var resultData []InsertData

	var tempDx []InsertData
	var tempLt []InsertData
	var tempYd []InsertData

	var pData []InsertData

	tempMap := make(map[string][]string, 0)

	for _, date := range dates {
		if date.TargetIsp == date.MonitorIsp {
			tempData = append(tempData, date)
		}
	}

	// 聚合并去重
	for _, value := range tempData {
		//	制造key
		key := fmt.Sprintf("%s-%s", strings.ReplaceAll(value.MonitorNode, " ", ""), strings.ReplaceAll(value.MonitorIsp, " ", ""))
		if _, ok := tempMap[key]; !ok {
			tempMap[key] = value.TargetNode
		} else {
			//	去重
			//	1、拿出value
			var result []string

			tempArr1 := tempMap[key]
			tempArr2 := value.TargetNode

			for _, t := range tempArr1 {
				tempArr2 = append(tempArr2, t)
			}
			result = removeDuplicates(tempArr2)
			tempMap[key] = result
		}
	}

	marshal2, err := json.Marshal(tempMap)
	if err != nil {
		return nil
	}
	fmt.Println("+++++++++")
	fmt.Println(string(marshal2))
	fmt.Println("+++++++++")

	// 拼接成合并后的格式
	for key, value := range tempMap {
		var ID InsertData
		//	 切分key
		keyArray := strings.Split(key, "-")
		monitorNode := keyArray[0]
		isp := keyArray[1]
		ID.MonitorIsp = isp
		ID.TargetIsp = isp
		ID.MonitorNode = monitorNode
		ID.TargetNode = value
		resultData = append(resultData, ID)
	}

	for _, datum := range resultData {
		switch datum.MonitorIsp {
		case "dx":
			tempDx = append(tempDx, datum)
		case "lt":
			tempLt = append(tempLt, datum)
		case "yd":
			tempYd = append(tempYd, datum)
		default:
			continue
		}
	}

	for _, dx := range tempDx {
		pData = append(pData, dx)
	}

	for _, lt := range tempLt {
		pData = append(pData, lt)
	}

	for _, yd := range tempYd {
		pData = append(pData, yd)
	}

	return pData
}

func TestNewDirector(t *testing.T) {
	manager := NewInstance()
	manager.reloadCfgData()
	manager.dealCfgData()
	lastData := manager.removeDup()
	Export(lastData)
	moTomo := OneToMore(lastData)
	//printjson(moTomo)
	iData := covToInsertData(moTomo)
	lData := removeDiffIsp(iData)
	fmt.Println("=============")
	marshal, err := json.Marshal(lData)
	if err != nil {
		return
	}
	fmt.Println(string(marshal))
	fmt.Println("=============")
	ExportInsertData(lData)
	//addIpingTask(lData)

	//tempData1 := InsertData{
	//	MonitorIsp:  "lt",
	//	MonitorNode: "henan",
	//	TargetIsp:   "dx",
	//	TargetNode: []string{
	//		"guizhou",
	//	},
	//}
	//
	//var tempdata2 []InsertData
	//tempdata2 = append(tempdata2, tempData1)
	//
	//addIpingTask(tempdata2)

}
