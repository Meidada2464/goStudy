/**
 * Package util
 * @Author fuqiang.li <fuqiang.li@baishan.com>
 * @Date 2023/4/17
 */

package util

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5Str(s string) string {
	return Md5Bytes([]byte(s))
}

func Md5Bytes(b []byte) string {
	h := md5.New()
	h.Write(b)
	return hex.EncodeToString(h.Sum(nil))
}
