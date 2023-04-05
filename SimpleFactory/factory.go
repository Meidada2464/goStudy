package SimpleFactory

import "fmt"

// 这里来做一个简单工厂模式的示例

// 接口

type Restaurant interface {
	// GetFood 提供一个抽象的方法
	GetFood()
}

// DinnerRoomFirst 具体工厂1
type DinnerRoomFirst struct {
}

// DinnerRoomerSecond 具体工厂2
type DinnerRoomerSecond struct {
}

// 对接口方法进行了实现

func (dF *DinnerRoomFirst) GetFood() {
	fmt.Println("DinnerRoomerFirst提供的食物")
}

func (dF *DinnerRoomerSecond) GetFood() {
	fmt.Println("DinnerRoomerSecond提供的食物")
}

// NewRestaurant 抽象工厂,输入需要的食物，并提供
func NewRestaurant(name string) Restaurant {
	switch name {
	case "first":
		return &DinnerRoomFirst{}
	case "second":
		return &DinnerRoomerSecond{}
	}
	return nil
}
