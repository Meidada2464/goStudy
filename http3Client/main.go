/**
 * Package http3Client
 * @Author fengfeng.mei <fengfeng.mei@baishan.com>
 * @Date 2024/10/22 17:18
 */

package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"github.com/quic-go/quic-go/qlog"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httptrace"
	"os"
	"time"
)

const addr = "https://localhost:9999/demo/text"

type httpDlInfo struct {
	status            int
	fileSize          int64
	dstIp             string
	connTime          float64
	dnsTime           float64
	tlsTime           float64
	firstByteTime     float64
	downloadTime      float64
	totalDownloadTime float64
	downloadRate      float64
}

func main() {
	keyLogFile := "keylogFile.log"
	var keyLog io.Writer
	f, err := os.Create(keyLogFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	keyLog = f

	pool, err := x509.SystemCertPool()
	if err != nil {
		log.Fatal(err)
	}
	//testdata.AddRootCA(pool)

	roundTripper := &http3.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs:            pool,
			KeyLogWriter:       keyLog,
			InsecureSkipVerify: false,
		},
		QUICConfig: &quic.Config{
			Tracer: qlog.DefaultConnectionTracer,
		},
	}

	defer roundTripper.Close()

	hclient := &http.Client{
		Transport: roundTripper,
	}

	// 创建一个request请求
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

	request, err := http.NewRequest("GET", addr, nil)
	if err != nil {
		log.Printf("Error creating request: %s", err)
		return
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
	request.WithContext(httptrace.WithClientTrace(ctx, trace))
	defer cancelFunc()

	rsp, err := hclient.Do(request)
	if err != nil {
		log.Fatal("hClient request Error: ", err)
		return
	}
	log.Printf("Got response for %s: %#v", addr, rsp)

	//body := &bytes.Buffer{}
	//_, err = io.Copy(body, rsp.Body)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//log.Printf("Response Body (%d bytes):\n%s", body.Len(), body.Bytes())

	// 计算下载速率 (字节每秒)
	b, err := io.ReadAll(rsp.Body)
	if err != nil {
		_ = fmt.Errorf("error reading response body %v", err)
	}

	// 记录响应完成时间
	end = time.Now()

	// 计算各项时间
	httpInfo.connTime = connEnd.Sub(connStart).Seconds()
	httpInfo.tlsTime = tlsEnd.Sub(tlsStart).Seconds()
	httpInfo.firstByteTime = firstByteEnd.Sub(connEnd).Seconds()
	httpInfo.downloadTime = end.Sub(firstByteEnd).Seconds()
	httpInfo.totalDownloadTime = end.Sub(start).Seconds()
	httpInfo.dnsTime = dnsEnd.Sub(dnsStart).Seconds()
	httpInfo.fileSize = int64(len(b))
	httpInfo.downloadRate = 0

	if httpInfo.downloadTime != 0 {
		httpInfo.downloadRate = float64(httpInfo.fileSize) / httpInfo.downloadTime
	}

	fmt.Println("httpInfo.fileSize:", httpInfo.fileSize)
	fmt.Println("httpInfo.dnsTime:", httpInfo.dnsTime)
	fmt.Println("httpInfo.connTime:", httpInfo.connTime)
	fmt.Println("httpInfo.tlsTime:", httpInfo.tlsTime)
	fmt.Println("httpInfo.firstByteTime:", httpInfo.firstByteTime)
	fmt.Println("httpInfo.downloadTime:", httpInfo.downloadTime)
	fmt.Println("httpInfo.totalDownloadTime:", httpInfo.totalDownloadTime)
	fmt.Println("httpInfo.downloadRate:", httpInfo.downloadRate)

	//fmt.Println("httpInfo:", httpInfo)

}
