package Builder

import (
	"fmt"
	"testing"
)

func TestNewDirector(t *testing.T) {
	// 先创造一个生成器，初始化产品
	builder := NewConcreteBuilder()
	// 创建一个主管，用于说明需要如何、怎样初始化一个产品
	director := NewDirector(&builder)
	// 根据规定好步骤去封装一个产品
	director.Construct()
	result := builder.GetResult()

	fmt.Println(result)
}
