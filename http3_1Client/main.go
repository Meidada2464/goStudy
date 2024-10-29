/**
 * Package http3Client
 * @Author fengfeng.mei <fengfeng.mei@baishan.com>
 * @Date 2024/10/21 17:49
 */

package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/quic-go/quic-go"
	"io"
	"time"
)

const addr = "127.0.0.1:9999"
const message = "ccc"

func main() {
	cbx := context.Background()
	session, err := quic.DialAddr(cbx, addr, &tls.Config{InsecureSkipVerify: true}, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	stream, err := session.OpenStreamSync(cbx)
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		fmt.Printf("Client: Sending '%s'\n", message)
		_, err = stream.Write([]byte(message))
		if err != nil {
			fmt.Println(err)
			return
		}
		buf := make([]byte, len(message))
		_, err = io.ReadFull(stream, buf)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Client: Got '%s'\n", buf)

		time.Sleep(2 * time.Second)
	}
}
