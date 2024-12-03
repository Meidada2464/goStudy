/**
 * Package TemplentTest
 * @Author fengfeng.mei <fengfeng.mei@baishan.com>
 * @Date 2024/11/15 16:33
 */

package TemplentTest

import (
	"fmt"
	"testing"
)

func TestMapErr(t *testing.T) {
	var alreadys = make(map[string]bool)
	alreadys["aa"] = true
	if _, ok := alreadys["aa"]; ok {
		fmt.Println("result1:", alreadys[""])
	}
	fmt.Println("result2:", alreadys[""])
}
