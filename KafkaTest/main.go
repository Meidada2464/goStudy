package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/segmentio/kafka-go"
	"time"
)

// 使用kafka-go操作kafka
func main() {
	ctx := context.Background()
	go writeKafka(ctx)
	readKafka(ctx)
}

var (
	reader *kafka.Reader
	topic  = "mallard2-store-bigdata"
)

func writeKafka(ctx context.Context) {
	// 新建kafka链接
	producer := &kafka.Writer{
		Addr:                   kafka.TCP("127.0.0.1:9092"),
		Topic:                  topic,
		WriteTimeout:           1 * time.Second,
		Balancer:               &kafka.Hash{},
		RequiredAcks:           kafka.RequireNone,
		AllowAutoTopicCreation: true,
	}
	defer producer.Close()

	// 3次尝试的机会，因为首次创建失败必定失败，因为会创建topic
	for {
		// 写入一个
		err := producer.WriteMessages(
			ctx,
			kafka.Message{Key: []byte("1"), Value: []byte("B")},
			kafka.Message{Key: []byte("2"), Value: []byte("I")},
			kafka.Message{Key: []byte("3"), Value: []byte("O")},
			kafka.Message{Key: []byte("4"), Value: []byte("P")},
			kafka.Message{Key: []byte("5"), Value: []byte("H")},
			kafka.Message{Key: []byte("6"), Value: []byte("I")},
			kafka.Message{Key: []byte("7"), Value: []byte("L")},
			kafka.Message{Key: []byte("8"), Value: []byte("I")},
			kafka.Message{Key: []byte("9"), Value: []byte("A")},
		)
		if err != nil {
			if errors.Is(err, kafka.LeaderNotAvailable) {
				time.Sleep(500 * time.Millisecond)
				fmt.Println("等待leader选举完成,err:", err)
			} else {
				fmt.Println("批量写入kafka失败,err:", err)
			}
		}
		fmt.Println("批量写入kafka成功")
		time.Sleep(60 * time.Second)
	}
}

func readKafka(ctx context.Context) {
	consumer := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{"127.0.0.1:9092"},
		Topic:          topic,
		GroupID:        "mallard2-store-bigdata",
		CommitInterval: 1 * time.Second,
		StartOffset:    kafka.FirstOffset,
	})

	for {
		if message, err := consumer.ReadMessage(ctx); err != nil {
			fmt.Println("读取kafka失败:", err)
			break
		} else {
			fmt.Println("读取kafka成功: topic:", message.Topic, "partition:", message.Partition, "offset:", message.Offset, "key:", string(message.Key), "value:", string(message.Value))
		}
	}
}
