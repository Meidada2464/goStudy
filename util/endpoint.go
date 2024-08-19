/**
 * Package util
 * @Author fuqiang.li <fuqiang.li@baishan.com>
 * @Date 2023/7/12
 */

package util

import (
	"os"
	"strings"
	"unicode"
)

const SpecialDefineHostnamePath = "/allconf/hostname.conf"

// GetHostname 获取主机名：
//
// 1. 从/allconf/hostname.conf获取
//
// 2. 从hostname获取
func GetHostname() string {
	hostname, _ := os.Hostname()
	rawData, err := os.ReadFile(SpecialDefineHostnamePath)
	if err != nil {
		return hostname
	}
	data := strings.ReplaceAll(string(rawData), "hostname=", "")
	data = strings.TrimSpace(data)
	if data != "" {
		return data
	}
	return hostname
}

func IsOverseasEp(s string) bool {
	if s == "" {
		return false
	}
	ss := strings.Split(s, "-")
	if len(ss) < 1 || len(ss[0]) == 0 {
		return false
	}
	for _, r := range ss[0] {
		if !unicode.IsLetter(r) || !unicode.IsUpper(r) {
			return false
		}
	}
	return true
}
