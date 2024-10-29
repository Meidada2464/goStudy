/**
 * Package netTest
 * @Author fengfeng.mei <fengfeng.mei@baishan.com>
 * @Date 2024/10/24 17:13
 */

package netTest

import (
	"fmt"
	"net/url"
	"testing"
)

func TestNetPaK(t *testing.T) {
	u, err := url.Parse("http://mon.local/150K.dat")
	if err != nil {
		return
	}
	fmt.Println("u", u)

	host := u.Hostname()
	portStr := u.Port()
	scheme := u.Scheme
	fmt.Println("host", host, "port", portStr, "scheme", scheme)
}
