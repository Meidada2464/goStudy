package SimpleFactory

import "testing"

func TestNewRestaurant(t *testing.T) {
	NewRestaurant("first").GetFood()
	NewRestaurant("second").GetFood()
}

// 简单工厂模式就多个具体对象是实现一个通用接口的方法，并实现其方法，类似于多态的模式
