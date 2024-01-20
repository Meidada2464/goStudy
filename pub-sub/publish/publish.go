package publish

import (
	"sync"
	"time"
)

var (
	once     sync.Once
	instance *Publisher
)

// 创建两种类型
type (
	subscriber chan interface{}         //订阅者为一个通道
	topicFunc  func(v interface{}) bool // 主题为一个过滤器
)

// Publisher 创建发布者
type Publisher struct {
	m           sync.RWMutex             // 读写锁
	buffer      int                      // 订阅队列的缓存大小
	timeout     time.Duration            // 发布的超时时间
	subscribers map[subscriber]topicFunc //订阅者信息
}

// NewPublisher 创建一个发布器
func NewPublisher(publishTimeout time.Duration, buffer int) *Publisher {
	once.Do(func() {
		instance = &Publisher{
			timeout:     publishTimeout,
			buffer:      buffer,
			subscribers: make(map[subscriber]topicFunc),
		}
	})
	return instance
}

// AddSubscribe 添加一个订阅者用于订阅全部的主题
func (p *Publisher) AddSubscribe() chan interface{} {
	return p.SubscribeTopic(nil)
}

// SubscribeTopic 添加一个新的订阅者，订阅指定的主题
func (p *Publisher) SubscribeTopic(topic topicFunc) chan interface{} {
	ch := make(chan interface{}, p.buffer)
	p.m.Lock()
	defer p.m.Unlock()
	p.subscribers[ch] = topic
	return ch
}

// Evict 退出订阅
func (p *Publisher) Evict(sub chan interface{}) {
	p.m.Lock()
	defer p.m.Unlock()

	delete(p.subscribers, sub)
	close(sub)
}

// Publish 发布一个主题
func (p *Publisher) Publish(v interface{}) {
	p.m.Lock()
	defer p.m.Unlock()

	var wg sync.WaitGroup
	for sub, topic := range p.subscribers {
		wg.Add(1)
		// 发布一个主题
		go p.sendTopic(sub, topic, v, &wg)
	}
	wg.Wait()
}

// 向所有订阅者中发布内容
func (p *Publisher) sendTopic(
	sub subscriber, topic topicFunc, v interface{}, wg *sync.WaitGroup,
) {
	defer wg.Done()
	if topic != nil && !topic(v) {
		return
	}

	select {
	case sub <- v:
	case <-time.After(p.timeout):
	}
}

func (p *Publisher) Close() {
	p.m.Lock()
	defer p.m.Unlock()

	for sub := range p.subscribers {
		delete(p.subscribers, sub)
	}
}
