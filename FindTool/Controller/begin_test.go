package Controller

import (
	"fmt"
	"goStudy/FindTool/Service"
	"testing"
)

func TestFindHost(t *testing.T) {
	Service.FindHost()
	fmt.Println("hello")
}
