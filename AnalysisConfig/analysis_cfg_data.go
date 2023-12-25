package AnalysisConfig

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/tealeg/xlsx"
	"log"
	"os"
	"strings"
	"sync"
)

var inFile = "./d_m.xlsx"
var cfgFile = "./iping_defined_task.json"
var once sync.Once
var instancd *Manager

type (
	Manager struct {
		Configs  detectCfg
		out      []outData
		tempData map[string]outData
		lastData []outData
	}

	detectCfg struct {
		IpingDefinedTask []detectElem `json:"iping_defined_task"`
	}

	outData struct {
		DetectType  string
		MonitorIsp  string
		MonitorNode string
		TargetIsp   string
		TargetNode  string
	}

	detectElem struct {
		MonitorType string `json:"monitor_type"`
		Task        struct {
			Monitor struct {
				Isp  string   `json:"isp"`
				Node []string `json:"node"`
			} `json:"monitor"`
			Target struct {
				Isp  string   `json:"isp"`
				Node []string `json:"node"`
			} `json:"target"`
		} `json:"task"`
	}
)

func NewInstance() *Manager {
	once.Do(func() {
		instancd = &Manager{
			tempData: map[string]outData{},
		}
	})
	return instancd
}

// 导入
func Import() []detectElem {
	//打开文件
	xlsxFile, err := xlsx.OpenFile(inFile)
	if err != nil {
		log.Println(err)
		return nil
	}

	detectCfgs := make([]detectElem, 0)

	//遍历sheet页
	for _, sheet := range xlsxFile.Sheets {
		fmt.Println("sheet name:", sheet.Name)
		//遍历行读取
		for i, row := range sheet.Rows {
			if i > 0 {
				var dc detectElem
				//遍历每行的列读取
				for k, cell := range row.Cells {
					//fmt.Println(k,cell.Value)
					switch k {
					case 0:
						elementKeyData := cell.Value
						elementKeyArray := strings.Split(elementKeyData, "-")
						monitorNode := elementKeyArray[0]
						isp := elementKeyArray[1]
						// 拼接monitor数据
						dc.MonitorType = "prov_to_prov"
						monitorNodeArray := []string{monitorNode}
						dc.Task.Monitor.Isp = isp
						dc.Task.Target.Isp = isp
						dc.Task.Monitor.Node = monitorNodeArray
					case 1:
						// 拼接target数据
						elementValueData := cell.Value
						elementValueArray := strings.Split(elementValueData, ",")
						dc.Task.Target.Node = elementValueArray
					default:
						return nil
					}
				}
				if dc.MonitorType != "" {
					detectCfgs = append(detectCfgs, dc)
				}
			}
		}
	}
	resp, err := json.Marshal(detectCfgs)
	if err != nil {
		log.Println(err)
		return nil
	}
	fmt.Println(string(resp))
	return detectCfgs
}

func (c *Manager) reloadCfgData() {

	var temp detectCfg
	var res detectCfg

	fileByte, err := os.ReadFile(cfgFile)
	if err != nil {
		fmt.Println("err-35:", err)
		return
	}

	err = json.Unmarshal(fileByte, &temp)
	if err != nil {
		fmt.Println("err-41:", err)
		return
	}

	elems := Import()
	if elems == nil {
		fmt.Println("读取xlsx出错")
		return
	}

	tempElems := temp.IpingDefinedTask

	for _, elem := range elems {
		tempElems = append(tempElems, elem)
	}

	res.IpingDefinedTask = tempElems
	c.Configs = res

	fmt.Println("configData:", c.Configs)
}

// 重新处理数据，将数据扁平化
func (c *Manager) dealCfgData() {
	for _, IDT := range c.Configs.IpingDefinedTask {
		if len(IDT.Task.Monitor.Node) != 0 && len(IDT.Task.Target.Node) != 0 {
			for _, MN := range IDT.Task.Monitor.Node {
				for _, TN := range IDT.Task.Target.Node {
					c.out = append(c.out, outData{
						DetectType:  IDT.MonitorType,
						MonitorIsp:  IDT.Task.Monitor.Isp,
						MonitorNode: MN,
						TargetIsp:   IDT.Task.Target.Isp,
						TargetNode:  TN,
					})
				}
			}
		}
	}
}

func (c *Manager) removeDup() []outData {

	println("len(out)", len(c.out))

	for _, data := range c.out {
		// 使用MD5哈希函数
		byteTemp, err := json.Marshal(data)
		if err != nil {
			fmt.Println()
			return nil
		}
		md5Hash := md5.Sum(byteTemp)
		md5Hex := hex.EncodeToString(md5Hash[:])
		c.tempData[md5Hex] = data
	}

	fmt.Println("len(tempData)", len(c.tempData))

	// 去重
	for _, v := range c.tempData {
		c.lastData = append(c.lastData, v)
	}

	fmt.Println("len(lastData)", len(c.lastData))
	return c.lastData
}
