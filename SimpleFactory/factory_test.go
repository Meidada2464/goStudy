package SimpleFactory

import "testing"

func TestNewRestaurant(t *testing.T) {
	// 创建一个工厂
	var CF ConcreteFactory
	CF.FactoryMethod().Use()

	// 汽车产品
	var CarFactory CarFactory
	CarFactory.FactoryMethod().Use()

}
