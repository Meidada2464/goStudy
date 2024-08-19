/**
 * Package util
 * @Author fuqiang.li <fuqiang.li@baishan.com>
 * @Date 2023/7/13
 */

package util

import "testing"

type isOverseasTest struct {
	arg      string
	expected bool
}

var isOverseasTests = []isOverseasTest{
	{"AE-xx", true},
	{"aE-xx", false},
	{"aa-xx", false},
	{"", false},
	{"-", false},
	{"1-", false},
	{"A-", true},
}

func TestIsOverseasEp(t *testing.T) {
	for _, test := range isOverseasTests {
		if output := IsOverseasEp(test.arg); output != test.expected {
			t.Errorf("Output %t not equal to expected %t", output, test.expected)
		}
	}
}
