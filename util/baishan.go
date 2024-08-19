/**
 * Package util
 * @Author fuqiang.li <fuqiang.li@baishan.com>
 * @Date 2023/6/14
 */

package util

import (
	"net"
	"strings"
	"unicode"
)

// IsValidEp 是否合法的主机名
func IsValidEp(s string) bool {
	for _, r := range s {
		if unicode.Is(unicode.Latin, r) ||
			unicode.IsDigit(r) ||
			r == '-' {
			return true
		}
	}
	return false
}

// IsValidMultiIPv4 是否合法的ipv4，用,拼接
func IsValidMultiIPv4(s string) bool {
	for _, ip := range strings.Split(s, ",") {
		if net.ParseIP(ip).To4() == nil {
			return false
		}
	}
	return false
}

// IsValidMultiIP 是否合法的ip，用,拼接
func IsValidMultiIP(s string) bool {
	for _, ip := range strings.Split(s, ",") {
		if net.ParseIP(ip) == nil {
			return false
		}
	}
	return true
}

// IsValidMultiSvcType 是否合法的服务类型，用,拼接
func IsValidMultiSvcType(s string) bool {
	for _, r := range s {
		if unicode.Is(unicode.Latin, r) ||
			unicode.IsDigit(r) ||
			r == '-' || r == '_' || r == ',' {
			return true
		}
	}
	return false
}

// IsValidMultiISP 是否合法的isp，用,拼接
func IsValidMultiISP(s string) bool {
	for _, r := range s {
		if unicode.Is(unicode.Latin, r) ||
			unicode.IsDigit(r) || r == ',' {
			return true
		}
	}
	return false
}

// IsValidMultiProv 是否合法的省份
func IsValidMultiProv(s string) bool {
	for _, r := range s {
		if unicode.Is(unicode.Latin, r) ||
			unicode.IsDigit(r) || r == ' ' {
			return true
		}
	}
	return false
}
