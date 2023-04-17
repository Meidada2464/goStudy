package Singleton

import (
	"fmt"
	"testing"
)

func TestGetInstance(t *testing.T) {
	s1 := GetInstance("singleton1")
	s1.setPrice()
	s1.setEngine()
	s1.setBrand()
	fmt.Println("创建了一个单例对象：", s1)

	// 再创建一个该结构体的对象：
	s2 := GetInstance("singleton1")
	fmt.Println("又创建了一个s1单例对象，但是不对值进行初始化：", s2)

	s3 := GetInstance("singleton2")
	s3.setPrice()
	s3.setEngine()
	s3.setBrand()
	fmt.Println("创建了一个s3单例对象：", s3)

	s4 := GetInstance("singleton2")
	fmt.Println("又创建了一个s3单例对象，但是不对值进行初始化：", s4)

	s5 := GetInstance("singleton3")
	fmt.Println("创建了一个s5单例对象：", s5)

	// 结果：
	// 创建了一个单例对象： &{玉皇大帝 v11 ￥209999}
	// 又创建了一个单例对象，但是不对值进行初始化： &{玉皇大帝 v11 ￥209999}

	// 原因： 由sync.Once保证了创建的对象的唯一性

}
