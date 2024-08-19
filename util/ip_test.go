/**
 * Package util
 * @Author fuqiang.li <fuqiang.li@baishan.com>
 * @Date 2024/3/12
 */

package util

import (
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func TestGetIpFamily(t *testing.T) {
	t.Run("IPv4Address", func(t *testing.T) {
		result := GetIpFamily(net.ParseIP("192.168.1.1"), "IPv4", "IPv6", "Unknown")
		assert.Equal(t, "IPv4", result)
	})

	t.Run("IPv6Address", func(t *testing.T) {
		result := GetIpFamily(net.ParseIP("2001:db8::2:1"), "IPv4", "IPv6", "Unknown")
		assert.Equal(t, "IPv6", result)
	})

	t.Run("UnknownAddress", func(t *testing.T) {
		result := GetIpFamily(net.ParseIP("unknown"), "IPv4", "IPv6", "Unknown")
		assert.Equal(t, "Unknown", result)
	})
}

func TestConvertIpLibIsp(t *testing.T) {
	t.Run("ChinaTelecom", func(t *testing.T) {
		result := ConvertIpLibIsp("ChinaTelecom")
		assert.Equal(t, "dx", result)
	})

	t.Run("ChinaUnicom", func(t *testing.T) {
		result := ConvertIpLibIsp("ChinaUnicom")
		assert.Equal(t, "lt", result)
	})

	t.Run("ChinaMobile", func(t *testing.T) {
		result := ConvertIpLibIsp("ChinaMobile")
		assert.Equal(t, "yd", result)
	})

	t.Run("Other", func(t *testing.T) {
		result := ConvertIpLibIsp("Other")
		assert.Equal(t, "Other", result)
	})
}

func TestIsIPv4(t *testing.T) {
	t.Run("ValidIPv4", func(t *testing.T) {
		result := IsIPv4("192.168.1.1")
		assert.True(t, result)
	})

	t.Run("InvalidIPv4", func(t *testing.T) {
		result := IsIPv4("2001:db8::2:1")
		assert.False(t, result)
	})
}

func TestIsIPv6(t *testing.T) {
	t.Run("ValidIPv6", func(t *testing.T) {
		result := IsIPv6("2001:db8::2:1")
		assert.True(t, result)
	})

	t.Run("InvalidIPv6", func(t *testing.T) {
		result := IsIPv6("192.168.1.1")
		assert.False(t, result)
	})
}

func TestGetIPv4NetMask(t *testing.T) {
	t.Run("ValidNetMask", func(t *testing.T) {
		result := GetIPv4NetMask("192.168.1.1", 24)
		assert.NotNil(t, result)
	})

	t.Run("InvalidNetMask", func(t *testing.T) {
		result := GetIPv4NetMask("192.168.1.1", 33)
		assert.Nil(t, result)
	})

	t.Run("InvalidIP", func(t *testing.T) {
		result := GetIPv4NetMask("2001:db8::2:1", 24)
		assert.Nil(t, result)
	})
}

func TestIPv4ToUnt32(t *testing.T) {
	t.Run("ValidIPv4", func(t *testing.T) {
		result := IPv4ToUnt32(net.ParseIP("192.168.1.1"))
		assert.Equal(t, uint32(3232235777), result)
	})

	t.Run("InvalidIPv4", func(t *testing.T) {
		result := IPv4ToUnt32(net.ParseIP("2001:db8::2:1"))
		assert.Equal(t, uint32(0), result)
	})
}

func TestUint32ToIPv4(t *testing.T) {
	t.Run("ValidUint32", func(t *testing.T) {
		result := Uint32ToIPv4(3232235777)
		assert.Equal(t, net.ParseIP("192.168.1.1"), result)
	})
}
