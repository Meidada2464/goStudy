/**
 * Package util
 * @Author fuqiang.li <fuqiang.li@baishan.com>
 * @Date 2023/7/17
 */

package util

import (
	"fmt"
	"net"
	"strings"
)

// GetIpFamily 获取ip类型
func GetIpFamily[V any](ip net.IP, ipv4, ipv6, unknown V) V {
	if ip.To4() != nil {
		return ipv4
	}
	if ip.To4() == nil && ip.To16() != nil {
		return ipv6
	}
	return unknown
}

// ConvertIpLibIsp ip库运营商转化
func ConvertIpLibIsp(s string) string {
	switch s {
	case "ChinaTelecom":
		return "dx"
	case "ChinaUnicom":
		return "lt"
	case "ChinaMobile":
		return "yd"
	default:
		return s
	}
}

type IP interface {
	string | net.IP | *net.IP
}

// IsIPv4 是否是ipv4地址
func IsIPv4[T IP](ip T) bool {
	switch v := (interface{})(ip).(type) {
	case string:
		return net.ParseIP(v).To4() != nil
	case net.IP:
		return v.To4() != nil
	case *net.IP:
		return v.To4() != nil
	}
	return false
}

func isIPv6(ip net.IP) bool {
	if ip == nil {
		return false
	}
	return ip.To4() == nil && ip.To16() != nil && strings.Contains(ip.String(), ":")
}

// IsIPv6 是否是ipv6地址
func IsIPv6[T IP](ip T) bool {
	switch v := (interface{})(ip).(type) {
	case string:
		return isIPv6(net.ParseIP(v))
	case net.IP:
		return isIPv6(v)
	case *net.IP:
		return isIPv6(*v)
	}
	return false
}

// GetIPv4NetMask 获取ipv4子网掩码对象
func GetIPv4NetMask[T IP](ip T, n uint8) *net.IPNet {
	if n < 0 || n > 32 {
		return nil
	}
	var nip net.IP
	switch v := (interface{})(ip).(type) {
	case string:
		nip = net.ParseIP(v)
	case net.IP:
		nip = v
	case *net.IP:
		nip = *v
	default:
		return nil
	}
	if nip.To4() == nil {
		return nil
	}
	_, ipNet, _ := net.ParseCIDR(fmt.Sprintf("%s/%d", nip.To4().String(), n))
	return ipNet
}

// IPv4ToUnt32 ipv4转uint32
//
// exp:
// 1.1.1.1
// 0b 00000001 00000001 00000001 00000001
// decimal: 16843009
func IPv4ToUnt32(ip net.IP) uint32 {
	if ip == nil {
		return 0
	}
	ipv4 := ip.To4()
	if ipv4 == nil {
		return 0
	}
	return uint32(ipv4[0])<<24 + uint32(ipv4[1])<<16 + uint32(ipv4[2])<<8 + uint32(ipv4[3])
}

// Uint32ToIPv4 uint32转ipv4
func Uint32ToIPv4(v uint32) net.IP {
	return net.IPv4(byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
}
