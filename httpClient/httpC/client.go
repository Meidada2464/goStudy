/**
 * Package httpC
 * @Author fengfeng.mei <fengfeng.mei@baishan.com>
 * @Date 2024/12/3 14:21
 */

package httpC

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptrace"
	"strings"
	"sync"
	"time"
)

func httpRequestClient(param httpTaskParam, ip, port string, wg *sync.WaitGroup) (httpDlInfo, error) {
	if wg != nil {
		defer wg.Done()
	}

	var (
		start        time.Time
		connStart    time.Time
		connEnd      time.Time
		tlsStart     time.Time
		tlsEnd       time.Time
		firstByteEnd time.Time
		end          time.Time
		dnsStart     time.Time
		dnsEnd       time.Time
		resolvedIPs  []net.IPAddr
		httpInfo     httpDlInfo
	)

	httpInfo.status = HTTP_FAILED

	// 创建 HTTP 客户端
	client := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
			TLSClientConfig:   &tls.Config{InsecureSkipVerify: param.skipVerify},
			ForceAttemptHTTP2: param.enableHttp2,
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				start = time.Now()
				// 使用指定的IP地址和端口
				var d net.Dialer
				var conn net.Conn
				var err error

				// 指定出口ip
				isSrcIp := false
				if param.srcIp != "" {
					srcIp := net.ParseIP(param.srcIp)
					if srcIp != nil {
						d.LocalAddr = &net.TCPAddr{IP: srcIp}
						isSrcIp = true
					}
				}

				if ip == "" {
					conn, err = d.DialContext(ctx, network, addr)
				} else {
					httpInfo.dstIp = ip
					conn, err = d.DialContext(ctx, network, net.JoinHostPort(ip, port))
				}

				if err != nil {
					if !isSrcIp {
						return nil, err
					}

					var retryD net.Dialer
					if ip == "" {
						conn, err = retryD.DialContext(ctx, network, addr)
					} else {
						httpInfo.dstIp = ip
						conn, err = retryD.DialContext(ctx, network, net.JoinHostPort(ip, port))
					}
				}

				if err != nil {
					return nil, err
				}

				err = conn.SetDeadline(time.Now().Add(param.maxTimeout)) // 设置发送接收数据超时
				if err != nil {
					return nil, err
				}
				return conn, nil
			},
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if !param.followRedirects {
				return http.ErrUseLastResponse
			}

			if param.maxRedirects != nil && len(via) >= *param.maxRedirects {
				return http.ErrUseLastResponse
			}

			return nil
		},
	}

	// 创建 HTTP 请求
	var request *http.Request
	var err error
	if param.HttpMethod == "GET" {
		request, err = http.NewRequest(param.HttpMethod, param.url, nil)
	} else if param.HttpMethod == "POST" {
		request, err = http.NewRequest(param.HttpMethod, param.url, bytes.NewReader(param.payload))
	} else {
		// 暂时不允许其他的方式
		return httpInfo, fmt.Errorf("http method %s is not supported", param.HttpMethod)
	}

	if err != nil {
		return httpInfo, fmt.Errorf("new request error: %v", err)
	}

	// 添加请求头
	for key, value := range param.header {
		// host头部要用指定方法设置
		if key == "Host" || key == "host" {
			request.Host = value
			continue
		}
		request.Header.Add(key, value)
	}

	// 添加basic auth
	if param.basicAuth != nil {
		userName, passwd := getBasicAuth(*param.basicAuth)
		if userName != "" && passwd != "" {
			request.SetBasicAuth(userName, passwd)
		}
	}

	// 创建 HTTP trace
	trace := &httptrace.ClientTrace{
		DNSStart: func(info httptrace.DNSStartInfo) {
			dnsStart = time.Now()
		},
		DNSDone: func(info httptrace.DNSDoneInfo) {
			dnsEnd = time.Now()
			resolvedIPs = info.Addrs
			if len(resolvedIPs) == 0 {
				return
			}
			httpInfo.dstIp = resolvedIPs[0].IP.String()
		},
		TLSHandshakeStart: func() {
			tlsStart = time.Now()
		},
		TLSHandshakeDone: func(cs tls.ConnectionState, err error) {
			tlsEnd = time.Now()
		},
		GotFirstResponseByte: func() {
			firstByteEnd = time.Now()
		},
		ConnectStart: func(network, addr string) {
			connStart = time.Now()
		},
		ConnectDone: func(network, addr string, err error) {
			connEnd = time.Now()
		},
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	request = request.WithContext(httptrace.WithClientTrace(ctx, trace))

	// 发送请求
	res, err := client.Do(request)
	if err != nil {
		return httpInfo, fmt.Errorf("do request error: %v", err)
	}

	// 5xx retry
	if res.StatusCode == HTTP_CODE_502 || res.StatusCode == HTTP_CODE_503 {
		return httpInfo, fmt.Errorf("do request code:%d", res.StatusCode)
	}

	defer cancelFunc()
	defer res.Body.Close()

	expectedHttpCode := false
	if len(param.expectedHttpCodes) == 0 {
		expectedHttpCode = true
	} else {
		for _, v := range param.expectedHttpCodes {
			if res.StatusCode == v {
				expectedHttpCode = true
				break
			}
		}
	}

	// 非预期状态码直接返回
	if !expectedHttpCode {
		return httpInfo,
			fmt.Errorf("expected http is %v, but got %v", param.expectedHttpCodes, res.StatusCode)
	}

	// 计算下载速率 (字节每秒)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return httpInfo, fmt.Errorf("error reading response body %v", err)
	}

	httpInfo.status = res.StatusCode

	// 记录响应完成时间
	end = time.Now()

	// 计算各项时间
	httpInfo.connTime = connEnd.Sub(connStart).Seconds()
	httpInfo.tlsTime = tlsEnd.Sub(tlsStart).Seconds()
	httpInfo.firstByteTime = firstByteEnd.Sub(connEnd).Seconds()
	httpInfo.downloadTime = end.Sub(firstByteEnd).Seconds()
	httpInfo.totalDownloadTime = end.Sub(start).Seconds()
	httpInfo.dnsTime = dnsEnd.Sub(dnsStart).Seconds()
	httpInfo.fileSize = int64(len(body))
	httpInfo.downloadRate = 0

	if httpInfo.downloadTime != 0 {
		httpInfo.downloadRate = float64(httpInfo.fileSize) / httpInfo.downloadTime
	}

	return httpInfo, nil
}

func getBasicAuth(basicAuth string) (string, string) {
	item := strings.Split(basicAuth, ":")
	if len(item) != 2 {
		return "", ""
	}
	return item[0], item[1]
}
