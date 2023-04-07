package AbstractFactory

import (
	"fmt"
	"testing"
)

// 抽象工厂就是可以自定义做饭的工厂，在简单工厂多加了一层工厂层
func TestNewSimpleFactory(t *testing.T) {
	// 选工厂
	f1 := NewSimpleFactory("factory1")
	f1.CreateRice().Cook()
	f1.CreateTomato().Cook()

	fmt.Println("=========================")

	f2 := NewSimpleFactory("DongBei")
	f2.CreateRice().Cook()
	f2.CreateTomato().Cook()
}
