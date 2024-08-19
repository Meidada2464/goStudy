/**
 * Package util
 * @Author fuqiang.li <fuqiang.li@baishan.com>
 * @Date 2023/4/25
 */

package util

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

// HTTPRequest 发起http请求
func HTTPRequest(cli *http.Client, method, url string,
	headers map[string]string, validCodes []int,
	payload []byte) ([]byte, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	resp, err := cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	for _, code := range validCodes {
		if resp.StatusCode != code {
			return nil, fmt.Errorf("http code fail: %d", resp.StatusCode)
		}
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
