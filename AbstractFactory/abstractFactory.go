package AbstractFactory

import "fmt"

type Lunch interface {
	Cook()
}

type Rice struct {
}

func (r *Rice) Cook() {
	fmt.Println("This is Rice")
}

type Tomato struct {
}

func (t *Tomato) Cook() {
	fmt.Println("This is a nice tomato")
}

type LunchFactory interface {
	CreateRice() Lunch
	CreateTomato() Lunch
}

type SimpleFactory struct {
}

func (s *SimpleFactory) CreateRice() Lunch {
	return &Rice{}
}

func (s *SimpleFactory) CreateTomato() Lunch {
	return &Tomato{}
}

type SimpleFactory2 struct {
}

func (s *SimpleFactory2) CreateRice() Lunch {
	fmt.Println("factory2 做的东北大米")
	return &DongBeiRice{}
}

func (s *SimpleFactory2) CreateTomato() Lunch {
	return &DongBeiTomato{}
}

type DongBeiRice struct {
}

func (d *DongBeiRice) Cook() {
	fmt.Println("这是一碗香喷喷的东北大米！")
}

type DongBeiTomato struct {
}

func (d *DongBeiTomato) Cook() {
	fmt.Println("这是东北tomato！")
}

// NewSimpleFactory 返回值是一个接口，表示无论什么商家的，只要是实现了该接口的商家都可返回
func NewSimpleFactory(s string) LunchFactory {
	switch s {
	case "factory1":
		return &SimpleFactory{}
	case "DongBei":
		return &SimpleFactory2{}
	}
	return nil
}
