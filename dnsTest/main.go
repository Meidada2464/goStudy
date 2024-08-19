package main

import (
	"context"
	"fmt"
	"net"
	"time"
)

func dnsTest() {
	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: 10 * time.Second,
			}
			return d.DialContext(ctx, "udp", "119.29.29.29:53")
		},
	}

	ips, err := r.LookupHost(context.Background(), "mallard-transfer-open.bs58i.baishancdnx.com")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ips)
}

func main() {
	dnsTest()
}
