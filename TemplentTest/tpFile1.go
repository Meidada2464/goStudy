package TemplentTest

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
)

type UserInfo struct {
	Name   string
	Gender string
	Age    int
}

func SayHello1(w http.ResponseWriter, r *http.Request) {
	tpFunc := template.FuncMap{
		"kua": func(arg string) (string, error) {
			return arg + "真帅", nil
		},
		"func2": func(arg string) (string, error) {
			return arg + "靓仔", nil
		},
		"urlencode": func(value interface{}) string {
			if value == nil {
				return ""
			}
			if value.(string) == "" {
				return ""
			}
			return url.QueryEscape(value.(string))
		},
	}

	htmlByte, err := ioutil.ReadFile("./hello.tmpl")
	if err != nil {
		fmt.Println("read html failed, err:", err)
		return
	}

	// 采用链式操作在Parse之前调用Funcs添加自定义的kua函数
	tmpl, err := template.New("hello").Funcs(tpFunc).Parse(string(htmlByte))
	if err != nil {
		fmt.Println("create template failed, err:", err)
		return
	}

	user := UserInfo{
		Name:   "",
		Gender: "男",
		Age:    18,
	}

	// 使用user渲染模板，并将结果写入w
	err1 := tmpl.Execute(w, user)
	if err1 != nil {
		return
	}

}
