package SimpleFactory

import "fmt"

// 这里来做一个简单工厂模式的示例
// 简单工厂模式是指，由接口定义工厂和产品的接口。工厂接口的返回对象是产品接口，只有实现了产品接口的对象会被返回
// 优点：解除了工厂和产品的耦合度，使产品和工厂解耦

// 实例

type Factory interface {
	// FactoryMethod 简单的工厂，解耦了工厂和产品
	FactoryMethod() Product
}

type Product interface {
	// Use 产品必选要有的特征
	Use()
}

// ConcreteFactory 开始实例化一个工厂
type ConcreteFactory struct {
}

type ConcreteProduct struct {
}

// CarFactory 再创建一个工厂，用来生产汽车
type CarFactory struct {
}

// CarProduct 汽车产品
type CarProduct struct {
}

func (c *CarFactory) FactoryMethod() Product {
	return &CarProduct{}
}

func (c CarProduct) Use() {
	fmt.Println("汽车产品！")
}

// ConcreteFactory具体工厂实现Factory接口

func (c *ConcreteFactory) FactoryMethod() Product {
	// 这里为什么一定要是地址？
	return &ConcreteProduct{}
}

// Use 定义产品，产品也实现product接口
func (c *ConcreteProduct) Use() {
	fmt.Println("具体的产品")
}
