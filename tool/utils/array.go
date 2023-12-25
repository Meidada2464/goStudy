/**
 * Package utils
 * @Author fuqiang.li <fuqiang.li@baishan.com>
 * @Date 2022/3/21
 */

package utils

// IntFind 确认数组中是否包含对应的int
func IntFind(ss []int, subStr int) bool {
	for _, s := range ss {
		if s == subStr {
			return true
		}
	}
	return false
}
