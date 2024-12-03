/**
 * Package pingByNic
 * @Author fengfeng.mei <fengfeng.mei@baishan.com>
 * @Date 2024/11/14 17:46
 */

package main

import (
	"fmt"
	"net"
)

func main() {
	v4Nics, v6Nics := GetALlNic()
	fmt.Println("==========V4 Nics ==========")
	for k, v := range v4Nics {
		fmt.Println(k, v)
	}

	//fmt.Println("==========V6 Nics ==========")
	//
	//for k, v := range v6Nics {
	//	fmt.Println(k, v)
	//}
}

func GetALlNic() (map[string]string, map[string]string) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, nil
	}

	var Ipv4nics = make(map[string]string)
	var Ipv6nics = make(map[string]string)

	for _, intf := range interfaces {

		addrs, err := intf.Addrs()
		if err != nil {
			return nil, nil
		}

		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}

			if ipNet.IP.IsLoopback() || ipNet.IP.IsMulticast() {
				continue
			}

			if ipNet.IP.To4() != nil {
				Ipv4nics[intf.Name] = ipNet.IP.String()
			} else {
				Ipv6nics[intf.Name] = ipNet.IP.String()
			}
		}
	}

	return Ipv4nics, Ipv6nics
}
