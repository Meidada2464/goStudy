/**
 * Package http3_1Server
 * @Author fengfeng.mei <fengfeng.mei@baishan.com>
 * @Date 2024/10/21 17:57
 */
package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/quic-go/quic-go"
	"io"
	"math/big"
)

const saddr = "127.0.0.1:9999"

func main() {
	ctx := context.Background()
	listener, err := quic.ListenAddr(saddr, generateTLSConfig(), nil)
	if err != nil {
		fmt.Println(err)
	}
	for {
		sess, err := listener.Accept(ctx)
		if err != nil {
			fmt.Println(err)
		} else {
			go dealSession(sess)
		}
	}
}
func dealSession(sess quic.Connection) {
	ctx := context.Background()
	stream, err := sess.AcceptStream(ctx)
	if err != nil {
		panic(err)
	} else {
		for {
			_, err = io.Copy(loggingWriter{stream}, stream)
		}
	}
}

type loggingWriter struct{ io.Writer }

func (w loggingWriter) Write(b []byte) (int, error) {
	fmt.Printf("Server: Got '%s'\n", string(b))
	return w.Writer.Write(b)
}

// Setup a bare-bones TLS config for the server
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
