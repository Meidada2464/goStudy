package handleTest

import (
	"fmt"
	"testing"
)

type ChinaPeople func(face string, age int) string

func (c ChinaPeople) name() {

}

type people interface {
	name()
}

func Address(p people) {

}

func Close() {
	fmt.Println("Close")
}
func Close2() {
	fmt.Println("Close2")
}
func Close3() string {
	fmt.Println("Close3")
	return ""
}

//func TestHandle(t *testing.T) {
//	defer Close3()
//	defer Close()
//	defer Close2()
//
//	recover()
//	panic("error")
//}

func sum(max int) int {
	total := 0
	for i := 0; i < max; i++ {
		total += i
	}
	return total
}

func fooWithDefer() {
	defer func() {
		sum(10)
	}()
}

func fooWithoutDefer() {
	sum(10)
}

func BenchmarkFooWithDefer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fooWithDefer()
	}
}

func BenchmarkFooWithoutDefer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fooWithoutDefer()
	}
}
