package passwordTest

import (
	"fmt"
	"regexp"
	"testing"
	"unicode"
)

func TestPsw(t *testing.T) {
	password := "Aini19990518"
	if ValidatePassword(password) {
		fmt.Println("密码有效")
	} else {
		fmt.Println("密码在8到20之间且须包含大小写字母和数字, 不能包含中文、全角字符、特殊字符")
	}
}

// ValidatePassword 验证密码是否符合规则
func ValidatePassword(password string) bool {
	// 检查是否包含中文字符
	isMatchHz := containsChinese(password)

	// 检查是否包含大小写字母和数字，且长度在8到20之间
	isMatch, _ := regexp.MatchString(`(?=.*[a-z])(?=.*[A-Z])(?=.*\d).{8,20}$`, password)

	// 检查是否不包含全角字符
	if isMatch && !isMatchHz && ContainQj(password) {
		return true
	} else {
		return false
	}
}

// containsChinese 检查字符串是否包含中文字符
func containsChinese(str string) bool {
	for _, r := range str {
		if unicode.Is(unicode.Han, r) {
			return true
		}
	}
	return false
}

// ContainQj 检查字符串是否包含全角字符 (返回 true 表示没有全角字符)
func ContainQj(str string) bool {
	for _, r := range str {
		if (r >= 65248 && r <= 65535) || r == 12288 {
			return false
		}
	}
	return true
}
