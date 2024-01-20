package test

import (
	"fmt"
	"goStudy/pub-sub/publish"
	"strings"
	"testing"
	"time"
)

// 进行订阅/发布者模式的测试
func TestPubSub(t *testing.T) {
	// 创建一个发布者
	p := publish.NewPublisher(100*time.Millisecond, 10)
	defer p.Close()

	// 添加订阅者，订阅所有主题
	all := p.AddSubscribe()
	// 添加接收具体主题的订阅者
	golang := p.SubscribeTopic(func(v interface{}) bool {
		if s, ok := v.(string); ok {
			return strings.Contains(s, "golang")
		}
		return false
	})

	// 添加接收具体主题的订阅者
	greet := p.SubscribeTopic(func(v interface{}) bool {
		if s, ok := v.(string); ok {
			return strings.Contains(s, "are you ok")
		}
		return false
	})
	// 开始发布内容

	p.Publish("hello, word!")
	p.Publish("hello golang, are you ok?")

	// 输出订阅者接收到的内容
	go func() {
		for msg := range all {
			fmt.Println("all:", msg)
		}
	}()

	go func() {
		for msg := range golang {
			fmt.Println("golang:", msg)
		}
	}()

	go func() {
		for msg := range greet {
			fmt.Println("greet:", msg)
		}
	}()

	time.Sleep(3 * time.Second)
}
