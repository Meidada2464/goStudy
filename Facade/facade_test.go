package Facade

import (
	"fmt"
	"testing"
)

func TestNewCarFacade(t *testing.T) {
	car1 := NewCarFacade()
	car1.CreateNewCompleteCar("凯迪拉克", "V12", "后驱车")

	car2 := NewCarFacade()
	car2.CreateNewCompleteCar("别克", "V6", "CT5")

	fmt.Println("car1:", car1, "car2:", car2)
}
