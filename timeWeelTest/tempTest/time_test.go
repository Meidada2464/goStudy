package tempTest

import (
	"fmt"
	"testing"
	"time"
)

func TestTimeTest(t *testing.T) {
	t1 := time.Duration(60) * time.Millisecond
	fmt.Println("t1:", t1)

	fmt.Println("t1.Seconds():", t1.Seconds())

	i1 := int(t1.Seconds())

	t2 := time.Second.Seconds()
	fmt.Println("t2:", t2)

	i2 := 0 % 1

	fmt.Println("i1:", i1)
	fmt.Println("i2:", i2)

}
