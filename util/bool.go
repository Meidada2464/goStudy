/**
 * Package util
 * @Author fuqiang.li <fuqiang.li@baishan.com>
 * @Date 2023/7/17
 */

package util

func NewTrue() *bool {
	v := true
	return &v
}

func NewFalse() *bool {
	v := false
	return &v
}
