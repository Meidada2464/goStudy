/**
 * Package util
 * @Author fuqiang.li <fuqiang.li@baishan.com>
 * @Date 2023/3/22
 */

package util

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// StructDeepCopy 结构体深拷贝，可以看看是否有比json更高效的方法
func StructDeepCopy(src, dst interface{}) error {
	if src == nil || dst == nil || reflect.ValueOf(dst).IsZero() || reflect.ValueOf(dst).IsNil() {
		return fmt.Errorf("must no nil")
	}
	data, err := json.Marshal(&src)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &dst)
}

// Struct2Map 结构体转map结构，比json效率高好几倍
func Struct2Map(data interface{}) map[string]string {
	defer func() {
		recover()
	}()
	if data == nil {
		return nil
	}
	var (
		t   = reflect.TypeOf(data)
		res = make(map[string]string)
	)
	// 遍历结构体的字段
	for i := 0; i < t.NumField(); i++ {
		var (
			field = t.Field(i)
			value = reflect.ValueOf(data).Field(i)
			name  = field.Name
			// 不支持name,omitempty这种，后续按需扩展
			tagName = strings.TrimSpace(field.Tag.Get("json"))
		)
		if tagName != "" {
			name = tagName
		}
		res[name] = value.String()
	}
	return res
}
