package SyncOAData

import (
	"fmt"
	"github.com/imroc/req"
	"goStudy/SyncOAData/Common"
	"goStudy/SyncOAData/Model"
)

func SyncAgentIntegrationData() {
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
		"https://aegis.bs58i.baishancloud.com:10037/api/v1/config/queryAppConfigDetails?appName=one-agent-integration-config&env=pro&showHistory=true",
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
	OAIData := res.Data
	fmt.Println("OAData", OAIData)

	var PlugRegisterList []Model.PlugRegister

	for i := 0; i < len(OAIData); i++ {
		var PlugOldVersionList []Model.PlugOldVersion

		// 拼装子表数据
		if len(OAIData[i].Versions) != 0 {
			for j := 0; j < len(OAIData[i].Versions); j++ {
				PlugOldVersionList = append(PlugOldVersionList, Model.PlugOldVersion{
					PlugName:      OAIData[i].Key,
					PlugVersion:   OAIData[i].Versions[j].Version,
					PlugParameter: OAIData[i].Versions[j].Config,
					Introduce:     OAIData[i].Versions[j].Remark,
					CreatedUser:   "梅峰峰",
					UpdatedUser:   "梅峰峰",
				})
			}
		}

		PlugRegisterList = append(PlugRegisterList, Model.PlugRegister{
			PlugName:       OAIData[i].Key,
			PlugVersion:    OAIData[i].Version,
			PlugType:       "官方插件",
			Introduce:      OAIData[i].Remark,
			Configuration:  OAIData[i].Config,
			CreatedUser:    "梅峰峰",
			UpdatedUser:    "梅峰峰",
			PlugOldVersion: PlugOldVersionList,
		})
	}
	fmt.Println("插入的数据结构结果：", PlugRegisterList)

	err = db.Table("plug_register").Create(&PlugRegisterList).Error
	if err != nil {
		fmt.Println("批量插入失败：", err)
	}

}
