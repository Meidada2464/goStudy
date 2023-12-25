package ActualCombat

import (
	"fmt"
	"testing"
)

func TestClothes(t *testing.T) {
	var fc GYClothesStore
	c1 := fc.MakeClothes("anta")
	c2 := fc.MakeClothes("peak")

	fmt.Println("clothes1:", c1)
	fmt.Println("clothes2:", c2)
}
