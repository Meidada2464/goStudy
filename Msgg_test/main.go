package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	ContentTypeJson = "Content-Type: application/json"
	Wechat          = 1
)

type (
	MsggResponse struct {
		Errno int `json:"errno"`
		Data  struct {
			Code int    `json:"code"`
			Msg  string `json:"msg"`
			Data MData  `json:"data"`
		} `json:"data"`
		Message string `json:"message"`
	}

	MData struct {
		Mail   Res `json:"mail"`
		Wechat Res `json:"wechat"`
		Call   Res `json:"call"`
		SMS    Res `json:"sms"`
	}

	Res struct {
		Code int    `json:"code"`
		Id   string `json:"id"`
	}
)

func main() {
	context := `{"token":"f9a69f17d44c6f7bce24b5e7ee79b805","params":{"to":["fengfeng.mei"],"content":"【探测系统高优先级通知-test】\n尊敬的%!s(*string=0x1400047a948)探测用户您好,您的探测任务%!s(*string=0x1400047a958)已到期自动关闭，如有需要继续使用请立即前往探测系统处理。\n【我来处理】https://detect-web.bs58i.baishancdnx.com/#/network/task"}}`

	err := HttpPost("https://msgg.bs58i.baishancdnx.com/api/app/1.0/msgg/submitwechat", context)
	if err != nil {
		return
	}
}

func HttpPost2(url, body string) error {
	response, err := http.Post(url, ContentTypeJson, strings.NewReader(body))
	// request error
	if err != nil {
		return err
	}
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 {
		return errors.New(response.Status)
	}

	// read body error
	if err != nil {
		return err
	}

	var msggResponse MsggResponse
	err = json.Unmarshal(responseBody, &msggResponse)
	if msggResponse.Errno != 0 {
		return errors.New(msggResponse.Message)
	}
	return err
}

func HttpPost(url, context string) error {
	request, err := http.NewRequest("POST", url, strings.NewReader(context))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}

	if err != nil {
		return err
	}

	var msggResponse MsggResponse
	err = json.Unmarshal(responseBody, &msggResponse)
	if msggResponse.Errno != 0 {
		return errors.New(msggResponse.Message)
	}
	return err
}
