package Singleton

import (
	"fmt"
	"sync"
)

// 这是一个单例模式，属于对象创建模式

// 该模式的重点是：
// 1、某个类只能有一个实例 2、该类只能是自己创建的（因此要提供一个访问该实例的方法） 3、该类必须向整个系统提供使用

// 单例模式包含的角色： 1、Singleton:单例

// Singleton 提供一个接口供单例的实例结构体进行实现，并且内部的实现方法都是私有的
type Singleton interface {
	setBrand()
	setEngine()
	setPrice()
}

// 定义一个私有的、服务于全局的结构体，里面的属性是私有的，外部无法赋值，需要实现接口中的方法赋值
type singleton struct {
	brand  string
	engine string
	price  string
}

func (s *singleton) setBrand() {
	s.brand = "玉皇大帝"
}

func (s *singleton) setEngine() {
	s.engine = "v11"
}

func (s *singleton) setPrice() {
	s.price = "￥209999"
}

type singleton2 struct {
	brand  string
	engine string
	price  string
}

func (s *singleton2) setBrand() {
	s.brand = "四驱达车"
}

func (s *singleton2) setEngine() {
	s.engine = "v13"
}

func (s *singleton2) setPrice() {
	s.price = "￥209998"
}

// 提供一个实现该犯法的对外接口，并且使用 sync.Once保证其之初始化一次

var (
	instance1 *singleton
	instance2 *singleton2
	once1     sync.Once
	once2     sync.Once
)

// GetInstance 实例化对象的统一对外的方法,这里返回值使用的是借口，可以实现创建多个需要创建的单例对象
func GetInstance(s string) Singleton {
	// once 对象只能用一次
	if s == "singleton1" {
		once1.Do(func() {
			instance1 = &singleton{}
		})
		return instance1
	}

	if s == "singleton2" {
		once2.Do(func() {
			instance2 = &singleton2{}
		})
		return instance2
	}

	fmt.Println("未能找到需要单例的对象")
	return nil
}
