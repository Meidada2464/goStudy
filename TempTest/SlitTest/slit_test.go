package SlitTest

import (
	"encoding/json"
	"fmt"
	"goStudy/imroc/req"
	"testing"
)

type Router struct {
	middlewareChain []string
}

type (
	ProvRes struct {
		Code      int          `json:"code,omitempty"`
		Data      []ProResData `json:"data"`
		Msg       string       `json:"msg,omitempty"`
		RequestId string       `json:"request_id,omitempty"`
	}

	ProResData struct {
		ID        int64  `json:"id"`
		CountryID int    `json:"country_id"`
		RegionID  int    `json:"region_id"`
		Name      string `json:"name"`
		Ename     string `json:"ename"`
		Code      string `json:"code"`
		Status    int    `json:"status"`
		Remark    string `json:"remark"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}
)

func TestUncode(t *testing.T) {
	//str := `{"code": 0,"request_id": "1701748710.5442_782933","data": [{"id": 1,"country_id": 1,"region_id": 4,"name": "北京","ename": "beijing","code": "beijing","status": 1,"remark": "beijing","created_at": "2020-02-17 14:57:32","updated_at": "2020-02-17 14:57:32"}]}`
	//
	//var result ProvRes
	//
	//err := json.Unmarshal([]byte(str), &result)
	//if err != nil {
	//	fmt.Println("Error decoding JSON:", err)
	//	return
	//}

	var res ProvRes
	provinceIdMap := make(map[int64]string, 10)
	resp, err := req.Get("https://service-taishan.bs58i.baishancloud.com/api/location/province")
	if err != nil {
		fmt.Printf("%+v\n", err)
	}

	statusCode := resp.Response().StatusCode
	if statusCode == 200 {
		err = resp.ToJSON(&res)
		if err != nil {
			fmt.Printf("%+v\n", err)

		}
	}

	if statusCode != 200 {
		fmt.Printf("%+v\n", err)

		return
	}

	for _, d := range res.Data {
		provinceIdMap[d.ID] = d.Name
	}

	fmt.Printf("%+v\n", provinceIdMap)

}

func TestMapNull(t *testing.T) {
	ipInfo := MapT()
	ipInfoStr, _ := json.Marshal(ipInfo)
	fmt.Println("aaa", string(ipInfoStr))
}

func MapT() (aInfo map[string]int) {
	aInfo = map[string]int{}
	return
}
