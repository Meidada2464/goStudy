/**
 * Package producer
 * @Author fengfeng.mei <fengfeng.mei@baishan.com>
 * @Date 2024/10/23 14:20
 */

package main

import (
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 配置 Sarama AsyncProducer
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true // 设置为 true 以便接收成功消息的回调

	// 连接 Kafka 代理
	brokers := []string{"localhost:9092"} // 替换为实际的 Kafka broker 地址
	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("无法创建 Kafka 异步生产者: %v", err)
	}
	defer producer.Close()

	// 捕获中断信号以便优雅退出
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// 发送消息到 Kafka
	topic := "test-topic"
	go func() {
		for i := 0; i < 10; i++ {
			message := fmt.Sprintf("消息编号: %d", i)
			producer.Input() <- &sarama.ProducerMessage{
				Topic: topic,
				Value: sarama.StringEncoder(message),
			}
			fmt.Printf("发送消息: %s\n", message)
		}
	}()

	// 处理发送成功的信息
	//go func() {
	//	for success := range producer.Successes() {
	//		fmt.Printf("消息成功发送，分区: %d, 偏移量: %d, 内容: %s\n, topic.",
	//			success.Partition, success.Offset, success.Value)
	//	}
	//}()

	// 处理发送失败的错误
	go func() {
		for err := range producer.Errors() {
			fmt.Printf("发送消息出错: %v\n", err)
		}
	}()

	go func() {
		for {
			select {
			case ok := <-producer.Successes():
				fmt.Println("aaaaaaaaaaaaaaaaa")
				fmt.Println("topic", ok.Topic)
				if ok != nil && ok.Topic == "" {
					fmt.Printf("消息成功发送，分区: %d, 偏移量: %d, 内容: %s\n",
						ok.Partition, ok.Offset, ok.Value)
					fmt.Println("bbbbbbbbbbbbbbbbbb")
				}
			}
		}
	}()

	// 等待信号优雅退出
	<-signals
	fmt.Println("生产者已退出")
}
