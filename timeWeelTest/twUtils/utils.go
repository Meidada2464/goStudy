/**
 * Package twUtils
 * @Author fengfeng.mei <fengfeng.mei@baishan.com>
 * @Date 2024/11/30 21:00
 */

package twUtils

import "time"

func GetInterval(t int) time.Duration {
	return time.Duration(t) * time.Second
}
