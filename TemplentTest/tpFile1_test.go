package TemplentTest

import (
	"fmt"
	"net/http"
	"testing"
)

func TestName1(t *testing.T) {
	http.HandleFunc("/", SayHello1)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		fmt.Println("HTTP server failed,err:", err)
		return
	}
}
