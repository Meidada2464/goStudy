package FileUpload

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestNewCarFacade(t *testing.T) {
	// 初始化连接池
	client := http.Client{
		Transport: &http.Transport{
			MaxIdleConns:          10,
			MaxIdleConnsPerHost:   5,
			ResponseHeaderTimeout: 120 * time.Second,
			IdleConnTimeout:       120 * 3 * time.Second,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: 60 * time.Second,
	}
	req, err := http.NewRequest("GET", "http://bsgtm-hexcenter.dolfindnsx.net/src/bsip_ipdb/2023.10.24/bgpdata.zip", bytes.NewBufferString(""))
	if err != nil {
		fmt.Println("http request fail, error :", err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "cc650ce2dedf53d85f135a224dbf2cdc")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("client.Do fail", err)
		return
	}
	if err != nil {
		return
	}
	defer resp.Body.Close()
	fmt.Println("resp:", resp)
	if resp.StatusCode != http.StatusOK {
		return
	}
	// 创建一个临时文件
	out, err := os.Create("./bgpdata.zip")
	if err != nil {
		fmt.Println("os.Create fail, error :", err)
		return
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println("io.Copy fail, error :", err)
		return
	}
}
