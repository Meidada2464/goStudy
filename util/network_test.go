/**
 * Package util
 * @Author fuqiang.li <fuqiang.li@baishan.com>
 * @Date 2024/3/7
 */

package util

import "testing"

func TestGetNICs(t *testing.T) {
	_, err := GetNICs()
	if err != nil {
		t.Errorf("GetNICs() failed, err: %v", err)
	}
}
