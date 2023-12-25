package Timetest

import (
	"fmt"
	"sync"
	"testing"
)

func TestIsInTime(t *testing.T) {
	var s = &RuleStrategy{
		RunBegin: "00:01,01:01",
		RunEnd:   "00:14,02:01",
	}

	res := s.IsInTime(1682356538)

	fmt.Println("res:", res)

}

type testStruct struct {
	Age *int
}

func TestSlice(t *testing.T) {

	res := make(map[int]testStruct)
	ageTmp := 18

	per := testStruct{
		Age: &ageTmp,
	}

	for i := 0; i < 5; i++ {
		per.Age = &i
		res[i] = per
		fmt.Printf("lastRes: %p \n", &per)
	}

	fmt.Println("lastRes:", res)
}

func TestRegexp(t *testing.T) {
	var map_test sync.Map
	map_test.Store("bgp-beijing-beijing-1-123-59-102-14", "ddddddddddddddddddddddddd")
	actual, loaded := map_test.LoadOrStore("bgp-beijing-beijing-1-123-59-102-14", "d5a54edd8c80411593eb1420dc13b774")
	map_test.Store("bgp-beijing-beijing-1-123-59-102-14", "d5a54edd8c80411593eb1420dc13b774")

	fmt.Println("actual", actual)
	fmt.Println("loaded", loaded)
}
