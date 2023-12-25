package SyncOAData

import (
	"fmt"
	"github.com/imroc/req"
	"goStudy/SyncOAData/Common"
	"goStudy/SyncOAData/Model"
	"strings"
)

func SyncOneAgentData() {
	db, err := InitDataBase()
	if err != nil {
		fmt.Println("链接失败：", err)
		return
	}
	// 获取OA数据
	params := req.Param{
		"appToken": "93rjewk8vmp4tzwuqnq9z37bqlm4kyh2",
	}
	resp, err := req.Get(
		"https://aegis.bs58i.baishancloud.com:10037/api/v1/config/queryAppConfigDetails?appName=one-agent&env=pro&showHistory=true",
		params,
	)
	if err != nil {
		return
	}
	//fmt.Println("resp", resp)
	var res Common.ResBody2
	err = resp.ToJSON(&res)
	if err != nil {
		return
	}

	if res.Code != 200 {
		err = fmt.Errorf("response body code is not 200:%v", res)
		return
	}
	OAData := res.Data
	fmt.Println("OAData", OAData)

	var PlugInstallList []Model.PlugInstall

	for i := 0; i < len(OAData); i++ {
		var PlugInstallVersionList []Model.PlugInstallVersion
		var types string
		types = "单机"

		if strings.Contains(OAData[i].Key, "@") {
			types = "机器组"
		}
		if OAData[i].Key == "default" {
			types = "默认类型"
		}

		if len(OAData[i].Versions) != 0 {
			for j := 0; j < len(OAData[i].Versions); j++ {
				PlugInstallVersionList = append(PlugInstallVersionList, Model.PlugInstallVersion{
					Types:     types,
					Name:      OAData[i].Key,
					Version:   OAData[i].Versions[j].Version,
					Parameter: OAData[i].Versions[j].Config,
				})
			}
		}

		tstemp := fmt.Sprintf("XN%d", i+1)

		PlugInstallList = append(PlugInstallList, Model.PlugInstall{
			TaskName:               tstemp,
			TaskGroup:              "效能管理研发",
			PlugName:               OAData[i].AppName,
			PlugVersion:            OAData[i].Version,
			PlugType:               "官方插件",
			Types:                  types,
			Name:                   OAData[i].Key,
			Parameter:              OAData[i].Config,
			Notes:                  OAData[i].Remark,
			PlugInstallVersionList: PlugInstallVersionList,
			CreatedUser:            "梅峰峰",
			UpdatedUser:            "梅峰峰",
		})
	}
	fmt.Println("插入的数据结构结果：", PlugInstallList)

	err = db.Table("plug_install").Create(&PlugInstallList).Error
	if err != nil {
		fmt.Println("批量插入失败：", err)
	}

}
