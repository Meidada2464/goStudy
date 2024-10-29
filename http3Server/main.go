/**
 * Package http3Server
 * @Author fengfeng.mei <fengfeng.mei@baishan.com>
 * @Date 2024/10/21 14:39
 */

package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"github.com/quic-go/quic-go/qlog"
	"goStudy/http3Server/testdata"
	"io"
	"math/big"
	"net/http"
	"time"
)

const (
	server = "localhost:9999"
)

// setupHandler 设置开启的http3的服务路径
func setupHandler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/demo/tile", func(w http.ResponseWriter, r *http.Request) {
		// Small 40x40 png
		w.Write([]byte{
			0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x00, 0x00, 0x0d,
			0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x28, 0x00, 0x00, 0x00, 0x28,
			0x01, 0x03, 0x00, 0x00, 0x00, 0xb6, 0x30, 0x2a, 0x2e, 0x00, 0x00, 0x00,
			0x03, 0x50, 0x4c, 0x54, 0x45, 0x5a, 0xc3, 0x5a, 0xad, 0x38, 0xaa, 0xdb,
			0x00, 0x00, 0x00, 0x0b, 0x49, 0x44, 0x41, 0x54, 0x78, 0x01, 0x63, 0x18,
			0x61, 0x00, 0x00, 0x00, 0xf0, 0x00, 0x01, 0xe2, 0xb8, 0x75, 0x22, 0x00,
			0x00, 0x00, 0x00, 0x49, 0x45, 0x4e, 0x44, 0xae, 0x42, 0x60, 0x82,
		})
	})

	mux.HandleFunc("/demo/tiles", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "<html><head><style>img{width:40px;height:40px;}</style></head><body>")
		for i := 0; i < 200; i++ {
			fmt.Fprintf(w, `<img src="/demo/tile?cachebust=%d">`, i)
		}
		io.WriteString(w, "</body></html>")
	})

	mux.HandleFunc("/demo/echo", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("error reading body while handling /echo: %s\n", err.Error())
		}
		w.Write(body)
	})

	mux.HandleFunc("/demo/text", func(w http.ResponseWriter, r *http.Request) {
		// 设置响应头，指定内容类型为 text/html
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		// HTML 页面内容
		htmlContent := `
    <!DOCTYPE html>
    <html>
    <head>
        <title>ChatGPT 介绍</title>
    </head>
    <body>
        <h1>ChatGPT 简介</h1>
        <p>ChatGPT 是由 OpenAI 开发的一种基于大型语言模型的对话生成系统。它使用了先进的深度学习技术，能够生成具有上下文和逻辑连贯性的自然语言对话。</p>
        <h2>ChatGPT 的特点：</h2>
        <ul>
            <li>生成自然、连贯的对话内容</li>
            <li>支持多种语言和对话主题</li>
            <li>可用于聊天机器人、智能助手等应用</li>
            <li>提供了丰富的 API 和工具，方便开发和集成</li>
        </ul>
            <p>想要了解更多关于 ChatGPT 的信息，请访问官方网站：<a 
    href="https://www.openai.com/gpt">OpenAI GPT</a></p>
    </body>
    </html>
    `
		// 将 HTML 内容写入响应体
		fmt.Fprintf(w, htmlContent)
		time.Sleep(5 * time.Second)
	})
	return mux

}

func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{Certificates: []tls.Certificate{tlsCert}}
}

func main() {
	// 开启http3服务
	fmt.Println("listening on ", server)
	var (
		err error
	)
	handler := setupHandler()
	certFile, keyFile := testdata.GetCertificatePaths()
	s := http3.Server{
		Handler: handler,
		Addr:    server,
		QUICConfig: &quic.Config{
			Tracer: qlog.DefaultConnectionTracer,
		},
	}
	err = s.ListenAndServeTLS(certFile, keyFile)
	if err != nil {
		fmt.Println(err)
	}
}
