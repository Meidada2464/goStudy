package utils

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLocalIP(t *testing.T) {
	Convey("test-localip", t, func() {
		ip := LocalIP()
		So(ip, ShouldNotBeEmpty)
		So(ip, ShouldNotEqual, "127.0.0.1")
		t.Log(ip)

		ips := LocalIPs()
		So(len(ips), ShouldBeGreaterThanOrEqualTo, 1)
		for _, ip := range ips {
			So(ip, ShouldNotEqual, "127.0.0.1")
		}
		t.Log(ips)
	})
}
