/**
 * Package utils
 * @Author fuqiang.li <fuqiang.li@baishan.com>
 * @Date 2022/5/16
 */

package utils

import "os"

// IsFileExist 判断文件是否存在
func IsFileExist(path string) bool {
	return IsPathExist(path)
}

// IsPathExist 判断目录是否存在
func IsPathExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
