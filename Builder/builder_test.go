package Builder

import (
	"fmt"
	"testing"
)

func TestNewDirector(t *testing.T) {
	// 创建一个空的需要创建的复杂结构体
	p1 := NewProduct1()
	fmt.Println("创建p1模具：", p1)
	// 创建一个指挥官，即定义需要创建这个空结构体p1的步骤的指挥官
	d1 := NewDirector(p1)
	fmt.Println("p1开始组装：", p1)
	d1.Construct("凯迪拉克", "v12", "199999")
	fmt.Println("p1组装完成：", p1)

	// 创建一个空的需要创建的复杂结构体
	p2 := NewProduct2()
	fmt.Println("创建p2模具：", p2)
	// 创建一个指挥官，即定义需要创建这个空结构体p2的步骤的指挥官
	d2 := NewDirector(p2)
	fmt.Println("p2开始组装：", p2)
	d2.Construct("小鹏汽车", "超级电动引擎", "299998")
	fmt.Println("p2组装完成：", p2)
}
