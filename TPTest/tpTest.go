package TPTest

//
//import (
//	"fmt"
//	"html/template"
//	"net/http"
//	"net/url"
//)
//
//type UserInfo struct {
//	Name stringl
//}
//
//func SayHello(w http.ResponseWriter, r *http.Request) {
//
//	tpFunc := template.FuncMap{
//		"kua": func(arg string) (string, error) {
//			return arg + "真帅", nil
//		},
//		"func2": func(arg string) (string, error) {
//			return arg + "靓仔", nil
//		},
//		"urlencode": func(value interface{}) string {
//			return url.QueryEscape(value.(string))
//		},
//	}
//	//解析指定文件生成模板对象
//
//	tmpl, err := template.ParseFiles("./hello.tmpl")
//	if err != nil {
//		fmt.Println("create template failed, err:", err)
//		return
//	}
//	// 利用给定数据渲染模板，并将结果写入w
//	err2 := tmpl.Execute(w, "沙河小王子")
//	if err2 != nil {
//		return
//	}
//}
