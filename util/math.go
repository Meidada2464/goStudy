/**
 * Package util
 * @Author fuqiang.li <fuqiang.li@baishan.com>
 * @Date 2024/1/29
 */

package util

import "math"

type float interface {
	~float32 | ~float64
}

// ToFixed rounds down to a float to a given number of decimal places
//
// digits must be greater than or equal to 0
func ToFixed[T float](n T, digits int) T {
	if digits < 0 {
		digits = 0
	}
	pow := math.Pow10(digits)
	i := int64(float64(n) * pow)
	return T(float64(i) / pow)
}
